package main

import (
	_ "github.com/lib/pq"
	"web-server/internal/admin"
	"web-server/internal/api"
	"web-server/internal/api/auth"
	"web-server/internal/api/chirps"
	"web-server/internal/api/polka"
	"web-server/internal/api/users"
	"web-server/internal/global"
)

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	c := global.LoadConfig()

	handler := http.NewServeMux()

	handler.Handle("GET /app/", appHandler(c))

	handler.Handle("GET /api/health", api.HealthHandler())
	handler.Handle("POST /api/login", auth.LoginHandler(c))
	handler.Handle("POST /api/refresh", auth.RefreshHandler(c))
	handler.Handle("POST /api/revoke", auth.RevokeHandler(c))
	handler.Handle("POST /api/users", users.CreateHandler(c))
	handler.Handle("PUT /api/users", users.UpdateHandler(c))
	handler.Handle("GET /api/chirps", chirps.ReadHandler(c))
	handler.Handle("GET /api/chirps/{chirpID}", chirps.ReadByIDHandler(c, "chirpID"))
	handler.Handle("DELETE /api/chirps/{chirpID}", chirps.DeleteByIDHandler(c, "chirpID"))
	handler.Handle("POST /api/chirps", chirps.CreateHandler(c))
	handler.Handle("POST /api/polka/webhooks", polka.WebhooksHandler(c))

	handler.Handle("GET /admin/metrics", admin.MetricsHandler(c))
	handler.Handle("POST /admin/reset", admin.ResetHandler(c))

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	log.Println(fmt.Sprintf("Starting server on port %v!", server.Addr))
	log.Fatal(server.ListenAndServe())
}

func appHandler(c *global.Config) http.Handler {
	app := http.StripPrefix("/app/", http.FileServer(http.Dir("./")))
	return c.MiddlewareMetricsInc(app)
}
