package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/julis-bw-nw/julis-service-register-app/backend/db"
	"github.com/julis-bw-nw/julis-service-register-app/backend/user"
)

type encrypter struct {
	gcm   cipher.AEAD
	nonce []byte
}

func newEncrypter(secret []byte) (*encrypter, error) {
	c, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	return &encrypter{
		gcm:   gcm,
		nonce: nonce,
	}, nil
}

func (e encrypter) Encrypt(data string) string {
	return string(e.gcm.Seal(e.nonce, e.nonce, []byte(data), nil))
}

func (e encrypter) Decrypt(data string) (string, error) {
	nonceSize := len(e.nonce)
	b := []byte(data)
	nonce, ciphertext := b[:nonceSize], b[nonceSize:]

	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (e encrypter) EncryptUserData(u user.DataDTO) (db.EncryptedUserData, error) {
	return db.EncryptedUserData{
		FirstName: e.Encrypt(u.FirstName),
		LastName:  e.Encrypt(u.LastName),
		Email:     e.Encrypt(u.Email),
		Password:  e.Encrypt(u.Password),
	}, nil
}
