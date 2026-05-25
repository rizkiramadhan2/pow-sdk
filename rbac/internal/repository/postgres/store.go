package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rizkiramadhan2/rbac/internal/entity"
	"github.com/rizkiramadhan2/rbac/internal/utils"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{db: db}
}

func mapErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return utils.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return utils.ErrAlreadyExists
		case "23503":
			return utils.ErrNotFound
		}
	}

	return err
}

func (s *Store) CreateWorkspace(ctx context.Context, workspace entity.Workspace) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_workspaces (id, name)
        VALUES ($1, $2)
    `, workspace.ID, workspace.Name)

	return mapErr(err)
}

func (s *Store) GetWorkspace(ctx context.Context, workspaceID string) (*entity.Workspace, error) {
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

func (s *Store) CreateUser(ctx context.Context, user entity.User) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_users (id, workspace_id, email, name)
        VALUES ($1, $2, $3, $4)
    `, user.ID, user.WorkspaceID, user.Email, user.Name)

	return mapErr(err)
}

func (s *Store) GetUser(ctx context.Context, workspaceID, userID string) (*entity.User, error) {
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

func (s *Store) CreateRole(ctx context.Context, role entity.Role) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_roles (id, workspace_id, name, description)
        VALUES ($1, $2, $3, $4)
    `, role.ID, role.WorkspaceID, role.Name, role.Description)

	return mapErr(err)
}

func (s *Store) GetRole(ctx context.Context, workspaceID, roleID string) (*entity.Role, error) {
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

func (s *Store) ListRoles(ctx context.Context, workspaceID string) ([]entity.Role, error) {
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

func (s *Store) CreatePolicy(ctx context.Context, policy entity.Policy) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_policies (id, workspace_id, resource, action, effect)
        VALUES ($1, $2, $3, $4, $5)
    `, policy.ID, policy.WorkspaceID, policy.Resource, policy.Action, string(policy.Effect))

	return mapErr(err)
}

func (s *Store) GetPolicy(ctx context.Context, workspaceID, policyID string) (*entity.Policy, error) {
	var policy entity.Policy

	err := s.db.QueryRow(ctx, `
        SELECT id, workspace_id, resource, action, effect, created_at, updated_at
        FROM rbac_policies
        WHERE workspace_id = $1 AND id = $2
    `, workspaceID, policyID).Scan(
		&policy.ID,
		&policy.WorkspaceID,
		&policy.Resource,
		&policy.Action,
		&policy.Effect,
		&policy.CreatedAt,
		&policy.UpdatedAt,
	)
	if err != nil {
		return nil, mapErr(err)
	}

	return &policy, nil
}

func (s *Store) ListPolicies(ctx context.Context, workspaceID string) ([]entity.Policy, error) {
	rows, err := s.db.Query(ctx, `
        SELECT id, workspace_id, resource, action, effect, created_at, updated_at
        FROM rbac_policies
        WHERE workspace_id = $1
        ORDER BY resource ASC, action ASC
    `, workspaceID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	var policies []entity.Policy

	for rows.Next() {
		var policy entity.Policy

		err := rows.Scan(
			&policy.ID,
			&policy.WorkspaceID,
			&policy.Resource,
			&policy.Action,
			&policy.Effect,
			&policy.CreatedAt,
			&policy.UpdatedAt,
		)
		if err != nil {
			return nil, mapErr(err)
		}

		policies = append(policies, policy)
	}

	return policies, mapErr(rows.Err())
}

func (s *Store) AssignRoleToUser(ctx context.Context, assignment entity.UserRoleAssignment) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_user_roles (workspace_id, user_id, role_id)
        VALUES ($1, $2, $3)
    `, assignment.WorkspaceID, assignment.UserID, assignment.RoleID)

	return mapErr(err)
}

func (s *Store) RemoveRoleFromUser(ctx context.Context, workspaceID, userID, roleID string) error {
	result, err := s.db.Exec(ctx, `
        DELETE FROM rbac_user_roles
        WHERE workspace_id = $1 AND user_id = $2 AND role_id = $3
    `, workspaceID, userID, roleID)
	if err != nil {
		return mapErr(err)
	}

	if result.RowsAffected() == 0 {
		return utils.ErrNotFound
	}

	return nil
}

func (s *Store) ListUserRoles(ctx context.Context, workspaceID, userID string) ([]entity.Role, error) {
	rows, err := s.db.Query(ctx, `
        SELECT r.id, r.workspace_id, r.name, r.description, r.created_at, r.updated_at
        FROM rbac_roles r
        INNER JOIN rbac_user_roles ur
            ON ur.workspace_id = r.workspace_id
            AND ur.role_id = r.id
        WHERE ur.workspace_id = $1
        AND ur.user_id = $2
        ORDER BY r.name ASC
    `, workspaceID, userID)
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

func (s *Store) AttachPolicyToRole(ctx context.Context, relation entity.RolePolicy) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_role_policies (workspace_id, role_id, policy_id)
        VALUES ($1, $2, $3)
    `, relation.WorkspaceID, relation.RoleID, relation.PolicyID)

	return mapErr(err)
}

