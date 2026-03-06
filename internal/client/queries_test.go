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

func TestSearchIssues(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issueSearch": map[string]any{
						"nodes": []map[string]any{
							{
								"id":         "i1",
								"identifier": "ENG-1",
								"title":      "Search result issue",
								"updatedAt":  now.Format(time.RFC3339),
								"priority":   2,
								"state":      map[string]any{"name": "In Progress", "type": "started"},
								"assignee":   nil,
								"team":       map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							},
						},
						"pageInfo": map[string]any{
							"hasNextPage": false,
							"endCursor":   "",
						},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		issues, pageInfo, err := c.SearchIssues("search result", 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Fatalf("expected 1 issue, got %d", len(issues))
		}
		if issues[0].Identifier != "ENG-1" {
			t.Errorf("expected ENG-1, got %q", issues[0].Identifier)
		}
		if pageInfo == nil {
			t.Error("expected pageInfo to be non-nil")
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"issueSearch": map[string]any{
						"nodes":    []map[string]any{},
						"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		issues, _, err := c.SearchIssues("nothing", 10)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 0 {
			t.Fatalf("expected 0 issues, got %d", len(issues))
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

func TestGetProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"project": map[string]any{
				"id":          "proj1",
				"name":        "My Project",
				"description": "Project description",
				"state":       "started",
				"startDate":   "2024-01-01",
				"targetDate":  "2024-06-01",
				"url":         "https://linear.app/team/project/my-project",
				"lead":        map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
			},
		}

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": responseData})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		project, err := c.GetProject("proj1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if project.ID != "proj1" {
			t.Errorf("expected id proj1, got %q", project.ID)
		}
		if project.Name != "My Project" {
			t.Errorf("expected name 'My Project', got %q", project.Name)
		}
		if project.StartDate != "2024-01-01" {
			t.Errorf("expected startDate '2024-01-01', got %q", project.StartDate)
		}
		if project.Lead == nil || project.Lead.Name != "Alice" {
			t.Errorf("expected lead Alice")
		}
	})

	t.Run("not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"project": nil}})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		_, err := c.GetProject("nonexistent")
		if err == nil {
			t.Fatal("expected error for not found project")
		}
	})
}

func TestListProjectMilestones(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"project": map[string]any{
						"projectMilestones": map[string]any{
							"nodes": []map[string]any{
								{"id": "ms1", "name": "Alpha Release", "targetDate": "2024-03-01", "description": "First milestone"},
							},
						},
					},
				},
			})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		milestones, err := c.ListProjectMilestones("proj1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(milestones) != 1 {
			t.Fatalf("expected 1 milestone, got %d", len(milestones))
		}
		if milestones[0].Name != "Alpha Release" {
			t.Errorf("expected name 'Alpha Release', got %q", milestones[0].Name)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"project": map[string]any{
						"projectMilestones": map[string]any{"nodes": []any{}},
					},
				},
			})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		milestones, err := c.ListProjectMilestones("proj1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(milestones) != 0 {
			t.Errorf("expected 0 milestones, got %d", len(milestones))
		}
	})

	t.Run("project not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"project": nil}})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		_, err := c.ListProjectMilestones("nonexistent")
		if err == nil {
			t.Fatal("expected error for not found project")
		}
	})
}

func TestListProjectLabels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"projectLabels": map[string]any{
						"nodes": []map[string]any{
							{"id": "pl1", "name": "Frontend", "color": "#ff0000", "description": "Frontend work"},
							{"id": "pl2", "name": "Backend", "color": "#0000ff", "description": ""},
						},
					},
				},
			})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		labels, err := c.ListProjectLabels()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(labels) != 2 {
			t.Fatalf("expected 2 labels, got %d", len(labels))
		}
		if labels[0].ID != "pl1" {
			t.Errorf("expected id pl1, got %q", labels[0].ID)
		}
		if labels[0].Name != "Frontend" {
			t.Errorf("expected name Frontend, got %q", labels[0].Name)
		}
		if labels[0].Color != "#ff0000" {
			t.Errorf("expected color #ff0000, got %q", labels[0].Color)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"projectLabels": map[string]any{"nodes": []map[string]any{}},
				},
			})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		labels, err := c.ListProjectLabels()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(labels) != 0 {
			t.Errorf("expected 0 labels, got %d", len(labels))
		}
	})
}

