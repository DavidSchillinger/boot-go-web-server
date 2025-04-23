package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization header")
	}

	value := strings.Split(auth, " ")
	if len(value) != 2 {
		return "", errors.New("invalid authorization header")
	}
	if value[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return value[1], nil
}
