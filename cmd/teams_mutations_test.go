package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestTeamsCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"CreateTeam": map[string]any{
			"teamCreate": map[string]any{
				"success": true,
				"team":    map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = teamsCreateCmd.Flags().Set("name", "Engineering")
	_ = teamsCreateCmd.Flags().Set("key", "ENG")
	defer func() {
		_ = teamsCreateCmd.Flags().Set("name", "")
		_ = teamsCreateCmd.Flags().Set("key", "")
	}()

	var buf bytes.Buffer
	teamsCreateCmd.SetOut(&buf)

	if err := teamsCreateCmd.RunE(teamsCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "создана") {
		t.Errorf("expected 'создана' in output, got: %s", out)
	}
}

func TestTeamsCreateNoToken(t *testing.T) {
	token = ""
	_ = teamsCreateCmd.Flags().Set("name", "Engineering")
	_ = teamsCreateCmd.Flags().Set("key", "ENG")
	defer func() {
		_ = teamsCreateCmd.Flags().Set("name", "")
		_ = teamsCreateCmd.Flags().Set("key", "")
	}()

	err := teamsCreateCmd.RunE(teamsCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestTeamsCreateNoName(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = teamsCreateCmd.Flags().Set("name", "")
	_ = teamsCreateCmd.Flags().Set("key", "ENG")
	defer func() { _ = teamsCreateCmd.Flags().Set("key", "") }()

	err := teamsCreateCmd.RunE(teamsCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no name")
	}
}

func TestTeamsCreateJSON(t *testing.T) {
	responses := map[string]any{
		"CreateTeam": map[string]any{
			"teamCreate": map[string]any{
				"success": true,
				"team":    map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = teamsCreateCmd.Flags().Set("name", "Engineering")
	_ = teamsCreateCmd.Flags().Set("key", "ENG")
	defer func() {
		_ = teamsCreateCmd.Flags().Set("name", "")
		_ = teamsCreateCmd.Flags().Set("key", "")
	}()

	var buf bytes.Buffer
	teamsCreateCmd.SetOut(&buf)

	if err := teamsCreateCmd.RunE(teamsCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestTeamsUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"UpdateTeam": map[string]any{
			"teamUpdate": map[string]any{
				"success": true,
				"team":    map[string]any{"id": "t1", "key": "ENG", "name": "Engineering Updated"},
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

	_ = teamsUpdateCmd.Flags().Set("name", "Engineering Updated")
	defer func() { _ = teamsUpdateCmd.Flags().Set("name", "") }()

	var buf bytes.Buffer
	teamsUpdateCmd.SetOut(&buf)

	if err := teamsUpdateCmd.RunE(teamsUpdateCmd, []string{"t1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "обновлена") {
		t.Errorf("expected 'обновлена' in output, got: %s", out)
	}
}

func TestTeamsUpdateNoToken(t *testing.T) {
	token = ""
	err := teamsUpdateCmd.RunE(teamsUpdateCmd, []string{"t1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestTeamsUpdateJSON(t *testing.T) {
	responses := map[string]any{
		"UpdateTeam": map[string]any{
			"teamUpdate": map[string]any{
				"success": true,
				"team":    map[string]any{"id": "t1", "key": "ENG", "name": "Engineering Updated"},
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
	teamsUpdateCmd.SetOut(&buf)

	if err := teamsUpdateCmd.RunE(teamsUpdateCmd, []string{"t1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestTeamsDeleteSuccess(t *testing.T) {
	responses := map[string]any{
		"DeleteTeam": map[string]any{
			"teamDelete": map[string]any{"success": true},
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
	teamsDeleteCmd.SetOut(&buf)

	if err := teamsDeleteCmd.RunE(teamsDeleteCmd, []string{"t1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "удалена") {
		t.Errorf("expected 'удалена' in output, got: %s", out)
	}
}

func TestTeamsDeleteNoToken(t *testing.T) {
	token = ""
	err := teamsDeleteCmd.RunE(teamsDeleteCmd, []string{"t1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
