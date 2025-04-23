package users

import (
	"encoding/json"
	"log"
	"net/http"
	"web-server/internal/auth"
	"web-server/internal/database"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func UpdateHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.AuthenticateRequest(c, r.Header)
		if err != nil {
			log.Printf("error authenticating request: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		params := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&params); err != nil {
			log.Printf("error decoding user: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		password, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("error hashing password: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		user, err := c.DB.UpdateUserByID(r.Context(), database.UpdateUserByIDParams{
			ID:             userID,
			Email:          params.Email,
			HashedPassword: password,
		})
		if err != nil {
			log.Printf("error creating user: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusOK,
			Body: User{
				ID:          user.ID,
				CreatedAt:   user.CreatedAt,
				UpdatedAt:   user.UpdatedAt,
				Email:       user.Email,
				IsChirpyRed: user.IsChirpyRed,
			},
		})
	})
}
