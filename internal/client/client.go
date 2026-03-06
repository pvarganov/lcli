package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const linearAPIURL = "https://api.linear.app/graphql"

// LinearClient — интерфейс для взаимодействия с Linear GraphQL API.
type LinearClient interface {
	Do(query string, variables map[string]any, result any) error
}

// graphqlRequest — структура тела запроса.
type graphqlRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

// graphqlResponse — обёртка вокруг ответа API.
type graphqlResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []graphqlError  `json:"errors"`
}

type graphqlError struct {
	Message string `json:"message"`
}

// Client — реализация LinearClient.
type Client struct {
	httpClient *http.Client
	token      string
	url        string
}

// New создаёт новый Client с заданным токеном.
func New(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		url:        linearAPIURL,
	}
}

// NewWithURL создаёт Client с произвольным URL (для тестов).
func NewWithURL(token, url string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		url:        url,
	}
}

// Do выполняет GraphQL-запрос и декодирует результат в result.
func (c *Client) Do(query string, variables map[string]any, result any) error {
	body, err := json.Marshal(graphqlRequest{Query: query, Variables: variables})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %s", resp.Status)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var gqlResp graphqlResponse
	if err := json.Unmarshal(raw, &gqlResp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}

	if len(gqlResp.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", gqlResp.Errors[0].Message)
	}

	if result != nil && gqlResp.Data != nil {
		if err := json.Unmarshal(gqlResp.Data, result); err != nil {
			return fmt.Errorf("unmarshal data: %w", err)
		}
	}

	return nil
}
