package application

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/identity/internal/application/commands"
	"github.com/mcorrigan89/identity/internal/application/queries"
	"github.com/mcorrigan89/identity/internal/domain/entities"
	"github.com/mcorrigan89/identity/internal/domain/services"
	"github.com/mcorrigan89/identity/internal/infrastructure/config"
	"github.com/mcorrigan89/identity/internal/infrastructure/oauth"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"

	"github.com/rs/zerolog"
)

type UserApplicationService interface {
	GetUserByID(ctx context.Context, query queries.UserByIDQuery) (*entities.UserEntity, error)
	GetUserBySessionToken(ctx context.Context, query queries.UserBySessionTokenQuery) (*entities.UserEntity, error)
	CreatePasswordUser(ctx context.Context, cmd commands.CreateNewPasswordUserCommand) (*entities.UserContextEntity, error)
	UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entities.UserEntity, error)
	OAuthLinkURL(ctx context.Context, query queries.OauthUrlQuery) string
	OAuthLogin(ctx context.Context, cmd commands.OauthLoginCallbackCommand) (*entities.UserContextEntity, error)
}

type userApplicationService struct {
	config       *config.Config
	wg           *sync.WaitGroup
	logger       *zerolog.Logger
	db           *pgxpool.Pool
	queries      models.Querier
	userService  services.UserService
	oauthService oauth.OAuthService
}

func NewUserApplicationService(db *pgxpool.Pool, wg *sync.WaitGroup, cfg *config.Config, logger *zerolog.Logger, userService services.UserService, oauthService oauth.OAuthService) *userApplicationService {
	dbQueries := models.New(db)
	return &userApplicationService{
		db:           db,
		config:       cfg,
		wg:           wg,
		logger:       logger,
		queries:      dbQueries,
		userService:  userService,
		oauthService: oauthService,
	}
}

func (app *userApplicationService) GetUserByID(ctx context.Context, query queries.UserByIDQuery) (*entities.UserEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting user by ID")

	user, err := app.userService.GetUserByID(ctx, app.queries, query.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by ID")
		return nil, err
	}

	return user, nil
}

func (app *userApplicationService) GetUserBySessionToken(ctx context.Context, query queries.UserBySessionTokenQuery) (*entities.UserEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Getting user by sessionToken")

	userContext, err := app.userService.GetUserContextBySessionToken(ctx, app.queries, query.SessionToken)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by sessionToken")
		return nil, err
	}

	return userContext.User, nil
}

func (app *userApplicationService) CreatePasswordUser(ctx context.Context, cmd commands.CreateNewPasswordUserCommand) (*entities.UserContextEntity, error) {
	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	app.logger.Info().Ctx(ctx).Msg("Creating new user")

	userEntity, password := cmd.ToDomain()

	createdUser, err := app.userService.CreatePasswordUser(ctx, qtx, userEntity, password)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create new user")
		return nil, err
	}

	userWithAuth, err := app.userService.GetUserByID(ctx, qtx, createdUser.ID)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get user by ID after creation")
		return nil, err
	}

	userSession, err := app.userService.CreateSession(ctx, qtx, userWithAuth)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create new session")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return userSession, nil
}

func (app *userApplicationService) UpdateUser(ctx context.Context, cmd commands.UpdateUserCommand) (*entities.UserEntity, error) {
	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)

	app.logger.Info().Ctx(ctx).Msg("Updating user")

	userEntity := cmd.ToDomain()

	updatedUser, err := app.userService.UpdateUser(ctx, qtx, userEntity)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to update user")
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
		return nil, err
	}

	return updatedUser, nil
}

func (app *userApplicationService) OAuthLinkURL(ctx context.Context, query queries.OauthUrlQuery) string {
	app.logger.Info().Ctx(ctx).Msg("Getting Oauth URL")

	oauthUrl := app.oauthService.GetGoogleOAuthLinkURL(ctx)

	return oauthUrl
}

func (app *userApplicationService) OAuthLogin(ctx context.Context, cmd commands.OauthLoginCallbackCommand) (*entities.UserContextEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Logging in with Oauth")

	switch cmd.Provider {
	case "google":
		return app.googleOAuthLogin(ctx, cmd)
	default:
		app.logger.Err(nil).Ctx(ctx).Msg("Invalid provider")
		return nil, entities.ErrInvalidProvider
	}

}

func (app *userApplicationService) googleOAuthLogin(ctx context.Context, cmd commands.OauthLoginCallbackCommand) (*entities.UserContextEntity, error) {
	app.logger.Info().Ctx(ctx).Msg("Logging in with Google Oauth")

	tx, cancel, err := postgres.CreateTransaction(ctx, app.db)
	defer cancel()
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to create transaction")
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := models.New(app.db).WithTx(tx)
	googleUser, googleData, err := app.oauthService.LoginWithGoogle(ctx, cmd.AuthCode)
	if err != nil {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to login with Google")
		return nil, err
	}

	userEntity, err := app.userService.GetUserByProviderID(ctx, qtx, googleUser.ID, "google")
	if err != nil && err != entities.ErrUserNotFound {
		app.logger.Err(err).Ctx(ctx).Msg("Failed to get User by Google ID")
		return nil, err
	}

	// Create new user if not found
	if userEntity != nil {
		userWithContext, err := app.userService.CreateSession(ctx, qtx, userEntity)
		if err != nil {
			app.logger.Err(err).Ctx(ctx).Msg("Failed to create new session")
			return nil, err
		}
		err = tx.Commit(ctx)
		if err != nil {
			app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
			return nil, err
		}
		return userWithContext, nil
	} else {
		// Find user with same email and create new auth provider
		userEntity, err := app.userService.GetUserByEmail(ctx, qtx, googleUser.Email)
		if err != nil && err != entities.ErrUserNotFound {
			app.logger.Err(err).Ctx(ctx).Msg("Failed to get User by email")
			return nil, err
		}
		if userEntity != nil {
			userAuth, err := app.userService.AddAuthProviderToUser(ctx, qtx, userEntity, "google", googleUser.ID, "1234", nil)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to add auth provider to user")
				return nil, err
			}

			userWithContext, err := app.userService.CreateSession(ctx, qtx, userAuth)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to create new session")
				return nil, err
			}
			err = tx.Commit(ctx)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
				return nil, err
			}
			return userWithContext, nil
		} else {

			userEntityToCreate := entities.UserEntity{
				ID:           uuid.New(),
				GivenName:    googleUser.GivenName,
				FamilyName:   googleUser.FamilyName,
				PrimaryEmail: &googleUser.Email,
			}

			googleDataJson, err := json.Marshal(googleData)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to marshal Google Token")
				return nil, err
			}

			userEntity, err = app.userService.CreateOAuthUser(ctx, qtx, &userEntityToCreate, "google", "1234", googleUser.ID, googleDataJson)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to create User from Google")
				return nil, err
			}

			userSession, err := app.userService.CreateSession(ctx, qtx, userEntity)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to create new session")
				return nil, err
			}

			err = tx.Commit(ctx)
			if err != nil {
				app.logger.Err(err).Ctx(ctx).Msg("Failed to commit transaction")
				return nil, err
			}

			return userSession, nil
		}
	}
}
