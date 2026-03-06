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

func newInitiativesTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestInitiativesListOutput(t *testing.T) {
	responseData := map[string]any{
		"initiatives": map[string]any{
			"nodes": []map[string]any{
				{
					"id":     "init1",
					"name":   "Platform Initiative",
					"status": "planned",
					"owner":  map[string]any{"id": "user1", "name": "Alice"},
				},
			},
		},
	}

	srv := newInitiativesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := initiativesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"init1", "Platform Initiative", "planned", "Alice"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestInitiativesListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"initiatives": map[string]any{
			"nodes": []map[string]any{
				{"id": "init1", "name": "Platform Initiative"},
			},
		},
	}

	srv := newInitiativesTestServer(t, responseData)
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

	cmd := initiativesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"init1"`) {
		t.Errorf("expected init1 in JSON output, got:\n%s", out)
	}
}

func TestInitiativesListEmpty(t *testing.T) {
	responseData := map[string]any{
		"initiatives": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newInitiativesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := initiativesListCmd
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

func TestInitiativesView(t *testing.T) {
	responseData := map[string]any{
		"initiative": map[string]any{
			"id":          "init1",
			"name":        "Platform Initiative",
			"status":      "planned",
			"description": "Scale the platform",
		},
	}

	srv := newInitiativesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := initiativesViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"init1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"init1", "Platform Initiative", "planned", "Scale the platform"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestInitiativesCreate(t *testing.T) {
	responseData := map[string]any{
		"initiativeCreate": map[string]any{
			"success": true,
			"initiative": map[string]any{
				"id":     "init1",
				"name":   "New Initiative",
				"status": "planned",
			},
		},
	}

	srv := newInitiativesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := initiativesCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "New Initiative")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "init1") {
		t.Errorf("expected init1 in output, got:\n%s", out)
	}
}

func TestInitiativesUpdate(t *testing.T) {
	responseData := map[string]any{
		"initiativeUpdate": map[string]any{
			"success": true,
			"initiative": map[string]any{
				"id":   "init1",
				"name": "Updated Initiative",
			},
		},
	}

	srv := newInitiativesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := initiativesUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "Updated Initiative")

	if err := cmd.RunE(cmd, []string{"init1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "init1") {
		t.Errorf("expected init1 in output, got:\n%s", out)
	}
}

func TestInitiativesArchive(t *testing.T) {
	responseData := map[string]any{
		"initiativeArchive": map[string]any{
			"success": true,
		},
	}

	srv := newInitiativesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := initiativesArchiveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"init1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "init1") {
		t.Errorf("expected init1 in output, got:\n%s", out)
	}
}
