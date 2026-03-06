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

func newReleasesTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestReleasesListOutput(t *testing.T) {
	responseData := map[string]any{
		"releases": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "rel1",
					"name":      "v1.0.0",
					"createdAt": "2024-01-01T00:00:00Z",
					"pipeline":  map[string]any{"id": "pipe1", "name": "Main"},
					"stage":     map[string]any{"id": "stage1", "name": "Production", "color": "#00ff00"},
				},
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("expected 'v1.0.0' in output, got: %s", output)
	}
	if !strings.Contains(output, "Main") {
		t.Errorf("expected 'Main' in output, got: %s", output)
	}
}

func TestReleasesListJSON(t *testing.T) {
	responseData := map[string]any{
		"releases": map[string]any{
			"nodes": []map[string]any{
				{
					"id":   "rel1",
					"name": "v1.0.0",
				},
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, `"id"`) {
		t.Errorf("expected JSON output with 'id', got: %s", output)
	}
}

func TestReleasesListEmpty(t *testing.T) {
	responseData := map[string]any{
		"releases": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesListCmd
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

func TestReleasesViewOutput(t *testing.T) {
	responseData := map[string]any{
		"release": map[string]any{
			"id":          "rel1",
			"name":        "v1.0.0",
			"description": "First release",
			"createdAt":   "2024-01-01T00:00:00Z",
			"commitSha":   "abc123",
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"rel1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("expected 'v1.0.0' in output, got: %s", output)
	}
	if !strings.Contains(output, "abc123") {
		t.Errorf("expected 'abc123' in output, got: %s", output)
	}
}

func TestReleasesCreateOutput(t *testing.T) {
	responseData := map[string]any{
		"releaseCreate": map[string]any{
			"success": true,
			"release": map[string]any{
				"id":        "rel1",
				"name":      "v1.0.0",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("name", "v1.0.0")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "v1.0.0") {
		t.Errorf("expected 'v1.0.0' in output, got: %s", output)
	}
}

func TestReleasesCompleteOutput(t *testing.T) {
	responseData := map[string]any{
		"releaseComplete": map[string]any{
			"success": true,
			"release": map[string]any{
				"id":          "rel1",
				"name":        "v1.0.0",
				"completedAt": "2024-06-01T00:00:00Z",
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesCompleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"rel1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "завершён") {
		t.Errorf("expected completion message, got: %s", output)
	}
}

func TestReleasesPipelinesListOutput(t *testing.T) {
	responseData := map[string]any{
		"releasePipelines": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "pipe1",
					"name":      "Main Pipeline",
					"slugId":    "main-pipeline",
					"type":      "github",
					"createdAt": "2024-01-01T00:00:00Z",
				},
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesPipelinesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Main Pipeline") {
		t.Errorf("expected 'Main Pipeline' in output, got: %s", output)
	}
}

func TestReleasesPipelinesListJSON(t *testing.T) {
	responseData := map[string]any{
		"releasePipelines": map[string]any{
			"nodes": []map[string]any{
				{
					"id":   "pipe1",
					"name": "Main Pipeline",
				},
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesPipelinesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, `"id"`) {
		t.Errorf("expected JSON output with 'id', got: %s", output)
	}
}

func TestReleasesPipelinesCreateOutput(t *testing.T) {
	responseData := map[string]any{
		"releasePipelineCreate": map[string]any{
			"success": true,
			"releasePipeline": map[string]any{
				"id":        "pipe1",
				"name":      "Main Pipeline",
				"slugId":    "main-pipeline",
				"type":      "github",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := newReleasesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := releasesPipelinesCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("name", "Main Pipeline")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Main Pipeline") {
		t.Errorf("expected 'Main Pipeline' in output, got: %s", output)
	}
}
