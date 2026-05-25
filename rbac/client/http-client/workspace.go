package client

import (
    "context"
    "net/http"
)

func (c *Client) CreateWorkspace(ctx context.Context, workspace Workspace) error {
    return c.do(ctx, http.MethodPost, "/v1/workspaces", workspace, nil)
}

func (c *Client) GetWorkspace(ctx context.Context, workspaceID string) (*Workspace, error) {
    var workspace Workspace

    err := c.do(ctx, http.MethodGet, "/v1/workspaces/"+workspaceID, nil, &workspace)
    if err != nil {
        return nil, err
    }

    return &workspace, nil
}