package regkey

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s Service) postCreateRegisterKeyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto registerKeyDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := dto.validate(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		if dto.KeyValue == "" {
			keyValue, err := generateRandomKey()
			if err != nil {
				log.Printf("failed to generate register key: %s", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			dto.KeyValue = keyValue
		}

		regKey := mapRegisterKeyDTOToData(dto)
		if err := s.DataService.CreateRegisterKey(&regKey); err != nil {
			log.Printf("failed to create regkey: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dto = mapRegisterKeyDataToDTO(regKey)
		if err := json.NewEncoder(w).Encode(&dto); err != nil {
			log.Printf("failed to marshal regkey to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
