package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/mcorrigan89/identity/internal/config"
)

func configDB(cfg *config.Config) *pgxpool.Config {
	const defaultMaxConns = int32(4)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}

func openDBPool(cfg *config.Config) (*pgxpool.Pool, error) {
	dbConfigurationOptions := configDB(cfg)

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbConfigurationOptions)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	connection, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	err = dbpool.Ping(ctx)

	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
