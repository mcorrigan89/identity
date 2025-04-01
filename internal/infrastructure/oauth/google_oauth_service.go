package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mcorrigan89/identity/internal/infrastructure/config"
	"github.com/rs/zerolog"
)

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	IDToken      string `json:"id_token"`
}

type GoogleOAuthService struct {
	config *config.Config
	logger *zerolog.Logger
}

func NewGoogleOAuthService(logger *zerolog.Logger, config *config.Config) *GoogleOAuthService {
	return &GoogleOAuthService{logger: logger, config: config}
}

func (service *GoogleOAuthService) GetGoogleLoginURL(clientUrl, clientID string) string {

	redirectUrl := fmt.Sprintf("%s/callback/google", clientUrl)

	queryValues := url.Values{
		"scope":                  {"openid profile email"},
		"access_type":            {"offline"},
		"include_granted_scopes": {"true"},
		"response_type":          {"code"},
		"state":                  {"google"},
		"redirect_uri":           {redirectUrl},
		"client_id":              {clientID},
	}

	authUrl := fmt.Sprintf("https://accounts.google.com/o/oauth2/v2/auth?%s", queryValues.Encode())

	return authUrl

}

func (service *GoogleOAuthService) loginWithGoogleCode(ctx context.Context, code string) (*GoogleUser, *GoogleTokenResponse, error) {
	tokenResponse, err := service.getGoogleTokenFromCode(ctx, code)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to get Google Token")
		return nil, nil, err
	}

	user, err := service.greateGoogleUser(ctx, tokenResponse)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to create Google User")
		return nil, nil, err
	}

	return user, tokenResponse, nil
}

func (service *GoogleOAuthService) getGoogleTokenFromCode(ctx context.Context, code string) (*GoogleTokenResponse, error) {
	service.logger.Info().Ctx(ctx).Str("code", code).Msg("Getting Google Auth with Code")
	clientID := service.config.OAuth.Google.ClientID
	clientSecret := service.config.OAuth.Google.ClientSecret
	redirectUrl := fmt.Sprintf("%s/callback/google", service.config.CientURL)

	url := fmt.Sprintf("https://oauth2.googleapis.com/token?code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, clientID, clientSecret, redirectUrl)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to get Google Token")
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to read Google Token response")
		return nil, err
	}

	var tokenResponse GoogleTokenResponse
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to unmarshal Google Token response")
		return nil, err
	}

	return &tokenResponse, nil
}

type GoogleUser struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	VerifiedEmail bool    `json:"verified_email"`
	Name          *string `json:"name,omitempty"`
	GivenName     *string `json:"given_name,omitempty"`
	FamilyName    *string `json:"family_name,omitempty"`
	Picture       *string `json:"picture,omitempty"`
}

func (service *GoogleOAuthService) greateGoogleUser(ctx context.Context, tokenResponse *GoogleTokenResponse) (*GoogleUser, error) {
	service.logger.Info().Ctx(ctx).Msg("Creating Google User")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo", nil)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to create Google User request")
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tokenResponse.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to get Google User")
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to read Google User response")
		return nil, err
	}

	var user GoogleUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		service.logger.Err(err).Ctx(ctx).Msg("Failed to unmarshal Google User response")
		return nil, err
	}

	return &user, nil
}
