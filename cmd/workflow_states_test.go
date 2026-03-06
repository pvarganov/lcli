package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestWorkflowStatesListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListWorkflowStates": map[string]any{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{
					{
						"id":    "ws1",
						"name":  "Todo",
						"type":  "unstarted",
						"color": "#e2e2e2",
						"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					},
					{
						"id":    "ws2",
						"name":  "In Progress",
						"type":  "started",
						"color": "#f2c94c",
						"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					},
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

	_ = workflowStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = workflowStatesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	workflowStatesListCmd.SetOut(&buf)

	if err := workflowStatesListCmd.RunE(workflowStatesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Todo") {
		t.Errorf("expected 'Todo' in output, got: %s", out)
	}
	if !strings.Contains(out, "In Progress") {
		t.Errorf("expected 'In Progress' in output, got: %s", out)
	}
	if !strings.Contains(out, "unstarted") {
		t.Errorf("expected 'unstarted' in output, got: %s", out)
	}
}

func TestWorkflowStatesListJSON(t *testing.T) {
	responses := map[string]any{
		"ListWorkflowStates": map[string]any{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{
					{
						"id":    "ws1",
						"name":  "Todo",
						"type":  "unstarted",
						"color": "#e2e2e2",
						"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					},
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

	origFormat := outputFormat
	outputFormat = "json"
	defer func() { outputFormat = origFormat }()

	_ = workflowStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = workflowStatesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	workflowStatesListCmd.SetOut(&buf)

	if err := workflowStatesListCmd.RunE(workflowStatesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestWorkflowStatesListEmpty(t *testing.T) {
	responses := map[string]any{
		"ListWorkflowStates": map[string]any{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{},
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

	_ = workflowStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = workflowStatesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	workflowStatesListCmd.SetOut(&buf)

	if err := workflowStatesListCmd.RunE(workflowStatesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestWorkflowStatesListNoToken(t *testing.T) {
	token = ""
	_ = workflowStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = workflowStatesListCmd.Flags().Set("team", "") }()

	err := workflowStatesListCmd.RunE(workflowStatesListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestWorkflowStatesListNoTeam(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = workflowStatesListCmd.Flags().Set("team", "")

	err := workflowStatesListCmd.RunE(workflowStatesListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no team")
	}
}

func TestWorkflowStatesCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"CreateWorkflowState": map[string]any{
			"workflowStateCreate": map[string]any{
				"success": true,
				"workflowState": map[string]any{
					"id":    "ws3",
					"name":  "Review",
					"type":  "started",
					"color": "#f2c94c",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = workflowStatesCreateCmd.Flags().Set("team", "t1")
	_ = workflowStatesCreateCmd.Flags().Set("name", "Review")
	_ = workflowStatesCreateCmd.Flags().Set("type", "started")
	defer func() {
		_ = workflowStatesCreateCmd.Flags().Set("team", "")
		_ = workflowStatesCreateCmd.Flags().Set("name", "")
		_ = workflowStatesCreateCmd.Flags().Set("type", "")
	}()

	var buf bytes.Buffer
	workflowStatesCreateCmd.SetOut(&buf)

	if err := workflowStatesCreateCmd.RunE(workflowStatesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "создан") {
		t.Errorf("expected 'создан' in output, got: %s", out)
	}
}

func TestWorkflowStatesCreateNoToken(t *testing.T) {
	token = ""
	_ = workflowStatesCreateCmd.Flags().Set("team", "t1")
	_ = workflowStatesCreateCmd.Flags().Set("name", "Review")
	_ = workflowStatesCreateCmd.Flags().Set("type", "started")
	defer func() {
		_ = workflowStatesCreateCmd.Flags().Set("team", "")
		_ = workflowStatesCreateCmd.Flags().Set("name", "")
		_ = workflowStatesCreateCmd.Flags().Set("type", "")
	}()

	err := workflowStatesCreateCmd.RunE(workflowStatesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestWorkflowStatesCreateJSON(t *testing.T) {
	responses := map[string]any{
		"CreateWorkflowState": map[string]any{
			"workflowStateCreate": map[string]any{
				"success": true,
				"workflowState": map[string]any{
					"id":    "ws3",
					"name":  "Review",
					"type":  "started",
					"color": "#f2c94c",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	origFormat := outputFormat
	outputFormat = "json"
	defer func() { outputFormat = origFormat }()

	_ = workflowStatesCreateCmd.Flags().Set("team", "t1")
	_ = workflowStatesCreateCmd.Flags().Set("name", "Review")
	_ = workflowStatesCreateCmd.Flags().Set("type", "started")
	defer func() {
		_ = workflowStatesCreateCmd.Flags().Set("team", "")
		_ = workflowStatesCreateCmd.Flags().Set("name", "")
		_ = workflowStatesCreateCmd.Flags().Set("type", "")
	}()

	var buf bytes.Buffer
	workflowStatesCreateCmd.SetOut(&buf)

	if err := workflowStatesCreateCmd.RunE(workflowStatesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestWorkflowStatesUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"UpdateWorkflowState": map[string]any{
			"workflowStateUpdate": map[string]any{
				"success": true,
				"workflowState": map[string]any{
					"id":    "ws1",
					"name":  "In Review",
					"type":  "started",
					"color": "#f2c94c",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = workflowStatesUpdateCmd.Flags().Set("name", "In Review")
	defer func() { _ = workflowStatesUpdateCmd.Flags().Set("name", "") }()

	var buf bytes.Buffer
	workflowStatesUpdateCmd.SetOut(&buf)

	if err := workflowStatesUpdateCmd.RunE(workflowStatesUpdateCmd, []string{"ws1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "обновлён") {
		t.Errorf("expected 'обновлён' in output, got: %s", out)
	}
}

func TestWorkflowStatesUpdateNoToken(t *testing.T) {
	token = ""
	err := workflowStatesUpdateCmd.RunE(workflowStatesUpdateCmd, []string{"ws1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestWorkflowStatesUpdateJSON(t *testing.T) {
	responses := map[string]any{
		"UpdateWorkflowState": map[string]any{
			"workflowStateUpdate": map[string]any{
				"success": true,
				"workflowState": map[string]any{
					"id":    "ws1",
					"name":  "In Review",
					"type":  "started",
					"color": "#f2c94c",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	origFormat := outputFormat
	outputFormat = "json"
	defer func() { outputFormat = origFormat }()

	var buf bytes.Buffer
	workflowStatesUpdateCmd.SetOut(&buf)

	if err := workflowStatesUpdateCmd.RunE(workflowStatesUpdateCmd, []string{"ws1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestWorkflowStatesArchiveSuccess(t *testing.T) {
	responses := map[string]any{
		"ArchiveWorkflowState": map[string]any{
			"workflowStateArchive": map[string]any{
				"success": true,
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

	var buf bytes.Buffer
	workflowStatesArchiveCmd.SetOut(&buf)

	if err := workflowStatesArchiveCmd.RunE(workflowStatesArchiveCmd, []string{"ws1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "архивирован") {
		t.Errorf("expected 'архивирован' in output, got: %s", out)
	}
}

func TestWorkflowStatesArchiveNoToken(t *testing.T) {
	token = ""
	err := workflowStatesArchiveCmd.RunE(workflowStatesArchiveCmd, []string{"ws1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
