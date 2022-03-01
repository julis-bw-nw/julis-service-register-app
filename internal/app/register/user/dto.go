package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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
