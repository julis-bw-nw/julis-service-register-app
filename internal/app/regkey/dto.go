package regkey

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type registrationKeyDTO struct {
	ID                  int64     `json:"id"`
	CreatedAt           time.Time `json:"createdAt"`
	ClaimedAt           time.Time `json:"claimedAt"`
	KeyValue            string    `json:"keyValue"`
	InstantRegistration bool      `json:"instantRegistration"`
}

func (dto registrationKeyDTO) validate() error {
	fieldRules := []*validation.FieldRules{
		validation.Field(&dto.KeyValue, is.Alphanumeric),
	}

	return validation.ValidateStruct(&dto, fieldRules...)
}
