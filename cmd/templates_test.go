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

func newTemplatesTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestTemplatesListOutput(t *testing.T) {
	responseData := map[string]any{
		"templates": []map[string]any{
			{
				"id":   "tpl1",
				"name": "Bug Report",
				"type": "issue",
				"team": map[string]any{"id": "team1", "name": "Engineering"},
			},
		},
	}

	srv := newTemplatesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := templatesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"tpl1", "Bug Report", "issue", "Engineering"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestTemplatesListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"templates": []map[string]any{
			{"id": "tpl1", "name": "Bug Report", "type": "issue"},
		},
	}

	srv := newTemplatesTestServer(t, responseData)
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

	cmd := templatesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"tpl1"`) {
		t.Errorf("expected tpl1 in JSON output, got:\n%s", out)
	}
}

func TestTemplatesListEmpty(t *testing.T) {
	responseData := map[string]any{
		"templates": []map[string]any{},
	}

	srv := newTemplatesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := templatesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}

func TestTemplatesView(t *testing.T) {
	responseData := map[string]any{
		"template": map[string]any{
			"id":          "tpl1",
			"name":        "Bug Report",
			"type":        "issue",
			"description": "Template for bug reports",
			"creator":     map[string]any{"id": "user1", "name": "Alice"},
		},
	}

	srv := newTemplatesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := templatesViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"tpl1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"tpl1", "Bug Report", "issue", "Template for bug reports", "Alice"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestTemplatesCreate(t *testing.T) {
	responseData := map[string]any{
		"templateCreate": map[string]any{
			"success": true,
			"template": map[string]any{
				"id":   "tpl1",
				"name": "New Template",
				"type": "issue",
			},
		},
	}

	srv := newTemplatesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := templatesCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "New Template")
	cmd.Flags().Set("type", "issue")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "tpl1") {
		t.Errorf("expected tpl1 in output, got:\n%s", out)
	}
}

func TestTemplatesUpdate(t *testing.T) {
	responseData := map[string]any{
		"templateUpdate": map[string]any{
			"success": true,
			"template": map[string]any{
				"id":   "tpl1",
				"name": "Updated Template",
				"type": "issue",
			},
		},
	}

	srv := newTemplatesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := templatesUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "Updated Template")

	if err := cmd.RunE(cmd, []string{"tpl1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "tpl1") {
		t.Errorf("expected tpl1 in output, got:\n%s", out)
	}
}

func TestTemplatesDelete(t *testing.T) {
	responseData := map[string]any{
		"templateDelete": map[string]any{
			"success": true,
		},
	}

	srv := newTemplatesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := templatesDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"tpl1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "tpl1") {
		t.Errorf("expected tpl1 in output, got:\n%s", out)
	}
}
