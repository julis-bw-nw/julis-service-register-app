package data

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/user"
	"golang.org/x/net/context"
)

func (s *Service) OnUserRegistration(f func(user.User)) {
	
}

func (s *Service) ClaimRegistrationKey(keyValue string, user user.User) (bool, error) {
	encryptedPwd := s.EncryptionService.Encrypt(user.Password)
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	var keyId int64
	var instantRegistration bool
	if err := s.QueryRow(ctx, `
SELECT id, instant_registration
FROM registration_keys
WHERE key_value = $1
AND claimed_at IS NULL;`, keyValue).Scan(&keyId, &instantRegistration); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	// Transaction to ensure that a registration key can only be claimed by one user
	return true, s.BeginFunc(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
INSERT INTO unregistered_users
(registration_key_id, approved_at, first_name, last_name, email, password)
VALUSE ($1, CASE WHEN $2 THEN CURRENT_TIMESTAMP ELSE NULL, $3, $4, $5, $6);
`, keyId, instantRegistration, user.FirstName, user.LastName, user.Email, encryptedPwd); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, `
INSERT INTO registration_keys
(claimed_at)
VALUSE (CURRENT_TIMESTAMP)
WHERE id = $1;
`, keyId); err != nil {
			return err
		}
		return nil
	})
}

func (s *Service) UsersToRegister() ([]user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	rows, err := s.Query(ctx, `
SELECT first_name, last_name, email, password
FROM unregistered_users
WHERE approved_at IS NOT NULL;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var encryptedPwd []byte
		var user user.User
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Email, &encryptedPwd); err != nil {
			return nil, err
		}
		pwd, err := s.EncryptionService.Decrypt(encryptedPwd)
		if err != nil {
			return nil, err
		}
		user.Password = pwd
		users = append(users, user)
	}
	return users, nil
}
