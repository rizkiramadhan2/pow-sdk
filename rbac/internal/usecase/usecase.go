package usecase

import (
	"context"

	"github.com/rizkiramadhan2/rbac/internal/entity"
	"github.com/rizkiramadhan2/rbac/internal/repository"
	"github.com/rizkiramadhan2/rbac/internal/utils"
)

type WorkspaceUsecase struct {
	store repository.Store
}

func NewWorkspaceUsecase(store repository.Store) *WorkspaceUsecase {
	return &WorkspaceUsecase{store: store}
}

func (u *WorkspaceUsecase) Create(ctx context.Context, workspace entity.Workspace) error {
	if workspace.ID == "" || workspace.Name == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.CreateWorkspace(ctx, workspace)
}

func (u *WorkspaceUsecase) Get(ctx context.Context, workspaceID string) (*entity.Workspace, error) {
	if workspaceID == "" {
		return nil, utils.ErrInvalidRequest
	}

	return u.store.GetWorkspace(ctx, workspaceID)
}

type UserUsecase struct {
	store repository.Store
}

func NewUserUsecase(store repository.Store) *UserUsecase {
	return &UserUsecase{store: store}
}

func (u *UserUsecase) Create(ctx context.Context, user entity.User) error {
	if user.ID == "" || user.WorkspaceID == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.CreateUser(ctx, user)
}

func (u *UserUsecase) Get(ctx context.Context, workspaceID, userID string) (*entity.User, error) {
	if workspaceID == "" || userID == "" {
		return nil, utils.ErrInvalidRequest
	}

	return u.store.GetUser(ctx, workspaceID, userID)
}

type RoleUsecase struct {
	store repository.Store
}

func NewRoleUsecase(store repository.Store) *RoleUsecase {
	return &RoleUsecase{store: store}
}

func (u *RoleUsecase) Create(ctx context.Context, role entity.Role) error {
	if role.ID == "" || role.WorkspaceID == "" || role.Name == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.CreateRole(ctx, role)
}

func (u *RoleUsecase) Get(ctx context.Context, workspaceID, roleID string) (*entity.Role, error) {
	if workspaceID == "" || roleID == "" {
		return nil, utils.ErrInvalidRequest
	}

	return u.store.GetRole(ctx, workspaceID, roleID)
}

func (u *RoleUsecase) List(ctx context.Context, workspaceID string) ([]entity.Role, error) {
	if workspaceID == "" {
		return nil, utils.ErrInvalidRequest
	}

	return u.store.ListRoles(ctx, workspaceID)
}

func (u *RoleUsecase) AssignToUser(ctx context.Context, assignment entity.UserRoleAssignment) error {
	if assignment.WorkspaceID == "" || assignment.UserID == "" || assignment.RoleID == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.AssignRoleToUser(ctx, assignment)
}

func (u *RoleUsecase) RemoveFromUser(ctx context.Context, workspaceID, userID, roleID string) error {
	if workspaceID == "" || userID == "" || roleID == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.RemoveRoleFromUser(ctx, workspaceID, userID, roleID)
}

type PolicyUsecase struct {
	store repository.Store
}

func NewPolicyUsecase(store repository.Store) *PolicyUsecase {
	return &PolicyUsecase{store: store}
}

func (u *PolicyUsecase) Create(ctx context.Context, policy entity.Policy) error {
	if policy.ID == "" ||
		policy.WorkspaceID == "" ||
		policy.Resource == "" ||
		policy.Action == "" ||
		policy.Effect == "" {
		return utils.ErrInvalidRequest
	}

	if policy.Effect != entity.EffectAllow && policy.Effect != entity.EffectDeny {
		return utils.ErrInvalidRequest
	}

	return u.store.CreatePolicy(ctx, policy)
}

func (u *PolicyUsecase) Get(ctx context.Context, workspaceID, policyID string) (*entity.Policy, error) {
	if workspaceID == "" || policyID == "" {
		return nil, utils.ErrInvalidRequest
	}

	return u.store.GetPolicy(ctx, workspaceID, policyID)
}

func (u *PolicyUsecase) List(ctx context.Context, workspaceID string) ([]entity.Policy, error) {
	if workspaceID == "" {
		return nil, utils.ErrInvalidRequest
	}

	return u.store.ListPolicies(ctx, workspaceID)
}

func (u *PolicyUsecase) AttachToRole(ctx context.Context, relation entity.RolePolicy) error {
	if relation.WorkspaceID == "" || relation.RoleID == "" || relation.PolicyID == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.AttachPolicyToRole(ctx, relation)
}

func (u *PolicyUsecase) DetachFromRole(ctx context.Context, workspaceID, roleID, policyID string) error {
	if workspaceID == "" || roleID == "" || policyID == "" {
		return utils.ErrInvalidRequest
	}

	return u.store.DetachPolicyFromRole(ctx, workspaceID, roleID, policyID)
}

type AuthzUsecase struct {
	store repository.Store
}

func NewAuthzUsecase(store repository.Store) *AuthzUsecase {
	return &AuthzUsecase{store: store}
}

func (u *AuthzUsecase) Can(ctx context.Context, req entity.CheckRequest) (*entity.CheckResult, error) {
	if req.WorkspaceID == "" || req.UserID == "" || req.Resource == "" || req.Action == "" {
		return nil, utils.ErrInvalidRequest
	}

	roles, err := u.store.ListUserRoles(ctx, req.WorkspaceID, req.UserID)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return &entity.CheckResult{
			Allowed: false,
			Reason:  "user has no role in workspace",
		}, nil
	}

	roleIDs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	policies, err := u.store.ListPoliciesByRoleIDs(ctx, req.WorkspaceID, roleIDs)
	if err != nil {
		return nil, err
	}

	allowed := false

	for _, policy := range policies {
		if !matchPolicy(policy, req.Resource, req.Action) {
			continue
		}

		if policy.Effect == entity.EffectDeny {
			return &entity.CheckResult{
				Allowed: false,
				Reason:  "explicit deny policy matched",
			}, nil
		}

		if policy.Effect == entity.EffectAllow {
			allowed = true
		}
	}

	if allowed {
		return &entity.CheckResult{
			Allowed: true,
			Reason:  "allow policy matched",
		}, nil
	}

	return &entity.CheckResult{
		Allowed: false,
		Reason:  "no matching allow policy",
	}, nil
}

func matchPolicy(policy entity.Policy, resource, action string) bool {
	resourceMatch := policy.Resource == resource || policy.Resource == "*"
	actionMatch := policy.Action == action || policy.Action == "*"

	return resourceMatch && actionMatch
}
