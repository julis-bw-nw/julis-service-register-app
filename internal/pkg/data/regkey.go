package data

import (
	"gorm.io/gorm"
)

type RegisterKeyService interface {
	CreateRegisterKey(regKey *RegisterKey) error
}

type RegisterKey struct {
	gorm.Model
	Users               []User
	MaxClaims           uint   `gorm:"default=1"`
	KeyValue            string `gorm:"unique"`
	InstantRegistration bool
}

func (db service) CreateRegisterKey(regKey *RegisterKey) error {
	return db.Create(regKey).Error
}
