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

func newMiscTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestRateLimitOutput(t *testing.T) {
	responseData := map[string]any{
		"rateLimitStatus": map[string]any{
			"identifier": "org-123",
			"kind":       "requestComplexity",
			"limits": []map[string]any{
				{
					"type":            "requestComplexity",
					"allowedAmount":   10000,
					"requestedAmount": 5,
					"remainingAmount": 9995,
					"period":          3600,
					"reset":           "2026-03-06T12:00:00Z",
				},
			},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := rateLimitCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "10000") {
		t.Errorf("expected maxComplexity '10000' in output, got: %s", output)
	}
}

func TestRateLimitJSON(t *testing.T) {
	responseData := map[string]any{
		"rateLimitStatus": map[string]any{
			"identifier": "org-abc",
			"kind":       "requestComplexity",
			"limits":     []map[string]any{},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := rateLimitCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "org-abc") {
		t.Errorf("expected 'org-abc' in JSON output, got: %s", output)
	}
}

func TestTimeSchedulesListOutput(t *testing.T) {
	responseData := map[string]any{
		"timeSchedules": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "ts1",
					"name":      "On-call",
					"timezone":  "UTC",
					"createdAt": "2026-01-01T00:00:00Z",
					"updatedAt": "2026-01-01T00:00:00Z",
				},
			},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := timeSchedulesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "On-call") {
		t.Errorf("expected 'On-call' in output, got: %s", output)
	}
}

func TestTimeSchedulesListJSON(t *testing.T) {
	responseData := map[string]any{
		"timeSchedules": map[string]any{
			"nodes": []map[string]any{
				{
					"id":        "ts2",
					"name":      "Weekend Schedule",
					"timezone":  "Europe/Moscow",
					"createdAt": "2026-01-01T00:00:00Z",
					"updatedAt": "2026-01-01T00:00:00Z",
				},
			},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := timeSchedulesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Weekend Schedule") {
		t.Errorf("expected 'Weekend Schedule' in JSON output, got: %s", output)
	}
}

func TestTimeSchedulesCreate(t *testing.T) {
	responseData := map[string]any{
		"timeScheduleCreate": map[string]any{
			"success": true,
			"timeSchedule": map[string]any{
				"id":        "ts3",
				"name":      "New Schedule",
				"timezone":  "UTC",
				"createdAt": "2026-01-01T00:00:00Z",
				"updatedAt": "2026-01-01T00:00:00Z",
			},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := timeSchedulesCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("name", "New Schedule")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "New Schedule") {
		t.Errorf("expected 'New Schedule' in output, got: %s", output)
	}
}

func TestTimeSchedulesDelete(t *testing.T) {
	responseData := map[string]any{
		"timeScheduleDelete": map[string]any{
			"success": true,
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := timeSchedulesDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	if err := cmd.RunE(cmd, []string{"ts1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "ts1") {
		t.Errorf("expected 'ts1' in output, got: %s", output)
	}
}

func TestTriageResponsibilitiesListOutput(t *testing.T) {
	responseData := map[string]any{
		"team": map[string]any{
			"triageResponsibility": map[string]any{
				"id":          "tr1",
				"createdAt":   "2026-01-01T00:00:00Z",
				"action":      "assignIssues",
				"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				"currentUser": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
			},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := triageResponsibilitiesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = ""
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("team", "t1")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "assignIssues") {
		t.Errorf("expected 'assignIssues' in output, got: %s", output)
	}
	if !strings.Contains(output, "Alice") {
		t.Errorf("expected 'Alice' in output, got: %s", output)
	}
}

func TestTriageResponsibilitiesListJSON(t *testing.T) {
	responseData := map[string]any{
		"team": map[string]any{
			"triageResponsibility": map[string]any{
				"id":          "tr2",
				"createdAt":   "2026-01-01T00:00:00Z",
				"action":      "notifyIssues",
				"team":        map[string]any{"id": "t2", "key": "OPS", "name": "Operations"},
				"currentUser": nil,
			},
		},
	}

	srv := newMiscTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := triageResponsibilitiesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	outputFormat = "json"
	defer func() { outputFormat = "" }()

	cmd.Flags().Set("team", "t2")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "notifyIssues") {
		t.Errorf("expected 'notifyIssues' in JSON output, got: %s", output)
	}
}
