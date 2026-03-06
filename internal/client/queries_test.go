package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestListIssues(t *testing.T) {
	now := time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC)
	responseData := map[string]any{
		"issues": map[string]any{
			"nodes": []map[string]any{
				{
					"id":         "i1",
					"identifier": "ENG-1",
					"title":      "Test issue",
					"updatedAt":  now.Format(time.RFC3339),
					"priority":   1,
					"state":      map[string]any{"name": "In Progress", "type": "started"},
					"assignee":   map[string]any{"id": "u1", "name": "Bob", "displayName": "Bob Smith", "email": "bob@test.com"},
					"team":       map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	issues, pageInfo, err := c.ListIssues(IssueFilter{Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Identifier != "ENG-1" {
		t.Errorf("expected identifier ENG-1, got %q", issues[0].Identifier)
	}
	if issues[0].Assignee == nil || issues[0].Assignee.DisplayName != "Bob Smith" {
		t.Errorf("expected assignee Bob Smith")
	}
	if pageInfo == nil {
		t.Error("expected pageInfo to be non-nil")
	}
}

func TestGetIssue(t *testing.T) {
	now := time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC)
	responseData := map[string]any{
		"issue": map[string]any{
			"id":          "i2",
			"identifier":  "ENG-2",
			"title":       "Another issue",
			"description": "Some description",
			"updatedAt":   now.Format(time.RFC3339),
			"priority":    4,
			"state":       map[string]any{"name": "Done", "type": "completed"},
			"assignee":    nil,
			"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	issue, err := c.GetIssue("ENG-2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Identifier != "ENG-2" {
		t.Errorf("expected ENG-2, got %q", issue.Identifier)
	}
	if issue.Description != "Some description" {
		t.Errorf("expected description, got %q", issue.Description)
	}
	if issue.Assignee != nil {
		t.Error("expected nil assignee")
	}
}

func TestListIssueLabels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issueLabels": map[string]any{
						"nodes": []map[string]any{
							{"id": "label1", "name": "Bug", "color": "#ff0000", "description": "A bug"},
							{"id": "label2", "name": "Feature", "color": "#00ff00", "description": ""},
						},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		labels, err := c.ListIssueLabels()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(labels) != 2 {
			t.Fatalf("expected 2 labels, got %d", len(labels))
		}
		if labels[0].Name != "Bug" {
			t.Errorf("expected Bug, got %q", labels[0].Name)
		}
		if labels[0].Color != "#ff0000" {
			t.Errorf("expected #ff0000, got %q", labels[0].Color)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issueLabels": map[string]any{
						"nodes": []map[string]any{},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		labels, err := c.ListIssueLabels()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(labels) != 0 {
			t.Fatalf("expected 0 labels, got %d", len(labels))
		}
	})
}

func TestGetIssueLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issueLabel": map[string]any{
						"id": "label1", "name": "Bug", "color": "#ff0000", "description": "A bug",
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		label, err := c.GetIssueLabel("label1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if label.ID != "label1" {
			t.Errorf("expected label1, got %q", label.ID)
		}
		if label.Name != "Bug" {
			t.Errorf("expected Bug, got %q", label.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issueLabel": nil,
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		_, err := c.GetIssueLabel("nonexistent")
		if err == nil {
			t.Fatal("expected error for not found label")
		}
	})
}

func TestListIssueRelations(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
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
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		rels, err := c.ListIssueRelations("i1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(rels) != 1 {
			t.Fatalf("expected 1 relation, got %d", len(rels))
		}
		if rels[0].ID != "rel1" {
			t.Errorf("expected rel1, got %s", rels[0].ID)
		}
		if rels[0].Type != "blocks" {
			t.Errorf("expected blocks, got %s", rels[0].Type)
		}
		if rels[0].RelatedIssue.Identifier != "ENG-2" {
			t.Errorf("expected ENG-2, got %s", rels[0].RelatedIssue.Identifier)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issue": map[string]any{
						"relations": map[string]any{
							"nodes": []map[string]any{},
						},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		rels, err := c.ListIssueRelations("i1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(rels) != 0 {
			t.Errorf("expected 0 relations, got %d", len(rels))
		}
	})

	t.Run("issue not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issue": nil,
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		_, err := c.ListIssueRelations("nonexistent")
		if err == nil {
			t.Fatal("expected error for nil issue")
		}
	})
}

func TestPriorityLabel(t *testing.T) {
	cases := []struct {
		p    int
		want string
	}{
		{0, "No priority"},
		{1, "Urgent"},
		{2, "High"},
		{3, "Medium"},
		{4, "Low"},
		{99, "No priority"},
	}
	for _, tc := range cases {
		if got := PriorityLabel(tc.p); got != tc.want {
			t.Errorf("PriorityLabel(%d) = %q, want %q", tc.p, got, tc.want)
		}
	}
}
