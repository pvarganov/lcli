package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/pavelvarganov/lcli/internal/client"
)

func TestProjectRelationsAddSuccess(t *testing.T) {
	responses := map[string]any{
		"projectRelationCreate": map[string]any{
			"projectRelationCreate": map[string]any{
				"success": true,
				"projectRelation": map[string]any{
					"id":   "pr1",
					"type": "related",
					"project":        map[string]any{"id": "p1", "name": "Project A"},
					"relatedProject": map[string]any{"id": "p2", "name": "Project B"},
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

	_ = projectRelationsAddCmd.Flags().Set("project", "p1")
	_ = projectRelationsAddCmd.Flags().Set("related", "p2")
	_ = projectRelationsAddCmd.Flags().Set("type", "related")
	defer func() {
		_ = projectRelationsAddCmd.Flags().Set("project", "")
		_ = projectRelationsAddCmd.Flags().Set("related", "")
		_ = projectRelationsAddCmd.Flags().Set("type", "related")
	}()

	var buf bytes.Buffer
	projectRelationsAddCmd.SetOut(&buf)

	if err := projectRelationsAddCmd.RunE(projectRelationsAddCmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Связь создана") {
		t.Errorf("expected 'Связь создана' in output, got: %s", out)
	}
	if !strings.Contains(out, "pr1") {
		t.Errorf("expected 'pr1' in output, got: %s", out)
	}
}

func TestProjectRelationsAddNoToken(t *testing.T) {
	token = ""
	err := projectRelationsAddCmd.RunE(projectRelationsAddCmd, []string{})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}

func TestProjectRelationsRemoveSuccess(t *testing.T) {
	responses := map[string]any{
		"projectRelationDelete": map[string]any{
			"projectRelationDelete": map[string]any{
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
	projectRelationsRemoveCmd.SetOut(&buf)

	if err := projectRelationsRemoveCmd.RunE(projectRelationsRemoveCmd, []string{"pr1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "Связь удалена") {
		t.Errorf("expected 'Связь удалена' in output, got: %s", out)
	}
	if !strings.Contains(out, "pr1") {
		t.Errorf("expected 'pr1' in output, got: %s", out)
	}
}

func TestProjectRelationsRemoveNoToken(t *testing.T) {
	token = ""
	err := projectRelationsRemoveCmd.RunE(projectRelationsRemoveCmd, []string{"pr1"})
	if err == nil {
		t.Fatal("expected error when no token")
	}
}
