package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
)

type registerDTO struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	RegistrationKey string `json:"registrationKey"`
}

func (dto registerDTO) validate() error {
	fieldRules := []*validation.FieldRules{
		validation.Field(&dto.FirstName, validation.Required, is.Alpha),
		validation.Field(&dto.LastName, validation.Required, is.Alpha),
		validation.Field(&dto.Email, validation.Required, is.Email),
	}

	return validation.ValidateStruct(&dto, fieldRules...)
}

type userFullDTO struct {
	ID        uint64    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
}

func mapUserFull(user data.User) userFullDTO {
	return userFullDTO{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
