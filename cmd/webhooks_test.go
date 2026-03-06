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

func newWebhooksTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestWebhooksListOutput(t *testing.T) {
	responseData := map[string]any{
		"webhooks": map[string]any{
			"nodes": []map[string]any{
				{
					"id":            "wh1",
					"url":           "https://example.com/hook",
					"enabled":       true,
					"secret":        "mysecret",
					"resourceTypes": []string{"Issue", "Comment"},
					"team":          map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"wh1", "https://example.com/hook", "yes", "ENG"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWebhooksListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"webhooks": map[string]any{
			"nodes": []map[string]any{
				{
					"id":      "wh1",
					"url":     "https://example.com/hook",
					"enabled": true,
				},
			},
		},
	}

	srv := newWebhooksTestServer(t, responseData)
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

	cmd := webhooksListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []map[string]any
	if err := json.Unmarshal([]byte(buf.String()), &result); err != nil {
		t.Fatalf("expected valid JSON: %v\noutput: %s", err, buf.String())
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(result))
	}
}

func TestWebhooksListEmpty(t *testing.T) {
	responseData := map[string]any{
		"webhooks": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected empty message, got:\n%s", buf.String())
	}
}

func TestWebhooksViewOutput(t *testing.T) {
	responseData := map[string]any{
		"webhook": map[string]any{
			"id":            "wh1",
			"url":           "https://example.com/hook",
			"enabled":       true,
			"secret":        "mysecret",
			"resourceTypes": []string{"Issue"},
			"team":          map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksViewCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"wh1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"wh1", "https://example.com/hook", "mysecret", "ENG"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWebhooksCreateOutput(t *testing.T) {
	responseData := map[string]any{
		"webhookCreate": map[string]any{
			"success": true,
			"webhook": map[string]any{
				"id":      "wh1",
				"url":     "https://example.com/hook",
				"enabled": true,
			},
		},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksCreateCmd
	cmd.Flags().Set("url", "https://example.com/hook")
	cmd.Flags().Set("team-id", "t1")

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "wh1") {
		t.Errorf("expected wh1 in output, got:\n%s", buf.String())
	}
}

func TestWebhooksUpdateOutput(t *testing.T) {
	responseData := map[string]any{
		"webhookUpdate": map[string]any{
			"success": true,
			"webhook": map[string]any{
				"id":      "wh1",
				"url":     "https://example.com/hook2",
				"enabled": false,
			},
		},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksUpdateCmd
	cmd.Flags().Set("url", "https://example.com/hook2")

	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"wh1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "wh1") {
		t.Errorf("expected wh1 in output, got:\n%s", buf.String())
	}
}

func TestWebhooksDeleteOutput(t *testing.T) {
	responseData := map[string]any{
		"webhookDelete": map[string]any{"success": true},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"wh1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "wh1") {
		t.Errorf("expected wh1 in output, got:\n%s", buf.String())
	}
}

func TestWebhooksRotateSecretOutput(t *testing.T) {
	responseData := map[string]any{
		"webhookRotateSecret": map[string]any{
			"success": true,
			"secret":  "newsecret",
		},
	}

	srv := newWebhooksTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := webhooksRotateSecretCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"wh1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"wh1", "newsecret"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWebhooksNoToken(t *testing.T) {
	token = ""
	err := webhooksListCmd.RunE(webhooksListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}
