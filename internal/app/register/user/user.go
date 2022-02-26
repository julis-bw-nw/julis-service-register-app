package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

type DataService interface {
	ClaimRegistrationKey(key string, u User) (bool, error)
}

type Service struct {
	DataService DataService
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Post("/", s.postRegisterUserHandler())
	router.ServeHTTP(w, r)
}
