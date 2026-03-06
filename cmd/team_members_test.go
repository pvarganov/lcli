package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestTeamMembersListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListTeamMembers": map[string]any{
			"team": map[string]any{
				"members": map[string]any{
					"nodes": []map[string]any{
						{
							"id":   "tm1",
							"user": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
							"team": map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							"role": "member",
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
	teamMembersListCmd.SetOut(&buf)

	if err := teamMembersListCmd.RunE(teamMembersListCmd, []string{"t1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alice Smith") {
		t.Errorf("expected 'Alice Smith' in output, got: %s", out)
	}
	if !strings.Contains(out, "member") {
		t.Errorf("expected 'member' in output, got: %s", out)
	}
}

func TestTeamMembersListJSON(t *testing.T) {
	responses := map[string]any{
		"ListTeamMembers": map[string]any{
			"team": map[string]any{
				"members": map[string]any{
					"nodes": []map[string]any{
						{
							"id":   "tm1",
							"user": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
							"team": map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							"role": "member",
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

	origFormat := outputFormat
	outputFormat = "json"
	defer func() { outputFormat = origFormat }()

	var buf bytes.Buffer
	teamMembersListCmd.SetOut(&buf)

	if err := teamMembersListCmd.RunE(teamMembersListCmd, []string{"t1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestTeamMembersListEmpty(t *testing.T) {
	responses := map[string]any{
		"ListTeamMembers": map[string]any{
			"team": map[string]any{
				"members": map[string]any{
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

	var buf bytes.Buffer
	teamMembersListCmd.SetOut(&buf)

	if err := teamMembersListCmd.RunE(teamMembersListCmd, []string{"t1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestTeamMembersListNoToken(t *testing.T) {
	token = ""
	err := teamMembersListCmd.RunE(teamMembersListCmd, []string{"t1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestTeamMembersAddSuccess(t *testing.T) {
	responses := map[string]any{
		"CreateTeamMembership": map[string]any{
			"teamMembershipCreate": map[string]any{
				"success": true,
				"teamMembership": map[string]any{
					"id":   "tm1",
					"user": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
					"team": map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					"role": "member",
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
	teamMembersAddCmd.SetOut(&buf)

	if err := teamMembersAddCmd.RunE(teamMembersAddCmd, []string{"t1", "u1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "добавлен") {
		t.Errorf("expected 'добавлен' in output, got: %s", out)
	}
}

func TestTeamMembersAddNoToken(t *testing.T) {
	token = ""
	err := teamMembersAddCmd.RunE(teamMembersAddCmd, []string{"t1", "u1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestTeamMembersRemoveSuccess(t *testing.T) {
	responses := map[string]any{
		"DeleteTeamMembership": map[string]any{
			"teamMembershipDelete": map[string]any{"success": true},
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
	teamMembersRemoveCmd.SetOut(&buf)

	if err := teamMembersRemoveCmd.RunE(teamMembersRemoveCmd, []string{"tm1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "удалено") {
		t.Errorf("expected 'удалено' in output, got: %s", out)
	}
}

func TestTeamMembersRemoveNoToken(t *testing.T) {
	token = ""
	err := teamMembersRemoveCmd.RunE(teamMembersRemoveCmd, []string{"tm1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
