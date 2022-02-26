package ldap

import (
	"net/http"

	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/user"
)

type DataService interface {
	UsersToRegister() ([]user.User, error)
}

type Service struct {
	DataService DataService
	Client      http.Client
}

func (s Service) RegisterUser(user user.User) error {
	return nil
}

func (s Service) createUser(user user.User) error {
	graphql := map[string]interface{}{
		"variables": map[string]interface{}{
			"user": map[string]interface{}{
				"id":          "hschlehlein",
				"email":       "hendrik.schlehlein@julis.de",
				"displayName": "Hendrik Schlehlein",
				"firstName":   "Hendrik",
				"lastName":    "Schlehlein",
			},
		},
		"query":         "mutation CreateUser($user: CreateUserInput!) {\n  createUser(user: $user) {\n    id\n    creationDate\n  }\n}\n",
		"operationName": "CreateUser",
	}

	req, err := http.NewRequest("")
}
