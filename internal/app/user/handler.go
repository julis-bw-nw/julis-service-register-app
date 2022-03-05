package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
)

func (s Service) getUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := s.DataService.Users()
		if err != nil {
			log.Printf("failed to get users form data service: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dtos := make([]userDTO, len(users))
		for i := range dtos {
			dtos[i] = mapUserDataToDTO(users[i])
		}

		if err := json.NewEncoder(w).Encode(dtos); err != nil {
			log.Printf("failed to marshal users to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s Service) postCreateUserHandler() http.HandlerFunc {
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

		user := mapRegisterDTOToData(dto)
		if err := s.DataService.RegisterUser(&user, dto.RegistrationKey); err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				log.Printf("failed to claim registration key: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		userDTO := mapUserDataToDTO(user)
		if err := json.NewEncoder(w).Encode(&userDTO); err != nil {
			log.Printf("failed to marshal user to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func (s Service) postRegisterUserInLDAPHandler() http.HandlerFunc {
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

		ldapUser := mapUserDataToLDAP(user)
		if err := s.LDAPService.CreateUser(ldapUser); err != nil {
			log.Printf("failed to create user with id %q in LDAP: %s", userId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
