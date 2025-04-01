package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/identity/internal/domain/entities"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, querier models.Querier, userId uuid.UUID) (*entities.UserEntity, error)
	GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error)
	GetUserByProviderID(ctx context.Context, querier models.Querier, provider string, providerID string) (*entities.UserEntity, error)
	GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error)
	CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	CreateUserProviderData(ctx context.Context, querier models.Querier, user *entities.UserEntity, provider string, providerID string, providerValue string, email *string, providerData []byte) error
	UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity, token string, expiresAt time.Time) (*entities.UserContextEntity, error)
}
