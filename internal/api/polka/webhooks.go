package polka

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"web-server/internal/auth"
	"web-server/internal/global"
	"web-server/internal/respond"
)

func WebhooksHandler(c *global.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := auth.GetAPIKey(r.Header)
		if err != nil {
			log.Printf("error authenticating request: %s", err)
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}
		if key != c.Env.PolkaKey {
			respond.WithStatus(w, http.StatusUnauthorized)
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := struct {
			Event string `json:"event"`
			Data  struct {
				UserId uuid.UUID `json:"user_id"`
			}
		}{}
		if err := decoder.Decode(&params); err != nil {
			log.Printf("error decoding chirp: %s", err)
			respond.WithStatus(w, http.StatusInternalServerError)
			return
		}

		if params.Event != "user.upgraded" {
			respond.WithStatus(w, http.StatusNoContent)
			return
		}

		if _, err := c.DB.UpgradeUserToChirpyRedByUserID(r.Context(), params.Data.UserId); err != nil {
			respond.WithStatus(w, http.StatusNotFound)
			return
		}

		respond.WithStatus(w, http.StatusNoContent)
	})
}
