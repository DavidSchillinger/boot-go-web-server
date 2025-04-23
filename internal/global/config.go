package global

import (
	"database/sql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"web-server/internal/database"
)

type Environment struct {
	Secret   string
	Platform string
	PolkaKey string
}

type Config struct {
	Env            Environment
	DB             *database.Queries
	FileserverHits atomic.Int32
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	secret := os.Getenv("SECRET")
	platform := os.Getenv("PLATFORM")
	polkaKey := os.Getenv("POLKA_KEY")
	databaseURL := os.Getenv("DATABASE_URL")

	if secret == "" || databaseURL == "" || platform == "" || polkaKey == "" {
		log.Fatal("error loading .env file")
	}

	environment := Environment{
		Secret:   secret,
		Platform: platform,
		PolkaKey: polkaKey,
	}

	return &Config{
		FileserverHits: atomic.Int32{},
		Env:            environment,
		DB:             connectToDatabase(databaseURL),
	}
}

func (c *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func connectToDatabase(dbURL string) *database.Queries {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	return database.New(db)
}
