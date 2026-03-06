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
