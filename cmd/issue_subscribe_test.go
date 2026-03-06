package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestIssueSubscribeSuccess(t *testing.T) {
	responses := []map[string]any{
		{"issueSubscribe": map[string]any{"success": true}},
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

	cmd := issueSubscribeCmd
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

func TestIssueSubscribeNoToken(t *testing.T) {
	token = ""
	err := issueSubscribeCmd.RunE(issueSubscribeCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestIssueSubscribeFailure(t *testing.T) {
	responses := []map[string]any{
		{"issueSubscribe": map[string]any{"success": false}},
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

	err := issueSubscribeCmd.RunE(issueSubscribeCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when success=false")
	}
}

func TestIssueUnsubscribeSuccess(t *testing.T) {
	responses := []map[string]any{
		{"issueUnsubscribe": map[string]any{"success": true}},
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

	cmd := issueUnsubscribeCmd
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

func TestIssueUnsubscribeNoToken(t *testing.T) {
	token = ""
	err := issueUnsubscribeCmd.RunE(issueUnsubscribeCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestIssueUnsubscribeFailure(t *testing.T) {
	responses := []map[string]any{
		{"issueUnsubscribe": map[string]any{"success": false}},
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

	err := issueUnsubscribeCmd.RunE(issueUnsubscribeCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when success=false")
	}
}
