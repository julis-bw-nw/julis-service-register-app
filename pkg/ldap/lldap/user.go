package lldap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julis-bw-nw/julis-service-register-app/pkg/ldap"
)

func (s *service) CreateUser(user ldap.User) error {
	graphql := map[string]interface{}{
		"variables": map[string]interface{}{
			"user": map[string]interface{}{
				"id":          user.ID(),
				"email":       user.Email,
				"displayName": user.DisplayName(),
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

	s.mu.RLock()
	addr := s.host
	s.mu.RUnlock()
	url := fmt.Sprintf("http://%s/%s", addr, "api/graphql")
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bb))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ldap.ErrFailedToCreateUser
	}

	return nil
}
