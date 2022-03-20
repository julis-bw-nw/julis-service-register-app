package user

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
	"github.com/julis-bw-nw/julis-service-register-app/pkg/ldap"
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

func mapRegisterDTOToData(dto registerDTO) data.User {
	return data.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
	}
}

type userDTO struct {
	ID         uint       `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	ApprovedAt *time.Time `json:"approvedAt,omitempty"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	Email      string     `json:"email"`
}

func mapUserDataToDTO(user data.User) userDTO {
	return userDTO{
		ID:         user.ID,
		CreatedAt:  user.CreatedAt,
		ApprovedAt: user.ApprovedAt,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Email:      user.Email,
	}
}

func mapUserDataToLDAP(user data.User) ldap.User {
	return ldap.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
}
