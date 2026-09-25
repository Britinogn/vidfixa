package middleware

import (
	"net/http"

	"github.com/go-chi/cors"

	"gitlab.com/britinogn/vidfixa/config"
)

// CORS defines which frontend can call this API and with what
// methods/headers. Origin comes from CLIENT_URL in .env, so it's
// swapped per environment (dev/prod) without touching code.
func CORS(cfg *config.Config) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{cfg.ClientURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
}
