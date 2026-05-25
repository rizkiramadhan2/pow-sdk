package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiramadhan2/rbac/internal/entity"
)

type RoleRepo struct {
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool) *RoleRepo {
	return &RoleRepo{db: db}
}

func (s *RoleRepo) CreateRole(ctx context.Context, role entity.Role) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_roles (id, workspace_id, name, description)
        VALUES ($1, $2, $3, $4)
    `, role.ID, role.WorkspaceID, role.Name, role.Description)

	return mapErr(err)
}

func (s *RoleRepo) GetRole(ctx context.Context, workspaceID, roleID string) (*entity.Role, error) {
	var role entity.Role

	err := s.db.QueryRow(ctx, `
        SELECT id, workspace_id, name, description, created_at, updated_at
        FROM rbac_roles
        WHERE workspace_id = $1 AND id = $2
    `, workspaceID, roleID).Scan(
		&role.ID,
		&role.WorkspaceID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, mapErr(err)
	}

	return &role, nil
}

func (s *RoleRepo) ListRoles(ctx context.Context, workspaceID string) ([]entity.Role, error) {
	rows, err := s.db.Query(ctx, `
        SELECT id, workspace_id, name, description, created_at, updated_at
        FROM rbac_roles
        WHERE workspace_id = $1
        ORDER BY name ASC
    `, workspaceID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	var roles []entity.Role

	for rows.Next() {
		var role entity.Role

		err := rows.Scan(
			&role.ID,
			&role.WorkspaceID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, mapErr(err)
		}

		roles = append(roles, role)
	}

	return roles, mapErr(rows.Err())
}
