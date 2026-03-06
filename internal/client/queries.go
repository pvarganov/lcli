package client

import "time"

// Comment представляет комментарий к задаче Linear.
type Comment struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
	User      *User     `json:"user"`
}

// Issue представляет задачу Linear.
type Issue struct {
	ID         string     `json:"id"`
	Identifier string     `json:"identifier"`
	Title      string     `json:"title"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	Priority   int        `json:"priority"`
	State      IssueState `json:"state"`
	Assignee   *User      `json:"assignee"`
	Description string    `json:"description"`
	Team       Team       `json:"team"`
}

type IssueState struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

type Team struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// IssueFilter — параметры фильтрации для списка задач.
type IssueFilter struct {
	Assignee string
	Status   string
	Team     string
	Limit    int
}

// IssueConnection — ответ на запрос списка задач.
type IssueConnection struct {
	Issues struct {
		Nodes []Issue `json:"nodes"`
	} `json:"issues"`
}

// IssueResult — ответ на запрос одной задачи.
type IssueResult struct {
	Issue Issue `json:"issue"`
}

// ListIssues возвращает список задач по фильтрам.
func (c *Client) ListIssues(filter IssueFilter) ([]Issue, error) {
	query := `
query ListIssues($first: Int, $filter: IssueFilter) {
  issues(first: $first, filter: $filter) {
    nodes {
      id
      identifier
      title
      updatedAt
      priority
      state { name type }
      assignee { id name displayName email }
      team { id key name }
    }
  }
}`

	variables := map[string]any{
		"first": filter.Limit,
	}

	issueFilter := map[string]any{}
	if filter.Team != "" {
		issueFilter["team"] = map[string]any{
			"key": map[string]any{"eq": filter.Team},
		}
	}
	if filter.Status != "" {
		issueFilter["state"] = map[string]any{
			"name": map[string]any{"eq": filter.Status},
		}
	}
	if filter.Assignee != "" {
		issueFilter["assignee"] = map[string]any{
			"displayName": map[string]any{"eq": filter.Assignee},
		}
	}
	if len(issueFilter) > 0 {
		variables["filter"] = issueFilter
	}

	var result IssueConnection
	if err := c.Do(query, variables, &result); err != nil {
		return nil, err
	}
	return result.Issues.Nodes, nil
}

// GetIssue возвращает задачу по ID.
func (c *Client) GetIssue(id string) (*Issue, error) {
	query := `
query GetIssue($id: String!) {
  issue(id: $id) {
    id
    identifier
    title
    description
    updatedAt
    priority
    state { name type }
    assignee { id name displayName email }
    team { id key name }
  }
}`

	variables := map[string]any{
		"id": id,
	}

	var result IssueResult
	if err := c.Do(query, variables, &result); err != nil {
		return nil, err
	}
	return &result.Issue, nil
}

// ListComments возвращает список комментариев к задаче по её ID.
func (c *Client) ListComments(issueID string) ([]Comment, error) {
	query := `
query ListComments($issueId: String!) {
  issue(id: $issueId) {
    comments {
      nodes {
        id
        body
        createdAt
        user { id name displayName email }
      }
    }
  }
}`

	var result struct {
		Issue struct {
			Comments struct {
				Nodes []Comment `json:"nodes"`
			} `json:"comments"`
		} `json:"issue"`
	}
	if err := c.Do(query, map[string]any{"issueId": issueID}, &result); err != nil {
		return nil, err
	}
	return result.Issue.Comments.Nodes, nil
}

// PriorityLabel возвращает текстовое обозначение приоритета.
func PriorityLabel(p int) string {
	switch p {
	case 1:
		return "Urgent"
	case 2:
		return "High"
	case 3:
		return "Medium"
	case 4:
		return "Low"
	default:
		return "No priority"
	}
}
