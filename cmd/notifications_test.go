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

func newNotificationsTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestNotificationsListOutput(t *testing.T) {
	now := time.Now().UTC()
	responseData := map[string]any{
		"notifications": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "n1",
					"type":      "issueAssignedToYou",
					"readAt":    nil,
					"createdAt": now.Format(time.RFC3339),
					"issue": map[string]any{
						"id": "i1", "identifier": "ENG-1", "title": "Test Issue",
					},
				},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
		},
	}

	srv := newNotificationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := notificationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"n1", "issueAssignedToYou", "ENG-1"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestNotificationsListJSONOutput(t *testing.T) {
	now := time.Now().UTC()
	responseData := map[string]any{
		"notifications": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "n1",
					"type":      "issueAssignedToYou",
					"readAt":    nil,
					"createdAt": now.Format(time.RFC3339),
				},
			},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
		},
	}

	srv := newNotificationsTestServer(t, responseData)
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

	cmd := notificationsListCmd
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
		t.Fatalf("expected 1 notification, got %d", len(result))
	}
}

func TestNotificationsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"notifications": map[string]any{
			"nodes":    []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
		},
	}

	srv := newNotificationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := notificationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected empty message, got:\n%s", buf.String())
	}
}

func TestNotificationsUnreadCountOutput(t *testing.T) {
	responseData := map[string]any{
		"notifications": map[string]any{
			"nodes": []map[string]any{
				{"id": "n1"},
				{"id": "n2"},
			},
		},
	}

	srv := newNotificationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := notificationsUnreadCountCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "2") {
		t.Errorf("expected count 2 in output, got:\n%s", buf.String())
	}
}

func TestNotificationsMarkReadOutput(t *testing.T) {
	responseData := map[string]any{
		"notificationUpdate": map[string]any{"success": true},
	}

	srv := newNotificationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := notificationsMarkReadCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"n1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "n1") {
		t.Errorf("expected n1 in output, got:\n%s", buf.String())
	}
}

func TestNotificationsMarkReadAll(t *testing.T) {
	responseData := map[string]any{
		"notificationMarkReadAll": map[string]any{"success": true},
	}

	srv := newNotificationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := notificationsMarkReadCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"all"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "Все") {
		t.Errorf("expected success message, got:\n%s", buf.String())
	}
}

func TestNotificationsArchiveOutput(t *testing.T) {
	responseData := map[string]any{
		"notificationArchive": map[string]any{"success": true},
	}

	srv := newNotificationsTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := notificationsArchiveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"n1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "n1") {
		t.Errorf("expected n1 in output, got:\n%s", buf.String())
	}
}

func TestNotificationsNoToken(t *testing.T) {
	token = ""
	err := notificationsListCmd.RunE(notificationsListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}
