package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterDTO struct {
	DataDTO
	RegisterKey string `json:"registerKey"`
}

type DataDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (dto RegisterDTO) Validate() error {
	fieldRules := []*validation.FieldRules{
		validation.Field(&dto.FirstName, validation.Required, is.Alpha),
		validation.Field(&dto.LastName, validation.Required, is.Alpha),
		validation.Field(&dto.Email, validation.Required, is.Email),
		// TODO: Password complexity validation
		validation.Field(&dto.Password, validation.Required, validation.Length(8, 40)),
	}

	return validation.ValidateStruct(&dto, fieldRules...)
}
