package commands

import (
	"github.com/google/uuid"
	"github.com/mcorrigan89/identity/internal/domain/entities"
)

type CreateNewPasswordUserCommand struct {
	Email      string  `json:"email" validate:"required,email"`
	GivenName  *string `json:"firstName" validate:"-"`
	FamilyName *string `json:"lastName" validate:"-"`
	Password   string  `json:"password" validate:"required,min=8"`
}

func (cmd *CreateNewPasswordUserCommand) ToDomain() (*entities.UserEntity, string) {
	return &entities.UserEntity{
		ID:           uuid.New(),
		GivenName:    cmd.GivenName,
		FamilyName:   cmd.FamilyName,
		PrimaryEmail: &cmd.Email,
	}, cmd.Password
}

type UpdateUserCommand struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	GivenName  *string   `json:"firstName" validate:"-"`
	FamilyName *string   `json:"lastName" validate:"-"`
}

func (cmd *UpdateUserCommand) ToDomain() *entities.UserEntity {
	return &entities.UserEntity{
		ID:         cmd.ID,
		GivenName:  cmd.GivenName,
		FamilyName: cmd.FamilyName,
	}
}

type OauthLoginCallbackCommand struct {
	Provider string
	AuthCode string
	Redirect string
}