func (s *Store) DetachPolicyFromRole(ctx context.Context, workspaceID, roleID, policyID string) error {
	result, err := s.db.Exec(ctx, `
        DELETE FROM rbac_role_policies
        WHERE workspace_id = $1 AND role_id = $2 AND policy_id = $3
    `, workspaceID, roleID, policyID)
	if err != nil {
		return mapErr(err)
	}

	if result.RowsAffected() == 0 {
		return utils.ErrNotFound
	}

	return nil
}

func (s *Store) ListRolePolicies(ctx context.Context, workspaceID, roleID string) ([]entity.Policy, error) {
	rows, err := s.db.Query(ctx, `
        SELECT p.id, p.workspace_id, p.resource, p.action, p.effect, p.created_at, p.updated_at
        FROM rbac_policies p
        INNER JOIN rbac_role_policies rp
            ON rp.workspace_id = p.workspace_id
            AND rp.policy_id = p.id
        WHERE rp.workspace_id = $1
        AND rp.role_id = $2
        ORDER BY p.resource ASC, p.action ASC
    `, workspaceID, roleID)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	var policies []entity.Policy

	for rows.Next() {
		var policy entity.Policy

		err := rows.Scan(
			&policy.ID,
			&policy.WorkspaceID,
			&policy.Resource,
			&policy.Action,
			&policy.Effect,
			&policy.CreatedAt,
			&policy.UpdatedAt,
		)
		if err != nil {
			return nil, mapErr(err)
		}

		policies = append(policies, policy)
	}

	return policies, mapErr(rows.Err())
}

func (s *Store) ListPoliciesByRoleIDs(ctx context.Context, workspaceID string, roleIDs []string) ([]entity.Policy, error) {
	if len(roleIDs) == 0 {
		return []entity.Policy{}, nil
	}

	rows, err := s.db.Query(ctx, `
        SELECT DISTINCT p.id, p.workspace_id, p.resource, p.action, p.effect, p.created_at, p.updated_at
        FROM rbac_policies p
        INNER JOIN rbac_role_policies rp
            ON rp.workspace_id = p.workspace_id
            AND rp.policy_id = p.id
        WHERE rp.workspace_id = $1
        AND rp.role_id = ANY($2)
    `, workspaceID, roleIDs)
	if err != nil {
		return nil, mapErr(err)
	}
	defer rows.Close()

	var policies []entity.Policy

	for rows.Next() {
		var policy entity.Policy

		err := rows.Scan(
			&policy.ID,
			&policy.WorkspaceID,
			&policy.Resource,
			&policy.Action,
			&policy.Effect,
			&policy.CreatedAt,
			&policy.UpdatedAt,
		)
		if err != nil {
			return nil, mapErr(err)
		}

		policies = append(policies, policy)
	}

	return policies, mapErr(rows.Err())
}
