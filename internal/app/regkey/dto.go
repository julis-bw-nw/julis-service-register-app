package regkey

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
)

type registerKeyDTO struct {
	ID                  uint      `json:"id"`
	CreatedAt           time.Time `json:"createdAt"`
	MaxClaims           int       `json:"maxClaims"`
	KeyValue            string    `json:"keyValue"`
	InstantRegistration bool      `json:"instantRegistration"`
}

func (dto registerKeyDTO) validate() error {
	fieldRules := []*validation.FieldRules{
		validation.Field(&dto.MaxClaims, validation.Required, validation.Min(1)),
		validation.Field(&dto.KeyValue, is.Alphanumeric),
	}

	return validation.ValidateStruct(&dto, fieldRules...)
}

func mapRegisterKeyDTOToData(dto registerKeyDTO) data.RegisterKey {
	return data.RegisterKey{
		MaxClaims:           uint(dto.MaxClaims),
		KeyValue:            dto.KeyValue,
		InstantRegistration: dto.InstantRegistration,
	}
}

func mapRegisterKeyDataToDTO(regKey data.RegisterKey) registerKeyDTO {
	return registerKeyDTO{
		ID:                  regKey.ID,
		CreatedAt:           regKey.CreatedAt,
		MaxClaims:           int(regKey.MaxClaims),
		KeyValue:            regKey.KeyValue,
		InstantRegistration: regKey.InstantRegistration,
	}
}
