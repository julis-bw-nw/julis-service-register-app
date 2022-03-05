package regkey

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
)

type Service struct {
	DataService data.Service
}

func (s Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := chi.NewRouter()
	router.Post("/", s.postCreateRegisterKeyHandler())
	router.ServeHTTP(w, r)
}

func generateRandomKey() (string, error) {
	keyBytes := make([]byte, 4)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}
