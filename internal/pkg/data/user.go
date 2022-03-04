package data

import "time"

type User struct {
	ID        uint64
	CreatedAt time.Time
	FirstName string
	LastName  string
	Email     string
}

type UserService interface {
	UserByID(id int64) (User, error)
	Users() ([]User, error)
	ClaimRegistrationKey(key string, user User) (keyExists bool, err error)
}
