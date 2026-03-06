package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestRelationsListSuccess(t *testing.T) {
	responses := []map[string]any{
		{
			"issue": map[string]any{
				"relations": map[string]any{
					"nodes": []map[string]any{
						{
							"id":   "rel1",
							"type": "blocks",
							"relatedIssue": map[string]any{
								"id":         "i2",
								"identifier": "ENG-2",
								"title":      "Blocked issue",
							},
						},
					},
				},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := relationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "rel1") {
		t.Errorf("expected rel1 in output, got: %s", out)
	}
	if !strings.Contains(out, "ENG-2") {
		t.Errorf("expected ENG-2 in output, got: %s", out)
	}
	if !strings.Contains(out, "blocks") {
		t.Errorf("expected blocks in output, got: %s", out)
	}
}

func TestRelationsListEmpty(t *testing.T) {
	responses := []map[string]any{
		{
			"issue": map[string]any{
				"relations": map[string]any{
					"nodes": []map[string]any{},
				},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := relationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "не найдены") {
		t.Errorf("expected 'не найдены' in output, got: %s", out)
	}
}

func TestRelationsListJSON(t *testing.T) {
	responses := []map[string]any{
		{
			"issue": map[string]any{
				"relations": map[string]any{
					"nodes": []map[string]any{
						{
							"id":   "rel1",
							"type": "related",
							"relatedIssue": map[string]any{
								"id":         "i3",
								"identifier": "ENG-3",
								"title":      "Related issue",
							},
						},
					},
				},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	outputFormat = "json"
	defer func() {
		token = ""
		outputFormat = ""
	}()

	cmd := relationsListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"ENG-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, `"rel1"`) {
		t.Errorf("expected json with rel1, got: %s", out)
	}
}

func TestRelationsListNoToken(t *testing.T) {
	token = ""
	err := relationsListCmd.RunE(relationsListCmd, []string{"ENG-1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestRelationsAddSuccess(t *testing.T) {
	responses := []map[string]any{
		{
			"issueRelationCreate": map[string]any{
				"success": true,
				"issueRelation": map[string]any{
					"id":   "rel-new",
					"type": "blocks",
					"relatedIssue": map[string]any{
						"id":         "i2",
						"identifier": "ENG-2",
						"title":      "Blocked",
					},
				},
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	relationIssueID = "i1"
	relationRelatedIssueID = "i2"
	relationRelType = "blocks"

	cmd := relationsAddCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "rel-new") {
		t.Errorf("expected rel-new in output, got: %s", out)
	}
}

func TestRelationsAddNoToken(t *testing.T) {
	token = ""
	err := relationsAddCmd.RunE(relationsAddCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestRelationsRemoveSuccess(t *testing.T) {
	responses := []map[string]any{
		{
			"issueRelationDelete": map[string]any{
				"success": true,
			},
		},
	}

	srv := newLabelAssignTestServer(t, responses)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := relationsRemoveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"rel1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "rel1") {
		t.Errorf("expected rel1 in output, got: %s", out)
	}
}

func TestRelationsRemoveNoToken(t *testing.T) {
	token = ""
	err := relationsRemoveCmd.RunE(relationsRemoveCmd, []string{"rel1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
