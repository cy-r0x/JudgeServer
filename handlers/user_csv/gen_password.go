package usercsv

import (
	"crypto/rand"
	"encoding/base64"
)

func generatePassword() string {
	const length = 8

	raw := make([]byte, length)
	_, err := rand.Read(raw)
	if err != nil {
		return "ChangeMe123!"
	}

	return base64.RawURLEncoding.EncodeToString(raw)
}
