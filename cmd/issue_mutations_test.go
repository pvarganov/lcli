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

// newCaptureMultiServer creates a server that routes by query keyword and captures issueCreate variables.
func newCaptureMultiServer(t *testing.T, captured *map[string]any, responses map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		if strings.Contains(req.Query, "issueCreate") && captured != nil {
			*captured = req.Variables
		}

		for keyword, data := range responses {
			if strings.Contains(req.Query, keyword) {
				w.Header().Set("Content-Type", "application/json")
				body, _ := json.Marshal(map[string]any{"data": data})
				w.Write(body)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
}

func baseCreateResponses() map[string]any {
	return map[string]any{
		"GetTeams": map[string]any{
			"teams": map[string]any{
				"nodes": []map[string]any{
					{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
		"issueCreate": issueCreateResponseData("ENG-1"),
	}
}

func resetCreateFlags() {
	createTitle = "Test issue"
	createTeam = "ENG"
	createDescription = ""
	createAssignee = ""
	createPriority = 0
	createDueDate = ""
	createEstimate = 0
	createLabels = ""
	createParent = ""
	createState = ""
	createCycleID = ""
	createProjectID = ""
	createMilestoneID = ""
	// Reset Changed state for flags that use cmd.Flags().Changed()
	if f := issueCreateCmd.Flags().Lookup("estimate"); f != nil {
		f.Changed = false
	}
}

func TestIssueCreateCmdDueDate(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createDueDate = "2026-04-01"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["dueDate"]; !ok || v != "2026-04-01" {
		t.Errorf("expected dueDate=2026-04-01, got %v", input["dueDate"])
	}
}

func TestIssueCreateCmdDueDateNotSentWhenEmpty(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if _, ok := input["dueDate"]; ok {
		t.Errorf("expected dueDate not sent when empty")
	}
}

func TestIssueCreateCmdEstimate(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	// Must use Flags().Set to mark flag as Changed
	if err := issueCreateCmd.Flags().Set("estimate", "5"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}
	defer func() {
		_ = issueCreateCmd.Flags().Set("estimate", "0")
		if f := issueCreateCmd.Flags().Lookup("estimate"); f != nil {
			f.Changed = false
		}
	}()

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; !ok {
		t.Errorf("expected estimate to be sent when flag is set")
	}
}

func TestIssueCreateCmdEstimateNotSentByDefault(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	// Do NOT set estimate flag — it should not appear in request
	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; ok {
		t.Errorf("expected estimate not sent when flag not set")
	}
}

func TestIssueCreateCmdParent(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createParent = "parent-abc"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["parentId"]; !ok || v != "parent-abc" {
		t.Errorf("expected parentId=parent-abc, got %v", input["parentId"])
	}
}

func TestIssueCreateCmdCycleID(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createCycleID = "cycle-xyz"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["cycleId"]; !ok || v != "cycle-xyz" {
		t.Errorf("expected cycleId=cycle-xyz, got %v", input["cycleId"])
	}
}

func TestIssueCreateCmdProjectID(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createProjectID = "proj-999"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectId"]; !ok || v != "proj-999" {
		t.Errorf("expected projectId=proj-999, got %v", input["projectId"])
	}
}

func TestIssueCreateCmdMilestoneID(t *testing.T) {
	var captured map[string]any
	srv := newCaptureMultiServer(t, &captured, baseCreateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createMilestoneID = "ms-007"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectMilestoneId"]; !ok || v != "ms-007" {
		t.Errorf("expected projectMilestoneId=ms-007, got %v", input["projectMilestoneId"])
	}
}

func TestIssueCreateCmdLabels(t *testing.T) {
	var captured map[string]any
	responses := baseCreateResponses()
	responses["issueLabels"] = map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "lbl-1", "name": "Bug"},
			},
		},
	}
	// The labels query returns issueLabels data
	responses["GetIssueLabels"] = map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "lbl-1", "name": "Bug"},
			},
		},
	}
	srv := newCaptureMultiServer(t, &captured, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createLabels = "Bug"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	rawLabels, ok := input["labelIds"]
	if !ok {
		t.Fatal("expected labelIds in request")
	}
	labels, _ := rawLabels.([]any)
	if len(labels) != 1 {
		t.Errorf("expected 1 label, got %v", labels)
	}
}

