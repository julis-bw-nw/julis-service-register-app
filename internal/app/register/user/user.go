package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
}

type DataService interface {
	ClaimRegistrationKey(key string, u User) (bool, error)
	UsersToRegister() ([]User, error)
	UserByID(id int64) (User, error)
}

type LDAPService interface {
	CreateUser(user User) error
}

type Service struct {
	DataService DataService
	LDAPService LDAPService
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Get("/", s.getUsersToRegisterHandler())
	router.Post("/", s.postRegisterUserHandler())
	router.Post("/{userId}", s.postCreateUserInLDAPHandler())
	router.ServeHTTP(w, r)
}
