package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkiramadhan2/rbac/internal/entity"
)

type PolicyRepo struct {
	db *pgxpool.Pool
}

func NewPolicyRepo(db *pgxpool.Pool) *PolicyRepo {
	return &PolicyRepo{db: db}
}

func (s *PolicyRepo) CreatePolicy(ctx context.Context, policy entity.Policy) error {
	_, err := s.db.Exec(ctx, `
        INSERT INTO rbac_policies (id, workspace_id, resource, action, effect)
        VALUES ($1, $2, $3, $4, $5)
    `, policy.ID, policy.WorkspaceID, policy.Resource, policy.Action, string(policy.Effect))

	return mapErr(err)
}

func (s *PolicyRepo) GetPolicy(ctx context.Context, workspaceID, policyID string) (*entity.Policy, error) {
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

func (s *PolicyRepo) ListPolicies(ctx context.Context, workspaceID string) ([]entity.Policy, error) {
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
