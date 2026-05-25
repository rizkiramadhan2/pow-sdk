package client

import (
	"context"
	"net/http"
)

func (c *Client) CreateUser(ctx context.Context, user User) error {
	return c.do(ctx, http.MethodPost, "/v1/users", user, nil)
}

func (c *Client) GetUser(ctx context.Context, workspaceID, userID string) (*User, error) {
	var user User

	path := "/v1/workspaces/" + workspaceID + "/users/" + userID

	err := c.do(ctx, http.MethodGet, path, nil, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
