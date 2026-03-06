package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectMilestonesListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListProjectMilestones": map[string]any{
			"project": map[string]any{
				"projectMilestones": map[string]any{
					"nodes": []map[string]any{
						{"id": "ms1", "name": "Alpha Release", "targetDate": "2024-03-01", "description": "First milestone"},
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
	projectMilestonesListCmd.SetOut(&buf)

	if err := projectMilestonesListCmd.RunE(projectMilestonesListCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alpha Release") {
		t.Errorf("expected 'Alpha Release' in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-03-01") {
		t.Errorf("expected targetDate in output, got: %s", out)
	}
}

func TestProjectMilestonesListJSON(t *testing.T) {
	responses := map[string]any{
		"ListProjectMilestones": map[string]any{
			"project": map[string]any{
				"projectMilestones": map[string]any{
					"nodes": []map[string]any{
						{"id": "ms1", "name": "Alpha", "targetDate": "", "description": ""},
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
	projectMilestonesListCmd.SetOut(&buf)

	if err := projectMilestonesListCmd.RunE(projectMilestonesListCmd, []string{"proj1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestProjectMilestonesListNoToken(t *testing.T) {
	token = ""
	err := projectMilestonesListCmd.RunE(projectMilestonesListCmd, []string{"proj1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectMilestonesCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectMilestoneCreate": map[string]any{
			"projectMilestoneCreate": map[string]any{
				"success": true,
				"projectMilestone": map[string]any{
					"id": "ms-new", "name": "New Milestone", "targetDate": "2024-06-01", "description": "",
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

	_ = projectMilestonesCreateCmd.Flags().Set("project-id", "proj1")
	_ = projectMilestonesCreateCmd.Flags().Set("name", "New Milestone")
	defer func() {
		_ = projectMilestonesCreateCmd.Flags().Set("project-id", "")
		_ = projectMilestonesCreateCmd.Flags().Set("name", "")
	}()

	var buf bytes.Buffer
	projectMilestonesCreateCmd.SetOut(&buf)

	if err := projectMilestonesCreateCmd.RunE(projectMilestonesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "New Milestone") {
		t.Errorf("expected 'New Milestone' in output, got: %s", out)
	}
}

func TestProjectMilestonesCreateNoToken(t *testing.T) {
	token = ""
	err := projectMilestonesCreateCmd.RunE(projectMilestonesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectMilestonesCreateMissingProjectID(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectMilestonesCreateCmd.Flags().Set("project-id", "")
	_ = projectMilestonesCreateCmd.Flags().Set("name", "Test")

	err := projectMilestonesCreateCmd.RunE(projectMilestonesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when project-id missing")
	}
}

func TestProjectMilestonesCreateMissingName(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = projectMilestonesCreateCmd.Flags().Set("project-id", "proj1")
	_ = projectMilestonesCreateCmd.Flags().Set("name", "")

	err := projectMilestonesCreateCmd.RunE(projectMilestonesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when name missing")
	}
}

func TestProjectMilestonesUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"projectMilestoneUpdate": map[string]any{
			"projectMilestoneUpdate": map[string]any{
				"success": true,
				"projectMilestone": map[string]any{
					"id": "ms1", "name": "Updated Milestone", "targetDate": "", "description": "",
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

	_ = projectMilestonesUpdateCmd.Flags().Set("name", "Updated Milestone")
	defer func() {
		_ = projectMilestonesUpdateCmd.Flags().Set("name", "")
	}()

	var buf bytes.Buffer
	projectMilestonesUpdateCmd.SetOut(&buf)

	if err := projectMilestonesUpdateCmd.RunE(projectMilestonesUpdateCmd, []string{"ms1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Updated Milestone") {
		t.Errorf("expected 'Updated Milestone' in output, got: %s", out)
	}
}

func TestProjectMilestonesUpdateNoToken(t *testing.T) {
	token = ""
	err := projectMilestonesUpdateCmd.RunE(projectMilestonesUpdateCmd, []string{"ms1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectMilestonesDeleteSuccess(t *testing.T) {
	responses := map[string]any{
		"projectMilestoneDelete": map[string]any{
			"projectMilestoneDelete": map[string]any{
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
	projectMilestonesDeleteCmd.SetOut(&buf)

	if err := projectMilestonesDeleteCmd.RunE(projectMilestonesDeleteCmd, []string{"ms1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ms1") {
		t.Errorf("expected milestone ID in output, got: %s", out)
	}
}

func TestProjectMilestonesDeleteNoToken(t *testing.T) {
	token = ""
	err := projectMilestonesDeleteCmd.RunE(projectMilestonesDeleteCmd, []string{"ms1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
