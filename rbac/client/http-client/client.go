package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	Host       string
	Secret     string
	HTTPClient *http.Client
}

type Client struct {
	host       string
	secret     string
	httpClient *http.Client
}

func New(config Config) *Client {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return &Client{
		host:       strings.TrimRight(config.Host, "/"),
		secret:     config.Secret,
		httpClient: httpClient,
	}
}

func (c *Client) do(ctx context.Context, method, path string, body any, out any) error {
	var reqBody *bytes.Reader

	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}

		reqBody = bytes.NewReader(payload)
	} else {
		reqBody = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.host+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.secret != "" {
		req.Header.Set("Authorization", "Bearer "+c.secret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp ErrorResponse
		_ = json.NewDecoder(resp.Body).Decode(&errResp)

		if errResp.Error != "" {
			return fmt.Errorf("rbac client: %s", errResp.Error)
		}

		return fmt.Errorf("rbac client: request failed with status %d", resp.StatusCode)
	}

	if out == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(out)
}
