package users

import "github.com/google/uuid"

type CreateUserDTO struct {
	ID          uuid.UUID `json:"id" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=8,max=100"`
	Name        string    `json:"name"  validate:"required,min=2,max=100"`
	Phone       *string   `json:"phone" validate:"omitempty,min=10,max=20"`
	Bio         *string   `json:"bio" validate:"omitempty,min=1,max=500"`
	Avatar      *string   `json:"avatar" validate:"omitempty,min=10,max=200"`
	Cover       *string   `json:"cover" validate:"omitempty,min=10,max=200"`
	DateOfBirth *string   `json:"date_of_birth" validate:"omitempty"`
	Sex         *string   `json:"sex" validate:"omitempty,oneof=male female other"`
	Goal        *string   `json:"goal" validate:"omitempty,oneof=maintenance weight_loss weight_gain"`
	HeightCm    *float64  `json:"height_cm" validate:"omitempty,min=50,max=300"`
	WeightKg    *float64  `json:"weight_kg" validate:"omitempty,min=20,max=500"`
	Status      *string   `json:"status" validate:"omitempty"`
}

type UpdateUserDTO struct {
	ID          uuid.UUID `json:"id"`
	Email       *string   `json:"email" validate:"omitempty,email"`
	Password    *string   `json:"password" validate:"omitempty,min=8,max=100"`
	Name        *string   `json:"name"  validate:"omitempty,min=2,max=100"`
	Phone       *string   `json:"phone" validate:"omitempty,min=10,max=20"`
	Bio         *string   `json:"bio" validate:"omitempty,min=1,max=500"`
	Avatar      *string   `json:"avatar" validate:"omitempty,min=10,max=200"`
	Cover       *string   `json:"cover" validate:"omitempty,min=10,max=200"`
	DateOfBirth *string   `json:"date_of_birth" validate:"omitempty"`
	Sex         *string   `json:"sex" validate:"omitempty,oneof=male female other"`
	Goal        *string   `json:"goal" validate:"omitempty,oneof=maintenance weight_loss weight_gain"`
	HeightCm    *float64  `json:"height_cm" validate:"omitempty,min=50,max=300"`
	WeightKg    *float64  `json:"weight_kg" validate:"omitempty,min=20,max=500"`
	Status      *string   `json:"status" validate:"omitempty"`
}

type GetUserDTO struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Phone       *string   `json:"phone"`
	Bio         *string   `json:"bio"`
	Avatar      *string   `json:"avatar"`
	Cover       *string   `json:"cover"`
	DateOfBirth *string   `json:"date_of_birth"`
	Sex         *string   `json:"sex"`
	Goal        *string   `json:"goal"`
	HeightCm    *float64  `json:"height_cm"`
	WeightKg    *float64  `json:"weight_kg"`
	Status      string    `json:"status"`
}
