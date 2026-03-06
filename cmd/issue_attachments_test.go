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

func newAttachmentsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestIssueAttachmentsListOutput(t *testing.T) {
	responseData := map[string]any{
		"issue": map[string]any{
			"attachments": map[string]any{
				"nodes": []map[string]any{
					{
						"id":         "att1",
						"title":      "GitHub PR #42",
						"url":        "https://github.com/org/repo/pull/42",
						"sourceType": "github",
						"subtitle":   "Open",
					},
				},
			},
		},
	}

	srv := newAttachmentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueAttachmentsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"issue1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"att1", "GitHub PR #42", "github", "Open"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestIssueAttachmentsListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"issue": map[string]any{
			"attachments": map[string]any{
				"nodes": []map[string]any{
					{
						"id":    "att1",
						"title": "GitHub PR #42",
						"url":   "https://github.com/org/repo/pull/42",
					},
				},
			},
		},
	}

	srv := newAttachmentsTestServer(t, responseData)
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

	cmd := issueAttachmentsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"issue1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"att1"`) {
		t.Errorf("expected att1 in JSON output, got:\n%s", out)
	}
}

func TestIssueAttachmentsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"issue": map[string]any{
			"attachments": map[string]any{
				"nodes": []map[string]any{},
			},
		},
	}

	srv := newAttachmentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueAttachmentsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"issue1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got:\n%s", out)
	}
}

func TestIssueAttachmentsLinkURL(t *testing.T) {
	responseData := map[string]any{
		"attachmentLinkURL": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att1",
				"title":      "Linear Docs",
				"url":        "https://linear.app/docs",
				"sourceType": "url",
				"subtitle":   "",
			},
		},
	}

	srv := newAttachmentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueAttachmentsLinkURLCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("url", "https://linear.app/docs")
	cmd.Flags().Set("title", "Linear Docs")

	if err := cmd.RunE(cmd, []string{"issue1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "att1") {
		t.Errorf("expected att1 in output, got:\n%s", out)
	}
}

func TestIssueAttachmentsLinkGitHubPR(t *testing.T) {
	responseData := map[string]any{
		"attachmentLinkGitHubPR": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att2",
				"title":      "Fix bug",
				"url":        "https://github.com/org/repo/pull/1",
				"sourceType": "github",
				"subtitle":   "Open",
			},
		},
	}

	srv := newAttachmentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueAttachmentsLinkGitHubPRCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("url", "https://github.com/org/repo/pull/1")
	cmd.Flags().Set("title", "Fix bug")

	if err := cmd.RunE(cmd, []string{"issue1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "att2") {
		t.Errorf("expected att2 in output, got:\n%s", out)
	}
}

func TestIssueAttachmentsDelete(t *testing.T) {
	responseData := map[string]any{
		"attachmentDelete": map[string]any{
			"success": true,
		},
	}

	srv := newAttachmentsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := issueAttachmentsDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"att1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "att1") {
		t.Errorf("expected att1 in output, got:\n%s", out)
	}
}
