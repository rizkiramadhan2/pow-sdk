package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiramadhan2/rbac/internal/entity"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) CreateUser(ctx context.Context, user entity.User) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_users (id, workspace_id, email, name)
        VALUES ($1, $2, $3, $4)
    `, user.ID, user.WorkspaceID, user.Email, user.Name)

	return mapErr(err)
}

func (s *UserRepo) GetUser(ctx context.Context, workspaceID, userID string) (*entity.User, error) {
	var user entity.User

	err := s.db.QueryRow(ctx, `
        SELECT id, workspace_id, email, name, created_at, updated_at
        FROM rbac_users
        WHERE workspace_id = $1 AND id = $2
    `, workspaceID, userID).Scan(
		&user.ID,
		&user.WorkspaceID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, mapErr(err)
	}

	return &user, nil
}
