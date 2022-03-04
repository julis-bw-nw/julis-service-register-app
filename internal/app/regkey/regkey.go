package regkey

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type RegistrationKey struct {
	ID                  int64
	CreatedAt           time.Time
	ClaimedAt           time.Time
	KeyValue            string
	InstantRegistration bool
}

type DataService interface {
	CreateRegistrationKey(keyValue string, instantRegistration bool) (RegistrationKey, error)
}

type Service struct {
	DataService DataService
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Post("/", s.postCreateRegisterKeyHandler())
	router.ServeHTTP(w, r)
}
