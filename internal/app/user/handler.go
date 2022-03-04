package user

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
)

func (s Service) getUsersToRegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.DataService.Users()
		if err != nil {
			log.Printf("failed to get users form data service: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dtos := make([]userFullDTO, len(users))
		for i := range dtos {
			dtos[i] = mapUserFull(users[i])
		}

		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("failed to marshal users to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

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

		user := data.User{
			FirstName: dto.FirstName,
			LastName:  dto.LastName,
			Email:     dto.Email,
		}

		instantRegistration, err := s.DataService.ClaimRegistrationKey(dto.RegistrationKey, user)
		if err != nil {
			log.Printf("failed to claim registration key: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if instantRegistration {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s Service) postCreateUserInLDAPHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := strconv.ParseInt(chi.URLParam(r, "userId"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user, err := s.DataService.UserByID(userId)
		if err != nil {
			log.Printf("failed to query user with id %q: %s", userId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := s.LDAPService.CreateUser(user); err != nil {
			log.Printf("failed to create user with id %q in LDAP: %s", userId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
