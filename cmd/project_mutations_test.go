package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectViewSuccess(t *testing.T) {
	responses := map[string]any{
		"GetProject": map[string]any{
			"project": map[string]any{
				"id":          "proj1",
				"name":        "My Project",
				"description": "A test project",
				"state":       "started",
				"startDate":   "2024-01-01",
				"targetDate":  "2024-06-01",
				"url":         "https://linear.app/proj",
				"lead":        map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
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
	projectViewCmd.SetOut(&buf)

	if err := projectViewCmd.RunE(projectViewCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "My Project") {
		t.Errorf("expected 'My Project' in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-01-01") {
		t.Errorf("expected startDate in output, got: %s", out)
	}
}

func TestProjectViewJSON(t *testing.T) {
	responses := map[string]any{
		"GetProject": map[string]any{
			"project": map[string]any{
				"id":          "proj1",
				"name":        "My Project",
				"description": "",
				"state":       "started",
				"startDate":   "",
				"targetDate":  "",
				"url":         "",
				"lead":        nil,
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
	projectViewCmd.SetOut(&buf)

	if err := projectViewCmd.RunE(projectViewCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestProjectViewNoToken(t *testing.T) {
	token = ""
	err := projectViewCmd.RunE(projectViewCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectCreate": map[string]any{
			"projectCreate": map[string]any{
				"success": true,
				"project": map[string]any{
					"id":          "proj-new",
					"name":        "New Project",
					"description": "",
					"state":       "planned",
					"startDate":   "",
					"targetDate":  "",
					"url":         "",
					"lead":        nil,
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

	_ = projectCreateCmd.Flags().Set("name", "New Project")
	_ = projectCreateCmd.Flags().Set("team-ids", "team1")
	defer func() {
		_ = projectCreateCmd.Flags().Set("name", "")
		_ = projectCreateCmd.Flags().Set("team-ids", "")
	}()

	var buf bytes.Buffer
	projectCreateCmd.SetOut(&buf)

	if err := projectCreateCmd.RunE(projectCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "New Project") {
		t.Errorf("expected 'New Project' in output, got: %s", out)
	}
}

func TestProjectCreateNoToken(t *testing.T) {
	token = ""
	err := projectCreateCmd.RunE(projectCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectCreateMissingName(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectCreateCmd.Flags().Set("name", "")
	_ = projectCreateCmd.Flags().Set("team-ids", "team1")

	err := projectCreateCmd.RunE(projectCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when name missing")
	}
}

func TestProjectUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectUpdate": map[string]any{
			"projectUpdate": map[string]any{
				"success": true,
				"project": map[string]any{
					"id":          "proj1",
					"name":        "Updated Project",
					"description": "",
					"state":       "started",
					"startDate":   "",
					"targetDate":  "",
					"url":         "",
					"lead":        nil,
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

	_ = projectUpdateCmd.Flags().Set("name", "Updated Project")
	defer func() {
		_ = projectUpdateCmd.Flags().Set("name", "")
	}()

	var buf bytes.Buffer
	projectUpdateCmd.SetOut(&buf)

	if err := projectUpdateCmd.RunE(projectUpdateCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Updated Project") {
		t.Errorf("expected 'Updated Project' in output, got: %s", out)
	}
}

func TestProjectUpdateNoToken(t *testing.T) {
	token = ""
	err := projectUpdateCmd.RunE(projectUpdateCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectDeleteSuccess(t *testing.T) {
	responses := map[string]any{
		"projectDelete": map[string]any{
			"projectDelete": map[string]any{
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
	projectDeleteCmd.SetOut(&buf)

	if err := projectDeleteCmd.RunE(projectDeleteCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "proj1") {
		t.Errorf("expected project ID in output, got: %s", out)
	}
}

func TestProjectDeleteNoToken(t *testing.T) {
	token = ""
	err := projectDeleteCmd.RunE(projectDeleteCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectArchiveSuccess(t *testing.T) {
	responses := map[string]any{
		"projectArchive": map[string]any{
			"projectArchive": map[string]any{
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
	projectArchiveCmd.SetOut(&buf)

	if err := projectArchiveCmd.RunE(projectArchiveCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "архивирован") {
		t.Errorf("expected 'архивирован' in output, got: %s", out)
	}
}

func TestProjectArchiveNoToken(t *testing.T) {
	token = ""
	err := projectArchiveCmd.RunE(projectArchiveCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectUnarchiveSuccess(t *testing.T) {
	responses := map[string]any{
		"projectUnarchive": map[string]any{
			"projectUnarchive": map[string]any{
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
	projectUnarchiveCmd.SetOut(&buf)

	if err := projectUnarchiveCmd.RunE(projectUnarchiveCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "разархивирован") {
		t.Errorf("expected 'разархивирован' in output, got: %s", out)
	}
}

func TestProjectUnarchiveNoToken(t *testing.T) {
	token = ""
	err := projectUnarchiveCmd.RunE(projectUnarchiveCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
