package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectsSearchSuccess(t *testing.T) {
	responses := map[string]any{
		"SearchProjects": map[string]any{
			"projects": map[string]any{
				"nodes": []map[string]any{
					{"id": "p1", "name": "Alpha Project", "description": "Test", "state": "started", "startDate": "", "targetDate": "", "url": "", "lead": nil},
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
	projectsSearchCmd.SetOut(&buf)

	if err := projectsSearchCmd.RunE(projectsSearchCmd, []string{"alpha"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Alpha Project") {
		t.Errorf("expected 'Alpha Project' in output, got: %s", out)
	}
}

func TestProjectsSearchJSON(t *testing.T) {
	responses := map[string]any{
		"SearchProjects": map[string]any{
			"projects": map[string]any{
				"nodes": []map[string]any{
					{"id": "p1", "name": "Alpha Project", "description": "", "state": "started", "startDate": "", "targetDate": "", "url": "", "lead": nil},
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
	projectsSearchCmd.SetOut(&buf)

	if err := projectsSearchCmd.RunE(projectsSearchCmd, []string{"alpha"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestProjectsSearchEmpty(t *testing.T) {
	responses := map[string]any{
		"SearchProjects": map[string]any{
			"projects": map[string]any{
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
	projectsSearchCmd.SetOut(&buf)

	if err := projectsSearchCmd.RunE(projectsSearchCmd, []string{"nonexistent"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected 'не найдены' in output, got: %s", out)
	}
}

func TestProjectsSearchNoToken(t *testing.T) {
	token = ""
	err := projectsSearchCmd.RunE(projectsSearchCmd, []string{"test"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
