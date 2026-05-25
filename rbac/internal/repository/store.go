package repository

import (
    "context"

    "github.com/rizkiramadhan2/rbac/internal/entity"
)

type Store interface {
    GetWorkspace(ctx context.Context, workspaceID string) (*entity.Workspace, error)
    CreateWorkspace(ctx context.Context, workspace entity.Workspace) error

    CreateUser(ctx context.Context, user entity.User) error
    GetUser(ctx context.Context, workspaceID, userID string) (*entity.User, error)

    CreateRole(ctx context.Context, role entity.Role) error
    GetRole(ctx context.Context, workspaceID, roleID string) (*entity.Role, error)
    ListRoles(ctx context.Context, workspaceID string) ([]entity.Role, error)

    CreatePolicy(ctx context.Context, policy entity.Policy) error
    GetPolicy(ctx context.Context, workspaceID, policyID string) (*entity.Policy, error)
    ListPolicies(ctx context.Context, workspaceID string) ([]entity.Policy, error)

    AssignRoleToUser(ctx context.Context, assignment entity.UserRoleAssignment) error
    RemoveRoleFromUser(ctx context.Context, workspaceID, userID, roleID string) error
    ListUserRoles(ctx context.Context, workspaceID, userID string) ([]entity.Role, error)

    AttachPolicyToRole(ctx context.Context, relation entity.RolePolicy) error
    DetachPolicyFromRole(ctx context.Context, workspaceID, roleID, policyID string) error
    ListRolePolicies(ctx context.Context, workspaceID, roleID string) ([]entity.Policy, error)
    ListPoliciesByRoleIDs(ctx context.Context, workspaceID string, roleIDs []string) ([]entity.Policy, error)
}