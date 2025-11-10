package middleware

import (
	"encoding/json"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// BindAndValidate binds JSON into dto and validates tags.
func BindAndValidate(c *fiber.Ctx, dto any) error {
	if err := c.BodyParser(dto); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	if err := validate.Struct(dto); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": err.Error()})
	}
	return nil
}

// BindMultipartForm binds multipart/form-data into dto and validates tags.
// It handles both form fields and file uploads.
// For JSON fields in multipart forms, it looks for a "payload" field containing JSON.
// For file uploads, it extracts the file with the given field name.
// If allowEmptyPayload is true, it will not require a payload field (useful for updates with just images).
func BindMultipartForm(c *fiber.Ctx, dto any, fileFieldName string, allowEmptyPayload ...bool) (*multipart.FileHeader, error) {
	contentType := string(c.Request().Header.ContentType())
	allowEmpty := len(allowEmptyPayload) > 0 && allowEmptyPayload[0]
	
	// Check if this is multipart/form-data
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		// If not multipart, try standard JSON binding
		if err := BindAndValidate(c, dto); err != nil {
			return nil, err
		}
		return nil, nil
	}

	var file *multipart.FileHeader
	var err error

	// Extract file if fileFieldName is provided
	if fileFieldName != "" {
		file, err = c.FormFile(fileFieldName)
		if err != nil {
			// File is optional - only return error if it's a real parsing error
			if !strings.Contains(err.Error(), "no such file") && 
			   !strings.Contains(err.Error(), "there is no uploaded file") &&
			   !strings.Contains(err.Error(), "bad request") {
				return nil, fiber.NewError(fiber.StatusBadRequest, "failed to parse file: "+err.Error())
			}
			file = nil
		}
		
		// Validate file size if file is provided
		if file != nil && file.Size == 0 {
			return nil, fiber.NewError(fiber.StatusBadRequest, "file is empty")
		}
	}

	// Try to parse JSON payload from "payload" field first
	jsonStr := c.FormValue("payload")
	if jsonStr != "" {
		// Decode JSON payload
		if err := json.Unmarshal([]byte(jsonStr), dto); err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid JSON payload: "+err.Error())
		}
		// Validate the DTO if payload was provided
		if err := validate.Struct(dto); err != nil {
			return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "validation failed: "+err.Error())
		}
	} else {
		// Fallback: try to parse all form fields directly into DTO
		// This allows direct form field binding (name, price, etc.)
		if err := c.BodyParser(dto); err != nil {
			// If allowEmptyPayload is true and we have a file, that's okay
			if allowEmpty && file != nil {
				// Just return the file, DTO can remain as-is (useful for image-only updates)
				return file, nil
			}
			// If BodyParser fails and we don't allow empty payload, return error
			if !allowEmpty {
				return nil, fiber.NewError(fiber.StatusBadRequest, "either provide 'payload' JSON field or valid form fields: "+err.Error())
			}
		} else {
			// Validate the DTO if we parsed form fields
			if err := validate.Struct(dto); err != nil {
				return nil, fiber.NewError(fiber.StatusUnprocessableEntity, "validation failed: "+err.Error())
			}
		}
	}

	return file, nil
}
