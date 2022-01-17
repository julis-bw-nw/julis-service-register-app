package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julis-bw-nw/julis-service-register-app/backend/db"
)

type EncryptionService interface {
	EncryptUserData(u DataDTO) (db.EncryptedUserData, error)
}

type DataService interface {
	ClaimRegistrationKey(keyValue string, u db.EncryptedUserData) (bool, error)
}

type Handler struct {
	DataService       DataService
	EncryptionService EncryptionService
}

func (h Handler) PostRegisterUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto RegisterDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := dto.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		encryptedUserData, err := h.EncryptionService.EncryptUserData(dto.DataDTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		keyExists, err := h.DataService.ClaimRegistrationKey(dto.RegistrationKey, encryptedUserData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !keyExists {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