func TestListProjectUpdates(t *testing.T) {
	now := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"project": map[string]any{
						"projectUpdates": map[string]any{
							"nodes": []map[string]any{
								{
									"id":        "pu1",
									"body":      "Week 1 progress",
									"health":    "onTrack",
									"createdAt": now.Format(time.RFC3339),
									"user":      map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
								},
							},
						},
					},
				},
			})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		updates, err := c.ListProjectUpdates("proj1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(updates) != 1 {
			t.Fatalf("expected 1 update, got %d", len(updates))
		}
		if updates[0].ID != "pu1" {
			t.Errorf("expected id pu1, got %q", updates[0].ID)
		}
		if updates[0].Health != "onTrack" {
			t.Errorf("expected health onTrack, got %q", updates[0].Health)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"project": map[string]any{
						"projectUpdates": map[string]any{"nodes": []map[string]any{}},
					},
				},
			})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		updates, err := c.ListProjectUpdates("proj1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(updates) != 0 {
			t.Errorf("expected 0 updates, got %d", len(updates))
		}
	})

	t.Run("project not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"project": nil}})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		_, err := c.ListProjectUpdates("nonexistent")
		if err == nil {
			t.Fatal("expected error for not found project")
		}
	})
}

func TestListProjectStatuses(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"projectStatuses": map[string]any{
				"nodes": []map[string]any{
					{"id": "ps1", "name": "Planned", "type": "planned", "color": "#0000ff", "description": "Project is planned", "position": 1.0},
					{"id": "ps2", "name": "In Progress", "type": "started", "color": "#00ff00", "description": "", "position": 2.0},
				},
			},
		}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": responseData})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		statuses, err := c.ListProjectStatuses()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(statuses) != 2 {
			t.Fatalf("expected 2 statuses, got %d", len(statuses))
		}
		if statuses[0].Name != "Planned" {
			t.Errorf("expected name 'Planned', got %q", statuses[0].Name)
		}
		if statuses[0].Type != "planned" {
			t.Errorf("expected type 'planned', got %q", statuses[0].Type)
		}
	})

	t.Run("empty", func(t *testing.T) {
		responseData := map[string]any{
			"projectStatuses": map[string]any{
				"nodes": []map[string]any{},
			},
		}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": responseData})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		statuses, err := c.ListProjectStatuses()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(statuses) != 0 {
			t.Errorf("expected 0 statuses, got %d", len(statuses))
		}
	})
}

func TestSearchProjects(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"projects": map[string]any{
				"nodes": []map[string]any{
					{
						"id":          "p1",
						"name":        "Alpha Project",
						"description": "Test project",
						"state":       "started",
						"startDate":   "",
						"targetDate":  "",
						"url":         "https://linear.app/p1",
						"lead":        nil,
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
		projects, err := c.SearchProjects("alpha")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(projects) != 1 {
			t.Fatalf("expected 1 project, got %d", len(projects))
		}
		if projects[0].Name != "Alpha Project" {
			t.Errorf("expected 'Alpha Project', got %q", projects[0].Name)
		}
	})

	t.Run("empty", func(t *testing.T) {
		responseData := map[string]any{
			"projects": map[string]any{
				"nodes": []map[string]any{},
			},
		}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": responseData})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		projects, err := c.SearchProjects("nonexistent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(projects) != 0 {
			t.Errorf("expected 0 projects, got %d", len(projects))
		}
	})
}

func TestListCycles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"cycles": map[string]any{
						"nodes": []map[string]any{
							{
								"id":          "cycle1",
								"number":      1,
								"name":        "Sprint 1",
								"startsAt":    now.Format(time.RFC3339),
								"endsAt":      end.Format(time.RFC3339),
								"completedAt": nil,
								"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							},
						},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		cycles, err := c.ListCycles("t1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cycles) != 1 {
			t.Fatalf("expected 1 cycle, got %d", len(cycles))
		}
		if cycles[0].Name != "Sprint 1" {
			t.Errorf("expected 'Sprint 1', got %q", cycles[0].Name)
		}
		if cycles[0].Number != 1 {
			t.Errorf("expected number 1, got %d", cycles[0].Number)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"cycles": map[string]any{
						"nodes": []map[string]any{},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		cycles, err := c.ListCycles("t1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(cycles) != 0 {
			t.Fatalf("expected 0 cycles, got %d", len(cycles))
		}
	})
}

func TestGetCycle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		end := time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC)
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"cycle": map[string]any{
						"id":          "cycle1",
						"number":      1,
						"name":        "Sprint 1",
						"startsAt":    now.Format(time.RFC3339),
						"endsAt":      end.Format(time.RFC3339),
						"completedAt": nil,
						"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		cycle, err := c.GetCycle("cycle1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cycle.ID != "cycle1" {
			t.Errorf("expected cycle1, got %q", cycle.ID)
		}
		if cycle.Team.Key != "ENG" {
			t.Errorf("expected ENG, got %q", cycle.Team.Key)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"cycle": nil,
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		_, err := c.GetCycle("nonexistent")
		if err == nil {
			t.Fatal("expected error for not found cycle")
		}
	})
}

