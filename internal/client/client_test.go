package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDo_Success(t *testing.T) {
	type response struct {
		Viewer struct {
			Name string `json:"name"`
		} `json:"viewer"`
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем заголовок авторизации
		if !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		// Проверяем тело запроса
		var req graphqlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Query == "" {
			t.Error("expected non-empty query")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"viewer": map[string]any{"name": "Test User"},
			},
		})
	}))
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	var result response
	err := c.Do(`{ viewer { name } }`, nil, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Viewer.Name != "Test User" {
		t.Errorf("expected 'Test User', got %q", result.Viewer.Name)
	}
}

func TestDo_GraphQLError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []map[string]any{
				{"message": "unauthorized"},
			},
		})
	}))
	defer srv.Close()

	c := NewWithURL("bad-token", srv.URL)
	err := c.Do(`{ viewer { name } }`, nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "unauthorized") {
		t.Errorf("expected 'unauthorized' in error, got %q", err.Error())
	}
}

func TestDo_NetworkError(t *testing.T) {
	// Сервер не запускаем — используем недоступный адрес
	c := NewWithURL("token", "http://127.0.0.1:19999")
	err := c.Do(`{ viewer { name } }`, nil, nil)
	if err == nil {
		t.Fatal("expected network error, got nil")
	}
}

func TestDo_WithVariables(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req graphqlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Variables["id"] != "issue-123" {
			t.Errorf("expected variable id='issue-123', got %v", req.Variables["id"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": nil})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	err := c.Do(`query($id: String!) { issue(id: $id) { id } }`, map[string]any{"id": "issue-123"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
