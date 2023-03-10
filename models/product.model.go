package models

import (
	"time"

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

type CreateProductSchema struct {
	Name      string `json:"name" validate:"required"`
	Price     int    `json:"price" validate:"required"`
	Category  string `json:"category,omitempty"`
	Published bool   `json:"published,omitempty"`
}

type UpdateProductSchema struct {
	Name      string `json:"name,omitempty"`
	Price     int    `json:"price,omitempty"`
	Category  string `json:"category,omitempty"`
	Published *bool  `json:"published,omitempty"`
}
