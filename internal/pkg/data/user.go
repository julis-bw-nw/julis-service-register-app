package data

import (
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/user"
	"golang.org/x/net/context"
)

func (db *Service) ClaimRegistrationKey(keyValue string, u user.Encrypted) (bool, error) {
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

	return true, db.BeginFunc(ctx, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `
INSERT INTO unregistered_users
(registration_key_id, first_name, last_name, email, password)
VALUSE ($1, $2, $3, $4, $5);
		`, keyId, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, `
INSERT INTO registration_keys
(claimed_at)
VALUSE (CURRENT_TIMESTAMP)
WHERE id = $1;`, keyId); err != nil {
			return err
		}
		return nil
	})
}
