package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectLabelsListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListProjectLabels": map[string]any{
			"projectLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "pl1", "name": "Frontend", "color": "#ff0000", "description": "Frontend work"},
					{"id": "pl2", "name": "Backend", "color": "#0000ff", "description": ""},
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
	projectLabelsListCmd.SetOut(&buf)

	if err := projectLabelsListCmd.RunE(projectLabelsListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Frontend") {
		t.Errorf("expected 'Frontend' in output, got: %s", out)
	}
	if !strings.Contains(out, "#ff0000") {
		t.Errorf("expected color '#ff0000' in output, got: %s", out)
	}
}

func TestProjectLabelsListJSON(t *testing.T) {
	responses := map[string]any{
		"ListProjectLabels": map[string]any{
			"projectLabels": map[string]any{
				"nodes": []map[string]any{
					{"id": "pl1", "name": "Frontend", "color": "#ff0000", "description": ""},
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
	projectLabelsListCmd.SetOut(&buf)

	if err := projectLabelsListCmd.RunE(projectLabelsListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestProjectLabelsListNoToken(t *testing.T) {
	token = ""
	err := projectLabelsListCmd.RunE(projectLabelsListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectLabelsListEmpty(t *testing.T) {
	responses := map[string]any{
		"ListProjectLabels": map[string]any{
			"projectLabels": map[string]any{
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
	projectLabelsListCmd.SetOut(&buf)

	if err := projectLabelsListCmd.RunE(projectLabelsListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected 'не найдены' in output, got: %s", out)
	}
}

func TestProjectLabelsCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectLabelCreate": map[string]any{
			"projectLabelCreate": map[string]any{
				"success": true,
				"projectLabel": map[string]any{
					"id": "pl-new", "name": "Design", "color": "#00ff00", "description": "",
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

	_ = projectLabelsCreateCmd.Flags().Set("name", "Design")
	_ = projectLabelsCreateCmd.Flags().Set("color", "#00ff00")
	defer func() {
		_ = projectLabelsCreateCmd.Flags().Set("name", "")
		_ = projectLabelsCreateCmd.Flags().Set("color", "")
	}()

	var buf bytes.Buffer
	projectLabelsCreateCmd.SetOut(&buf)

	if err := projectLabelsCreateCmd.RunE(projectLabelsCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Design") {
		t.Errorf("expected 'Design' in output, got: %s", out)
	}
}

func TestProjectLabelsCreateNoToken(t *testing.T) {
	token = ""
	err := projectLabelsCreateCmd.RunE(projectLabelsCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectLabelsCreateMissingName(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectLabelsCreateCmd.Flags().Set("name", "")

	err := projectLabelsCreateCmd.RunE(projectLabelsCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when name missing")
	}
}

func TestProjectLabelsUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectLabelUpdate": map[string]any{
			"projectLabelUpdate": map[string]any{
				"success": true,
				"projectLabel": map[string]any{
					"id": "pl1", "name": "Updated Label", "color": "#ff00ff", "description": "",
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

	_ = projectLabelsUpdateCmd.Flags().Set("name", "Updated Label")
	defer func() {
		_ = projectLabelsUpdateCmd.Flags().Set("name", "")
	}()

	var buf bytes.Buffer
	projectLabelsUpdateCmd.SetOut(&buf)

	if err := projectLabelsUpdateCmd.RunE(projectLabelsUpdateCmd, []string{"pl1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Updated Label") {
		t.Errorf("expected 'Updated Label' in output, got: %s", out)
	}
}

func TestProjectLabelsUpdateNoToken(t *testing.T) {
	token = ""
	err := projectLabelsUpdateCmd.RunE(projectLabelsUpdateCmd, []string{"pl1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectLabelsDeleteSuccess(t *testing.T) {
	responses := map[string]any{
		"projectLabelDelete": map[string]any{
			"projectLabelDelete": map[string]any{
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
	projectLabelsDeleteCmd.SetOut(&buf)

	if err := projectLabelsDeleteCmd.RunE(projectLabelsDeleteCmd, []string{"pl1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "pl1") {
		t.Errorf("expected label ID in output, got: %s", out)
	}
}

func TestProjectLabelsDeleteNoToken(t *testing.T) {
	token = ""
	err := projectLabelsDeleteCmd.RunE(projectLabelsDeleteCmd, []string{"pl1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
