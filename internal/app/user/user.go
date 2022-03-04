package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/ldap"
)

type Service struct {
	DataService data.UserService
	LDAPService ldap.Service
}

func (s Service) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", s.getUsersToRegisterHandler())
	r.Post("/", s.postRegisterUserHandler())
	r.Post("/{userId}", s.postCreateUserInLDAPHandler())
	return r
}
