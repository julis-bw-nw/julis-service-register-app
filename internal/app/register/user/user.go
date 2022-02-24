package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type RegisterDTO struct {
	DTO
	RegistrationKey string `json:"registrationKey"`
}

type DTO struct {
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

type Encrypted struct {
	FirstName []byte
	LastName  []byte
	Email     []byte
	Password  []byte
}

type EncryptionService interface {
	Encrypt(data string) []byte
}

type DataService interface {
	ClaimRegistrationKey(keyValue string, data Encrypted) (bool, error)
}

type Service struct {
	DataService       DataService
	EncryptionService EncryptionService
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Post("/", s.postRegisterUserHandler())
	router.ServeHTTP(w, r)
}

func (s Service) EncryptDTO(dto DTO) (Encrypted, error) {
	return Encrypted{
		FirstName: s.EncryptionService.Encrypt(dto.FirstName),
		LastName:  s.EncryptionService.Encrypt(dto.LastName),
		Email:     s.EncryptionService.Encrypt(dto.Email),
		Password:  s.EncryptionService.Encrypt(dto.Password),
	}, nil
}
