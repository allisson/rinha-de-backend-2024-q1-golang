package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	rinha "github.com/allisson/rinha-de-backend-2024-q1-golang"
)

func warnUp(ctx context.Context, pool *pgxpool.Pool) {
	for i := 1; i < 50; i++ {
		// nolint
		rinha.GetClient(ctx, pool, 9999)
	}
	for i := 1; i < 50; i++ {
		// nolint
		rinha.AddTransaction(ctx, pool, 9999, rinha.Transaction{Amount: 100, Type: rinha.CreditType, Description: "descricao", CreatedAt: time.Now().UTC()})
	}
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))
	cfg := rinha.NewConfig()
	ctx := context.Background()

	pool, err := rinha.SetupDatabaseConnection(ctx, cfg)
	if err != nil {
		slog.Error("Database connection error", "error", err.Error())
		os.Exit(1)
	}

	defer pool.Close()

	slog.Info("Starting the warnup")
	warnUp(ctx, pool)

	slog.Info("Starting the server")
	rinha.RunServer(cfg, pool)
}
