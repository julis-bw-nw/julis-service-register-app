package data

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrEmailAlreadyUsed = errors.New("email was already used")

type UserService interface {
	UserByID(id int64) (User, error)
	Users() ([]User, error)
	RegisterUser(user *User, regKeyValue string) error
	ApproveUser(user *User) error
}

type User struct {
	gorm.Model
	RegisterKeyID uint
	RegisterKey   RegisterKey
	ApprovedAt    *time.Time
	FirstName     string
	LastName      string
	Email         string `gorm:"unique"`
}

func (user *User) AfterCreate(tx *gorm.DB) error {
	v, ok := tx.Get("instantRegistration")
	if !ok {
		return nil
	}

	isInstantRegistration, ok := v.(bool)
	if !ok {
		return nil
	}

	if !isInstantRegistration {
		return nil
	}

	return tx.Model(user).Update("approved_at", time.Now()).Error
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
		var regKey RegisterKey
		result := tx.Raw(`
SELECT *
FROM register_keys k
WHERE k.key_value = ?
AND max_claims > (
	SELECT COUNT(u.id)
	FROM users u
	WHERE k.id = u.register_key_id
);`, regKeyValue).Scan(&regKey)
		if result.Error != nil {
			return result.Error
		} else if result.RowsAffected == 0 {
			return ErrRecordNotFound
		}

		if tx.First(&User{}, "email = ?", user.Email).RowsAffected > 0 {
			return ErrEmailAlreadyUsed
		}

		user.RegisterKeyID = regKey.ID
		return db.Set("instantRegistration", regKey.InstantRegistration).Create(user).Error
	})
}

func (db service) ApproveUser(user *User) error {
	return db.Model(user).Update("approved_at", time.Now()).Error
}
