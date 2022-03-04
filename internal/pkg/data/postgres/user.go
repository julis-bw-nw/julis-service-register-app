package postgres

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/julis-bw-nw/julis-service-register-app/internal/pkg/data"
	"golang.org/x/net/context"
)

func (s *Service) UserByID(id int64) (data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	var user data.User
	if err := s.QueryRow(ctx, `
SELECT first_name, last_name, email
FROM unregistered_users
WHERE id = $1;
`, id).Scan(&user.FirstName, &user.LastName, &user.Email); err != nil {
		return user, err
	}
	return user, nil
}

func (s *Service) Users() ([]data.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	rows, err := s.Query(ctx, `
SELECT first_name, last_name, email
FROM unregistered_users;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []data.User
	for rows.Next() {
		var user data.User
		if err := rows.Scan(&user.FirstName, &user.LastName, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *Service) ClaimRegistrationKey(keyValue string, user data.User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	var keyId int64
	var instantRegistration bool
	if err := s.QueryRow(ctx, `
SELECT id, instant_registration
FROM registration_keys
WHERE key_value = $1
AND claimed_at IS NULL;
`, keyValue).Scan(&keyId, &instantRegistration); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	// Transaction to ensure that a registration key can only be claimed by one user
	return true, s.BeginFunc(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
INSERT INTO unregistered_users
(registration_key_id, first_name, last_name, email)
VALUES ($1, $2, $3, $4);
`, keyId, user.FirstName, user.LastName, user.Email); err != nil {
			return err
		}

		_, err := tx.Exec(ctx, `
UPDATE registration_keys
SET claimed_at = CURRENT_TIMESTAMP
WHERE id = $1;
`, keyId)
		return err
	})
}
