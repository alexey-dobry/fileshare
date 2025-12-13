package model

import (
	"github.com/alexey-dobry/fileshare/pkg/validator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Credentials struct {
	gorm.Model
	UUID         uuid.UUID `gorm:"not null" validate:"required"`
	Email        string    `gorm:"not null,uniqueIndex" validate:"required"`
	PasswordHash string    `gorm:"not null" validate:"required"`
	Role         string    `gorm:"not null" validate:"required"`
}

func (u *Credentials) Validate() error {
	return validator.V.Struct(u)
}
