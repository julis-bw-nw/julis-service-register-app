package ldap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/user"
)

var ErrResponseNotOK = errors.New("response is not 200 OK")

type DataService interface {
	UsersToRegister() ([]user.User, error)
}

type Service struct {
	DataService DataService
	Client      http.Client
}

func (s Service) CreateUser(url string, user user.User) error {
	graphql := map[string]interface{}{
		"variables": map[string]interface{}{
			"user": map[string]interface{}{
				"id":          userId(user),
				"email":       user.Email,
				"displayName": displayName(user),
				"firstName":   user.FirstName,
				"lastName":    user.LastName,
			},
		},
		"query":         "mutation CreateUser($user: CreateUserInput!) {\n  createUser(user: $user) {\n    id\n    creationDate\n  }\n}\n",
		"operationName": "CreateUser",
	}

	bb, err := json.Marshal(graphql)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bb))
	if err != nil {
		return err
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrResponseNotOK
	}

	return nil
}

func userId(user user.User) string {
	return strings.ToLower(string(user.FirstName[0]) + user.LastName)
}

func displayName(user user.User) string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}
