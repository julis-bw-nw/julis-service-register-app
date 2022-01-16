package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/julis-bw-nw/julis-service-register-app/backend/db"
	"github.com/julis-bw-nw/julis-service-register-app/backend/user"
)

type encrypter struct {
	secret string
}

func (e encrypter) Encrypt(data string) (string, error) {
	block, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return "", err
	}

	b := []byte(data)
	cfb := cipher.NewCFBEncrypter(block, b)
	cipherText := make([]byte, len(b))
	cfb.XORKeyStream(cipherText, b)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (e encrypter) Decrypt(data string) (string, error) {
	block, err := aes.NewCipher([]byte(e.secret))
	if err != nil {
		return "", err
	}

	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, b)
	plainText := make([]byte, len(b))
	cfb.XORKeyStream(plainText, b)
	return string(plainText), nil
}

func (e encrypter) EncryptUserData(u user.DataDTO) (db.EncryptedUserData, error) {
	firstName, err := e.Encrypt(u.FirstName)
	if err != nil {
		return db.EncryptedUserData{}, err
	}

	lastName, err := e.Encrypt(u.LastName)
	if err != nil {
		return db.EncryptedUserData{}, err
	}

	email, err := e.Encrypt(u.Email)
	if err != nil {
		return db.EncryptedUserData{}, err
	}

	password, err := e.Encrypt(u.Password)
	if err != nil {
		return db.EncryptedUserData{}, err
	}

	return db.EncryptedUserData{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}, nil
}
