package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectUpdatesListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListProjectUpdates": map[string]any{
			"project": map[string]any{
				"projectUpdates": map[string]any{
					"nodes": []map[string]any{
						{
							"id":        "pu1",
							"body":      "Week 1 progress report",
							"health":    "onTrack",
							"createdAt": "2024-06-01T10:00:00Z",
							"user":      map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
						},
					},
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
	projectUpdatesListCmd.SetOut(&buf)

	if err := projectUpdatesListCmd.RunE(projectUpdatesListCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Week 1 progress") {
		t.Errorf("expected body in output, got: %s", out)
	}
	if !strings.Contains(out, "onTrack") {
		t.Errorf("expected health in output, got: %s", out)
	}
}

func TestProjectUpdatesListJSON(t *testing.T) {
	responses := map[string]any{
		"ListProjectUpdates": map[string]any{
			"project": map[string]any{
				"projectUpdates": map[string]any{
					"nodes": []map[string]any{
						{"id": "pu1", "body": "Update", "health": "onTrack", "createdAt": "2024-06-01T10:00:00Z", "user": nil},
					},
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
	projectUpdatesListCmd.SetOut(&buf)

	if err := projectUpdatesListCmd.RunE(projectUpdatesListCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestProjectUpdatesListNoToken(t *testing.T) {
	token = ""
	err := projectUpdatesListCmd.RunE(projectUpdatesListCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectUpdatesCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectUpdateCreate": map[string]any{
			"projectUpdateCreate": map[string]any{
				"success": true,
				"projectUpdate": map[string]any{
					"id":        "pu-new",
					"body":      "New update",
					"health":    "onTrack",
					"createdAt": "2024-06-01T10:00:00Z",
					"user":      nil,
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

	_ = projectUpdatesCreateCmd.Flags().Set("project-id", "proj1")
	_ = projectUpdatesCreateCmd.Flags().Set("body", "New update")
	defer func() {
		_ = projectUpdatesCreateCmd.Flags().Set("project-id", "")
		_ = projectUpdatesCreateCmd.Flags().Set("body", "")
	}()

	var buf bytes.Buffer
	projectUpdatesCreateCmd.SetOut(&buf)

	if err := projectUpdatesCreateCmd.RunE(projectUpdatesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "pu-new") {
		t.Errorf("expected update ID in output, got: %s", out)
	}
}

func TestProjectUpdatesCreateNoToken(t *testing.T) {
	token = ""
	err := projectUpdatesCreateCmd.RunE(projectUpdatesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectUpdatesCreateMissingProjectID(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectUpdatesCreateCmd.Flags().Set("project-id", "")
	_ = projectUpdatesCreateCmd.Flags().Set("body", "Some body")

	err := projectUpdatesCreateCmd.RunE(projectUpdatesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when project-id missing")
	}
}

func TestProjectUpdatesCreateMissingBody(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectUpdatesCreateCmd.Flags().Set("project-id", "proj1")
	_ = projectUpdatesCreateCmd.Flags().Set("body", "")

	err := projectUpdatesCreateCmd.RunE(projectUpdatesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when body missing")
	}
}

func TestProjectUpdatesUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectUpdateUpdate": map[string]any{
			"projectUpdateUpdate": map[string]any{
				"success": true,
				"projectUpdate": map[string]any{
					"id":        "pu1",
					"body":      "Updated body",
					"health":    "atRisk",
					"createdAt": "2024-06-01T10:00:00Z",
					"user":      nil,
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

	_ = projectUpdatesUpdateCmd.Flags().Set("body", "Updated body")
	defer func() {
		_ = projectUpdatesUpdateCmd.Flags().Set("body", "")
	}()

	var buf bytes.Buffer
	projectUpdatesUpdateCmd.SetOut(&buf)

	if err := projectUpdatesUpdateCmd.RunE(projectUpdatesUpdateCmd, []string{"pu1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "pu1") {
		t.Errorf("expected update ID in output, got: %s", out)
	}
}

func TestProjectUpdatesUpdateNoToken(t *testing.T) {
	token = ""
	err := projectUpdatesUpdateCmd.RunE(projectUpdatesUpdateCmd, []string{"pu1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectUpdatesDeleteSuccess(t *testing.T) {
	responses := map[string]any{
		"projectUpdateArchive": map[string]any{
			"projectUpdateArchive": map[string]any{
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
	projectUpdatesDeleteCmd.SetOut(&buf)

	if err := projectUpdatesDeleteCmd.RunE(projectUpdatesDeleteCmd, []string{"pu1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "pu1") {
		t.Errorf("expected update ID in output, got: %s", out)
	}
}

func TestProjectUpdatesDeleteNoToken(t *testing.T) {
	token = ""
	err := projectUpdatesDeleteCmd.RunE(projectUpdatesDeleteCmd, []string{"pu1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
