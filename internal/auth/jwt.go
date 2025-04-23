package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"time"
	"web-server/internal/global"
)

func CreateJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String(),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(
	tokenString,
	tokenSecret string,
) (uuid.UUID, error) {
	key := func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil }
	tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, key)
	if err != nil {
		return uuid.Nil, err
	}

	sub, err := tok.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func AuthenticateRequest(c *global.Config, h http.Header) (uuid.UUID, error) {
	token, err := GetBearerToken(h)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := ValidateJWT(token, c.Env.Secret)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
