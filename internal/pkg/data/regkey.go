package data

import (
	"gorm.io/gorm"
)

type RegisterKeyService interface {
	RegisterKeys() ([]RegisterKey, error)
	CreateRegisterKey(regKey *RegisterKey) error
}

type RegisterKey struct {
	gorm.Model
	Users               []User
	MaxClaims           uint   `gorm:"default=1"`
	KeyValue            string `gorm:"unique"`
	InstantRegistration bool
}

func (db service) RegisterKeys() ([]RegisterKey, error) {
	var regKeys []RegisterKey
	return regKeys, db.Preload("Users").Find(&regKeys).Error
}

func (db service) CreateRegisterKey(regKey *RegisterKey) error {
	return db.Create(regKey).Error
}
