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

func newCustomViewsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestViewsList(t *testing.T) {
	responseData := map[string]any{
		"customViews": map[string]any{
			"nodes": []map[string]any{
				{"id": "cv1", "name": "My View", "description": "A view", "icon": "eye", "color": "#ff0000"},
			},
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := viewsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cv1") {
		t.Errorf("expected cv1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "My View") {
		t.Errorf("expected 'My View' in output, got:\n%s", out)
	}
}

func TestViewsListJSON(t *testing.T) {
	responseData := map[string]any{
		"customViews": map[string]any{
			"nodes": []map[string]any{
				{"id": "cv1", "name": "My View"},
			},
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
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

	cmd := viewsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output, got:\n%s", out)
	}
}

func TestViewsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"customViews": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := viewsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}

func TestViewsView(t *testing.T) {
	responseData := map[string]any{
		"customView": map[string]any{
			"id":          "cv1",
			"name":        "My View",
			"description": "A custom view",
			"icon":        "eye",
			"color":       "#ff0000",
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := viewsViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"cv1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cv1") {
		t.Errorf("expected cv1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "My View") {
		t.Errorf("expected 'My View' in output, got:\n%s", out)
	}
}

func TestViewsCreate(t *testing.T) {
	responseData := map[string]any{
		"customViewCreate": map[string]any{
			"success": true,
			"customView": map[string]any{
				"id":   "cv1",
				"name": "My View",
			},
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := viewsCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "My View")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cv1") {
		t.Errorf("expected cv1 in output, got:\n%s", out)
	}
}

func TestViewsUpdate(t *testing.T) {
	responseData := map[string]any{
		"customViewUpdate": map[string]any{
			"success": true,
			"customView": map[string]any{
				"id":   "cv1",
				"name": "Updated View",
			},
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := viewsUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "Updated View")

	if err := cmd.RunE(cmd, []string{"cv1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cv1") {
		t.Errorf("expected cv1 in output, got:\n%s", out)
	}
}

func TestViewsDelete(t *testing.T) {
	responseData := map[string]any{
		"customViewDelete": map[string]any{
			"success": true,
		},
	}

	srv := newCustomViewsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := viewsDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"cv1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cv1") {
		t.Errorf("expected cv1 in output, got:\n%s", out)
	}
}
