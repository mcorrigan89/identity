package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/mcorrigan89/identity/internal/domain/entities"
	"github.com/mcorrigan89/identity/internal/domain/repositories"
	"github.com/mcorrigan89/identity/internal/infrastructure/config"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error)
	GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error)
	GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error)
	GetUserByProviderID(ctx context.Context, querier models.Querier, provider string, providerID string) (*entities.UserEntity, error)
	CreatePasswordUser(ctx context.Context, querier models.Querier, user *entities.UserEntity, password string) (*entities.UserEntity, error)
	CreateOAuthUser(ctx context.Context, querier models.Querier, user *entities.UserEntity, provider string, providerID string, providerValue string, providerData []byte) (*entities.UserEntity, error)
	AddAuthProviderToUser(ctx context.Context, querier models.Querier, user *entities.UserEntity, provider string, providerID string, providerValue string, providerData []byte) (*entities.UserEntity, error)
	UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error)
	CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserContextEntity, error)
}

type userService struct {
	config   *config.Config
	userRepo repositories.UserRepository
}

func NewUserService(config *config.Config, userRepo repositories.UserRepository) *userService {
	return &userService{config: config, userRepo: userRepo}
}

func (s *userService) GetUserByID(ctx context.Context, querier models.Querier, userID uuid.UUID) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetUserByID(ctx, querier, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, querier models.Querier, email string) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, querier, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) GetUserContextBySessionToken(ctx context.Context, querier models.Querier, sessionToken string) (*entities.UserContextEntity, error) {
	user, err := s.userRepo.GetUserContextBySessionToken(ctx, querier, sessionToken)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByProviderID(ctx context.Context, querier models.Querier, provider string, providerID string) (*entities.UserEntity, error) {
	user, err := s.userRepo.GetUserByProviderID(ctx, querier, provider, providerID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) CreatePasswordUser(ctx context.Context, querier models.Querier, user *entities.UserEntity, password string) (*entities.UserEntity, error) {
	user, err := s.userRepo.CreateUser(ctx, querier, user)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	hashedPasswordString := string(hashedPassword)

	err = s.userRepo.CreateUserProviderData(ctx, querier, user, "password", s.config.ServerToken, hashedPasswordString, user.PrimaryEmail, nil)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateOAuthUser(ctx context.Context, querier models.Querier, user *entities.UserEntity, provider string, providerID string, providerValue string, providerData []byte) (*entities.UserEntity, error) {
	user, err := s.userRepo.CreateUser(ctx, querier, user)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.CreateUserProviderData(ctx, querier, user, provider, providerID, providerValue, user.PrimaryEmail, providerData)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) AddAuthProviderToUser(ctx context.Context, querier models.Querier, user *entities.UserEntity, provider string, providerID string, providerValue string, providerData []byte) (*entities.UserEntity, error) {
	err := s.userRepo.CreateUserProviderData(ctx, querier, user, provider, providerID, providerValue, user.PrimaryEmail, providerData)
	if err != nil {
		return nil, err
	}

	userWithAuth, err := s.userRepo.GetUserByID(ctx, querier, user.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, entities.ErrUserNotFound
		}
		return nil, err
	}

	return userWithAuth, nil
}

func (s *userService) UpdateUser(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserEntity, error) {
	user, err := s.userRepo.UpdateUser(ctx, querier, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateSession(ctx context.Context, querier models.Querier, user *entities.UserEntity) (*entities.UserContextEntity, error) {

	token := xid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24 * 30)

	userSession, err := s.userRepo.CreateSession(ctx, querier, user, token, expiresAt)
	if err != nil {
		return nil, err
	}

	return userSession, nil
}
