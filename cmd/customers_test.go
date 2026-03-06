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

func newCustomersTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestCustomersListOutput(t *testing.T) {
	responseData := map[string]any{
		"customers": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "c1",
					"name":      "Acme Corp",
					"createdAt": "2024-01-01T00:00:00Z",
					"status":    map[string]any{"id": "s1", "name": "active", "displayName": "Active", "color": "#00ff00"},
					"tier":      map[string]any{"id": "t1", "name": "enterprise", "displayName": "Enterprise", "color": "#gold"},
				},
			},
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Acme Corp") {
		t.Errorf("expected 'Acme Corp' in output, got: %s", output)
	}
	if !strings.Contains(output, "Active") {
		t.Errorf("expected 'Active' in output, got: %s", output)
	}
}

func TestCustomersListOutputJSON(t *testing.T) {
	responseData := map[string]any{
		"customers": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "c1",
					"name":      "Acme Corp",
					"createdAt": "2024-01-01T00:00:00Z",
				},
			},
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, `"Acme Corp"`) {
		t.Errorf("expected JSON with Acme Corp, got: %s", output)
	}
}

func TestCustomersListEmpty(t *testing.T) {
	responseData := map[string]any{
		"customers": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "не найдены") {
		t.Errorf("expected empty message, got: %s", output)
	}
}

func TestCustomersViewOutput(t *testing.T) {
	responseData := map[string]any{
		"customer": map[string]any{
			"id":        "c1",
			"name":      "Acme Corp",
			"createdAt": "2024-01-01T00:00:00Z",
			"status":    map[string]any{"id": "s1", "name": "active", "displayName": "Active", "color": "#00ff00"},
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"c1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Acme Corp") {
		t.Errorf("expected 'Acme Corp' in output, got: %s", output)
	}
}

func TestCustomersCreateOutput(t *testing.T) {
	responseData := map[string]any{
		"customerCreate": map[string]any{
			"success": true,
			"customer": map[string]any{
				"id":        "c1",
				"name":      "New Corp",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.Flags().Set("name", "New Corp"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "New Corp") {
		t.Errorf("expected 'New Corp' in output, got: %s", output)
	}
}

func TestCustomersUpdateOutput(t *testing.T) {
	responseData := map[string]any{
		"customerUpdate": map[string]any{
			"success": true,
			"customer": map[string]any{
				"id":        "c1",
				"name":      "Updated Corp",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.Flags().Set("name", "Updated Corp"); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	if err := cmd.RunE(cmd, []string{"c1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Updated Corp") {
		t.Errorf("expected 'Updated Corp' in output, got: %s", output)
	}
}

func TestCustomersDeleteOutput(t *testing.T) {
	responseData := map[string]any{
		"customerDelete": map[string]any{
			"success": true,
		},
	}

	srv := newCustomersTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := customersDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"c1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "удалён") {
		t.Errorf("expected deletion message, got: %s", output)
	}
}
