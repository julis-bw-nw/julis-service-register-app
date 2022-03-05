package data

import (
	"gorm.io/gorm"
)

type UserService interface {
	UserByID(id int64) (User, error)
	Users() ([]User, error)
	RegisterUser(user *User, regKeyValue string) error
}

type User struct {
	gorm.Model
	RegisterKeyID uint
	RegisterKey   RegisterKey
	FirstName     string
	LastName      string
	Email         string
}

func (db service) UserByID(id int64) (User, error) {
	var user User
	return user, db.First(&user, id).Error
}

func (db service) Users() ([]User, error) {
	var users []User
	return users, db.Find(&users).Error
}

func (db service) RegisterUser(user *User, regKeyValue string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := tx.Raw(`
SELECT id
FROM register_keys k
WHERE k.key_value = ?
AND max_claims > (
	SELECT COUNT(u.id)
	FROM users u
	WHERE k.id = u.register_key_id
);`, regKeyValue).Scan(&user.RegisterKeyID)
		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return ErrRecordNotFound
		}

		return db.Create(user).Error
	})
}
