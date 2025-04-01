package oauth

import (
	"context"

	"github.com/mcorrigan89/identity/internal/infrastructure/config"
	"github.com/rs/zerolog"
)

type OAuthService interface {
	LoginWithGoogle(ctx context.Context, code string) (*GoogleUser, *GoogleTokenResponse, error)
	GetGoogleOAuthLinkURL(ctx context.Context) string
}

type oAuthService struct {
	config      *config.Config
	logger      *zerolog.Logger
	googleOauth *GoogleOAuthService
}

func NewOAuthService(logger *zerolog.Logger, config *config.Config) *oAuthService {
	googleOauth := NewGoogleOAuthService(logger, config)
	return &oAuthService{logger: logger, config: config, googleOauth: googleOauth}
}

func (service *oAuthService) LoginWithGoogle(ctx context.Context, code string) (*GoogleUser, *GoogleTokenResponse, error) {
	googleUser, googleData, err := service.googleOauth.loginWithGoogleCode(ctx, code)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to login with Google")
		return nil, nil, err
	}

	return googleUser, googleData, nil
}

func (service *oAuthService) GetGoogleOAuthLinkURL(ctx context.Context) string {
	service.logger.Info().Ctx(ctx).Msgf("Google Auth URL")
	authUrl := service.googleOauth.GetGoogleLoginURL(service.config.CientURL, service.config.OAuth.Google.ClientID)
	return authUrl
}
