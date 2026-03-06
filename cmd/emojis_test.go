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

func newEmojisTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestEmojisList(t *testing.T) {
	responseData := map[string]any{
		"emojis": map[string]any{
			"nodes": []map[string]any{
				{
					"id":   "emoji1",
					"name": "smile",
					"url":  "https://example.com/smile.png",
				},
			},
		},
	}

	srv := newEmojisTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := emojisListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "emoji1") {
		t.Errorf("expected emoji1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "smile") {
		t.Errorf("expected smile in output, got:\n%s", out)
	}
}

func TestEmojisListJSON(t *testing.T) {
	responseData := map[string]any{
		"emojis": map[string]any{
			"nodes": []map[string]any{
				{"id": "emoji1", "name": "smile", "url": "https://example.com/smile.png"},
			},
		},
	}

	srv := newEmojisTestServer(t, responseData)
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

	cmd := emojisListCmd
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

func TestEmojisListEmpty(t *testing.T) {
	responseData := map[string]any{
		"emojis": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newEmojisTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := emojisListCmd
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

func TestEmojisCreate(t *testing.T) {
	responseData := map[string]any{
		"emojiCreate": map[string]any{
			"success": true,
			"emoji": map[string]any{
				"id":   "emoji1",
				"name": "smile",
				"url":  "https://example.com/smile.png",
			},
		},
	}

	srv := newEmojisTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := emojisCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("name", "smile")
	cmd.Flags().Set("url", "https://example.com/smile.png")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "smile") {
		t.Errorf("expected smile in output, got:\n%s", out)
	}
}

func TestEmojisDelete(t *testing.T) {
	responseData := map[string]any{
		"emojiDelete": map[string]any{
			"success": true,
		},
	}

	srv := newEmojisTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := emojisDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"emoji1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "emoji1") {
		t.Errorf("expected emoji1 in output, got:\n%s", out)
	}
}
