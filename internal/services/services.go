package services

import (
	"sync"

	"github.com/mcorrigan89/identity/internal/config"
	"github.com/mcorrigan89/identity/internal/repositories"
	"github.com/rs/zerolog"
)

type ServicesUtils struct {
	logger *zerolog.Logger
	wg     *sync.WaitGroup
	config *config.Config
}

type Services struct {
	utils        ServicesUtils
	UserService  *UserService
	OAuthService *OAuthService
}

func NewServices(repositories *repositories.Repositories, cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup) Services {
	utils := ServicesUtils{
		logger: logger,
		wg:     wg,
		config: cfg,
	}

	userService := NewUserService(utils, repositories.UserRepository)
	oAuthService := NewOAuthService(utils, userService, repositories.UserRepository)

	return Services{
		utils:        utils,
		UserService:  userService,
		OAuthService: oAuthService,
	}
}
