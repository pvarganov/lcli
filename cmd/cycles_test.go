package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestCyclesListSuccess(t *testing.T) {
	responses := map[string]any{
		"ListCycles": map[string]any{
			"cycles": map[string]any{
				"nodes": []map[string]any{
					{
						"id":          "cycle1",
						"number":      1,
						"name":        "Sprint 1",
						"startsAt":    "2024-01-01T00:00:00Z",
						"endsAt":      "2024-01-14T00:00:00Z",
						"completedAt": nil,
						"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					},
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

	_ = cyclesListCmd.Flags().Set("team", "t1")
	defer func() { _ = cyclesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	cyclesListCmd.SetOut(&buf)

	if err := cyclesListCmd.RunE(cyclesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Sprint 1") {
		t.Errorf("expected 'Sprint 1' in output, got: %s", out)
	}
	if !strings.Contains(out, "2024-01-01") {
		t.Errorf("expected start date in output, got: %s", out)
	}
}

func TestCyclesListJSON(t *testing.T) {
	responses := map[string]any{
		"ListCycles": map[string]any{
			"cycles": map[string]any{
				"nodes": []map[string]any{
					{
						"id":          "cycle1",
						"number":      1,
						"name":        "Sprint 1",
						"startsAt":    "2024-01-01T00:00:00Z",
						"endsAt":      "2024-01-14T00:00:00Z",
						"completedAt": nil,
						"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					},
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

	_ = cyclesListCmd.Flags().Set("team", "t1")
	defer func() { _ = cyclesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	cyclesListCmd.SetOut(&buf)

	if err := cyclesListCmd.RunE(cyclesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON with 'id', got: %s", out)
	}
}

func TestCyclesListEmpty(t *testing.T) {
	responses := map[string]any{
		"ListCycles": map[string]any{
			"cycles": map[string]any{
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

	_ = cyclesListCmd.Flags().Set("team", "t1")
	defer func() { _ = cyclesListCmd.Flags().Set("team", "") }()

	var buf bytes.Buffer
	cyclesListCmd.SetOut(&buf)

	if err := cyclesListCmd.RunE(cyclesListCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestCyclesListNoToken(t *testing.T) {
	token = ""
	_ = cyclesListCmd.Flags().Set("team", "t1")
	defer func() { _ = cyclesListCmd.Flags().Set("team", "") }()

	err := cyclesListCmd.RunE(cyclesListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestCyclesListNoTeam(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = cyclesListCmd.Flags().Set("team", "")

	err := cyclesListCmd.RunE(cyclesListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no team")
	}
}

func TestCyclesViewSuccess(t *testing.T) {
	responses := map[string]any{
		"GetCycle": map[string]any{
			"cycle": map[string]any{
				"id":          "cycle1",
				"number":      2,
				"name":        "Sprint 2",
				"startsAt":    "2024-01-15T00:00:00Z",
				"endsAt":      "2024-01-28T00:00:00Z",
				"completedAt": nil,
				"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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
	cyclesViewCmd.SetOut(&buf)

	if err := cyclesViewCmd.RunE(cyclesViewCmd, []string{"cycle1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Sprint 2") {
		t.Errorf("expected 'Sprint 2' in output, got: %s", out)
	}
	if !strings.Contains(out, "ENG") {
		t.Errorf("expected 'ENG' in output, got: %s", out)
	}
}

func TestCyclesViewJSON(t *testing.T) {
	responses := map[string]any{
		"GetCycle": map[string]any{
			"cycle": map[string]any{
				"id":          "cycle1",
				"number":      1,
				"name":        "Sprint 1",
				"startsAt":    "2024-01-01T00:00:00Z",
				"endsAt":      "2024-01-14T00:00:00Z",
				"completedAt": nil,
				"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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
	cyclesViewCmd.SetOut(&buf)

	if err := cyclesViewCmd.RunE(cyclesViewCmd, []string{"cycle1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output with 'id', got: %s", out)
	}
}

func TestCyclesViewNoToken(t *testing.T) {
	token = ""
	err := cyclesViewCmd.RunE(cyclesViewCmd, []string{"cycle1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestCyclesCreateSuccess(t *testing.T) {
	responses := map[string]any{
		"CreateCycle": map[string]any{
			"cycleCreate": map[string]any{
				"success": true,
				"cycle": map[string]any{
					"id":          "cycle1",
					"number":      1,
					"name":        "Sprint 1",
					"startsAt":    "2024-01-01T00:00:00Z",
					"endsAt":      "2024-01-14T00:00:00Z",
					"completedAt": nil,
					"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = cyclesCreateCmd.Flags().Set("team", "t1")
	_ = cyclesCreateCmd.Flags().Set("starts-at", "2024-01-01T00:00:00Z")
	_ = cyclesCreateCmd.Flags().Set("ends-at", "2024-01-14T00:00:00Z")
	_ = cyclesCreateCmd.Flags().Set("name", "Sprint 1")
	defer func() {
		_ = cyclesCreateCmd.Flags().Set("team", "")
		_ = cyclesCreateCmd.Flags().Set("starts-at", "")
		_ = cyclesCreateCmd.Flags().Set("ends-at", "")
		_ = cyclesCreateCmd.Flags().Set("name", "")
	}()

	var buf bytes.Buffer
	cyclesCreateCmd.SetOut(&buf)

	if err := cyclesCreateCmd.RunE(cyclesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "создан") {
		t.Errorf("expected 'создан' in output, got: %s", out)
	}
}

func TestCyclesCreateNoToken(t *testing.T) {
	token = ""
	_ = cyclesCreateCmd.Flags().Set("team", "t1")
	_ = cyclesCreateCmd.Flags().Set("starts-at", "2024-01-01T00:00:00Z")
	_ = cyclesCreateCmd.Flags().Set("ends-at", "2024-01-14T00:00:00Z")
	defer func() {
		_ = cyclesCreateCmd.Flags().Set("team", "")
		_ = cyclesCreateCmd.Flags().Set("starts-at", "")
		_ = cyclesCreateCmd.Flags().Set("ends-at", "")
	}()

	err := cyclesCreateCmd.RunE(cyclesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestCyclesCreateNoTeam(t *testing.T) {
	token = "test-token"
	defer func() { token = "" }()

	_ = cyclesCreateCmd.Flags().Set("team", "")
	_ = cyclesCreateCmd.Flags().Set("starts-at", "2024-01-01T00:00:00Z")
	_ = cyclesCreateCmd.Flags().Set("ends-at", "2024-01-14T00:00:00Z")
	defer func() {
		_ = cyclesCreateCmd.Flags().Set("starts-at", "")
		_ = cyclesCreateCmd.Flags().Set("ends-at", "")
	}()

	err := cyclesCreateCmd.RunE(cyclesCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no team")
	}
}

func TestCyclesCreateJSON(t *testing.T) {
	responses := map[string]any{
		"CreateCycle": map[string]any{
			"cycleCreate": map[string]any{
				"success": true,
				"cycle": map[string]any{
					"id":          "cycle1",
					"number":      1,
					"name":        "Sprint 1",
					"startsAt":    "2024-01-01T00:00:00Z",
					"endsAt":      "2024-01-14T00:00:00Z",
					"completedAt": nil,
					"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = cyclesCreateCmd.Flags().Set("team", "t1")
	_ = cyclesCreateCmd.Flags().Set("starts-at", "2024-01-01T00:00:00Z")
	_ = cyclesCreateCmd.Flags().Set("ends-at", "2024-01-14T00:00:00Z")
	defer func() {
		_ = cyclesCreateCmd.Flags().Set("team", "")
		_ = cyclesCreateCmd.Flags().Set("starts-at", "")
		_ = cyclesCreateCmd.Flags().Set("ends-at", "")
	}()

	var buf bytes.Buffer
	cyclesCreateCmd.SetOut(&buf)

	if err := cyclesCreateCmd.RunE(cyclesCreateCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestCyclesUpdateSuccess(t *testing.T) {
	responses := map[string]any{
		"UpdateCycle": map[string]any{
			"cycleUpdate": map[string]any{
				"success": true,
				"cycle": map[string]any{
					"id":          "cycle1",
					"number":      1,
					"name":        "Sprint 1 Updated",
					"startsAt":    "2024-01-01T00:00:00Z",
					"endsAt":      "2024-01-14T00:00:00Z",
					"completedAt": nil,
					"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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

	_ = cyclesUpdateCmd.Flags().Set("name", "Sprint 1 Updated")
	defer func() { _ = cyclesUpdateCmd.Flags().Set("name", "") }()

	var buf bytes.Buffer
	cyclesUpdateCmd.SetOut(&buf)

	if err := cyclesUpdateCmd.RunE(cyclesUpdateCmd, []string{"cycle1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "обновлён") {
		t.Errorf("expected 'обновлён' in output, got: %s", out)
	}
}

func TestCyclesUpdateNoToken(t *testing.T) {
	token = ""
	err := cyclesUpdateCmd.RunE(cyclesUpdateCmd, []string{"cycle1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestCyclesUpdateJSON(t *testing.T) {
	responses := map[string]any{
		"UpdateCycle": map[string]any{
			"cycleUpdate": map[string]any{
				"success": true,
				"cycle": map[string]any{
					"id":          "cycle1",
					"number":      1,
					"name":        "Sprint 1 Updated",
					"startsAt":    "2024-01-01T00:00:00Z",
					"endsAt":      "2024-01-14T00:00:00Z",
					"completedAt": nil,
					"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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
	cyclesUpdateCmd.SetOut(&buf)

	if err := cyclesUpdateCmd.RunE(cyclesUpdateCmd, []string{"cycle1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"id"`) {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestCyclesArchiveSuccess(t *testing.T) {
	responses := map[string]any{
		"ArchiveCycle": map[string]any{
			"cycleArchive": map[string]any{
				"success": true,
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
	cyclesArchiveCmd.SetOut(&buf)

	if err := cyclesArchiveCmd.RunE(cyclesArchiveCmd, []string{"cycle1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "архивирован") {
		t.Errorf("expected 'архивирован' in output, got: %s", out)
	}
}

func TestCyclesArchiveNoToken(t *testing.T) {
	token = ""
	err := cyclesArchiveCmd.RunE(cyclesArchiveCmd, []string{"cycle1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
