package rinha

import (
	"context"
	"log/slog"

	"github.com/goccy/go-json"

	"github.com/allisson/pgxutil/v2"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDatabaseConnection(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	databaseURL := cfg.DatabaseURL
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		slog.Error("database config error", "error", err.Error())
		return nil, err
	}
	config.MinConns = int32(cfg.DatabaseMinConns)
	config.MaxConns = int32(cfg.DatabaseMaxConns)

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("database pool error", "error", err.Error())
		return nil, err
	}

	return pool, nil
}

func GetClient(ctx context.Context, pool *pgxpool.Pool, clientID uint) (Client, error) {
	client := Client{}
	options := pgxutil.NewFindOptions().WithFilter("id", clientID)
	err := pgxutil.Get(ctx, pool, "clientes", options, &client)
	return client, parseDatabaseError(err)
}

func AddTransaction(ctx context.Context, pool *pgxpool.Pool, clientID uint, transaction Transaction) (Balance, error) {
	jsonTransaction, err := json.Marshal(&transaction)
	if err != nil {
		return Balance{}, err
	}

	result := map[string]Balance{}

	if err := pgxscan.Get(ctx, pool, &result, "SELECT add_transaction($1, $2)", int(clientID), string(jsonTransaction)); err != nil {
		return Balance{}, parseDatabaseError(err)
	}

	return result["add_transaction"], nil
}
