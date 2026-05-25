package entity

import "time"

type Effect string

const (
    EffectAllow Effect = "allow"
    EffectDeny  Effect = "deny"
)

type Workspace struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
    ID          string    `json:"id"`
    WorkspaceID string    `json:"workspace_id"`
    Email       string    `json:"email"`
    Name        string    `json:"name"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Role struct {
    ID          string    `json:"id"`
    WorkspaceID string    `json:"workspace_id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Policy struct {
    ID          string    `json:"id"`
    WorkspaceID string    `json:"workspace_id"`
    Resource    string    `json:"resource"`
    Action      string    `json:"action"`
    Effect      Effect    `json:"effect"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type UserRoleAssignment struct {
    WorkspaceID string    `json:"workspace_id"`
    UserID      string    `json:"user_id"`
    RoleID      string    `json:"role_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type RolePolicy struct {
    WorkspaceID string    `json:"workspace_id"`
    RoleID      string    `json:"role_id"`
    PolicyID    string    `json:"policy_id"`
    CreatedAt   time.Time `json:"created_at"`
}

type CheckRequest struct {
    WorkspaceID string `json:"workspace_id"`
    UserID      string `json:"user_id"`
    Resource    string `json:"resource"`
    Action      string `json:"action"`
}

type CheckResult struct {
    Allowed bool   `json:"allowed"`
    Reason  string `json:"reason"`
}