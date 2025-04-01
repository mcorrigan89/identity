package middleware

import (
	"context"
	"net/http"

	"github.com/mcorrigan89/identity/internal/domain/entities"
)

func GetUserFromContext(ctx context.Context) *entities.UserContextEntity {
	userContext, ok := ctx.Value(currentUserContextKey).(*entities.UserContextEntity)
	if !ok {
		return nil
	}

	return userContext
}

func (m *middleware) Authorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userContext, ok := ctx.Value(currentUserContextKey).(*entities.UserContextEntity)
		if !ok {
			m.logger.Error().Ctx(ctx).Msg("User context not found in request context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if userContext.IsExpired() {
			m.logger.Error().Ctx(ctx).Msg("User context has expired")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
