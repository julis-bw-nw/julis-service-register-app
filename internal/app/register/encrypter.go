package register

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

type Encrypter struct {
	gcm   cipher.AEAD
	nonce []byte
}

func NewEncrypter(secret []byte) (*Encrypter, error) {
	b, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(b)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	return &Encrypter{
		gcm:   gcm,
		nonce: nonce,
	}, nil
}

func (e Encrypter) Encrypt(data string) []byte {
	return e.gcm.Seal(e.nonce, e.nonce, []byte(data), nil)
}

func (e Encrypter) Decrypt(data []byte) (string, error) {
	nonceSize := len(e.nonce)
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := e.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