// ---- UpdateIssue client-layer tests ----

func issueUpdateResponseData() any {
	return map[string]any{
		"issueUpdate": map[string]any{
			"success": true,
			"issue": map[string]any{
				"id":          "existing-id",
				"identifier":  "ENG-42",
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

func TestUpdateIssueInputDescription(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{Description: "new desc"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["description"]; !ok || v != "new desc" {
		t.Errorf("expected description=new desc, got %v", input["description"])
	}
}

func TestUpdateIssueInputDueDate(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{DueDate: "2026-05-01"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["dueDate"]; !ok || v != "2026-05-01" {
		t.Errorf("expected dueDate=2026-05-01, got %v", input["dueDate"])
	}
}

func TestUpdateIssueInputEstimate(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	est := 5
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{Estimate: &est})

	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; !ok {
		t.Errorf("expected estimate to be sent")
	}
}

func TestUpdateIssueInputEstimateNotSentWhenNil(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{Title: "x"})

	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; ok {
		t.Errorf("expected estimate not sent when nil")
	}
}

func TestUpdateIssueInputParentID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{ParentID: "parent-xyz"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["parentId"]; !ok || v != "parent-xyz" {
		t.Errorf("expected parentId=parent-xyz, got %v", input["parentId"])
	}
}

func TestUpdateIssueInputCycleID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{CycleID: "cycle-456"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["cycleId"]; !ok || v != "cycle-456" {
		t.Errorf("expected cycleId=cycle-456, got %v", input["cycleId"])
	}
}

func TestUpdateIssueInputProjectID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{ProjectID: "proj-111"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectId"]; !ok || v != "proj-111" {
		t.Errorf("expected projectId=proj-111, got %v", input["projectId"])
	}
}

func TestUpdateIssueInputMilestoneID(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{MilestoneID: "ms-222"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectMilestoneId"]; !ok || v != "ms-222" {
		t.Errorf("expected projectMilestoneId=ms-222, got %v", input["projectMilestoneId"])
	}
}

func TestUpdateIssueInputAddedLabelIDs(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{AddedLabelIDs: []string{"lbl-1", "lbl-2"}})

	input, _ := captured["input"].(map[string]any)
	raw, ok := input["addedLabelIds"]
	if !ok {
		t.Fatal("expected addedLabelIds in request")
	}
	labels, _ := raw.([]any)
	if len(labels) != 2 {
		t.Errorf("expected 2 addedLabelIds, got %v", labels)
	}
}

func TestUpdateIssueInputRemovedLabelIDs(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{RemovedLabelIDs: []string{"lbl-3"}})

	input, _ := captured["input"].(map[string]any)
	raw, ok := input["removedLabelIds"]
	if !ok {
		t.Fatal("expected removedLabelIds in request")
	}
	labels, _ := raw.([]any)
	if len(labels) != 1 {
		t.Errorf("expected 1 removedLabelId, got %v", labels)
	}
}

func TestUpdateIssueInputSnoozedUntilAt(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{SnoozedUntilAt: "2026-04-01T10:00:00Z"})

	input, _ := captured["input"].(map[string]any)
	if v, ok := input["snoozedUntilAt"]; !ok || v != "2026-04-01T10:00:00Z" {
		t.Errorf("expected snoozedUntilAt=2026-04-01T10:00:00Z, got %v", input["snoozedUntilAt"])
	}
}

