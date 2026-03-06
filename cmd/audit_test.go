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

func newAuditTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestAuditListOutput(t *testing.T) {
	responseData := map[string]any{
		"auditEntries": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "ae1",
					"type":      "issueCreate",
					"actorId":   "u1",
					"createdAt": "2024-01-01T00:00:00Z",
					"ip":        "1.2.3.4",
					"country":   "US",
				},
			},
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   "",
			},
		},
	}

	srv := newAuditTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := auditListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "issueCreate") {
		t.Errorf("expected 'issueCreate' in output, got: %s", output)
	}
	if !strings.Contains(output, "1.2.3.4") {
		t.Errorf("expected IP in output, got: %s", output)
	}
}

func TestAuditListJSON(t *testing.T) {
	responseData := map[string]any{
		"auditEntries": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "ae2",
					"type":      "teamCreate",
					"actorId":   "u2",
					"createdAt": "2024-01-02T00:00:00Z",
					"ip":        "5.6.7.8",
					"country":   "DE",
				},
			},
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   "",
			},
		},
	}

	srv := newAuditTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := auditListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "teamCreate") {
		t.Errorf("expected 'teamCreate' in JSON output, got: %s", output)
	}
}

func TestAuditListEmpty(t *testing.T) {
	responseData := map[string]any{
		"auditEntries": map[string]any{
			"nodes": []map[string]any{},
			"pageInfo": map[string]any{
				"hasNextPage": false,
				"endCursor":   "",
			},
		},
	}

	srv := newAuditTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := auditListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestAuditTypesOutput(t *testing.T) {
	responseData := map[string]any{
		"auditEntryTypes": []map[string]any{
			{
				"type":        "issueCreate",
				"description": "Issue created",
			},
			{
				"type":        "issueDelete",
				"description": "Issue deleted",
			},
		},
	}

	srv := newAuditTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := auditTypesCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "issueCreate") {
		t.Errorf("expected 'issueCreate' in output, got: %s", output)
	}
	if !strings.Contains(output, "Issue created") {
		t.Errorf("expected description in output, got: %s", output)
	}
}

func TestAuditTypesJSON(t *testing.T) {
	responseData := map[string]any{
		"auditEntryTypes": []map[string]any{
			{
				"type":        "userCreate",
				"description": "User created",
			},
		},
	}

	srv := newAuditTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := auditTypesCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "userCreate") {
		t.Errorf("expected 'userCreate' in JSON output, got: %s", output)
	}
}
