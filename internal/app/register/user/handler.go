package user

import (
	"encoding/json"
	"net/http"
)

func (s Service) postRegisterUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto registerDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := dto.validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		user := User{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Email:     dto.Email,
		}

		keyExists, err := s.DataService.ClaimRegistrationKey(dto.RegistrationKey, user)
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

func (s Service) getRegisteredUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}
