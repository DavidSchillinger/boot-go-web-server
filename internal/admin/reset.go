package admin

import (
	"log"
	"net/http"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func ResetHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.Env.Platform != "dev" {
			respond.WithStatus(w, http.StatusForbidden)
			return
		}

		c.FileserverHits.Store(0)

		err := c.DB.DeleteUsers(r.Context())
		if err != nil {
			log.Printf("Error deleting users: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
