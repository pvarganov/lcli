package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestIssueCommentOutput(t *testing.T) {
	now := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"commentCreate": map[string]any{
					"success": true,
					"comment": map[string]any{
						"id":        "cmt-1",
						"body":      "Test comment",
						"createdAt": now.Format(time.RFC3339),
						"user": map[string]any{
							"id":          "u1",
							"name":        "John",
							"displayName": "John Doe",
							"email":       "john@example.com",
						},
					},
				},
			},
		})
		w.Write(body)
	}))
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	commentBody = "Test comment"
	defer func() { commentBody = "" }()

	cmd := issueCommentCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "cmt-1") {
		t.Errorf("expected comment id in output, got: %s", out)
	}
	if !strings.Contains(out, "Комментарий добавлен") {
		t.Errorf("expected success message in output, got: %s", out)
	}
}

func TestIssueCommentNoToken(t *testing.T) {
	token = ""
	commentBody = "some text"
	defer func() { commentBody = "" }()
	err := issueCommentCmd.RunE(issueCommentCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestIssueCommentsOutput(t *testing.T) {
	now := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"issue": map[string]any{
					"comments": map[string]any{
						"nodes": []map[string]any{
							{
								"id":        "cmt-1",
								"body":      "First comment",
								"createdAt": now.Format(time.RFC3339),
								"user": map[string]any{
									"id":          "u1",
									"name":        "Alice",
									"displayName": "Alice Smith",
									"email":       "alice@example.com",
								},
							},
							{
								"id":        "cmt-2",
								"body":      "Second comment",
								"createdAt": now.Add(time.Hour).Format(time.RFC3339),
								"user": map[string]any{
									"id":          "u2",
									"name":        "Bob",
									"displayName": "Bob Jones",
									"email":       "bob@example.com",
								},
							},
						},
					},
				},
			},
		})
		w.Write(body)
	}))
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueCommentsCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "First comment") {
		t.Errorf("expected 'First comment' in output, got: %s", out)
	}
	if !strings.Contains(out, "Second comment") {
		t.Errorf("expected 'Second comment' in output, got: %s", out)
	}
	if !strings.Contains(out, "Alice Smith") {
		t.Errorf("expected 'Alice Smith' in output, got: %s", out)
	}
}

func TestIssueCommentsEmpty(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"issue": map[string]any{
					"comments": map[string]any{
						"nodes": []map[string]any{},
					},
				},
			},
		})
		w.Write(body)
	}))
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueCommentsCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Комментариев нет") {
		t.Errorf("expected empty message in output, got: %s", out)
	}
}

func TestIssueCommentsNoToken(t *testing.T) {
	token = ""
	err := issueCommentsCmd.RunE(issueCommentsCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestIssueCommentsJSONOutput(t *testing.T) {
	now := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{
			"data": map[string]any{
				"issue": map[string]any{
					"comments": map[string]any{
						"nodes": []map[string]any{
							{
								"id":        "cmt-1",
								"body":      "JSON comment",
								"createdAt": now.Format(time.RFC3339),
								"user":      nil,
							},
						},
					},
				},
			},
		})
		w.Write(body)
	}))
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

	cmd := issueCommentsCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var comments []map[string]any
	if err := json.Unmarshal(buf.Bytes(), &comments); err != nil {
		t.Fatalf("expected valid JSON, got: %v\noutput: %s", err, buf.String())
	}
	if len(comments) != 1 || comments[0]["id"] != "cmt-1" {
		t.Errorf("unexpected JSON output: %v", comments)
	}
}
