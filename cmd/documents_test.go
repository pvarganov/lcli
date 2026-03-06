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

func newDocumentsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestDocumentsListOutput(t *testing.T) {
	responseData := map[string]any{
		"documents": map[string]any{
			"nodes": []map[string]any{
				{
					"id":      "doc1",
					"title":   "Getting Started",
					"content": "# Hello",
					"project": map[string]any{"id": "proj1", "name": "My Project"},
					"creator": map[string]any{"id": "user1", "name": "Alice"},
				},
			},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"doc1", "Getting Started", "My Project", "Alice"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestDocumentsListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"documents": map[string]any{
			"nodes": []map[string]any{
				{"id": "doc1", "title": "Getting Started"},
			},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
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

	cmd := documentsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"doc1"`) {
		t.Errorf("expected doc1 in JSON output, got:\n%s", out)
	}
}

func TestDocumentsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"documents": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsListCmd
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

func TestDocumentsView(t *testing.T) {
	responseData := map[string]any{
		"document": map[string]any{
			"id":      "doc1",
			"title":   "Getting Started",
			"content": "# Hello World",
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"doc1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"doc1", "Getting Started", "# Hello World"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestDocumentsCreate(t *testing.T) {
	responseData := map[string]any{
		"documentCreate": map[string]any{
			"success": true,
			"document": map[string]any{
				"id":      "doc1",
				"title":   "New Doc",
				"content": "# Content",
			},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("title", "New Doc")
	cmd.Flags().Set("content", "# Content")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "doc1") {
		t.Errorf("expected doc1 in output, got:\n%s", out)
	}
}

func TestDocumentsUpdate(t *testing.T) {
	responseData := map[string]any{
		"documentUpdate": map[string]any{
			"success": true,
			"document": map[string]any{
				"id":    "doc1",
				"title": "Updated Doc",
			},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("title", "Updated Doc")

	if err := cmd.RunE(cmd, []string{"doc1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "doc1") {
		t.Errorf("expected doc1 in output, got:\n%s", out)
	}
}

func TestDocumentsDelete(t *testing.T) {
	responseData := map[string]any{
		"documentDelete": map[string]any{
			"success": true,
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"doc1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "doc1") {
		t.Errorf("expected doc1 in output, got:\n%s", out)
	}
}

func TestDocumentsSearch(t *testing.T) {
	responseData := map[string]any{
		"searchDocuments": map[string]any{
			"nodes": []map[string]any{
				{
					"id":    "doc1",
					"title": "Getting Started",
				},
			},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsSearchCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"getting"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"doc1", "Getting Started"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestDocumentsSearchEmpty(t *testing.T) {
	responseData := map[string]any{
		"searchDocuments": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newDocumentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := documentsSearchCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"nonexistent"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}
