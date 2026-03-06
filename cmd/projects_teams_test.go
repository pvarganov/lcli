package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func newProjectsTeamsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestProjectsListOutput(t *testing.T) {
	responseData := map[string]any{
		"projects": map[string]any{
			"nodes": []map[string]any{
				{
					"id":          "proj1",
					"name":        "Mobile App",
					"description": "iOS and Android app",
					"state":       "started",
				},
			},
		},
	}

	srv := newProjectsTeamsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := projectsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"Mobile App", "started", "iOS and Android app"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestProjectsListNoToken(t *testing.T) {
	token = ""
	err := projectsListCmd.RunE(projectsListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestProjectsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"projects": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newProjectsTeamsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := projectsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}

func TestTeamsListOutput(t *testing.T) {
	responseData := map[string]any{
		"teams": map[string]any{
			"nodes": []map[string]any{
				{
					"id":   "team1",
					"key":  "ENG",
					"name": "Engineering",
				},
				{
					"id":   "team2",
					"key":  "DES",
					"name": "Design",
				},
			},
		},
	}

	srv := newProjectsTeamsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := teamsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"ENG", "Engineering", "DES", "Design"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestTeamsListNoToken(t *testing.T) {
	token = ""
	err := teamsListCmd.RunE(teamsListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestProjectsListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"projects": map[string]any{
			"nodes": []map[string]any{
				{
					"id":          "proj1",
					"name":        "Mobile App",
					"description": "iOS and Android app",
					"state":       "started",
				},
			},
		},
	}

	srv := newProjectsTeamsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	outputFormat = "json"
	defer func() { outputFormat = "" }()

	cmd := projectsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var projects []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &projects); err != nil {
		t.Fatalf("expected valid JSON, got: %v\noutput: %s", err, buf.String())
	}
	if len(projects) != 1 || projects[0]["name"] != "Mobile App" {
		t.Errorf("unexpected JSON output: %v", projects)
	}
}

func TestTeamsListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"teams": map[string]any{
			"nodes": []map[string]any{
				{"id": "team1", "key": "ENG", "name": "Engineering"},
			},
		},
	}

	srv := newProjectsTeamsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	outputFormat = "json"
	defer func() { outputFormat = "" }()

	cmd := teamsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var teams []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &teams); err != nil {
		t.Fatalf("expected valid JSON, got: %v\noutput: %s", err, buf.String())
	}
	if len(teams) != 1 || teams[0]["key"] != "ENG" {
		t.Errorf("unexpected JSON output: %v", teams)
	}
}

func TestTeamsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"teams": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newProjectsTeamsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := teamsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected empty message, got: %s", buf.String())
	}
}
