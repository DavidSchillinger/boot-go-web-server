package chirps

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"web-server/internal/auth"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func DeleteByIDHandler(c *global.Config, pathField string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := auth.AuthenticateRequest(c, r.Header)
		if err != nil {
			log.Printf("error authenticating request: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		id, err := uuid.Parse(r.PathValue(pathField))
		if err != nil {
			log.Printf("error parsing chirp ID: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		chirp, err := c.DB.GetChirpById(r.Context(), id)
		if err != nil {
			log.Printf("error retrieving chirp: %s", err)
			respond.WithStatus(w, http.StatusNotFound)
			return
		}

		if chirp.UserID != userID {
			respond.WithStatus(w, http.StatusForbidden)
			return
		}

		if err := c.DB.DeleteChirpById(r.Context(), id); err != nil {
			log.Printf("error deleting chirp: %s", err)
			respond.WithStatus(w, http.StatusNotFound)
			return
		}

		respond.WithStatus(w, http.StatusNoContent)
	})
}
