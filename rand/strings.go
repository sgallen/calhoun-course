package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const rememberTokenBytes = 32

func genRandBytes(n int) ([]byte, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func String(n int) (string, error) {
	b, err := genRandBytes(n)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func RememberToken() (string, error) {
	return String(rememberTokenBytes)
}
