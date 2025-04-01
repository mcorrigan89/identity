package handlers

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/mcorrigan89/identity/internal/application"
	"github.com/mcorrigan89/identity/internal/application/commands"
	"github.com/mcorrigan89/identity/internal/application/queries"
	"github.com/mcorrigan89/identity/internal/interfaces/http/dto"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	logger         *zerolog.Logger
	userAppService application.UserApplicationService
}

func NewUserHandler(logger *zerolog.Logger, userAppService application.UserApplicationService) *UserHandler {
	return &UserHandler{
		logger:         logger,
		userAppService: userAppService,
	}
}

func (h *UserHandler) GetUserByID(ctx context.Context, input *struct {
	ID uuid.UUID `path:"id"`
}) (*dto.GetUserByIDResponse, error) {

	query := queries.UserByIDQuery{
		ID: input.ID,
	}

	user, err := h.userAppService.GetUserByID(ctx, query)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to get user by ID", err)
	}

	userDto := dto.NewUserDtoFromEntity(user)

	return &dto.GetUserByIDResponse{
		Body: userDto,
	}, nil
}

func (h *UserHandler) CreatePasswordUser(ctx context.Context, input *dto.CreatePasswordUserRequest) (*dto.CreateUserResponse, error) {
	cmd := commands.CreateNewPasswordUserCommand{
		Email:      input.Body.Email,
		GivenName:  input.Body.GivenName,
		FamilyName: input.Body.FamilyName,
		Password:   input.Body.Password,
	}

	userSessionEntity, err := h.userAppService.CreatePasswordUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to create user", err)
	}

	userDto := dto.NewUserDtoFromEntity(userSessionEntity.User)
	sessionDto := dto.SessionDto{
		Token:     userSessionEntity.SessionToken,
		ExpiresAt: userSessionEntity.ExpiresAt(),
	}

	resp := dto.CreateUserResponse{}

	resp.Body.User = userDto
	resp.Body.SessionDto = &sessionDto

	return &resp, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, input *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	cmd := commands.UpdateUserCommand{
		ID:         input.ID,
		GivenName:  input.Body.GivenName,
		FamilyName: input.Body.FamilyName,
	}

	userEntity, err := h.userAppService.UpdateUser(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to update user", err)
	}

	userDto := dto.NewUserDtoFromEntity(userEntity)

	resp := dto.UpdateUserResponse{}

	resp.Body.User = userDto

	return &resp, nil
}

func (h *UserHandler) OAuthLinkURL(ctx context.Context, input *dto.GetOAuthLinkURLRequest) (*dto.GetOAuthLinkURLResponse, error) {

	query := queries.OauthUrlQuery{
		Provider: input.Provider,
	}

	oauthLinkURL := h.userAppService.OAuthLinkURL(ctx, query)

	if oauthLinkURL == "" {
		return nil, huma.Error400BadRequest("Invalid provider", nil)
	}

	resp := dto.GetOAuthLinkURLResponse{}

	resp.Body.URL = oauthLinkURL
	return &resp, nil
}

func (h *UserHandler) OAuthLogin(ctx context.Context, input *dto.OAuthLoginRequest) (*dto.OAuthLoginResponse, error) {
	cmd := commands.OauthLoginCallbackCommand{
		AuthCode: input.AuthCode,
		Provider: input.Provider,
	}

	userSessionEntity, err := h.userAppService.OAuthLogin(ctx, cmd)
	if err != nil {
		return nil, huma.Error500InternalServerError("Failed to login with OAuth", err)
	}

	userDto := dto.NewUserDtoFromEntity(userSessionEntity.User)
	sessionDto := dto.SessionDto{
		Token:     userSessionEntity.SessionToken,
		ExpiresAt: userSessionEntity.ExpiresAt(),
	}

	resp := dto.OAuthLoginResponse{}

	resp.Body.User = userDto
	resp.Body.SessionDto = &sessionDto

	return &resp, nil
}
