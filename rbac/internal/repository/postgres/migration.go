package postgres

import (
	"context"
	_ "embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations.sql
var migrationsSQL string

func Migrate(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, migrationsSQL)
	return err
}
