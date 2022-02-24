package user

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s Service) postRegisterUserHandler() http.HandlerFunc {
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

		encrypted, err := s.EncryptDTO(dto.DTO)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		keyExists, err := s.DataService.ClaimRegistrationKey(dto.RegistrationKey, encrypted)
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
