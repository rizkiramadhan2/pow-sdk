package client

import (
	"context"
	"net/http"
)

func (c *Client) CreateRole(ctx context.Context, role Role) error {
	return c.do(ctx, http.MethodPost, "/v1/roles", role, nil)
}

func (c *Client) GetRole(ctx context.Context, workspaceID, roleID string) (*Role, error) {
	var role Role

	path := "/v1/workspaces/" + workspaceID + "/roles/" + roleID

	err := c.do(ctx, http.MethodGet, path, nil, &role)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (c *Client) ListRoles(ctx context.Context, workspaceID string) ([]Role, error) {
	var roles []Role

	path := "/v1/workspaces/" + workspaceID + "/roles"

	err := c.do(ctx, http.MethodGet, path, nil, &roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (c *Client) AssignRoleToUser(ctx context.Context, assignment UserRoleAssignment) error {
	return c.do(ctx, http.MethodPost, "/v1/roles/assign", assignment, nil)
}

func (c *Client) RemoveRoleFromUser(ctx context.Context, workspaceID, userID, roleID string) error {
	path := "/v1/workspaces/" + workspaceID + "/users/" + userID + "/roles/" + roleID

	return c.do(ctx, http.MethodDelete, path, nil, nil)
}
