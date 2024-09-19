package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
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

type TxFn func(*pgx.Tx) error

func (utils *ServicesUtils) WithTransaction(ctx context.Context, fn TxFn) error {
	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	tx, err := utils.db.Begin(ctx)
	if err != nil {
		utils.logger.Err(err).Ctx(ctx).Msg("Begin transaction")
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				utils.logger.Err(rbErr).Ctx(ctx).Msg("Rollback failed")
			}
		}
	}()

	err = fn(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)

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
