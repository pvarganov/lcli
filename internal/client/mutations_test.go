package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newMutationTestServer(t *testing.T, responseData any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body, _ := json.Marshal(map[string]any{"data": responseData})
		w.Write(body)
	}))
}

func TestCreateIssue(t *testing.T) {
	now := time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC)
	responseData := map[string]any{
		"issueCreate": map[string]any{
			"success": true,
			"issue": map[string]any{
				"id":          "new-issue-id",
				"identifier":  "ENG-100",
				"title":       "New test issue",
				"description": "Test description",
				"updatedAt":   now.Format(time.RFC3339),
				"priority":    2,
				"state":       map[string]any{"name": "Todo", "type": "unstarted"},
				"assignee":    nil,
				"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	issue, err := c.CreateIssue(CreateIssueInput{
		TeamID:      "t1",
		Title:       "New test issue",
		Description: "Test description",
		Priority:    2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Identifier != "ENG-100" {
		t.Errorf("expected identifier ENG-100, got %s", issue.Identifier)
	}
	if issue.Title != "New test issue" {
		t.Errorf("expected title 'New test issue', got %s", issue.Title)
	}
}

func TestUpdateIssue(t *testing.T) {
	now := time.Date(2024, 2, 1, 12, 0, 0, 0, time.UTC)
	responseData := map[string]any{
		"issueUpdate": map[string]any{
			"success": true,
			"issue": map[string]any{
				"id":          "existing-issue-id",
				"identifier":  "ENG-42",
				"title":       "Updated title",
				"description": "",
				"updatedAt":   now.Format(time.RFC3339),
				"priority":    3,
				"state":       map[string]any{"name": "In Progress", "type": "started"},
				"assignee":    nil,
				"team":        map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	prio := 3
	issue, err := c.UpdateIssue("ENG-42", UpdateIssueInput{
		Title:    "Updated title",
		Priority: &prio,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if issue.Identifier != "ENG-42" {
		t.Errorf("expected identifier ENG-42, got %s", issue.Identifier)
	}
	if issue.Title != "Updated title" {
		t.Errorf("expected title 'Updated title', got %s", issue.Title)
	}
}

func TestGetTeamByKey(t *testing.T) {
	responseData := map[string]any{
		"teams": map[string]any{
			"nodes": []map[string]any{
				{"id": "t1", "key": "ENG", "name": "Engineering"},
				{"id": "t2", "key": "PROD", "name": "Product"},
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	team, err := c.GetTeamByKey("ENG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team == nil {
		t.Fatal("expected team, got nil")
	}
	if team.ID != "t1" {
		t.Errorf("expected team ID t1, got %s", team.ID)
	}

	missing, err := c.GetTeamByKey("MISSING")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if missing != nil {
		t.Errorf("expected nil for missing team, got %+v", missing)
	}
}

func TestFindWorkflowStateByName(t *testing.T) {
	responseData := map[string]any{
		"workflowStates": map[string]any{
			"nodes": []map[string]any{
				{"id": "s1", "name": "Todo"},
				{"id": "s2", "name": "In Progress"},
				{"id": "s3", "name": "Done"},
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	id, err := c.FindWorkflowStateByName("t1", "In Progress")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "s2" {
		t.Errorf("expected state ID s2, got %s", id)
	}

	missing, err := c.FindWorkflowStateByName("t1", "Unknown")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if missing != "" {
		t.Errorf("expected empty string for missing state, got %s", missing)
	}
}

func TestCreateIssueLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueLabelCreate": map[string]any{
				"success": true,
				"issueLabel": map[string]any{
					"id":          "label-new",
					"name":        "My Label",
					"color":       "#aabbcc",
					"description": "A test label",
				},
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		label, err := c.CreateIssueLabel(CreateIssueLabelInput{Name: "My Label", Color: "#aabbcc", Description: "A test label"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if label.ID != "label-new" {
			t.Errorf("expected label-new, got %s", label.ID)
		}
		if label.Name != "My Label" {
			t.Errorf("expected My Label, got %s", label.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueLabelCreate": map[string]any{
				"success":    false,
				"issueLabel": nil,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateIssueLabel(CreateIssueLabelInput{Name: "Bad", Color: "#fff"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateIssueLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueLabelUpdate": map[string]any{
				"success": true,
				"issueLabel": map[string]any{
					"id":          "label1",
					"name":        "Updated Label",
					"color":       "#112233",
					"description": "",
				},
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		label, err := c.UpdateIssueLabel("label1", UpdateIssueLabelInput{Name: "Updated Label", Color: "#112233"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if label.Name != "Updated Label" {
			t.Errorf("expected Updated Label, got %s", label.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueLabelUpdate": map[string]any{
				"success":    false,
				"issueLabel": nil,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateIssueLabel("label1", UpdateIssueLabelInput{Name: "X"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteIssueLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueLabelDelete": map[string]any{
				"success": true,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteIssueLabel("label1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueLabelDelete": map[string]any{
				"success": false,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteIssueLabel("label1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestAddLabelToIssue(t *testing.T) {
	t.Run("success adds new label", func(t *testing.T) {
		requestNum := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			requestNum++
			var data any
			if requestNum == 1 {
				data = map[string]any{
					"issue": map[string]any{
						"labels": map[string]any{
							"nodes": []map[string]any{
								{"id": "label-existing"},
							},
						},
					},
				}
			} else {
				data = map[string]any{
					"issueUpdate": map[string]any{"success": true},
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"data": data})
		}))
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.AddLabelToIssue("issue1", "label-new"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("no-op if label already present", func(t *testing.T) {
		requestNum := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			requestNum++
			data := map[string]any{
				"issue": map[string]any{
					"labels": map[string]any{
						"nodes": []map[string]any{
							{"id": "label-existing"},
						},
					},
				},
			}
			json.NewEncoder(w).Encode(map[string]any{"data": data})
		}))
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.AddLabelToIssue("issue1", "label-existing"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if requestNum != 1 {
			t.Errorf("expected 1 request (no update needed), got %d", requestNum)
		}
	})

	t.Run("issueUpdate failure", func(t *testing.T) {
		requestNum := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			requestNum++
			var data any
			if requestNum == 1 {
				data = map[string]any{
					"issue": map[string]any{
						"labels": map[string]any{"nodes": []map[string]any{}},
					},
				}
			} else {
				data = map[string]any{
					"issueUpdate": map[string]any{"success": false},
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"data": data})
		}))
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.AddLabelToIssue("issue1", "label-new"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestRemoveLabelFromIssue(t *testing.T) {
	t.Run("success removes label", func(t *testing.T) {
		requestNum := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			requestNum++
			var data any
			if requestNum == 1 {
				data = map[string]any{
					"issue": map[string]any{
						"labels": map[string]any{
							"nodes": []map[string]any{
								{"id": "label-to-remove"},
								{"id": "label-keep"},
							},
						},
					},
				}
			} else {
				data = map[string]any{
					"issueUpdate": map[string]any{"success": true},
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"data": data})
		}))
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.RemoveLabelFromIssue("issue1", "label-to-remove"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("issueUpdate failure", func(t *testing.T) {
		requestNum := 0
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			requestNum++
			var data any
			if requestNum == 1 {
				data = map[string]any{
					"issue": map[string]any{
						"labels": map[string]any{
							"nodes": []map[string]any{{"id": "label1"}},
						},
					},
				}
			} else {
				data = map[string]any{
					"issueUpdate": map[string]any{"success": false},
				}
			}
			json.NewEncoder(w).Encode(map[string]any{"data": data})
		}))
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.RemoveLabelFromIssue("issue1", "label1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateIssueRelation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueRelationCreate": map[string]any{
				"success": true,
				"issueRelation": map[string]any{
					"id":   "rel1",
					"type": "blocks",
					"relatedIssue": map[string]any{
						"id":         "i2",
						"identifier": "ENG-2",
						"title":      "Related issue",
					},
				},
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		rel, err := c.CreateIssueRelation(CreateIssueRelationInput{
			IssueID:        "i1",
			RelatedIssueID: "i2",
			Type:           "blocks",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rel.ID != "rel1" {
			t.Errorf("expected rel1, got %s", rel.ID)
		}
		if rel.Type != "blocks" {
			t.Errorf("expected blocks, got %s", rel.Type)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueRelationCreate": map[string]any{
				"success":       false,
				"issueRelation": nil,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateIssueRelation(CreateIssueRelationInput{IssueID: "i1", RelatedIssueID: "i2", Type: "related"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteIssueRelation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueRelationDelete": map[string]any{
				"success": true,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteIssueRelation("rel1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueRelationDelete": map[string]any{
				"success": false,
			},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteIssueRelation("rel1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestFindUserByName(t *testing.T) {
	responseData := map[string]any{
		"users": map[string]any{
			"nodes": []map[string]any{
				{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@example.com"},
				{"id": "u2", "name": "Bob", "displayName": "Bob Jones", "email": "bob@example.com"},
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	user, err := c.FindUserByName("Alice Smith")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != "u1" {
		t.Errorf("expected user ID u1, got %s", user.ID)
	}

	byEmail, err := c.FindUserByName("bob@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if byEmail == nil || byEmail.ID != "u2" {
		t.Errorf("expected user u2 by email, got %+v", byEmail)
	}
}

func TestArchiveIssue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueArchive": map[string]any{"success": true},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveIssue("issue1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueArchive": map[string]any{"success": false},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveIssue("issue1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUnarchiveIssue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueUnarchive": map[string]any{"success": true},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.UnarchiveIssue("issue1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueUnarchive": map[string]any{"success": false},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.UnarchiveIssue("issue1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteIssue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueDelete": map[string]any{"success": true},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteIssue("issue1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueDelete": map[string]any{"success": false},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteIssue("issue1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestSubscribeToIssue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueSubscribe": map[string]any{"success": true},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.SubscribeToIssue("issue1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueSubscribe": map[string]any{"success": false},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.SubscribeToIssue("issue1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUnsubscribeFromIssue(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueUnsubscribe": map[string]any{"success": true},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.UnsubscribeFromIssue("issue1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"issueUnsubscribe": map[string]any{"success": false},
		})
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.UnsubscribeFromIssue("issue1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestBatchCreateIssues(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"issueBatchCreate": map[string]any{
				"success": true,
				"issues": []map[string]any{
					{"id": "issue1", "identifier": "ENG-1", "title": "First", "priority": 0, "state": map[string]any{"name": "Todo", "type": "unstarted"}, "team": map[string]any{"id": "t1", "key": "ENG", "name": "Eng"}},
					{"id": "issue2", "identifier": "ENG-2", "title": "Second", "priority": 0, "state": map[string]any{"name": "Todo", "type": "unstarted"}, "team": map[string]any{"id": "t1", "key": "ENG", "name": "Eng"}},
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		issues, err := c.BatchCreateIssues(BatchCreateIssuesInput{
			Issues: []CreateIssueInput{
				{TeamID: "t1", Title: "First"},
				{TeamID: "t1", Title: "Second"},
			},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 2 {
			t.Fatalf("expected 2 issues, got %d", len(issues))
		}
		if issues[0].Identifier != "ENG-1" {
			t.Errorf("expected ENG-1, got %s", issues[0].Identifier)
		}
	})

	t.Run("success_false", func(t *testing.T) {
		responseData := map[string]any{
			"issueBatchCreate": map[string]any{"success": false, "issues": []any{}},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		_, err := c.BatchCreateIssues(BatchCreateIssuesInput{
			Issues: []CreateIssueInput{{TeamID: "t1", Title: "Test"}},
		})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestBatchUpdateIssues(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"issueBatchUpdate": map[string]any{
				"success": true,
				"issues": []map[string]any{
					{"id": "issue1", "identifier": "ENG-1", "title": "Updated", "priority": 0, "state": map[string]any{"name": "In Progress", "type": "started"}, "team": map[string]any{"id": "t1", "key": "ENG", "name": "Eng"}},
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		issues, err := c.BatchUpdateIssues(BatchUpdateIssuesInput{
			IDs:    []string{"issue1"},
			Update: UpdateIssueInput{Title: "Updated"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(issues) != 1 {
			t.Fatalf("expected 1 issue, got %d", len(issues))
		}
		if issues[0].Title != "Updated" {
			t.Errorf("expected 'Updated', got %s", issues[0].Title)
		}
	})

	t.Run("success_false", func(t *testing.T) {
		responseData := map[string]any{
			"issueBatchUpdate": map[string]any{"success": false, "issues": []any{}},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		_, err := c.BatchUpdateIssues(BatchUpdateIssuesInput{
			IDs:    []string{"issue1"},
			Update: UpdateIssueInput{Title: "Test"},
		})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"commentUpdate": map[string]any{
				"success": true,
				"comment": map[string]any{
					"id":        "cm1",
					"body":      "Updated body",
					"createdAt": "2024-01-01T00:00:00Z",
					"user":      map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice", "email": "alice@example.com"},
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		comment, err := c.UpdateComment("cm1", "Updated body")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if comment.ID != "cm1" {
			t.Errorf("expected id 'cm1', got %s", comment.ID)
		}
		if comment.Body != "Updated body" {
			t.Errorf("expected body 'Updated body', got %s", comment.Body)
		}
	})

	t.Run("success_false", func(t *testing.T) {
		responseData := map[string]any{
			"commentUpdate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateComment("cm1", "body")
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"commentDelete": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteComment("cm1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("success_false", func(t *testing.T) {
		responseData := map[string]any{
			"commentDelete": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteComment("cm1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestResolveComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"commentResolve": map[string]any{
				"success": true,
				"comment": map[string]any{
					"id":        "cm1",
					"body":      "Some comment",
					"createdAt": "2024-01-01T00:00:00Z",
					"user":      nil,
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		comment, err := c.ResolveComment("cm1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if comment.ID != "cm1" {
			t.Errorf("expected id 'cm1', got %s", comment.ID)
		}
	})

	t.Run("success_false", func(t *testing.T) {
		responseData := map[string]any{
			"commentResolve": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if _, err := c.ResolveComment("cm1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUnresolveComment(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"commentUnresolve": map[string]any{
				"success": true,
				"comment": map[string]any{
					"id":        "cm1",
					"body":      "Some comment",
					"createdAt": "2024-01-01T00:00:00Z",
					"user":      nil,
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		comment, err := c.UnresolveComment("cm1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if comment.ID != "cm1" {
			t.Errorf("expected id 'cm1', got %s", comment.ID)
		}
	})

	t.Run("success_false", func(t *testing.T) {
		responseData := map[string]any{
			"commentUnresolve": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if _, err := c.UnresolveComment("cm1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"projectCreate": map[string]any{
				"success": true,
				"project": map[string]any{
					"id":          "proj1",
					"name":        "My Project",
					"description": "A project",
					"state":       "started",
					"startDate":   "2024-01-01",
					"targetDate":  "2024-06-01",
					"url":         "https://linear.app/team/project/my-project",
					"lead":        nil,
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		project, err := c.CreateProject(CreateProjectInput{
			Name:        "My Project",
			Description: "A project",
			TeamIDs:     []string{"team1"},
			StartDate:   "2024-01-01",
			TargetDate:  "2024-06-01",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if project.Name != "My Project" {
			t.Errorf("expected name 'My Project', got %q", project.Name)
		}
		if project.StartDate != "2024-01-01" {
			t.Errorf("expected startDate '2024-01-01', got %q", project.StartDate)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectCreate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.CreateProject(CreateProjectInput{Name: "x", TeamIDs: []string{"t1"}}); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"projectUpdate": map[string]any{
				"success": true,
				"project": map[string]any{
					"id":          "proj1",
					"name":        "Updated Project",
					"description": "Updated",
					"state":       "started",
					"startDate":   "",
					"targetDate":  "",
					"url":         "",
					"lead":        nil,
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		project, err := c.UpdateProject("proj1", UpdateProjectInput{Name: "Updated Project"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if project.Name != "Updated Project" {
			t.Errorf("expected 'Updated Project', got %q", project.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.UpdateProject("proj1", UpdateProjectInput{Name: "x"}); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectDelete": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProject("proj1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectDelete": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProject("proj1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestArchiveProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectArchive": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveProject("proj1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectArchive": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveProject("proj1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUnarchiveProject(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUnarchive": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.UnarchiveProject("proj1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUnarchive": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.UnarchiveProject("proj1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateProjectMilestone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectMilestoneCreate": map[string]any{
				"success": true,
				"projectMilestone": map[string]any{
					"id": "ms1", "name": "Alpha", "targetDate": "2024-03-01", "description": "First milestone",
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		ms, err := c.CreateProjectMilestone(CreateProjectMilestoneInput{
			ProjectID:  "proj1",
			Name:       "Alpha",
			TargetDate: "2024-03-01",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ms.Name != "Alpha" {
			t.Errorf("expected name 'Alpha', got %q", ms.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectMilestoneCreate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateProjectMilestone(CreateProjectMilestoneInput{ProjectID: "proj1", Name: "Alpha"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateProjectMilestone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectMilestoneUpdate": map[string]any{
				"success": true,
				"projectMilestone": map[string]any{
					"id": "ms1", "name": "Beta", "targetDate": "2024-06-01", "description": "",
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		ms, err := c.UpdateProjectMilestone("ms1", UpdateProjectMilestoneInput{Name: "Beta"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ms.Name != "Beta" {
			t.Errorf("expected name 'Beta', got %q", ms.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectMilestoneUpdate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateProjectMilestone("ms1", UpdateProjectMilestoneInput{Name: "Beta"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteProjectMilestone(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectMilestoneDelete": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProjectMilestone("ms1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectMilestoneDelete": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProjectMilestone("ms1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateProjectUpdate(t *testing.T) {
	now := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdateCreate": map[string]any{
				"success": true,
				"projectUpdate": map[string]any{
					"id":        "pu-new",
					"body":      "Progress update",
					"health":    "onTrack",
					"createdAt": now.Format(time.RFC3339),
					"user":      nil,
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		update, err := c.CreateProjectUpdate(CreateProjectUpdateInput{
			ProjectID: "proj1",
			Body:      "Progress update",
			Health:    "onTrack",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if update.ID != "pu-new" {
			t.Errorf("expected id pu-new, got %s", update.ID)
		}
		if update.Body != "Progress update" {
			t.Errorf("expected body 'Progress update', got %s", update.Body)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdateCreate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateProjectUpdate(CreateProjectUpdateInput{ProjectID: "proj1", Body: "test"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateProjectUpdate(t *testing.T) {
	now := time.Date(2024, 6, 1, 10, 0, 0, 0, time.UTC)

	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdateUpdate": map[string]any{
				"success": true,
				"projectUpdate": map[string]any{
					"id":        "pu1",
					"body":      "Updated body",
					"health":    "atRisk",
					"createdAt": now.Format(time.RFC3339),
					"user":      nil,
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		update, err := c.UpdateProjectUpdate("pu1", UpdateProjectUpdateInput{
			Body:   "Updated body",
			Health: "atRisk",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if update.Health != "atRisk" {
			t.Errorf("expected health atRisk, got %s", update.Health)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdateUpdate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateProjectUpdate("pu1", UpdateProjectUpdateInput{Body: "test"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateProjectLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectLabelCreate": map[string]any{
				"success": true,
				"projectLabel": map[string]any{
					"id": "pl-new", "name": "Design", "color": "#00ff00", "description": "Design work",
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		label, err := c.CreateProjectLabel(CreateProjectLabelInput{
			Name:        "Design",
			Color:       "#00ff00",
			Description: "Design work",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if label.ID != "pl-new" {
			t.Errorf("expected id pl-new, got %s", label.ID)
		}
		if label.Name != "Design" {
			t.Errorf("expected name Design, got %s", label.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectLabelCreate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateProjectLabel(CreateProjectLabelInput{Name: "Test"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateProjectLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectLabelUpdate": map[string]any{
				"success": true,
				"projectLabel": map[string]any{
					"id": "pl1", "name": "Updated Label", "color": "#ff00ff", "description": "",
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		label, err := c.UpdateProjectLabel("pl1", UpdateProjectLabelInput{
			Name:  "Updated Label",
			Color: "#ff00ff",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if label.Name != "Updated Label" {
			t.Errorf("expected name 'Updated Label', got %s", label.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectLabelUpdate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateProjectLabel("pl1", UpdateProjectLabelInput{Name: "Test"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteProjectLabel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectLabelDelete": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProjectLabel("pl1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectLabelDelete": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProjectLabel("pl1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestArchiveProjectUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdateArchive": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveProjectUpdate("pu1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectUpdateArchive": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveProjectUpdate("pu1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateProjectStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectStatusCreate": map[string]any{
				"success": true,
				"projectStatus": map[string]any{
					"id": "pst1", "name": "Planned", "type": "planned", "color": "#0000ff", "description": "", "position": 1.0,
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		status, err := c.CreateProjectStatus(CreateProjectStatusInput{
			Name:     "Planned",
			Type:     "planned",
			Color:    "#0000ff",
			Position: 1.0,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status.Name != "Planned" {
			t.Errorf("expected name 'Planned', got %q", status.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectStatusCreate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateProjectStatus(CreateProjectStatusInput{Name: "X", Type: "planned", Color: "#fff", Position: 1.0})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateProjectStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectStatusUpdate": map[string]any{
				"success": true,
				"projectStatus": map[string]any{
					"id": "pst1", "name": "Updated", "type": "started", "color": "#00ff00", "description": "", "position": 1.0,
				},
			},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		status, err := c.UpdateProjectStatus("pst1", UpdateProjectStatusInput{Name: "Updated"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if status.Name != "Updated" {
			t.Errorf("expected name 'Updated', got %q", status.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectStatusUpdate": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateProjectStatus("pst1", UpdateProjectStatusInput{Name: "X"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestArchiveProjectStatus(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectStatusArchive": map[string]any{"success": true},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveProjectStatus("pst1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		srv := newMutationTestServer(t, map[string]any{
			"projectStatusArchive": map[string]any{"success": false},
		})
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveProjectStatus("pst1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateProjectRelation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"projectRelationCreate": map[string]any{
				"success": true,
				"projectRelation": map[string]any{
					"id":   "pr1",
					"type": "related",
					"project": map[string]any{
						"id": "p1", "name": "Project A",
					},
					"relatedProject": map[string]any{
						"id": "p2", "name": "Project B",
					},
				},
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		rel, err := c.CreateProjectRelation(CreateProjectRelationInput{
			ProjectID:        "p1",
			RelatedProjectID: "p2",
			Type:             "related",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if rel.ID != "pr1" {
			t.Errorf("expected id 'pr1', got %q", rel.ID)
		}
		if rel.Type != "related" {
			t.Errorf("expected type 'related', got %q", rel.Type)
		}
		if rel.Project.Name != "Project A" {
			t.Errorf("expected project name 'Project A', got %q", rel.Project.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"projectRelationCreate": map[string]any{
				"success":         false,
				"projectRelation": nil,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.CreateProjectRelation(CreateProjectRelationInput{
			ProjectID:        "p1",
			RelatedProjectID: "p2",
			Type:             "related",
		}); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteProjectRelation(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"projectRelationDelete": map[string]any{
				"success": true,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProjectRelation("pr1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"projectRelationDelete": map[string]any{
				"success": false,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteProjectRelation("pr1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateCycle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
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
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		cycle, err := c.CreateCycle("t1", "Sprint 1", "2024-01-01T00:00:00Z", "2024-01-14T00:00:00Z")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cycle.ID != "cycle1" {
			t.Errorf("expected cycle ID 'cycle1', got '%s'", cycle.ID)
		}
		if cycle.Name != "Sprint 1" {
			t.Errorf("expected cycle name 'Sprint 1', got '%s'", cycle.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"cycleCreate": map[string]any{
				"success": false,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.CreateCycle("t1", "", "2024-01-01T00:00:00Z", "2024-01-14T00:00:00Z"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateCycle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
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
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		cycle, err := c.UpdateCycle("cycle1", "Sprint 1 Updated", "", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cycle.Name != "Sprint 1 Updated" {
			t.Errorf("expected updated name, got '%s'", cycle.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"cycleUpdate": map[string]any{
				"success": false,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.UpdateCycle("cycle1", "name", "", ""); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestArchiveCycle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"cycleArchive": map[string]any{
				"success": true,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveCycle("cycle1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"cycleArchive": map[string]any{
				"success": false,
			},
		}

		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveCycle("cycle1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateWorkflowState(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"workflowStateCreate": map[string]any{
				"success": true,
				"workflowState": map[string]any{
					"id":    "ws1",
					"name":  "Review",
					"type":  "started",
					"color": "#f2c94c",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		ws, err := c.CreateWorkflowState("t1", "Review", "started", "#f2c94c")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ws.Name != "Review" {
			t.Errorf("expected 'Review', got %q", ws.Name)
		}
		if ws.Type != "started" {
			t.Errorf("expected 'started', got %q", ws.Type)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"workflowStateCreate": map[string]any{
				"success": false,
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if _, err := c.CreateWorkflowState("t1", "Review", "started", ""); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateWorkflowState(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"workflowStateUpdate": map[string]any{
				"success": true,
				"workflowState": map[string]any{
					"id":    "ws1",
					"name":  "In Review",
					"type":  "started",
					"color": "#f2c94c",
					"team":  map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		ws, err := c.UpdateWorkflowState("ws1", "In Review", "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ws.Name != "In Review" {
			t.Errorf("expected 'In Review', got %q", ws.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"workflowStateUpdate": map[string]any{
				"success": false,
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if _, err := c.UpdateWorkflowState("ws1", "name", ""); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestArchiveWorkflowState(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"workflowStateArchive": map[string]any{
				"success": true,
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveWorkflowState("ws1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"workflowStateArchive": map[string]any{
				"success": false,
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()
		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveWorkflowState("ws1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"teamCreate": map[string]any{
				"success": true,
				"team":    map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		team, err := c.CreateTeam("Engineering", "ENG")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if team.Name != "Engineering" {
			t.Errorf("expected name 'Engineering', got %s", team.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"teamCreate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.CreateTeam("Engineering", "ENG"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"teamUpdate": map[string]any{
				"success": true,
				"team":    map[string]any{"id": "t1", "key": "ENG", "name": "Engineering Updated"},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		team, err := c.UpdateTeam("t1", UpdateTeamInput{Name: "Engineering Updated"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if team.Name != "Engineering Updated" {
			t.Errorf("expected updated name, got %s", team.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"teamUpdate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.UpdateTeam("t1", UpdateTeamInput{Name: "X"}); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteTeam(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"teamDelete": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteTeam("t1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"teamDelete": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteTeam("t1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateTeamMembership(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"teamMembershipCreate": map[string]any{
				"success": true,
				"teamMembership": map[string]any{
					"id":   "tm1",
					"user": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
					"team": map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					"role": "member",
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		m, err := c.CreateTeamMembership("t1", "u1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.User.Name != "Alice" {
			t.Errorf("expected user Alice, got %s", m.User.Name)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"teamMembershipCreate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.CreateTeamMembership("t1", "u1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteTeamMembership(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"teamMembershipDelete": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteTeamMembership("tm1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"teamMembershipDelete": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteTeamMembership("tm1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateTeamMembership(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"teamMembershipUpdate": map[string]any{
				"success": true,
				"teamMembership": map[string]any{
					"id":   "tm1",
					"user": map[string]any{"id": "u1", "name": "Alice", "displayName": "Alice Smith", "email": "alice@test.com"},
					"team": map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
					"role": "admin",
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		m, err := c.UpdateTeamMembership("tm1", "admin")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if m.Role != "admin" {
			t.Errorf("expected role admin, got %s", m.Role)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"teamMembershipUpdate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if _, err := c.UpdateTeamMembership("tm1", "admin"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestMarkNotificationRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"notificationUpdate": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.MarkNotificationRead("n1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"notificationUpdate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.MarkNotificationRead("n1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestMarkAllNotificationsRead(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"notificationMarkReadAll": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.MarkAllNotificationsRead(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"notificationMarkReadAll": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.MarkAllNotificationsRead(); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestArchiveNotification(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"notificationArchive": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveNotification("n1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"notificationArchive": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.ArchiveNotification("n1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestCreateWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"webhookCreate": map[string]any{
				"success": true,
				"webhook": map[string]any{
					"id":            "wh1",
					"url":           "https://example.com/hook",
					"enabled":       true,
					"secret":        "mysecret",
					"resourceTypes": []string{"Issue"},
					"team":          map[string]any{"id": "t1", "key": "ENG", "name": "Engineering"},
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		wh, err := c.CreateWebhook(WebhookCreateInput{
			URL:           "https://example.com/hook",
			TeamID:        "t1",
			Enabled:       true,
			ResourceTypes: []string{"Issue"},
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if wh.ID != "wh1" {
			t.Errorf("expected id wh1, got %s", wh.ID)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"webhookCreate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.CreateWebhook(WebhookCreateInput{URL: "https://example.com/hook", TeamID: "t1"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestUpdateWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"webhookUpdate": map[string]any{
				"success": true,
				"webhook": map[string]any{
					"id":      "wh1",
					"url":     "https://example.com/hook2",
					"enabled": false,
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		wh, err := c.UpdateWebhook("wh1", map[string]any{"url": "https://example.com/hook2", "enabled": false})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if wh.ID != "wh1" {
			t.Errorf("expected id wh1, got %s", wh.ID)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"webhookUpdate": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.UpdateWebhook("wh1", map[string]any{"url": "x"})
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestDeleteWebhook(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"webhookDelete": map[string]any{"success": true},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteWebhook("wh1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"webhookDelete": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		if err := c.DeleteWebhook("wh1"); err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestRotateWebhookSecret(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		responseData := map[string]any{
			"webhookRotateSecret": map[string]any{
				"success": true,
				"webhook": map[string]any{
					"id":     "wh1",
					"url":    "https://example.com/hook",
					"secret": "newsecret",
				},
			},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		wh, err := c.RotateWebhookSecret("wh1")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if wh.ID != "wh1" {
			t.Errorf("expected id wh1, got %s", wh.ID)
		}
	})

	t.Run("failure", func(t *testing.T) {
		responseData := map[string]any{
			"webhookRotateSecret": map[string]any{"success": false},
		}
		srv := newMutationTestServer(t, responseData)
		defer srv.Close()

		c := NewWithURL("test-token", srv.URL)
		_, err := c.RotateWebhookSecret("wh1")
		if err == nil {
			t.Fatal("expected error when success=false")
		}
	})
}

func TestAttachmentLinkURL(t *testing.T) {
	responseData := map[string]any{
		"attachmentLinkURL": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att1",
				"title":      "Linear Docs",
				"url":        "https://linear.app/docs",
				"sourceType": "url",
				"subtitle":   "",
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	att, err := c.AttachmentLinkURL("issue1", "https://linear.app/docs", "Linear Docs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.ID != "att1" {
		t.Errorf("expected id att1, got %s", att.ID)
	}
}

func TestAttachmentLinkGitHubPR(t *testing.T) {
	responseData := map[string]any{
		"attachmentLinkGitHubPR": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att2",
				"title":      "Fix bug",
				"url":        "https://github.com/org/repo/pull/1",
				"sourceType": "github",
				"subtitle":   "Open",
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	att, err := c.AttachmentLinkGitHubPR("issue1", "https://github.com/org/repo/pull/1", "Fix bug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.ID != "att2" {
		t.Errorf("expected id att2, got %s", att.ID)
	}
}

func TestAttachmentLinkGitHubIssue(t *testing.T) {
	responseData := map[string]any{
		"attachmentLinkGitHubIssue": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att3",
				"title":      "Bug report",
				"url":        "https://github.com/org/repo/issues/5",
				"sourceType": "github",
				"subtitle":   "",
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	att, err := c.AttachmentLinkGitHubIssue("issue1", "https://github.com/org/repo/issues/5", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.ID != "att3" {
		t.Errorf("expected id att3, got %s", att.ID)
	}
}

func TestAttachmentLinkGitLabMR(t *testing.T) {
	responseData := map[string]any{
		"attachmentLinkGitLabMR": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att4",
				"title":      "MR !42",
				"url":        "https://gitlab.com/org/repo/-/merge_requests/42",
				"sourceType": "gitlab",
				"subtitle":   "",
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	att, err := c.AttachmentLinkGitLabMR("issue1", "https://gitlab.com/org/repo/-/merge_requests/42", 42, "org/repo", "MR !42")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.ID != "att4" {
		t.Errorf("expected id att4, got %s", att.ID)
	}
}

func TestAttachmentDelete(t *testing.T) {
	responseData := map[string]any{
		"attachmentDelete": map[string]any{
			"success": true,
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	err := c.AttachmentDelete("att1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAttachmentUpdate(t *testing.T) {
	responseData := map[string]any{
		"attachmentUpdate": map[string]any{
			"success": true,
			"attachment": map[string]any{
				"id":         "att1",
				"title":      "Updated title",
				"url":        "https://linear.app/docs",
				"sourceType": "url",
				"subtitle":   "New subtitle",
			},
		},
	}

	srv := newMutationTestServer(t, responseData)
	defer srv.Close()

	c := NewWithURL("test-token", srv.URL)
	att, err := c.AttachmentUpdate("att1", AttachmentUpdateInput{Title: "Updated title", Subtitle: "New subtitle"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.Title != "Updated title" {
		t.Errorf("expected title 'Updated title', got %s", att.Title)
	}
}
