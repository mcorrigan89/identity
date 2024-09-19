package api

import (
	"context"
	"errors"
	"sync"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	identityv1 "github.com/mcorrigan89/identity/internal/api/serviceapis/identity/v1"
	"github.com/mcorrigan89/identity/internal/config"
	"github.com/mcorrigan89/identity/internal/entities"
	"github.com/mcorrigan89/identity/internal/services"

	"github.com/rs/zerolog"
)

type IdentityServerV1 struct {
	config   *config.Config
	wg       *sync.WaitGroup
	logger   *zerolog.Logger
	services *services.Services
}

func newIdentityProtoUrlServer(cfg *config.Config, logger *zerolog.Logger, wg *sync.WaitGroup, services *services.Services) *IdentityServerV1 {
	return &IdentityServerV1{
		config:   cfg,
		wg:       wg,
		logger:   logger,
		services: services,
	}
}

func (s *IdentityServerV1) GetUserById(ctx context.Context, req *connect.Request[identityv1.GetUserByIdRequest]) (*connect.Response[identityv1.GetUserByIdResponse], error) {
	userID := req.Msg.Id

	userUuid, err := uuid.Parse(userID)
	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error parsing user ID")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	userEntity, err := s.services.UserService.GetUserByID(ctx, userUuid)
	if err != nil {
		if err == entities.ErrUserNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		s.logger.Err(err).Ctx(ctx).Msg("Error getting user by ID")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&identityv1.GetUserByIdResponse{
		User: &identityv1.User{
			Id:         userEntity.ID.String(),
			GivenName:  userEntity.GivenName,
			FamilyName: userEntity.FamilyName,
			Email:      userEntity.Email,
		},
	})
	res.Header().Set("Identity-Version", "v1")
	return res, nil
}

func (s *IdentityServerV1) GetUserBySessionToken(ctx context.Context, req *connect.Request[identityv1.GetUserBySessionTokenRequest]) (*connect.Response[identityv1.GetUserBySessionTokenResponse], error) {
	token := req.Msg.Token
	if token == "" {
		err := errors.New("token is required")
		s.logger.Err(err).Ctx(ctx).Msg("Token is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	userEntity, _, err := s.services.UserService.GetUserBySessionToken(ctx, token)
	if err != nil {
		if err == entities.ErrUserNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		s.logger.Err(err).Ctx(ctx).Msg("Error getting user by session token")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&identityv1.GetUserBySessionTokenResponse{
		User: &identityv1.User{
			Id:         userEntity.ID.String(),
			GivenName:  userEntity.GivenName,
			FamilyName: userEntity.FamilyName,
			Email:      userEntity.Email,
		},
	})

	res.Header().Set("Identity-Version", "v1")
	return res, nil
}

func (s *IdentityServerV1) AuthenticateWithGoogleCode(ctx context.Context, req *connect.Request[identityv1.AuthenticateWithGoogleCodeRequest]) (*connect.Response[identityv1.AuthenticateWithGoogleCodeResponse], error) {
	code := req.Msg.Code
	if code == "" {
		err := errors.New("code is required")
		s.logger.Err(err).Ctx(ctx).Msg("Code is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	sessionEntity, err := s.services.OAuthService.LoginWithGoogleCode(ctx, code)

	if err != nil {
		s.logger.Err(err).Ctx(ctx).Msg("Error authenticating with Google code")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&identityv1.AuthenticateWithGoogleCodeResponse{
		Session: &identityv1.UserSession{
			Token:     sessionEntity.Token,
			ExpiresAt: sessionEntity.ExpiresAt.Format(time.RFC3339),
		},
	})

	res.Header().Set("Identity-Version", "v1")
	return res, nil
}

func (s *IdentityServerV1) AuthenticateWithPassword(ctx context.Context, req *connect.Request[identityv1.AuthenticateWithPasswordRequest]) (*connect.Response[identityv1.AuthenticateWithPasswordResponse], error) {
	email := req.Msg.Email
	password := req.Msg.Password
	if email == "" {
		err := errors.New("email is required")
		s.logger.Err(err).Ctx(ctx).Msg("Email is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	if password == "" {
		err := errors.New("password is required")
		s.logger.Err(err).Ctx(ctx).Msg("Password is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	sessionEntity, err := s.services.UserService.AuthenticateWithPassword(ctx, email, password)

	if err != nil {
		if err == entities.ErrUserNotFound {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		if err == entities.ErrInvalidCredentials {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		s.logger.Err(err).Ctx(ctx).Msg("Error authenticating with password")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&identityv1.AuthenticateWithPasswordResponse{
		Session: &identityv1.UserSession{
			Token:     sessionEntity.Token,
			ExpiresAt: sessionEntity.ExpiresAt.Format(time.RFC3339),
		},
	})

	res.Header().Set("Identity-Version", "v1")
	return res, nil

}

func (s *IdentityServerV1) CreateUser(ctx context.Context, req *connect.Request[identityv1.CreateUserRequest]) (*connect.Response[identityv1.CreateUserResponse], error) {
	givenName := req.Msg.GivenName
	familyName := req.Msg.FamilyName
	email := req.Msg.Email
	password := req.Msg.Password

	if email == "" {
		err := errors.New("email is required")
		s.logger.Err(err).Ctx(ctx).Msg("Email is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if password == "" {
		err := errors.New("password is required")
		s.logger.Err(err).Ctx(ctx).Msg("Password is required")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	userEntity, err := s.services.UserService.CreateUser(ctx, services.CreateUserArgs{
		GivenName:  givenName,
		FamilyName: familyName,
		Email:      email,
		Password:   password,
	})

	if err != nil {
		if err == entities.ErrDuplicateEmail {
			return nil, connect.NewError(connect.CodeAlreadyExists, err)
		}
		if err == entities.ErrUserAlreadyCreated {
			return nil, connect.NewError(connect.CodeUnavailable, err)
		}
		s.logger.Err(err).Ctx(ctx).Msg("Error creating user")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&identityv1.CreateUserResponse{
		User: &identityv1.User{
			Id:         userEntity.ID.String(),
			GivenName:  userEntity.GivenName,
			FamilyName: userEntity.FamilyName,
			Email:      userEntity.Email,
		},
	})

	res.Header().Set("Identity-Version", "v1")
	return res, nil
}
