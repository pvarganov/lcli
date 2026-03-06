package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

// newLabelAssignTestServer создаёт тестовый сервер с очередью ответов.
func newLabelAssignTestServer(t *testing.T, responses []map[string]any) *httptest.Server {
	t.Helper()
	requestNum := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idx := requestNum
		requestNum++
		if idx >= len(responses) {
			idx = len(responses) - 1
		}
		json.NewEncoder(w).Encode(map[string]any{"data": responses[idx]})
	}))
}

func TestLabelAddSuccess(t *testing.T) {
	responses := []map[string]any{
		// 1: ListIssueLabels
		{
			"issueLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "label-bug", "name": "Bug", "color": "#ff0000", "description": ""},
				},
			},
		},
		// 2: getIssueLabelIDs
		{
			"issue": map[string]any{
				"labels": map[string]any{
					"nodes": []map[string]any{},
				},
			},
		},
		// 3: issueUpdate
		{
			"issueUpdate": map[string]any{"success": true},
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

	cmd := labelAddCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1", "Bug"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Bug") {
		t.Errorf("expected 'Bug' in output, got: %s", out)
	}
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected 'ENG-1' in output, got: %s", out)
	}
}

func TestLabelAddLabelNotFound(t *testing.T) {
	responses := []map[string]any{
		{
			"issueLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "label-bug", "name": "Bug", "color": "#ff0000", "description": ""},
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

	err := labelAddCmd.RunE(labelAddCmd, []string{"ENG-1", "NonExistentLabel"})
	if err == nil {
		t.Fatal("expected error for non-existent label")
	}
	if !strings.Contains(err.Error(), "не найдена") {
		t.Errorf("expected 'не найдена' in error, got: %s", err.Error())
	}
}

func TestLabelAddNoToken(t *testing.T) {
	token = ""
	err := labelAddCmd.RunE(labelAddCmd, []string{"ENG-1", "Bug"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestLabelRemoveSuccess(t *testing.T) {
	responses := []map[string]any{
		// 1: ListIssueLabels
		{
			"issueLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "label-bug", "name": "Bug", "color": "#ff0000", "description": ""},
				},
			},
		},
		// 2: getIssueLabelIDs
		{
			"issue": map[string]any{
				"labels": map[string]any{
					"nodes": []map[string]any{
						{"id": "label-bug"},
					},
				},
			},
		},
		// 3: issueUpdate
		{
			"issueUpdate": map[string]any{"success": true},
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

	cmd := labelRemoveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1", "Bug"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Bug") {
		t.Errorf("expected 'Bug' in output, got: %s", out)
	}
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected 'ENG-1' in output, got: %s", out)
	}
}

func TestLabelRemoveLabelNotFound(t *testing.T) {
	responses := []map[string]any{
		{
			"issueLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "label-bug", "name": "Bug", "color": "#ff0000", "description": ""},
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

	err := labelRemoveCmd.RunE(labelRemoveCmd, []string{"ENG-1", "NonExistentLabel"})
	if err == nil {
		t.Fatal("expected error for non-existent label")
	}
	if !strings.Contains(err.Error(), "не найдена") {
		t.Errorf("expected 'не найдена' in error, got: %s", err.Error())
	}
}

func TestLabelRemoveNoToken(t *testing.T) {
	token = ""
	err := labelRemoveCmd.RunE(labelRemoveCmd, []string{"ENG-1", "Bug"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
