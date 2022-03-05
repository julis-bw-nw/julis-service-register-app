package lldap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/julis-bw-nw/julis-service-register-app/pkg/ldap"
)

type service struct {
	client *http.Client

	mu   sync.RWMutex
	host string
}

func New(c *http.Client, host string, options ...Option) ldap.Service {
	for _, opt := range options {
		opt(c, host)
	}

	return &service{
		client: c,
		host:   host,
	}
}

type Option func(c *http.Client, addr string)

func WithAuthenticatorTransport(username, password string) Option {
	return func(c *http.Client, addr string) {
		at := authenticatorTransport{
			client: &http.Client{
				Timeout: c.Timeout,
			},
			addr:     addr,
			username: username,
			password: password,
		}
		c.Transport = &at
	}
}

type simpleLoginDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponseDTO struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type refreshResponseDTO struct {
	Token string `json:"token"`
}

type authenticatorTransport struct {
	client *http.Client

	mu           sync.RWMutex
	addr         string
	username     string
	password     string
	token        string
	refreshToken string
}

func (at *authenticatorTransport) authenticate() (string, error) {
	at.mu.RLock()
	addr, username, password := at.addr, at.username, at.password
	at.mu.RUnlock()

	url := fmt.Sprintf("http://%s/%s", addr, "auth/simple/login")
	creds, err := json.Marshal(simpleLoginDTO{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(creds))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := at.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ldap.ErrAuthenticationFailed
	}

	var tokens loginResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&tokens); err != nil {
		return "", err
	}

	at.mu.Lock()
	at.token, at.refreshToken = tokens.Token, tokens.RefreshToken
	at.mu.Unlock()
	return tokens.Token, nil
}

func (at *authenticatorTransport) refresh() (string, error) {
	at.mu.RLock()
	addr, refreshToken := at.addr, at.refreshToken
	at.mu.RUnlock()

	url := fmt.Sprintf("http://%s/%s", addr, "auth/refresh")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("refresh-token", refreshToken)

	resp, err := at.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to refresh token")
	}

	var token refreshResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	at.mu.Lock()
	at.token = token.Token
	at.mu.Unlock()
	return token.Token, nil
}

func (at *authenticatorTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	at.mu.RLock()
	token, refreshToken := at.token, at.refreshToken
	at.mu.RUnlock()

	var err error
	if token == "" || refreshToken == "" {
		token, err = at.authenticate()
		if err != nil {
			return nil, err
		}
	} else {
		token, err = at.refresh()
		if err != nil {
			return nil, err
		}
	}

	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	return http.DefaultTransport.RoundTrip(r)
}
