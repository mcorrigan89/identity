package queries

import (
	"github.com/google/uuid"
)

type UserByIDQuery struct {
	ID uuid.UUID
}

type UserBySessionTokenQuery struct {
	SessionToken string
}

type OauthUrlQuery struct {
	Provider string
}
