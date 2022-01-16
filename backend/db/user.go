package db

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"golang.org/x/net/context"
)

type EncryptedUserData struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (db *DB) ClaimRegistrationKey(keyValue string, u EncryptedUserData) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.Timeout)
	defer cancel()

	var keyId int64
	if err := db.QueryRow(ctx, `
SELECT id
FROM registration_keys
WHERE key_value = $1
AND claimed_at IS NULL;
`, keyValue).Scan(&keyId); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	_, err := db.Exec(ctx, `
INSERT INTO unregistered_users
(registration_key_id, first_name, last_name, email, password)
VALUSE ($1, $2, $3, $4, $5);
`, keyId, u.FirstName, u.LastName, u.Email, u.Password)
	return true, err
}
