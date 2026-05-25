package client

import (
	"context"
	"net/http"
)

func (c *Client) CreatePolicy(ctx context.Context, policy Policy) error {
	return c.do(ctx, http.MethodPost, "/v1/policies", policy, nil)
}

func (c *Client) GetPolicy(ctx context.Context, workspaceID, policyID string) (*Policy, error) {
	var policy Policy

	path := "/v1/workspaces/" + workspaceID + "/policies/" + policyID

	err := c.do(ctx, http.MethodGet, path, nil, &policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (c *Client) ListPolicies(ctx context.Context, workspaceID string) ([]Policy, error) {
	var policies []Policy

	path := "/v1/workspaces/" + workspaceID + "/policies"

	err := c.do(ctx, http.MethodGet, path, nil, &policies)
	if err != nil {
		return nil, err
	}

	return policies, nil
}

func (c *Client) AttachPolicyToRole(ctx context.Context, relation RolePolicy) error {
	return c.do(ctx, http.MethodPost, "/v1/roles/policies/attach", relation, nil)
}

func (c *Client) DetachPolicyFromRole(ctx context.Context, workspaceID, roleID, policyID string) error {
	path := "/v1/workspaces/" + workspaceID + "/roles/" + roleID + "/policies/" + policyID

	return c.do(ctx, http.MethodDelete, path, nil, nil)
}
