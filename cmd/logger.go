package main

import (
	"os"
	"time"

	"github.com/mcorrigan89/identity/internal/infrastructure/middleware"
	"github.com/rs/zerolog"
)

type TracingHook struct{}

func (h TracingHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	correlationId := middleware.GetCorrelationIdFromContext(ctx)
	if correlationId != "" {
		e.Str("correlation_id", correlationId)
	}
	ip := middleware.GetIPFromContext(ctx)
	if ip != "" {
		e.Str("ip_address", ip)
	}
	sessionToken := middleware.GetSessionTokenFromContext(ctx)
	if sessionToken != "" {
		e.Str("session_token", sessionToken)
	}
}

func getLogger() zerolog.Logger {

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		With().
		Timestamp().
		Caller().
		Stack().
		Logger().
		Hook(TracingHook{})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return logger

}
