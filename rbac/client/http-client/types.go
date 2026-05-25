package client

type ErrorResponse struct {
	Error string `json:"error"`
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

type Workspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
}

type Role struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Policy struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspace_id"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Effect      string `json:"effect"`
}

type UserRoleAssignment struct {
	WorkspaceID string `json:"workspace_id"`
	UserID      string `json:"user_id"`
	RoleID      string `json:"role_id"`
}

type RolePolicy struct {
	WorkspaceID string `json:"workspace_id"`
	RoleID      string `json:"role_id"`
	PolicyID    string `json:"policy_id"`
}
