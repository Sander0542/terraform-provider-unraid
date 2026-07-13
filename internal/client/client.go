// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

// Client is a simple GraphQL HTTP client for the Unraid API.
type Client struct {
	endpoint   string
	apiToken   string
	httpClient *http.Client
}

// New creates a new Unraid API client.
func New(endpoint string, apiToken string) *Client {
	return &Client{
		endpoint:   endpoint,
		apiToken:   apiToken,
		httpClient: &http.Client{},
	}
}

// MakeRequest implements graphql.Client.
func (c *Client) MakeRequest(ctx context.Context, req *graphql.Request, resp *graphql.Response) error {
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiToken)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", httpResp.Status)
	}

	if err := json.NewDecoder(httpResp.Body).Decode(resp); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}
