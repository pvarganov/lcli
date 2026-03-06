package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestLabelsListTableOutput(t *testing.T) {
	responseData := map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "label1", "name": "Bug", "color": "#ff0000", "description": "Bug label"},
				{"id": "label2", "name": "Feature", "color": "#00ff00", "description": ""},
			},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := labelsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"Bug", "#ff0000", "Feature", "Bug label"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestLabelsListJSONOutput(t *testing.T) {
	responseData := map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{
				{"id": "label1", "name": "Bug", "color": "#ff0000", "description": ""},
			},
		},
	}

	srv := newIssuesTestServer(t, responseData)
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

	cmd := labelsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var labels []map[string]any
	if err := json.Unmarshal([]byte(buf.String()), &labels); err != nil {
		t.Fatalf("expected valid JSON, got: %v\noutput: %s", err, buf.String())
	}
	if len(labels) != 1 {
		t.Fatalf("expected 1 label, got %d", len(labels))
	}
	if labels[0]["name"] != "Bug" {
		t.Errorf("expected Bug, got %v", labels[0]["name"])
	}
}

func TestLabelsListEmpty(t *testing.T) {
	responseData := map[string]any{
		"issueLabels": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := labelsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "не найдены") {
		t.Errorf("expected 'не найдены' in output, got: %s", buf.String())
	}
}

func TestLabelsListNoToken(t *testing.T) {
	token = ""
	err := labelsListCmd.RunE(labelsListCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestLabelsCreateOutput(t *testing.T) {
	responseData := map[string]any{
		"issueLabelCreate": map[string]any{
			"success": true,
			"issueLabel": map[string]any{
				"id":          "new-label",
				"name":        "My Label",
				"color":       "#aabbcc",
				"description": "",
			},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	labelName = "My Label"
	labelColor = "#aabbcc"
	labelDescription = ""
	labelTeam = ""

	cmd := labelsCreateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "My Label") {
		t.Errorf("expected 'My Label' in output, got: %s", out)
	}
	if !strings.Contains(out, "new-label") {
		t.Errorf("expected 'new-label' in output, got: %s", out)
	}
}

func TestLabelsCreateNoToken(t *testing.T) {
	token = ""
	labelName = "Test"
	err := labelsCreateCmd.RunE(labelsCreateCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestLabelsUpdateOutput(t *testing.T) {
	responseData := map[string]any{
		"issueLabelUpdate": map[string]any{
			"success": true,
			"issueLabel": map[string]any{
				"id":          "label1",
				"name":        "Updated Label",
				"color":       "#112233",
				"description": "",
			},
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	labelName = "Updated Label"
	labelColor = "#112233"
	labelDescription = ""

	cmd := labelsUpdateCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"label1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "Updated Label") {
		t.Errorf("expected 'Updated Label' in output, got: %s", out)
	}
}

func TestLabelsUpdateNoToken(t *testing.T) {
	token = ""
	err := labelsUpdateCmd.RunE(labelsUpdateCmd, []string{"label1"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}

func TestLabelsDeleteOutput(t *testing.T) {
	responseData := map[string]any{
		"issueLabelDelete": map[string]any{
			"success": true,
		},
	}

	srv := newIssuesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := labelsDeleteCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"label1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "label1") {
		t.Errorf("expected 'label1' in output, got: %s", out)
	}
}

func TestLabelsDeleteNoToken(t *testing.T) {
	token = ""
	err := labelsDeleteCmd.RunE(labelsDeleteCmd, []string{"label1"})
	if err == nil {
		t.Fatal("expected error when no token, got nil")
	}
}
