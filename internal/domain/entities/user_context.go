package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"
)

type UserContextEntity struct {
	SessionToken   string
	UserID         uuid.UUID
	User           *UserEntity
	userExpired    bool
	expiresAt      time.Time
	impersonatorID *uuid.UUID
}

func NewUserContextEntity(user *UserEntity, userSessionModel models.UserSession) *UserContextEntity {
	return &UserContextEntity{
		SessionToken:   userSessionModel.Token,
		UserID:         user.ID,
		User:           user,
		userExpired:    userSessionModel.UserExpired,
		expiresAt:      userSessionModel.ExpiresAt,
		impersonatorID: userSessionModel.ImpersonatorID,
	}
}

func (uc *UserContextEntity) IsExpired() bool {
	return uc.userExpired || uc.expiresAt.Before(time.Now())
}

func (uc *UserContextEntity) ExpiresAt() time.Time {
	return uc.expiresAt
}