func TestUpdateIssueInputEmptyFieldsNotSent(t *testing.T) {
	var captured map[string]any
	srv := captureRequestServer(t, &captured, issueUpdateResponseData())
	defer srv.Close()

	c := client.NewWithURL("tok", srv.URL)
	_, _ = c.UpdateIssue("ENG-42", client.UpdateIssueInput{Title: "only-title"})

	input, _ := captured["input"].(map[string]any)
	for _, field := range []string{"description", "dueDate", "estimate", "parentId", "cycleId", "projectId", "projectMilestoneId", "addedLabelIds", "removedLabelIds", "snoozedUntilAt"} {
		if _, ok := input[field]; ok {
			t.Errorf("expected %s not sent when empty, but it was present", field)
		}
	}
}

// ---- IssueUpdateCmd cmd-layer tests ----

func issueUpdateCmdResponseData() any {
	return map[string]any{
		"issueUpdate": map[string]any{
			"success": true,
			"issue": map[string]any{
				"id":          "existing-id",
				"identifier":  "ENG-42",
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

func getIssueResponseData() any {
	return map[string]any{
		"issue": map[string]any{
			"id":          "existing-id",
			"identifier":  "ENG-42",
			"title":       "Test",
			"description": "",
			"updatedAt":   "2024-01-01T00:00:00Z",
			"priority":    0,
			"state":       map[string]any{"name": "Todo", "type": "unstarted"},
			"assignee":    nil,
			"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}
}

// newCaptureUpdateMultiServer routes by query keyword and captures issueUpdate variables.
func newCaptureUpdateMultiServer(t *testing.T, captured *map[string]any, responses map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Query     string         `json:"query"`
			Variables map[string]any `json:"variables"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)

		if strings.Contains(req.Query, "issueUpdate") && captured != nil {
			*captured = req.Variables
		}

		for keyword, data := range responses {
			if strings.Contains(req.Query, keyword) {
				w.Header().Set("Content-Type", "application/json")
				body, _ := json.Marshal(map[string]any{"data": data})
				w.Write(body)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{}}`))
	}))
}

func resetUpdateFlags() {
	updateTitle = ""
	updateStatus = ""
	updateAssignee = ""
	updatePriority = 0
	updateDescription = ""
	updateDueDate = ""
	updateEstimate = 0
	updateParent = ""
	updateCycleID = ""
	updateProjectID = ""
	updateMilestoneID = ""
	updateAddLabels = ""
	updateRemoveLabels = ""
	updateSnoozeUntil = ""
	for _, name := range []string{"priority", "estimate"} {
		if f := issueUpdateCmd.Flags().Lookup(name); f != nil {
			f.Changed = false
		}
	}
}

func baseUpdateResponses() map[string]any {
	return map[string]any{
		"issueUpdate": issueUpdateCmdResponseData(),
	}
}

func TestIssueUpdateCmdDescription(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateDescription = "New description text"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["description"]; !ok || v != "New description text" {
		t.Errorf("expected description=New description text, got %v", input["description"])
	}
}

func TestIssueUpdateCmdDueDate(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateDueDate = "2026-05-01"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["dueDate"]; !ok || v != "2026-05-01" {
		t.Errorf("expected dueDate=2026-05-01, got %v", input["dueDate"])
	}
}

func TestIssueUpdateCmdEstimate(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	if err := issueUpdateCmd.Flags().Set("estimate", "5"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}
	defer func() {
		_ = issueUpdateCmd.Flags().Set("estimate", "0")
		if f := issueUpdateCmd.Flags().Lookup("estimate"); f != nil {
			f.Changed = false
		}
	}()

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; !ok {
		t.Errorf("expected estimate to be sent when flag is set")
	}
}

func TestIssueUpdateCmdEstimateNotSentByDefault(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateTitle = "some title"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if _, ok := input["estimate"]; ok {
		t.Errorf("expected estimate not sent when flag not set")
	}
}

func TestIssueUpdateCmdParent(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateParent = "parent-id-abc"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["parentId"]; !ok || v != "parent-id-abc" {
		t.Errorf("expected parentId=parent-id-abc, got %v", input["parentId"])
	}
}

func TestIssueUpdateCmdCycleID(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateCycleID = "cycle-xyz"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["cycleId"]; !ok || v != "cycle-xyz" {
		t.Errorf("expected cycleId=cycle-xyz, got %v", input["cycleId"])
	}
}

func TestIssueUpdateCmdProjectID(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateProjectID = "proj-001"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectId"]; !ok || v != "proj-001" {
		t.Errorf("expected projectId=proj-001, got %v", input["projectId"])
	}
}

