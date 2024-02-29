package rinha

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrClientNotFound      = errors.New("client not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

func parseDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return ErrClientNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Message == "cliente_not_found" {
			return ErrClientNotFound
		}
		if pgErr.SQLState() == "23514" {
			return ErrInsufficientBalance
		}
	}

	return err
}
