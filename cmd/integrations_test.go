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

func newIntegrationsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestIntegrationsListOutput(t *testing.T) {
	responseData := map[string]any{
		"integrations": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "int1",
					"service":   "github",
					"createdAt": "2024-01-01T00:00:00Z",
					"team":      map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
	}

	srv := newIntegrationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := integrationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "github") {
		t.Errorf("expected 'github' in output, got: %s", output)
	}
	if !strings.Contains(output, "ENG") {
		t.Errorf("expected 'ENG' in output, got: %s", output)
	}
}

func TestIntegrationsListJSON(t *testing.T) {
	responseData := map[string]any{
		"integrations": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "int1",
					"service":   "slack",
					"createdAt": "2024-01-01T00:00:00Z",
				},
			},
		},
	}

	srv := newIntegrationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := integrationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "slack") {
		t.Errorf("expected 'slack' in json output, got: %s", output)
	}
}

func TestIntegrationsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"integrations": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newIntegrationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := integrationsListCmd
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

func TestIntegrationsDeleteCmd(t *testing.T) {
	responseData := map[string]any{
		"integrationDelete": map[string]any{
			"success": true,
		},
	}

	srv := newIntegrationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := integrationsDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"int1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "int1") {
		t.Errorf("expected ID in output, got: %s", buf.String())
	}
}
