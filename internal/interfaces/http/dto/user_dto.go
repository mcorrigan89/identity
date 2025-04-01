package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/identity/internal/domain/entities"
)

type UserDto struct {
	ID         uuid.UUID   `json:"id"`
	GivenName  *string     `json:"given_name"`
	FamilyName *string     `json:"family_name"`
	Emails     []*EmailDto `json:"emails"`
}

type EmailDto struct {
	Email     string    `json:"email"`
	Verified  bool      `json:"verified"`
	Primary   bool      `json:"primary"`
	LinkedTo  []string  `json:"linked_to"`
	CreatedAt time.Time `json:"created_at"`
}

func NewUserDtoFromEntity(entity *entities.UserEntity) *UserDto {

	if entity == nil {
		return nil
	}

	emails := make([]*EmailDto, 0)

	for _, email := range entity.Emails() {
		emails = append(emails, &EmailDto{
			Email:     email.Email,
			Verified:  email.Verified,
			Primary:   email.Primary,
			LinkedTo:  email.LinkedTo,
			CreatedAt: email.CreatedAt,
		})
	}

	return &UserDto{
		ID:         entity.ID,
		GivenName:  entity.GivenName,
		FamilyName: entity.FamilyName,
		Emails:     emails,
	}
}

type SessionDto struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type GetUserByIDResponse struct {
	Body *UserDto `json:"body"`
}

type CreatePasswordUserRequest struct {
	Body struct {
		Email      string  `json:"email"`
		GivenName  *string `json:"given_name"`
		FamilyName *string `json:"family_name"`
		Password   string  `json:"password"`
	}
}

type CreateUserResponse struct {
	Body struct {
		User       *UserDto    `json:"user"`
		SessionDto *SessionDto `json:"session"`
	} `json:"body"`
}

type UpdateUserRequest struct {
	ID   uuid.UUID `path:"id"`
	Body struct {
		GivenName  *string `json:"given_name"`
		FamilyName *string `json:"family_name"`
	}
}

type UpdateUserResponse struct {
	Body struct {
		User *UserDto `json:"user"`
	} `json:"body"`
}

type GetOAuthLinkURLRequest struct {
	Provider string `path:"provider"`
}

type GetOAuthLinkURLResponse struct {
	Body struct {
		URL string `json:"url"`
	} `json:"body"`
}

type OAuthLoginRequest struct {
	Provider string `path:"provider"`
	AuthCode string `query:"code"`
}

type OAuthLoginResponse struct {
	Body struct {
		User       *UserDto    `json:"user"`
		SessionDto *SessionDto `json:"session"`
	} `json:"body"`
}
