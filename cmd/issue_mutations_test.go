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

// multiHandlerServer создаёт тестовый сервер, который возвращает разные ответы
// в зависимости от тела запроса (операции GraphQL).
func newMultiResponseServer(t *testing.T, responses map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query string `json:"query"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		// Ищем совпадение по ключевому слову в запросе
		for keyword, data := range responses {
			if strings.Contains(req.Query, keyword) {
				w.Header().Set("Content-Type", "application/json")
				body, _ := json.Marshal(map[string]any{"data": data})
				w.Write(body)
				return
			}
		}
		// Дефолтный пустой ответ
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
}

func TestIssueCreateOutput(t *testing.T) {
	now := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)

	responses := map[string]any{
		"GetTeams": map[string]any{
			"teams": map[string]any{
				"nodes": []map[string]any{
					{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
		"issueCreate": map[string]any{
			"issueCreate": map[string]any{
				"success": true,
				"issue": map[string]any{
					"id":          "new-id",
					"identifier":  "ENG-101",
					"title":       "My new issue",
					"description": "",
					"updatedAt":   now.Format(time.RFC3339),
					"priority":    0,
					"state":       map[string]any{"name": "Todo", "type": "unstarted"},
					"assignee":    nil,
					"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
	}

	srv := newMultiResponseServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	// Сбрасываем флаги перед тестом
	createTitle = "My new issue"
	createTeam = "ENG"
	createDescription = ""
	createAssignee = ""
	createPriority = 0

	cmd := issueCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-101") {
		t.Errorf("expected ENG-101 in output, got: %s", out)
	}
	if !strings.Contains(out, "My new issue") {
		t.Errorf("expected 'My new issue' in output, got: %s", out)
	}
}

func TestIssueCreateNoToken(t *testing.T) {
	token = ""
	createTitle = "Test"
	createTeam = "ENG"
	err := issueCreateCmd.RunE(issueCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestIssueCreateTeamNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"teams": map[string]any{
					"nodes": []map[string]any{},
				},
			},
		})
		w.Write(body)
	}))
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	createTitle = "Test issue"
	createTeam = "NONEXISTENT"

	err := issueCreateCmd.RunE(issueCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error for nonexistent team")
	}
	if !strings.Contains(err.Error(), "не найдена") {
		t.Errorf("expected 'не найдена' in error, got: %v", err)
	}
}

func TestIssueUpdateOutput(t *testing.T) {
	now := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)

	responses := map[string]any{
		"issueUpdate": map[string]any{
			"issueUpdate": map[string]any{
				"success": true,
				"issue": map[string]any{
					"id":          "existing-id",
					"identifier":  "ENG-42",
					"title":       "Updated issue title",
					"description": "",
					"updatedAt":   now.Format(time.RFC3339),
					"priority":    2,
					"state":       map[string]any{"name": "In Progress", "type": "started"},
					"assignee":    nil,
					"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
	}

	srv := newMultiResponseServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	updateTitle = "Updated issue title"
	updateStatus = ""
	updateAssignee = ""
	updatePriority = 0

	cmd := issueUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-42") {
		t.Errorf("expected ENG-42 in output, got: %s", out)
	}
	if !strings.Contains(out, "Updated issue title") {
		t.Errorf("expected 'Updated issue title' in output, got: %s", out)
	}
}

func TestIssueUpdateNoToken(t *testing.T) {
	token = ""
	err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}
