package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) GetDriverToken(ctx context.Context, driverID string) (string, error) {
	var token string

	err := r.db.QueryRow(ctx,
		`SELECT token FROM driver_push_tokens WHERE driver_id=$1`,
		driverID,
	).Scan(&token)

	if err != nil {
		return "", err
	}

	return token, nil
}