func TestListWorkflowStates(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"workflowStates": map[string]any{
						"nodes": []map[string]any{
							{
								"id":    "ws1",
								"name":  "Todo",
								"type":  "unstarted",
								"color": "#e2e2e2",
								"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							},
							{
								"id":    "ws2",
								"name":  "In Progress",
								"type":  "started",
								"color": "#f2c94c",
								"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							},
						},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		states, err := c.ListWorkflowStates("t1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(states) != 2 {
			t.Fatalf("expected 2 states, got %d", len(states))
		}
		if states[0].Name != "Todo" {
			t.Errorf("expected 'Todo', got %q", states[0].Name)
		}
		if states[1].Type != "started" {
			t.Errorf("expected 'started', got %q", states[1].Type)
		}
		if states[0].Team.Key != "ENG" {
			t.Errorf("expected 'ENG', got %q", states[0].Team.Key)
		}
	})

	t.Run("empty", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{
					"workflowStates": map[string]any{
						"nodes": []map[string]any{},
					},
				},
			})
		}))
		defer srv.Close()
		c := NewWithURL("token", srv.URL)
		states, err := c.ListWorkflowStates("t1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(states) != 0 {
			t.Fatalf("expected 0 states, got %d", len(states))
		}
	})
}

func TestListTeamMembers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"team": map[string]any{
				"members": map[string]any{
					"nodes": []map[string]any{
						{
							"id":   "tm1",
							"user": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
							"team": map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
							"role": "member",
						},
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
		members, err := c.ListTeamMembers("t1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(members) != 1 {
			t.Fatalf("expected 1 member, got %d", len(members))
		}
		if members[0].User.Name != "Alice" {
			t.Errorf("expected Alice, got %s", members[0].User.Name)
		}
		if members[0].Role != "member" {
			t.Errorf("expected role 'member', got %s", members[0].Role)
		}
	})

	t.Run("empty", func(t *testing.T) {
		responseData := map[string]any{
			"team": map[string]any{
				"members": map[string]any{
					"nodes": []map[string]any{},
				},
			},
		}
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{"data": responseData})
		}))
		defer srv.Close()

		c := NewWithURL("token", srv.URL)
		members, err := c.ListTeamMembers("t1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(members) != 0 {
			t.Fatalf("expected 0 members, got %d", len(members))
		}
	})
}

func TestListUsers(t *testing.T) {
	responseData := map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@example.com"},
				{"id": "u2", "name": "Bob", "displayName": "Bob Jones", "email": "bob@example.com"},
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	users, err := c.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].DisplayName != "Alice Smith" {
		t.Errorf("expected Alice Smith, got %q", users[0].DisplayName)
	}
}

func TestListUsersEmpty(t *testing.T) {
	responseData := map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	users, err := c.ListUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("expected 0 users, got %d", len(users))
	}
}

