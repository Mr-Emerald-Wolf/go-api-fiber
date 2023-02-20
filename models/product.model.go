package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Product struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id,omitempty"`
	Name      string    `gorm:"varchar(255);uniqueIndex;not null" json:"name,omitempty"`
	Price     int       `gorm:"not null" json:"price,omitempty"`
	Category  string    `gorm:"varchar(100)" json:"category,omitempty"`
	CreatedAt time.Time `gorm:"not null" json:"createdAt,omitempty"`
	UpdatedAt time.Time `gorm:"not null" json:"updatedAt,omitempty"`
}

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct[T any](payload T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

type CreateProductSchema struct {
	Name      string `json:"Name" validate:"required"`
	Price     int    `json:"price" validate:"required"`
	Category  string `json:"category,omitempty"`
	Published bool   `json:"published,omitempty"`
}

type UpdateProductSchema struct {
	Name      string `json:"Name,omitempty"`
	Price     int    `json:"price,omitempty"`
	Category  string `json:"category,omitempty"`
	Published *bool  `json:"published,omitempty"`
}
