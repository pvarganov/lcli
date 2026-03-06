package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestIssueArchiveSuccess(t *testing.T) {
	responses := []map[string]any{
		{"issueArchive": map[string]any{"success": true}},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueArchiveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected ENG-1 in output, got: %s", out)
	}
}

func TestIssueArchiveNoToken(t *testing.T) {
	token = ""
	err := issueArchiveCmd.RunE(issueArchiveCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestIssueArchiveFailure(t *testing.T) {
	responses := []map[string]any{
		{"issueArchive": map[string]any{"success": false}},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	err := issueArchiveCmd.RunE(issueArchiveCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when success=false")
	}
}

func TestIssueUnarchiveSuccess(t *testing.T) {
	responses := []map[string]any{
		{"issueUnarchive": map[string]any{"success": true}},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueUnarchiveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected ENG-1 in output, got: %s", out)
	}
}

func TestIssueUnarchiveNoToken(t *testing.T) {
	token = ""
	err := issueUnarchiveCmd.RunE(issueUnarchiveCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestIssueDeleteSuccess(t *testing.T) {
	responses := []map[string]any{
		{"issueDelete": map[string]any{"success": true}},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "ENG-1") {
		t.Errorf("expected ENG-1 in output, got: %s", out)
	}
}

func TestIssueDeleteNoToken(t *testing.T) {
	token = ""
	err := issueDeleteCmd.RunE(issueDeleteCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestIssueDeleteFailure(t *testing.T) {
	responses := []map[string]any{
		{"issueDelete": map[string]any{"success": false}},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	err := issueDeleteCmd.RunE(issueDeleteCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when success=false")
	}
}
