package middleware

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/mcorrigan89/identity/internal/domain/services"
	"github.com/mcorrigan89/identity/internal/infrastructure/config"
	"github.com/rs/zerolog"
)

type Middleware interface {
	ContextBuilder(next http.Handler) http.Handler
	RecoverPanic(next http.Handler) http.Handler
	EnabledCORS(next http.Handler) http.Handler
	Authorization(next http.HandlerFunc) http.HandlerFunc
}

type middleware struct {
	config      *config.Config
	logger      *zerolog.Logger
	db          *pgxpool.Pool
	userService services.UserService
}

func CreateMiddleware(config *config.Config, db *pgxpool.Pool, logger *zerolog.Logger, userService services.UserService) *middleware {
	return &middleware{config: config, logger: logger, db: db, userService: userService}
}
