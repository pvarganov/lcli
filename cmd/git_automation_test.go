package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestGitAutomationStatesListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListGitAutomationStates": map[string]any{
			"team": map[string]any{
				"gitAutomationStates": map[string]any{
					"nodes": []map[string]any{
						{
							"id":    "gas1",
							"event": "branchCreated",
							"state": map[string]any{"id": "ws1", "name": "In Progress", "type": "started", "color": "#f00"},
							"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = gitAutomationStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = gitAutomationStatesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	gitAutomationStatesListCmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := gitAutomationStatesListCmd.RunE(gitAutomationStatesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "branchCreated") {
		t.Errorf("expected 'branchCreated' in output, got: %s", out)
	}
	if !strings.Contains(out, "In Progress") {
		t.Errorf("expected 'In Progress' in output, got: %s", out)
	}
}

func TestGitAutomationStatesListJSON(t *testing.T) {
	responses := map[string]any{
		"ListGitAutomationStates": map[string]any{
			"team": map[string]any{
				"gitAutomationStates": map[string]any{
					"nodes": []map[string]any{
						{
							"id":    "gas1",
							"event": "branchMerged",
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

	_ = gitAutomationStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = gitAutomationStatesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	gitAutomationStatesListCmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := gitAutomationStatesListCmd.RunE(gitAutomationStatesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "branchMerged") {
		t.Errorf("expected 'branchMerged' in json output, got: %s", buf.String())
	}
}

func TestGitAutomationStatesListEmpty(t *testing.T) {
	responses := map[string]any{
		"ListGitAutomationStates": map[string]any{
			"team": map[string]any{
				"gitAutomationStates": map[string]any{
					"nodes": []map[string]any{},
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

	_ = gitAutomationStatesListCmd.Flags().Set("team", "t1")
	defer func() { _ = gitAutomationStatesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	gitAutomationStatesListCmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := gitAutomationStatesListCmd.RunE(gitAutomationStatesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestGitAutomationStatesCreateCmd(t *testing.T) {
	responses := map[string]any{
		"gitAutomationStateCreate": map[string]any{
			"gitAutomationStateCreate": map[string]any{
				"success": true,
				"gitAutomationState": map[string]any{
					"id":    "gas1",
					"event": "branchCreated",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = gitAutomationStatesCreateCmd.Flags().Set("team", "t1")
	_ = gitAutomationStatesCreateCmd.Flags().Set("event", "branchCreated")
	defer func() {
		_ = gitAutomationStatesCreateCmd.Flags().Set("team", "")
		_ = gitAutomationStatesCreateCmd.Flags().Set("event", "")
	}()

	var buf bytes.Buffer
	gitAutomationStatesCreateCmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := gitAutomationStatesCreateCmd.RunE(gitAutomationStatesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "gas1") {
		t.Errorf("expected ID in output, got: %s", buf.String())
	}
}

func TestGitAutomationStatesDeleteCmd(t *testing.T) {
	responses := map[string]any{
		"gitAutomationStateDelete": map[string]any{
			"gitAutomationStateDelete": map[string]any{
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
	gitAutomationStatesDeleteCmd.SetOut(&buf)

	if err := gitAutomationStatesDeleteCmd.RunE(gitAutomationStatesDeleteCmd, []string{"gas1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "gas1") {
		t.Errorf("expected ID in output, got: %s", buf.String())
	}
}

func TestGitAutomationBranchesCreateCmd(t *testing.T) {
	responses := map[string]any{
		"gitAutomationTargetBranchCreate": map[string]any{
			"gitAutomationTargetBranchCreate": map[string]any{
				"success": true,
				"gitAutomationTargetBranch": map[string]any{
					"id":            "tb1",
					"branchPattern": "main",
					"isRegex":       false,
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

	_ = gitAutomationBranchesCreateCmd.Flags().Set("team", "t1")
	_ = gitAutomationBranchesCreateCmd.Flags().Set("pattern", "main")
	defer func() {
		_ = gitAutomationBranchesCreateCmd.Flags().Set("team", "")
		_ = gitAutomationBranchesCreateCmd.Flags().Set("pattern", "")
	}()

	var buf bytes.Buffer
	gitAutomationBranchesCreateCmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := gitAutomationBranchesCreateCmd.RunE(gitAutomationBranchesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "main") {
		t.Errorf("expected 'main' in output, got: %s", buf.String())
	}
}

func TestGitAutomationBranchesDeleteCmd(t *testing.T) {
	responses := map[string]any{
		"gitAutomationTargetBranchDelete": map[string]any{
			"gitAutomationTargetBranchDelete": map[string]any{
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
	gitAutomationBranchesDeleteCmd.SetOut(&buf)

	if err := gitAutomationBranchesDeleteCmd.RunE(gitAutomationBranchesDeleteCmd, []string{"tb1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "tb1") {
		t.Errorf("expected ID in output, got: %s", buf.String())
	}
}
