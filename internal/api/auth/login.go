package auth

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"web-server/internal/api/users"
	"web-server/internal/auth"
	"web-server/internal/database"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func LoginHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("error decoding user: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		user, err := c.DB.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			log.Printf("error retrieving user: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		if err := auth.VerifyPasswordHash(user.HashedPassword, params.Password); err != nil {
			log.Printf("error verifying password: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		accessToken, err := auth.CreateJWT(user.ID, c.Env.Secret, time.Hour)
		if err != nil {
			log.Printf("error creating token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		refreshTokenRaw, err := auth.CreateRefreshToken()
		if err != nil {
			log.Printf("error creating refresh token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		refreshToken, err := c.DB.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:  refreshTokenRaw,
			UserID: user.ID,
		})
		if err != nil {
			log.Printf("error creating refresh token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusOK,
			Body: struct {
				users.User
				Token        string `json:"token"`
				RefreshToken string `json:"refresh_token"`
			}{
				User: users.User{
					ID:          user.ID,
					CreatedAt:   user.CreatedAt,
					UpdatedAt:   user.UpdatedAt,
					Email:       user.Email,
					IsChirpyRed: user.IsChirpyRed,
				},
				Token:        accessToken,
				RefreshToken: refreshToken.Token,
			},
		})
	})
}
