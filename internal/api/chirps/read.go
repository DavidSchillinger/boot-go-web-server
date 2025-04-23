package chirps

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"sort"
	"time"
	"web-server/internal/database"
	"web-server/internal/global"
	"web-server/internal/respond"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}

func ReadHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sortParam := r.URL.Query().Get("sort")
		authorIDParam := r.URL.Query().Get("author_id")

		var chirps []database.Chirp
		if authorIDParam != "" {
			authorID, err := uuid.Parse(authorIDParam)
			if err != nil {
				log.Printf("error parsing chirp ID: %s", err)
				respond.WithStatus(w, http.StatusBadRequest)
				return
			}

			chirps, err = c.DB.GetChirpsByUserID(r.Context(), authorID)
			if err != nil {
				log.Printf("error retrieving chirps: %s", err)
				respond.WithStatus(w, http.StatusInternalServerError)
				return
			}
		} else {
			var err error
			chirps, err = c.DB.GetChirps(r.Context())
			if err != nil {
				log.Printf("error retrieving chirps: %s", err)
				respond.WithStatus(w, http.StatusInternalServerError)
				return
			}
		}

		var body []Chirp
		for _, chirp := range chirps {
			body = append(body, Chirp{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}

		switch sortParam {
		case "desc":
			sort.Slice(body, func(i, j int) bool {
				return body[i].CreatedAt.After(body[j].CreatedAt)
			})
		case "asc":
			fallthrough
		default:
			sort.Slice(body, func(i, j int) bool {
				return body[i].CreatedAt.Before(body[j].CreatedAt)
			})
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusOK,
			Body: body,
		})
	})
}

func ReadByIDHandler(c *global.Config, pathField string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue(pathField))
		if err != nil {
			log.Printf("error parsing chirp ID: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		chirp, err := c.DB.GetChirpById(r.Context(), id)
		if err != nil {
			log.Printf("error retrieving chirp: %s", err)
			respond.WithError(w, respond.WithErrorOptions{
				Code:    http.StatusNotFound,
				Message: "chirp not found",
			})
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusOK,
			Body: Chirp{
				ID:        chirp.ID,
				UserID:    chirp.UserID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
			},
		})
	})
}
