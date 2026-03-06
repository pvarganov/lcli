package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pavelvarganov/lcli/internal/client"
)

func newIssuesTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestIssuesListOutput(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	responseData := map[string]any{
		"issues": map[string]any{
			"nodes": []map[string]any{
				{
					"id":         "abc123",
					"identifier": "ENG-42",
					"title":      "Fix login bug",
					"updatedAt":  now.Format(time.RFC3339),
					"priority":   2,
					"state":      map[string]any{"name": "In Progress", "type": "started"},
					"assignee":   map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@example.com"},
					"team":       map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ENG-42", "Fix login bug", "In Progress", "Alice Smith", "High"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestIssuesListNoToken(t *testing.T) {
	token = ""
	err := issuesListCmd.RunE(issuesListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestIssueViewOutput(t *testing.T) {
	now := time.Date(2024, 3, 5, 14, 30, 0, 0, time.UTC)
	responseData := map[string]any{
		"issue": map[string]any{
			"id":          "xyz789",
			"identifier":  "ENG-99",
			"title":       "Update docs",
			"description": "Add API documentation",
			"updatedAt":   now.Format(time.RFC3339),
			"priority":    3,
			"state":       map[string]any{"name": "Todo", "type": "unstarted"},
			"assignee":    nil,
			"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-99"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ENG-99", "Update docs", "Todo", "Medium", "Engineering", "Add API documentation"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestIssueViewNoToken(t *testing.T) {
	token = ""
	err := issueViewCmd.RunE(issueViewCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestIssuesListJSONOutput(t *testing.T) {
	now := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	responseData := map[string]any{
		"issues": map[string]any{
			"nodes": []map[string]any{
				{
					"id":         "abc123",
					"identifier": "ENG-42",
					"title":      "Fix login bug",
					"updatedAt":  now.Format(time.RFC3339),
					"priority":   2,
					"state":      map[string]any{"name": "In Progress", "type": "started"},
					"assignee":   map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@example.com"},
					"team":       map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	outputFormat = "json"
	defer func() { outputFormat = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	// Проверяем, что вывод — корректный JSON
	var issues []map[string]any
	if err := json.Unmarshal([]byte(out), &issues); err != nil {
		t.Fatalf("expected valid JSON output, got error: %v\noutput: %s", err, out)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue in JSON, got %d", len(issues))
	}
	if issues[0]["identifier"] != "ENG-42" {
		t.Errorf("expected identifier ENG-42, got %v", issues[0]["identifier"])
	}
}
