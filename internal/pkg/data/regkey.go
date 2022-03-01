package data

import (
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/regkey"
	"golang.org/x/net/context"
)

func (s *Service) CreateRegistrationKey(keyValue string, instantRegistration bool) (regkey.RegistrationKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	regKey := regkey.RegistrationKey{
		KeyValue:            keyValue,
		InstantRegistration: instantRegistration,
	}

	if err := s.QueryRow(ctx, `
INSERT INTO registration_keys
(key_value, instant_registration)
VALUES ($1, $2)
RETURNING id, created_at;
`, keyValue, instantRegistration).Scan(&regKey.ID, &regKey.CreatedAt); err != nil {
		return regKey, err
	}

	return regKey, nil
}
