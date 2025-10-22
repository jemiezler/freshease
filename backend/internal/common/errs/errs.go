package errs

import "errors"

var (
	NotFound         = errors.New("resource not found")
	Unauthorized     = errors.New("unauthorized")
	Conflict         = errors.New("conflict detected")
	NoFieldsToUpdate = errors.New("no fields to update")
)
