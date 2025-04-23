package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("no authorization header")
	}

	value := strings.Split(auth, " ")
	if len(value) != 2 {
		return "", errors.New("invalid authorization header")
	}
	if value[0] != "ApiKey" {
		return "", errors.New("invalid authorization header")
	}

	return value[1], nil
}