func TestIssueUpdateCmdMilestoneID(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateMilestoneID = "ms-999"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["projectMilestoneId"]; !ok || v != "ms-999" {
		t.Errorf("expected projectMilestoneId=ms-999, got %v", input["projectMilestoneId"])
	}
}

func TestIssueUpdateCmdSnoozeUntil(t *testing.T) {
	var captured map[string]any
	srv := newCaptureUpdateMultiServer(t, &captured, baseUpdateResponses())
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateSnoozeUntil = "2026-04-01T10:00:00Z"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["snoozedUntilAt"]; !ok || v != "2026-04-01T10:00:00Z" {
		t.Errorf("expected snoozedUntilAt=2026-04-01T10:00:00Z, got %v", input["snoozedUntilAt"])
	}
}

func TestIssueUpdateCmdAddLabels(t *testing.T) {
	var captured map[string]any
	responses := baseUpdateResponses()
	// Use unique substrings that won't conflict: "issue(id:" is in GetIssue but not GetIssueLabels
	responses["issue(id:"] = getIssueResponseData()
	responses["GetIssueLabels"] = map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "lbl-1", "name": "Bug"},
				{"id": "lbl-2", "name": "Feature"},
			},
		},
	}
	srv := newCaptureUpdateMultiServer(t, &captured, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateAddLabels = "Bug,Feature"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	raw, ok := input["addedLabelIds"]
	if !ok {
		t.Fatal("expected addedLabelIds in request")
	}
	labels, _ := raw.([]any)
	if len(labels) != 2 {
		t.Errorf("expected 2 addedLabelIds, got %v", labels)
	}
}

func TestIssueUpdateCmdRemoveLabels(t *testing.T) {
	var captured map[string]any
	responses := baseUpdateResponses()
	responses["issue(id:"] = getIssueResponseData()
	responses["GetIssueLabels"] = map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "lbl-3", "name": "Deprecated"},
			},
		},
	}
	srv := newCaptureUpdateMultiServer(t, &captured, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetUpdateFlags()
	updateRemoveLabels = "Deprecated"

	if err := issueUpdateCmd.RunE(issueUpdateCmd, []string{"ENG-42"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	raw, ok := input["removedLabelIds"]
	if !ok {
		t.Fatal("expected removedLabelIds in request")
	}
	labels, _ := raw.([]any)
	if len(labels) != 1 {
		t.Errorf("expected 1 removedLabelId, got %v", labels)
	}
}

// ---- IssueCreateCmd state test ----

func TestIssueCreateCmdState(t *testing.T) {
	var captured map[string]any
	responses := baseCreateResponses()
	responses["GetWorkflowStates"] = map[string]any{
		"workflowStates": map[string]any{
			"nodes": []map[string]any{
				{"id": "state-in-progress", "name": "In Progress"},
			},
		},
	}
	srv := newCaptureMultiServer(t, &captured, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client { return client.NewWithURL(tok, srv.URL) }
	defer func() { newLinearClient = origFactory }()
	token = "test-token"
	defer func() { token = "" }()

	resetCreateFlags()
	createState = "In Progress"

	if err := issueCreateCmd.RunE(issueCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input, _ := captured["input"].(map[string]any)
	if v, ok := input["stateId"]; !ok || v != "state-in-progress" {
		t.Errorf("expected stateId=state-in-progress, got %v", input["stateId"])
	}
}
