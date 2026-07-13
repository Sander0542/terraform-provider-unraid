// Copyright (c) Sander Jochems
// SPDX-License-Identifier: MIT

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphqlResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []graphqlError `json:"errors,omitempty"`
}

type graphqlError struct {
	Message string `json:"message"`
}

// Do executes a GraphQL query and decodes the response into T.
func Do[T any](ctx context.Context, c *Client, query string, variables map[string]any) (T, error) {
	var zero T

	body, err := json.Marshal(graphqlRequest{Query: query, Variables: variables})
	if err != nil {
		return zero, fmt.Errorf("marshaling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return zero, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return zero, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	var result graphqlResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return zero, fmt.Errorf("decoding response: %w", err)
	}

	if len(result.Errors) > 0 {
		return zero, fmt.Errorf("graphql error: %s", result.Errors[0].Message)
	}

	return result.Data, nil
}
