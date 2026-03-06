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

func newFavoritesTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestFavoritesList(t *testing.T) {
	responseData := map[string]any{
		"favorites": map[string]any{
			"nodes": []map[string]any{
				{
					"id":   "fav1",
					"type": "issue",
					"issue": map[string]any{
						"id":         "i1",
						"identifier": "ENG-1",
						"title":      "Test issue",
					},
				},
			},
		},
	}

	srv := newFavoritesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := favoritesListCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "fav1") {
		t.Errorf("expected fav1 in output, got:\n%s", out)
	}
	if !strings.Contains(out, "issue") {
		t.Errorf("expected 'issue' in output, got:\n%s", out)
	}
}

func TestFavoritesListJSON(t *testing.T) {
	responseData := map[string]any{
		"favorites": map[string]any{
			"nodes": []map[string]any{
				{"id": "fav1", "type": "issue"},
			},
		},
	}

	srv := newFavoritesTestServer(t, responseData)
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

	cmd := favoritesListCmd
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

func TestFavoritesListEmpty(t *testing.T) {
	responseData := map[string]any{
		"favorites": map[string]any{
			"nodes": []map[string]any{},
		},
	}

	srv := newFavoritesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := favoritesListCmd
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

func TestFavoritesAdd(t *testing.T) {
	responseData := map[string]any{
		"favoriteCreate": map[string]any{
			"success": true,
			"favorite": map[string]any{
				"id":   "fav1",
				"type": "issue",
			},
		},
	}

	srv := newFavoritesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := favoritesAddCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.Flags().Set("type", "issue")
	cmd.Flags().Set("id", "i1")

	if err := cmd.RunE(cmd, []string{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "fav1") {
		t.Errorf("expected fav1 in output, got:\n%s", out)
	}
}

func TestFavoritesRemove(t *testing.T) {
	responseData := map[string]any{
		"favoriteDelete": map[string]any{
			"success": true,
		},
	}

	srv := newFavoritesTestServer(t, responseData)
	defer srv.Close()

	origFactory := newLinearClient
	newLinearClient = func(tok string) *client.Client {
		return client.NewWithURL(tok, srv.URL)
	}
	defer func() { newLinearClient = origFactory }()

	token = "test-token"
	defer func() { token = "" }()

	cmd := favoritesRemoveCmd
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	if err := cmd.RunE(cmd, []string{"fav1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "fav1") {
		t.Errorf("expected fav1 in output, got:\n%s", out)
	}
}