func TestGetUser(t *testing.T) {
	responseData := map[string]any{
		"user": map[string]any{
			"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@example.com",
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	user, err := c.GetUser("u1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %q", user.Email)
	}
}

func TestGetUserNotFound(t *testing.T) {
	responseData := map[string]any{
		"user": nil,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	_, err := c.GetUser("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found user, got nil")
	}
}

func TestGetViewer(t *testing.T) {
	responseData := map[string]any{
		"viewer": map[string]any{
			"id": "v1", "name": "Me", "displayName": "Current User", "email": "me@example.com",
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	viewer, err := c.GetViewer()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if viewer.DisplayName != "Current User" {
		t.Errorf("expected Current User, got %q", viewer.DisplayName)
	}
}

func TestListNotifications(t *testing.T) {
	now := time.Date(2024, 5, 10, 8, 0, 0, 0, time.UTC)
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	notifications, pageInfo, err := c.ListNotifications(10, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifications) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(notifications))
	}
	if notifications[0].Type != "issueAssignedToYou" {
		t.Errorf("expected type issueAssignedToYou, got %q", notifications[0].Type)
	}
	if notifications[0].Issue == nil || notifications[0].Issue.Identifier != "ENG-1" {
		t.Error("expected issue with identifier ENG-1")
	}
	if pageInfo == nil {
		t.Error("expected pageInfo to be non-nil")
	}
}

func TestListNotificationsEmpty(t *testing.T) {
	responseData := map[string]any{
		"notifications": map[string]any{
			"nodes":    []map[string]any{},
			"pageInfo": map[string]any{"hasNextPage": false, "endCursor": ""},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	notifications, _, err := c.ListNotifications(10, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notifications) != 0 {
		t.Fatalf("expected 0 notifications, got %d", len(notifications))
	}
}

func TestGetNotificationsUnreadCount(t *testing.T) {
	responseData := map[string]any{
		"notifications": map[string]any{
			"nodes": []map[string]any{
				{"id": "n1"},
				{"id": "n2"},
				{"id": "n3"},
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	count, err := c.GetNotificationsUnreadCount()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 unread, got %d", count)
	}
}

func TestListWebhooks(t *testing.T) {
	responseData := map[string]any{
		"webhooks": map[string]any{
			"nodes": []map[string]any{
				{
					"id":            "wh1",
					"url":           "https://example.com/hook",
					"enabled":       true,
					"secret":        "mysecret",
					"resourceTypes": []string{"Issue", "Comment"},
					"team":          map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
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
	webhooks, err := c.ListWebhooks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(webhooks) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(webhooks))
	}
	if webhooks[0].ID != "wh1" {
		t.Errorf("expected id wh1, got %s", webhooks[0].ID)
	}
	if webhooks[0].URL != "https://example.com/hook" {
		t.Errorf("expected url https://example.com/hook, got %s", webhooks[0].URL)
	}
}

func TestListWebhooksEmpty(t *testing.T) {
	responseData := map[string]any{
		"webhooks": map[string]any{
			"nodes": []map[string]any{},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	webhooks, err := c.ListWebhooks()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(webhooks) != 0 {
		t.Fatalf("expected 0 webhooks, got %d", len(webhooks))
	}
}

func TestGetWebhook(t *testing.T) {
	responseData := map[string]any{
		"webhook": map[string]any{
			"id":            "wh1",
			"url":           "https://example.com/hook",
			"enabled":       true,
			"secret":        "mysecret",
			"resourceTypes": []string{"Issue"},
			"team":          map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	wh, err := c.GetWebhook("wh1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if wh.ID != "wh1" {
		t.Errorf("expected id wh1, got %s", wh.ID)
	}
}

func TestGetWebhookNotFound(t *testing.T) {
	responseData := map[string]any{
		"webhook": nil,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	_, err := c.GetWebhook("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found webhook")
	}
}

func TestListAttachments(t *testing.T) {
	responseData := map[string]any{
		"issue": map[string]any{
			"attachments": map[string]any{
				"nodes": []map[string]any{
					{
						"id":         "att1",
						"title":      "GitHub PR #42",
						"url":        "https://github.com/org/repo/pull/42",
						"sourceType": "github",
						"subtitle":   "Open",
					},
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
	attachments, err := c.ListAttachments("issue1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(attachments) != 1 {
		t.Fatalf("expected 1 attachment, got %d", len(attachments))
	}
	if attachments[0].ID != "att1" {
		t.Errorf("expected id att1, got %q", attachments[0].ID)
	}
	if attachments[0].Title != "GitHub PR #42" {
		t.Errorf("expected title 'GitHub PR #42', got %q", attachments[0].Title)
	}
}

func TestListAttachmentsEmpty(t *testing.T) {
	responseData := map[string]any{
		"issue": map[string]any{
			"attachments": map[string]any{
				"nodes": []map[string]any{},
			},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	attachments, err := c.ListAttachments("issue1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(attachments) != 0 {
		t.Errorf("expected 0 attachments, got %d", len(attachments))
	}
}

func TestListDocuments(t *testing.T) {
	responseData := map[string]any{
		"documents": map[string]any{
			"nodes": []map[string]any{
				{
					"id":      "doc1",
					"title":   "Getting Started",
					"content": "# Hello",
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
	docs, err := c.ListDocuments()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}
	if docs[0].ID != "doc1" {
		t.Errorf("expected id doc1, got %q", docs[0].ID)
	}
	if docs[0].Title != "Getting Started" {
		t.Errorf("expected title 'Getting Started', got %q", docs[0].Title)
	}
}

func TestListDocumentsEmpty(t *testing.T) {
	responseData := map[string]any{
		"documents": map[string]any{
			"nodes": []map[string]any{},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	docs, err := c.ListDocuments()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 0 {
		t.Errorf("expected 0 documents, got %d", len(docs))
	}
}

func TestGetDocument(t *testing.T) {
	responseData := map[string]any{
		"document": map[string]any{
			"id":      "doc1",
			"title":   "Getting Started",
			"content": "# Hello",
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	doc, err := c.GetDocument("doc1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doc.ID != "doc1" {
		t.Errorf("expected id doc1, got %q", doc.ID)
	}
}

func TestGetDocumentNotFound(t *testing.T) {
	responseData := map[string]any{
		"document": nil,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"data": responseData})
	}))
	defer srv.Close()

	c := NewWithURL("token", srv.URL)
	_, err := c.GetDocument("nonexistent")
	if err == nil {
		t.Fatal("expected error for not found document")
	}
}

func TestSearchDocuments(t *testing.T) {
	responseData := map[string]any{
		"searchDocuments": map[string]any{
			"nodes": []map[string]any{
				{
					"id":    "doc1",
					"title": "Getting Started",
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
	docs, err := c.SearchDocuments("getting")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}
	if docs[0].ID != "doc1" {
		t.Errorf("expected id doc1, got %q", docs[0].ID)
	}
}
