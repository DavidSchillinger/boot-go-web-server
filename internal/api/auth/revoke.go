package auth

import (
	"log"
	"net/http"
	"web-server/internal/auth"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func RevokeHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenRaw, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("error retrieving token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		if _, err := c.DB.RevokeRefreshToken(r.Context(), tokenRaw); err != nil {
			log.Printf("error revoking token: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		respond.WithJSON(w, respond.WithJSONOptions{
			Code: http.StatusNoContent,
			Body: nil,
		})
	})
}
