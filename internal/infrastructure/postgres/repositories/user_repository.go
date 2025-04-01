package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mcorrigan89/identity/internal/domain/entities"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"
)

type postgresUserRepository struct {
}

func NewPostgresUserRepository() *postgresUserRepository {
	return &postgresUserRepository{}
}

func (repo *postgresUserRepository) GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	userRow, err := querier.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userAuthRows, err := querier.GetUserAuthByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	userAuthModels := make([]models.UserAuth, len(userAuthRows))
	for i, auth := range userAuthRows {
		userAuthModels[i] = models.UserAuth{
			ID:            auth.UserAuth.ID,
			ProviderID:    auth.UserAuth.ProviderID,
			Provider:      auth.UserAuth.Provider,
			Email:         auth.UserAuth.Email,
			EmailVerified: auth.UserAuth.EmailVerified,
			CreatedAt:     auth.UserAuth.CreatedAt,
		}
	}

	return entities.NewUserEntity(userRow.User, userAuthModels), nil
}

func (repo *postgresUserRepository) GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserBySessionToken(ctx, sessionToken)
	if err != nil {
		return nil, err
	}

	userAuthRows, err := querier.GetUserAuthByUserID(ctx, row.User.ID)
	if err != nil {
		return nil, err
	}

	userAuthModels := make([]models.UserAuth, len(userAuthRows))
	for i, auth := range userAuthRows {
		userAuthModels[i] = models.UserAuth{
			ID:            auth.UserAuth.ID,
			ProviderID:    auth.UserAuth.ProviderID,
			Provider:      auth.UserAuth.Provider,
			Email:         auth.UserAuth.Email,
			EmailVerified: auth.UserAuth.EmailVerified,
			CreatedAt:     auth.UserAuth.CreatedAt,
		}
	}

	userEntity := entities.NewUserEntity(row.User, userAuthModels)

	return entities.NewUserContextEntity(userEntity, row.UserSession), nil
}

func (repo *postgresUserRepository) GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserByEmail(ctx, &email)
	if err != nil {
		return nil, err
	}

	userAuthRows, err := querier.GetUserAuthByUserID(ctx, row.User.ID)
	if err != nil {
		return nil, err
	}

	userAuthModels := make([]models.UserAuth, len(userAuthRows))
	for i, auth := range userAuthRows {
		userAuthModels[i] = models.UserAuth{
			ID:            auth.UserAuth.ID,
			ProviderID:    auth.UserAuth.ProviderID,
			Provider:      auth.UserAuth.Provider,
			Email:         auth.UserAuth.Email,
			EmailVerified: auth.UserAuth.EmailVerified,
			CreatedAt:     auth.UserAuth.CreatedAt,
		}
	}

	return entities.NewUserEntity(row.User, userAuthModels), nil
}

func (repo *postgresUserRepository) GetUserByProviderID(ctx context.Context, querier models.Querier, provider string, providerID string) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.GetUserByProviderID(ctx, models.GetUserByProviderIDParams{
		Provider:   provider,
		ProviderID: providerID,
	})
	if err != nil {
		return nil, err
	}

	userAuthRows, err := querier.GetUserAuthByUserID(ctx, row.User.ID)
	if err != nil {
		return nil, err
	}

	userAuthModels := make([]models.UserAuth, len(userAuthRows))
	for i, auth := range userAuthRows {
		userAuthModels[i] = models.UserAuth{
			ID:            auth.UserAuth.ID,
			ProviderID:    auth.UserAuth.ProviderID,
			Provider:      auth.UserAuth.Provider,
			Email:         auth.UserAuth.Email,
			EmailVerified: auth.UserAuth.EmailVerified,
			CreatedAt:     auth.UserAuth.CreatedAt,
		}
	}

	return entities.NewUserEntity(row.User, userAuthModels), nil
}

func (repo *postgresUserRepository) CreateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateUser(ctx, models.CreateUserParams{
		ID:           user.ID,
		GivenName:    user.GivenName,
		FamilyName:   user.FamilyName,
		PrimaryEmail: user.PrimaryEmail,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, entities.ErrEmailInUse
			case "users_user_handle_key":
				return nil, entities.ErrUserHandleInUse
			default:
				return nil, fmt.Errorf("unique constraint violation: %s", pgErr.ConstraintName)
			}
		}
		return nil, err
	}

	return entities.NewUserEntity(row, []models.UserAuth{}), nil
}

func (repo *postgresUserRepository) UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.UpdateUser(ctx, models.UpdateUserParams{
		ID:           user.ID,
		GivenName:    user.GivenName,
		FamilyName:   user.FamilyName,
		PrimaryEmail: user.PrimaryEmail,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "users_email_key":
				return nil, entities.ErrEmailInUse
			case "users_user_handle_key":
				return nil, entities.ErrUserHandleInUse
			default:
				return nil, fmt.Errorf("unique constraint violation: %s", pgErr.ConstraintName)
			}
		}
		return nil, err
	}

	userAuthRows, err := querier.GetUserAuthByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	userAuthModels := make([]models.UserAuth, len(userAuthRows))
	for i, auth := range userAuthRows {
		userAuthModels[i] = models.UserAuth{
			ID:            auth.UserAuth.ID,
			ProviderID:    auth.UserAuth.ProviderID,
			Provider:      auth.UserAuth.Provider,
			Email:         auth.UserAuth.Email,
			EmailVerified: auth.UserAuth.EmailVerified,
			CreatedAt:     auth.UserAuth.CreatedAt,
		}
	}

	return entities.NewUserEntity(row, userAuthModels), nil
}

func (repo *postgresUserRepository) CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity, token string, expiresAt time.Time) (*entities.UserContextEntity, error) {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	row, err := querier.CreateUserSession(ctx, models.CreateUserSessionParams{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}

	return entities.NewUserContextEntity(user, row), nil
}

func (repo *postgresUserRepository) CreateUserProviderData(ctx context.Context, querier models.Querier, user *entities.UserEntity, provider string, providerID string, providerValue string, email *string, providerData []byte) error {
	ctx, cancel := context.WithTimeout(ctx, postgres.DefaultTimeout)
	defer cancel()

	_, err := querier.CreateUserAuth(ctx, models.CreateUserAuthParams{
		UserID:        user.ID,
		Value:         providerValue,
		Provider:      provider,
		ProviderID:    providerID,
		ProviderData:  providerData,
		Email:         email,
		EmailVerified: false,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			switch pgErr.ConstraintName {
			case "user_auth_value_key":
				return entities.ErrUserClaimed
			default:
				return fmt.Errorf("unique constraint violation: %s", pgErr.ConstraintName)
			}
		}

		return err
	}
	return nil
}
