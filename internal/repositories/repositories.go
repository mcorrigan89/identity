package repositories

import (
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mcorrigan89/identity/internal/repositories/models"
	"github.com/rs/zerolog"
)

const defaultTimeout = 10 * time.Second

type ServicesUtils struct {
	logger *zerolog.Logger
	wg     *sync.WaitGroup
	db     *pgxpool.Pool
}

type Repositories struct {
	utils          ServicesUtils
	UserRepository *UserRepository
}

func NewRepositories(db *pgxpool.Pool, logger *zerolog.Logger, wg *sync.WaitGroup) Repositories {
	queries := models.New(db)
	utils := ServicesUtils{
		logger: logger,
		wg:     wg,
		db:     db,
	}

	userRepo := NewUserRepository(utils, db, queries)

	return Repositories{
		utils:          utils,
		UserRepository: userRepo,
	}
}
