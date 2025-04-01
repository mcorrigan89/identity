package middleware

import (
	"context"
	"net/http"

	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"
	"github.com/rs/xid"
)

type contextKey string

const (
	ipKey                 contextKey = "ip"
	correlationIDKey      contextKey = "correlation_id"
	sessionTokenKey       contextKey = "sessionTokenKey"
	currentUserContextKey contextKey = "currentUserContextKey"
)

var SessionTokenKey = "x-session-token"

func GetCorrelationIdFromContext(ctx context.Context) string {
	correlationId, ok := ctx.Value(correlationIDKey).(string)
	if !ok {
		return ""
	}
	return correlationId
}

func GetSessionTokenFromContext(ctx context.Context) string {
	sessionToken, ok := ctx.Value(sessionTokenKey).(string)
	if !ok {
		return ""
	}
	return sessionToken
}

func GetIPFromContext(ctx context.Context) string {
	ip, ok := ctx.Value(ipKey).(string)
	if !ok {
		return ""
	}
	return ip
}

func (m *middleware) ContextBuilder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		querier := models.New(m.db)

		ctx = context.WithValue(ctx, ipKey, r.RemoteAddr)
		correlationID := xid.New().String()
		ctx = context.WithValue(ctx, correlationIDKey, correlationID)

		sessionToken := r.Header.Get(SessionTokenKey)
		ctx = context.WithValue(ctx, sessionTokenKey, sessionToken)

		ctx = m.logger.WithContext(ctx)

		userContext, err := m.userService.GetUserContextBySessionToken(ctx, querier, sessionToken)
		if err == nil {
			ctx = context.WithValue(ctx, currentUserContextKey, userContext)
		}

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
