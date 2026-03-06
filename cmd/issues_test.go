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

// captureIssuesListServer создаёт тестовый сервер, захватывающий переменные запроса ListIssues.
func captureIssuesListServer(t *testing.T, captured *map[string]any, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		if strings.Contains(req.Query, "ListIssues") {
			*captured = req.Variables
		}
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func issuesListResponseData() any {
	return map[string]any{
		"issues": map[string]any{
			"nodes":    []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
		},
	}
}

func TestIssuesListFilterPriority(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesPriority = 2
	defer func() { issuesPriority = -1 }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	filter, _ := captured["filter"].(map[string]any)
	if filter == nil {
		t.Fatal("expected filter in request variables")
	}
	priority, _ := filter["priority"].(map[string]any)
	if priority == nil {
		t.Fatal("expected priority filter")
	}
	if priority["eq"] == nil {
		t.Errorf("expected priority.eq to be set, got %v", priority)
	}
}

func TestIssuesListFilterPriorityNotSentWhenDefault(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesPriority = -1
	defer func() { issuesPriority = -1 }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	filter, _ := captured["filter"].(map[string]any)
	if filter != nil {
		if _, ok := filter["priority"]; ok {
			t.Errorf("expected priority not sent when -1, but it was present")
		}
	}
}

func TestIssuesListFilterLabel(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesLabel = "Bug"
	defer func() { issuesLabel = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	filter, _ := captured["filter"].(map[string]any)
	if filter == nil {
		t.Fatal("expected filter in request variables")
	}
	labels, _ := filter["labels"].(map[string]any)
	if labels == nil {
		t.Fatal("expected labels filter")
	}
	some, _ := labels["some"].(map[string]any)
	if some == nil {
		t.Fatal("expected labels.some filter")
	}
	name, _ := some["name"].(map[string]any)
	if name == nil || name["eqIgnoreCase"] != "Bug" {
		t.Errorf("expected labels.some.name.eqIgnoreCase=Bug, got %v", some)
	}
}

func TestIssuesListFilterProjectID(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesProjectID = "proj-123"
	defer func() { issuesProjectID = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	filter, _ := captured["filter"].(map[string]any)
	if filter == nil {
		t.Fatal("expected filter in request variables")
	}
	project, _ := filter["project"].(map[string]any)
	if project == nil {
		t.Fatal("expected project filter")
	}
	id, _ := project["id"].(map[string]any)
	if id == nil || id["eq"] != "proj-123" {
		t.Errorf("expected project.id.eq=proj-123, got %v", project)
	}
}

func TestIssuesListFilterCycleID(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesCycleID = "cycle-456"
	defer func() { issuesCycleID = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	filter, _ := captured["filter"].(map[string]any)
	if filter == nil {
		t.Fatal("expected filter in request variables")
	}
	cycle, _ := filter["cycle"].(map[string]any)
	if cycle == nil {
		t.Fatal("expected cycle filter")
	}
	id, _ := cycle["id"].(map[string]any)
	if id == nil || id["eq"] != "cycle-456" {
		t.Errorf("expected cycle.id.eq=cycle-456, got %v", cycle)
	}
}

func TestIssuesListFilterCreator(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesCreator = "Alice"
	defer func() { issuesCreator = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	filter, _ := captured["filter"].(map[string]any)
	if filter == nil {
		t.Fatal("expected filter in request variables")
	}
	creator, _ := filter["creator"].(map[string]any)
	if creator == nil {
		t.Fatal("expected creator filter")
	}
	dn, _ := creator["displayName"].(map[string]any)
	if dn == nil || dn["eq"] != "Alice" {
		t.Errorf("expected creator.displayName.eq=Alice, got %v", creator)
	}
}

func TestIssuesListFilterOrderBy(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesOrderBy = "updatedAt"
	defer func() { issuesOrderBy = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	if captured["orderBy"] != "updatedAt" {
		t.Errorf("expected orderBy=updatedAt, got %v", captured["orderBy"])
	}
}

func TestIssuesListFilterOrderByNotSentWhenEmpty(t *testing.T) {
	var captured map[string]any
	srv := captureIssuesListServer(t, &captured, issuesListResponseData())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	issuesOrderBy = ""
	defer func() { issuesOrderBy = "" }()

	cmd := issuesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	_ = cmd.RunE(cmd, []string{})

	if _, ok := captured["orderBy"]; ok {
		t.Errorf("expected orderBy not sent when empty, but it was present")
	}
}

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
