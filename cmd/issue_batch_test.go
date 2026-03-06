package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestIssueBatchUpdateNoToken(t *testing.T) {
	token = ""
	issueBatchUpdateCmd.Flags().Set("ids", "issue-1")
	issueBatchUpdateCmd.Flags().Set("status", "In Progress")
	err := issueBatchUpdateCmd.RunE(issueBatchUpdateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestIssueBatchUpdateNoIDs(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()
	batchUpdateIDs = ""
	batchUpdateStatus = "In Progress"
	err := issueBatchUpdateCmd.RunE(issueBatchUpdateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no IDs")
	}
}

func TestIssueBatchUpdateNoStatus(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()
	batchUpdateIDs = "issue-1"
	batchUpdateStatus = ""
	err := issueBatchUpdateCmd.RunE(issueBatchUpdateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no status")
	}
}

func TestIssueBatchUpdateSuccess(t *testing.T) {
	responses := []map[string]any{
		// 1: GetIssue — получаем команду
		{
			"issue": map[string]any{
				"id":         "issue-uuid-1",
				"identifier": "ENG-1",
				"title":      "Test Issue",
				"updatedAt":  "2024-01-01T00:00:00Z",
				"priority":   0,
				"state":      map[string]any{"name": "Todo", "type": "unstarted"},
				"team":       map[string]any{"id": "team-1", "key": "ENG", "name": "Eng"},
			},
		},
		// 2: FindWorkflowStateByName
		{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{
					{"id": "state-in-progress", "name": "In Progress"},
				},
			},
		},
		// 3: BatchUpdateIssues
		{
			"issueBatchUpdate": map[string]any{
				"success": true,
				"issues": []map[string]any{
					{
						"id":         "issue-uuid-1",
						"identifier": "ENG-1",
						"title":      "Test Issue",
						"priority":   0,
						"state":      map[string]any{"name": "In Progress", "type": "started"},
						"team":       map[string]any{"id": "team-1", "key": "ENG", "name": "Eng"},
					},
				},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	batchUpdateIDs = "issue-uuid-1"
	batchUpdateStatus = "In Progress"

	cmd := issueBatchUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "1") {
		t.Errorf("expected count in output, got: %s", out)
	}
}

func TestIssueBatchUpdateStatusNotFound(t *testing.T) {
	responses := []map[string]any{
		// 1: GetIssue
		{
			"issue": map[string]any{
				"id":         "issue-uuid-1",
				"identifier": "ENG-1",
				"title":      "Test Issue",
				"updatedAt":  "2024-01-01T00:00:00Z",
				"priority":   0,
				"state":      map[string]any{"name": "Todo", "type": "unstarted"},
				"team":       map[string]any{"id": "team-1", "key": "ENG", "name": "Eng"},
			},
		},
		// 2: FindWorkflowStateByName — статус не найден
		{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	batchUpdateIDs = "issue-uuid-1"
	batchUpdateStatus = "NonExistentStatus"

	err := issueBatchUpdateCmd.RunE(issueBatchUpdateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when status not found")
	}
}

func TestIssueBatchUpdateJSONOutput(t *testing.T) {
	responses := []map[string]any{
		{
			"issue": map[string]any{
				"id":         "issue-uuid-1",
				"identifier": "ENG-1",
				"title":      "Test Issue",
				"updatedAt":  "2024-01-01T00:00:00Z",
				"priority":   0,
				"state":      map[string]any{"name": "Todo", "type": "unstarted"},
				"team":       map[string]any{"id": "team-1", "key": "ENG", "name": "Eng"},
			},
		},
		{
			"workflowStates": map[string]any{
				"nodes": []map[string]any{
					{"id": "state-in-progress", "name": "In Progress"},
				},
			},
		},
		{
			"issueBatchUpdate": map[string]any{
				"success": true,
				"issues": []map[string]any{
					{
						"id":         "issue-uuid-1",
						"identifier": "ENG-1",
						"title":      "Test Issue",
						"priority":   0,
						"state":      map[string]any{"name": "In Progress", "type": "started"},
						"team":       map[string]any{"id": "team-1", "key": "ENG", "name": "Eng"},
					},
				},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	batchUpdateIDs = "issue-uuid-1"
	batchUpdateStatus = "In Progress"
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	cmd := issueBatchUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected ENG-1 in JSON output, got: %s", out)
	}
}
