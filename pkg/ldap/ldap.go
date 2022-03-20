package ldap

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrAuthenticationFailed = errors.New("authentication failed")
	ErrFailedToCreateUser   = errors.New("failed to create user")
)

type Service interface {
	CreateUser(user User) error
}

type User struct {
	FirstName string
	LastName  string
	Email     string
}

func (user User) ID() string {
	return strings.ToLower(string(user.FirstName[0]) + user.LastName)
}

func (user User) DisplayName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}
