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

func newRoadmapsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestRoadmapsListOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmaps": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "rm1",
					"name":      "Q1 Roadmap",
					"createdAt": "2024-01-01T00:00:00Z",
					"owner":     map[string]any{"id": "u1", "name": "Alice"},
				},
			},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Q1 Roadmap") {
		t.Errorf("expected 'Q1 Roadmap' in output, got: %s", output)
	}
	if !strings.Contains(output, "Alice") {
		t.Errorf("expected 'Alice' in output, got: %s", output)
	}
}

func TestRoadmapsListOutputJSON(t *testing.T) {
	responseData := map[string]any{
		"roadmaps": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "rm1",
					"name":      "Q1 Roadmap",
					"createdAt": "2024-01-01T00:00:00Z",
				},
			},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Q1 Roadmap") {
		t.Errorf("expected JSON output with Q1 Roadmap")
	}
}

func TestRoadmapsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"roadmaps": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected 'не найдены' in output, got: %s", buf.String())
	}
}

func TestRoadmapsViewOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmap": map[string]any{
			"id":          "rm1",
			"name":        "Q1 Roadmap",
			"description": "First quarter goals",
			"createdAt":   "2024-01-01T00:00:00Z",
			"owner":       map[string]any{"id": "u1", "name": "Alice"},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"rm1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Q1 Roadmap") {
		t.Errorf("expected 'Q1 Roadmap' in output, got: %s", output)
	}
	if !strings.Contains(output, "Alice") {
		t.Errorf("expected 'Alice' in output, got: %s", output)
	}
}

func TestRoadmapsCreateOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmapCreate": map[string]any{
			"success": true,
			"roadmap": map[string]any{
				"id":        "rm1",
				"name":      "New Roadmap",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("name", "New Roadmap")
	defer cmd.Flags().Set("name", "")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "New Roadmap") {
		t.Errorf("expected 'New Roadmap' in output, got: %s", buf.String())
	}
}

func TestRoadmapsUpdateOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmapUpdate": map[string]any{
			"success": true,
			"roadmap": map[string]any{
				"id":        "rm1",
				"name":      "Updated Roadmap",
				"createdAt": "2024-01-01T00:00:00Z",
			},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("name", "Updated Roadmap")
	defer cmd.Flags().Set("name", "")

	if err := cmd.RunE(cmd, []string{"rm1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Updated Roadmap") {
		t.Errorf("expected 'Updated Roadmap' in output, got: %s", buf.String())
	}
}

func TestRoadmapsDeleteOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmapDelete": map[string]any{
			"success": true,
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"rm1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "rm1") {
		t.Errorf("expected 'rm1' in output, got: %s", buf.String())
	}
}

func TestRoadmapsAddProjectOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmapToProjectCreate": map[string]any{
			"success": true,
			"roadmapToProject": map[string]any{
				"id":      "rel1",
				"roadmap": map[string]any{"id": "rm1", "name": "Q1 Roadmap"},
				"project": map[string]any{"id": "p1", "name": "My Project"},
			},
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsAddProjectCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("roadmap", "rm1")
	cmd.Flags().Set("project", "p1")
	defer func() {
		cmd.Flags().Set("roadmap", "")
		cmd.Flags().Set("project", "")
	}()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Q1 Roadmap") {
		t.Errorf("expected 'Q1 Roadmap' in output, got: %s", output)
	}
}

func TestRoadmapsRemoveProjectOutput(t *testing.T) {
	responseData := map[string]any{
		"roadmapToProjectDelete": map[string]any{
			"success": true,
		},
	}

	srv := newRoadmapsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := roadmapsRemoveProjectCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"rel1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "rel1") {
		t.Errorf("expected 'rel1' in output, got: %s", buf.String())
	}
}
