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
	issues, err := c.ListIssues(IssueFilter{Limit: 10})
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
