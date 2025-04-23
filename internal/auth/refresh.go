package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func CreateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes), nil
}
