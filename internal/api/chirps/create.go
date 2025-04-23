package chirps

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
	"web-server/internal/auth"
	"web-server/internal/database"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func CreateHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.AuthenticateRequest(c, r.Header)
		if err != nil {
			log.Printf("error authenticating request: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := struct {
			Body string `json:"body"`
		}{}
		if err := decoder.Decode(&params); err != nil {
			log.Printf("error decoding chirp: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		if len(params.Body) > 140 {
			respond.WithError(w, respond.WithErrorOptions{
				Code:    http.StatusBadRequest,
				Message: "chirp is too long",
			})
			return
		}

		split := strings.Split(params.Body, " ")
		var cleaned []string
		for _, sub := range split {
			lower := strings.ToLower(sub)
			switch lower {
			case "kerfuffle":
				fallthrough
			case "sharbert":
				fallthrough
			case "fornax":
				cleaned = append(cleaned, "****")
			default:
				cleaned = append(cleaned, sub)
			}
		}

		chirp, err := c.DB.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   params.Body,
			UserID: userID,
		})
		if err != nil {
			log.Printf("error creating chirp: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusCreated,
			Body: struct {
				ID        uuid.UUID `json:"id"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Body      string    `json:"body"`
				UserID    uuid.UUID `json:"user_id"`
			}{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			},
		})
	})
}
