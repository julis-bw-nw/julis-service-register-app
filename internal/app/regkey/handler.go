package regkey

import (
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

func (s Service) postCreateRegisterKeyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto registrationKeyDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := dto.validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if dto.KeyValue == "" {
			keyBytes := make([]byte, 4)
			if _, err := rand.Read(keyBytes); err != nil {
				log.Printf("failed to generate registration key: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			dto.KeyValue = hex.EncodeToString(keyBytes)
		}

		regKey, err := s.DataService.CreateRegistrationKey(dto.KeyValue, dto.InstantRegistration)
		if err != nil {
			log.Printf("failed to create regkey: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dto.ID = regKey.ID
		dto.CreatedAt = regKey.CreatedAt

		bb, err := json.Marshal(dto)
		if err != nil {
			log.Printf("failed to marshal regkey to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(bb)
		w.WriteHeader(http.StatusCreated)
	}
}
