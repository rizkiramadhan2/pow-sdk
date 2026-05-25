package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiramadhan2/rbac/internal/entity"
)

type WorkspaceRepo struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepo(db *pgxpool.Pool) *WorkspaceRepo {
	return &WorkspaceRepo{db: db}
}

func (s *WorkspaceRepo) CreateWorkspace(ctx context.Context, workspace entity.Workspace) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_workspaces (id, name)
        VALUES ($1, $2)
    `, workspace.ID, workspace.Name)

	return mapErr(err)
}

func (s *WorkspaceRepo) GetWorkspace(ctx context.Context, workspaceID string) (*entity.Workspace, error) {
	var workspace entity.Workspace

	err := s.db.QueryRow(ctx, `
        SELECT id, name, created_at, updated_at
        FROM rbac_workspaces
        WHERE id = $1
    `, workspaceID).Scan(
		&workspace.ID,
		&workspace.Name,
		&workspace.CreatedAt,
		&workspace.UpdatedAt,
	)
	if err != nil {
		return nil, mapErr(err)
	}

	return &workspace, nil
}
