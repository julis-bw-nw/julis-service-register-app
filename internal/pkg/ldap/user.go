package ldap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/julis-bw-nw/julis-service-register-app/internal/app/register/user"
)

var ErrResponseNotOK = errors.New("response is not 200 OK")

type Service struct {
	Client  *http.Client
	BaseURL string
}

func (s Service) CreateUser(user user.User) error {
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
		"query":         "mutation CreateUser($user: CreateUserInput!) { createUser(user: $user) { id creationDate } }",
		"operationName": "CreateUser",
	}

	bb, err := json.Marshal(graphql)
	if err != nil {
		return err
	}

	log.Println(string(bb))

	// TODO: Could cause a race condition on concurrent read of s.BaseURL
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/api/graphql", bytes.NewReader(bb))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzUxMiJ9.eyJleHAiOiIyMDIyLTAzLTAyVDE4OjMzOjA5LjUyMDQxMTg1MloiLCJpYXQiOiIyMDIyLTAzLTAxVDE4OjMzOjA5LjUyMDQxMjYyMloiLCJ1c2VyIjoiYWRtaW4iLCJncm91cHMiOlsibGxkYXBfYWRtaW4iXX0.IIfZWI7iY2gU-txGDDu8eFaY_2YwD-qf64SFJqlEOusyJp8KyvWSZtjWdxvWdyRMKGPKwzwLi8N0P6g_C1GL5g")
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		log.Println(string(data))
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
