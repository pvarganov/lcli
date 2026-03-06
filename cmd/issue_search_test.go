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

func newSearchTestServer(t *testing.T, nodes []map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"issueSearch": map[string]any{
					"nodes":    nodes,
					"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
				},
			},
		})
	}))
}

func TestIssueSearchCmdTableOutput(t *testing.T) {
	now := time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC)
	nodes := []map[string]any{
		{
			"id":         "i1",
			"identifier": "ENG-1",
			"title":      "Search result",
			"updatedAt":  now.Format(time.RFC3339),
			"priority":   2,
			"state":      map[string]any{"name": "In Progress", "type": "started"},
			"assignee":   nil,
			"team":       map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}
	srv := newSearchTestServer(t, nodes)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	var buf bytes.Buffer
	issueSearchCmd.SetOut(&buf)

	if err := issueSearchCmd.RunE(issueSearchCmd, []string{"result"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected ENG-1 in output, got: %s", out)
	}
	if !strings.Contains(out, "Search result") {
		t.Errorf("expected title in output, got: %s", out)
	}
	if !strings.Contains(out, "In Progress") {
		t.Errorf("expected status in output, got: %s", out)
	}
}

func TestIssueSearchCmdJSONOutput(t *testing.T) {
	now := time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC)
	nodes := []map[string]any{
		{
			"id":         "i2",
			"identifier": "ENG-2",
			"title":      "JSON issue",
			"updatedAt":  now.Format(time.RFC3339),
			"priority":   0,
			"state":      map[string]any{"name": "Todo", "type": "unstarted"},
			"assignee":   nil,
			"team":       map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}
	srv := newSearchTestServer(t, nodes)
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

	var buf bytes.Buffer
	issueSearchCmd.SetOut(&buf)

	if err := issueSearchCmd.RunE(issueSearchCmd, []string{"json"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	var issues []map[string]any
	if err := json.Unmarshal([]byte(out), &issues); err != nil {
		t.Fatalf("expected valid JSON, got: %v\noutput: %s", err, out)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0]["identifier"] != "ENG-2" {
		t.Errorf("expected ENG-2, got %v", issues[0]["identifier"])
	}
}

func TestIssueSearchCmdEmptyResults(t *testing.T) {
	srv := newSearchTestServer(t, []map[string]any{})
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	var buf bytes.Buffer
	issueSearchCmd.SetOut(&buf)

	if err := issueSearchCmd.RunE(issueSearchCmd, []string{"noresults"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message in output, got: %s", out)
	}
}

func TestIssueSearchCmdNoToken(t *testing.T) {
	token = ""
	err := issueSearchCmd.RunE(issueSearchCmd, []string{"query"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}
