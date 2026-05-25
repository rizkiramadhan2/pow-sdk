package client

import (
    "context"
    "net/http"
)

func (c *Client) Can(ctx context.Context, req CheckRequest) (*CheckResult, error) {
    var result CheckResult

    err := c.do(ctx, http.MethodPost, "/v1/authz/check", req, &result)
    if err != nil {
        return nil, err
    }

    return &result, nil
}