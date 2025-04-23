package users

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"web-server/internal/auth"
	"web-server/internal/database"
	"web-server/internal/global"
	"web-server/internal/respond"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func CreateHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding user: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		password, err := auth.HashPassword(params.Password)
		if err != nil {
			log.Printf("Error hashing password: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		user, err := c.DB.CreateUser(r.Context(), database.CreateUserParams{
			Email:          params.Email,
			HashedPassword: password,
		})
		if err != nil {
			log.Printf("Error creating user: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusCreated,
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
