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

// captureRequestServer создаёт тестовый сервер, который захватывает тело GraphQL-запроса.
func captureRequestServer(t *testing.T, captured *map[string]any, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		*captured = req.Variables
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func issueCreateResponseData(identifier string) any {
	return map[string]any{
		"issueCreate": map[string]any{
			"success": true,
			"issue": map[string]any{
				"id":          "new-id",
				"identifier":  identifier,
				"title":       "Test",
				"description": "",
				"updatedAt":   "2024-01-01T00:00:00Z",
				"priority":    0,
				"state":       map[string]any{"name": "Todo", "type": "unstarted"},
				"assignee":    nil,
				"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
			},
		},
	}
}

func TestCreateIssueInputDueDate(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	est := 3
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:  "t1",
		Title:   "Test",
		DueDate: "2026-04-01",
		Estimate: &est,
	})

	input, _ := captured["input"].(map[string]any)
	if input == nil {
		t.Fatal("no input in captured request")
	}
	if v, ok := input["dueDate"]; !ok || v != "2026-04-01" {
		t.Errorf("expected dueDate=2026-04-01, got %v", input["dueDate"])
	}
	if v, ok := input["estimate"]; !ok || v == nil {
		t.Errorf("expected estimate to be set, got %v", input["estimate"])
	}
}

func TestCreateIssueInputEstimateZeroNotSent(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID: "t1",
		Title:  "Test",
		// Estimate is nil — should not be sent
	})

	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; ok {
		t.Errorf("expected estimate not sent when nil, but it was present")
	}
	if _, ok := input["dueDate"]; ok {
		t.Errorf("expected dueDate not sent when empty, but it was present")
	}
}

func TestCreateIssueInputLabelIDs(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:   "t1",
		Title:    "Test",
		LabelIDs: []string{"lbl-1", "lbl-2"},
	})

	input, _ := captured["input"].(map[string]any)
	rawLabels, ok := input["labelIds"]
	if !ok {
		t.Fatal("expected labelIds in request, not found")
	}
	labels, _ := rawLabels.([]any)
	if len(labels) != 2 {
		t.Errorf("expected 2 labelIds, got %v", labels)
	}
}

func TestCreateIssueInputLabelIDsNotSentWhenEmpty(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID: "t1",
		Title:  "Test",
	})

	input, _ := captured["input"].(map[string]any)
	if _, ok := input["labelIds"]; ok {
		t.Errorf("expected labelIds not sent when empty, but it was present")
	}
}

func TestCreateIssueInputParentID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:   "t1",
		Title:    "Test",
		ParentID: "parent-id-123",
	})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["parentId"]; !ok || v != "parent-id-123" {
		t.Errorf("expected parentId=parent-id-123, got %v", input["parentId"])
	}
}

func TestCreateIssueInputStateID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:  "t1",
		Title:   "Test",
		StateID: "state-id-xyz",
	})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["stateId"]; !ok || v != "state-id-xyz" {
		t.Errorf("expected stateId=state-id-xyz, got %v", input["stateId"])
	}
}

func TestCreateIssueInputCycleID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:  "t1",
		Title:   "Test",
		CycleID: "cycle-abc",
	})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["cycleId"]; !ok || v != "cycle-abc" {
		t.Errorf("expected cycleId=cycle-abc, got %v", input["cycleId"])
	}
}

func TestCreateIssueInputProjectID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:    "t1",
		Title:     "Test",
		ProjectID: "proj-001",
	})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectId"]; !ok || v != "proj-001" {
		t.Errorf("expected projectId=proj-001, got %v", input["projectId"])
	}
}

func TestCreateIssueInputMilestoneID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueCreateResponseData("ENG-1"))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.CreateIssue(client.CreateIssueInput{
		TeamID:      "t1",
		Title:       "Test",
		MilestoneID: "milestone-007",
	})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectMilestoneId"]; !ok || v != "milestone-007" {
		t.Errorf("expected projectMilestoneId=milestone-007, got %v", input["projectMilestoneId"])
	}
}

func TestFindLabelsByNames(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"issueLabels": map[string]any{
					"nodes": []map[string]any{
						{"id": "lbl-1", "name": "Bug"},
						{"id": "lbl-2", "name": "Feature"},
					},
				},
			},
		})
		w.Write(body)
	}))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	ids, err := c.FindLabelsByNames("team-1", []string{"Bug", "Feature"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 ids, got %d", len(ids))
	}
	if ids[0] != "lbl-1" || ids[1] != "lbl-2" {
		t.Errorf("unexpected ids: %v", ids)
	}
}

func TestFindLabelsByNamesNotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"issueLabels": map[string]any{
					"nodes": []map[string]any{
						{"id": "lbl-1", "name": "Bug"},
					},
				},
			},
		})
		w.Write(body)
	}))
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, err := c.FindLabelsByNames("team-1", []string{"NonExistent"})
	if err == nil {
		t.Fatal("expected error for not found label, got nil")
	}
	if !strings.Contains(err.Error(), "метка не найдена") {
		t.Errorf("expected 'метка не найдена' in error, got: %v", err)
	}
}
