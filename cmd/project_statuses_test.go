package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectStatusesListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListProjectStatuses": map[string]any{
			"projectStatuses": map[string]any{
				"nodes": []map[string]any{
					{"id": "ps1", "name": "Planned", "type": "planned", "color": "#0000ff", "description": "Planning phase", "position": 1.0},
					{"id": "ps2", "name": "In Progress", "type": "started", "color": "#00ff00", "description": "", "position": 2.0},
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

	var buf bytes.Buffer
	projectStatusesListCmd.SetOut(&buf)

	if err := projectStatusesListCmd.RunE(projectStatusesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Planned") {
		t.Errorf("expected 'Planned' in output, got: %s", out)
	}
	if !strings.Contains(out, "planned") {
		t.Errorf("expected type 'planned' in output, got: %s", out)
	}
	if !strings.Contains(out, "#0000ff") {
		t.Errorf("expected color '#0000ff' in output, got: %s", out)
	}
}

func TestProjectStatusesListJSON(t *testing.T) {
	responses := map[string]any{
		"ListProjectStatuses": map[string]any{
			"projectStatuses": map[string]any{
				"nodes": []map[string]any{
					{"id": "ps1", "name": "Planned", "type": "planned", "color": "#0000ff", "description": "", "position": 1.0},
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

	origFormat := outputFormat
	outputFormat = "json"
	defer func() { outputFormat = origFormat }()

	var buf bytes.Buffer
	projectStatusesListCmd.SetOut(&buf)

	if err := projectStatusesListCmd.RunE(projectStatusesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestProjectStatusesListNoToken(t *testing.T) {
	token = ""
	err := projectStatusesListCmd.RunE(projectStatusesListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectStatusesListEmpty(t *testing.T) {
	responses := map[string]any{
		"ListProjectStatuses": map[string]any{
			"projectStatuses": map[string]any{
				"nodes": []map[string]any{},
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
	projectStatusesListCmd.SetOut(&buf)

	if err := projectStatusesListCmd.RunE(projectStatusesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected 'не найдены' in output, got: %s", out)
	}
}

func TestProjectStatusesCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectStatusCreate": map[string]any{
			"projectStatusCreate": map[string]any{
				"success": true,
				"projectStatus": map[string]any{
					"id": "ps-new", "name": "Backlog", "type": "backlog", "color": "#aaaaaa", "description": "", "position": 0.0,
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

	_ = projectStatusesCreateCmd.Flags().Set("name", "Backlog")
	_ = projectStatusesCreateCmd.Flags().Set("type", "backlog")
	_ = projectStatusesCreateCmd.Flags().Set("color", "#aaaaaa")
	defer func() {
		_ = projectStatusesCreateCmd.Flags().Set("name", "")
		_ = projectStatusesCreateCmd.Flags().Set("type", "")
		_ = projectStatusesCreateCmd.Flags().Set("color", "")
	}()

	var buf bytes.Buffer
	projectStatusesCreateCmd.SetOut(&buf)

	if err := projectStatusesCreateCmd.RunE(projectStatusesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Backlog") {
		t.Errorf("expected 'Backlog' in output, got: %s", out)
	}
}

func TestProjectStatusesCreateNoToken(t *testing.T) {
	token = ""
	err := projectStatusesCreateCmd.RunE(projectStatusesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectStatusesCreateMissingName(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectStatusesCreateCmd.Flags().Set("name", "")
	_ = projectStatusesCreateCmd.Flags().Set("type", "planned")
	_ = projectStatusesCreateCmd.Flags().Set("color", "#fff")

	err := projectStatusesCreateCmd.RunE(projectStatusesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when name missing")
	}
}

func TestProjectStatusesCreateMissingType(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectStatusesCreateCmd.Flags().Set("name", "Test")
	_ = projectStatusesCreateCmd.Flags().Set("type", "")
	_ = projectStatusesCreateCmd.Flags().Set("color", "#fff")
	defer func() {
		_ = projectStatusesCreateCmd.Flags().Set("name", "")
		_ = projectStatusesCreateCmd.Flags().Set("color", "")
	}()

	err := projectStatusesCreateCmd.RunE(projectStatusesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when type missing")
	}
}

func TestProjectStatusesUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectStatusUpdate": map[string]any{
			"projectStatusUpdate": map[string]any{
				"success": true,
				"projectStatus": map[string]any{
					"id": "ps1", "name": "Updated Status", "type": "planned", "color": "#ff00ff", "description": "", "position": 1.0,
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

	_ = projectStatusesUpdateCmd.Flags().Set("name", "Updated Status")
	defer func() {
		_ = projectStatusesUpdateCmd.Flags().Set("name", "")
	}()

	var buf bytes.Buffer
	projectStatusesUpdateCmd.SetOut(&buf)

	if err := projectStatusesUpdateCmd.RunE(projectStatusesUpdateCmd, []string{"ps1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Updated Status") {
		t.Errorf("expected 'Updated Status' in output, got: %s", out)
	}
}

func TestProjectStatusesUpdateNoToken(t *testing.T) {
	token = ""
	err := projectStatusesUpdateCmd.RunE(projectStatusesUpdateCmd, []string{"ps1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
