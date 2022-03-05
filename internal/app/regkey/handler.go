package regkey

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s Service) getRegisterKeyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		regKeys, err := s.DataService.RegisterKeys()
		if err != nil {
			log.Printf("failed to get register keys form data service: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		dtos := make([]registerKeyDTO, len(regKeys))
		for i := range dtos {
			dtos[i] = mapRegisterKeyDataToDTO(regKeys[i])
		}

		if err := json.NewEncoder(w).Encode(dtos); err != nil {
			log.Printf("failed to marshal register keys to json: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

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
