package auth

import (
	"log"
	"net/http"
	"time"
	"web-server/internal/auth"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func RefreshHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRaw, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error retrieving token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		refreshToken, err := c.DB.GetRefreshToken(r.Context(), tokenRaw)
		if err != nil {
			log.Printf("error retrieving token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		if refreshToken.RevokedAt.Valid == true {
			log.Printf("token has been revoked: %s", refreshToken.Token)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		if refreshToken.ExpiresAt.Before(time.Now()) {
			log.Printf("token has expired: %s", refreshToken.Token)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		accessToken, err := auth.CreateJWT(refreshToken.UserID, c.Env.Secret, time.Hour)
		if err != nil {
			log.Printf("error creating token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusOK,
			Body: struct {
				Token string `json:"token"`
			}{
				Token: accessToken,
			},
		})
	})
}
