package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mcorrigan89/identity/internal/infrastructure/postgres/models"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrEmailInUse      = errors.New("email in use")
	ErrUserHandleInUse = errors.New("handle in use")
	ErrUserClaimed     = errors.New("user already claimed")
	ErrInvalidProvider = errors.New("invalid provider")
)

type UserEntity struct {
	ID           uuid.UUID
	GivenName    *string
	FamilyName   *string
	PrimaryEmail *string
	authEntities []*UserAuthEntity
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserAuthEntity struct {
	ID            uuid.UUID
	ProviderID    string
	Provider      string
	Email         *string
	EmailVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewUserEntity(userModel models.User, userAuth []models.UserAuth) *UserEntity {

	authEntities := make([]*UserAuthEntity, len(userAuth))
	for i, auth := range userAuth {
		authEntities[i] = &UserAuthEntity{
			ID:         auth.ID,
			ProviderID: auth.ProviderID,
			Provider:   auth.Provider,
			Email:      auth.Email,
			CreatedAt:  auth.CreatedAt,
			UpdatedAt:  auth.UpdatedAt,
		}
	}

	return &UserEntity{
		ID:           userModel.ID,
		GivenName:    userModel.GivenName,
		FamilyName:   userModel.FamilyName,
		PrimaryEmail: userModel.PrimaryEmail,
		authEntities: authEntities,
		CreatedAt:    userModel.CreatedAt,
		UpdatedAt:    userModel.UpdatedAt,
	}
}

func (u *UserEntity) FullName() *string {
	if u.GivenName != nil && u.FamilyName != nil {
		fullName := fmt.Sprintf("%s %s", *u.GivenName, *u.FamilyName)
		return &fullName
	}
	if u.GivenName != nil {
		return u.GivenName
	}

	if u.FamilyName != nil {
		return u.FamilyName
	}
	return nil
}

func (u *UserEntity) GetAuthEntities() []*UserAuthEntity {
	return u.authEntities
}

type EmailEntity struct {
	Email     string
	Verified  bool
	Primary   bool
	LinkedTo  []string
	CreatedAt time.Time
}

func (u *UserEntity) Emails() []*EmailEntity {
	emailMap := make(map[string]*EmailEntity)

	for _, auth := range u.authEntities {
		if auth.Email != nil {
			email := *auth.Email
			if _, exists := emailMap[email]; !exists {
				primary := false
				if u.PrimaryEmail != nil && email == *u.PrimaryEmail {
					primary = true
				}
				emailMap[email] = &EmailEntity{
					Email:     email,
					Verified:  false,
					Primary:   primary,
					LinkedTo:  []string{},
					CreatedAt: auth.CreatedAt,
				}
			}
			emailMap[email].LinkedTo = append(emailMap[email].LinkedTo, auth.Provider)
		}
	}
	emails := make([]*EmailEntity, 0, len(emailMap))
	for _, email := range emailMap {
		emails = append(emails, email)
	}

	return emails
}
