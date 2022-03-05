package data

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

type Service interface {
	UserService
	RegisterKeyService
}

type service struct {
	*gorm.DB
}

func NewPostgres(dsn string) (Service, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &service{
		DB: db,
	}, db.AutoMigrate(&User{}, &RegisterKey{})
}
