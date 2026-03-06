package client

import (
	"fmt"
	"time"
)

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
	After    string // курсор для пагинации
}

// PageInfo — информация о пагинации.
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

// IssueConnection — ответ на запрос списка задач.
type IssueConnection struct {
	Issues struct {
		Nodes    []Issue  `json:"nodes"`
		PageInfo PageInfo `json:"pageInfo"`
	} `json:"issues"`
}

// IssueResult — ответ на запрос одной задачи.
type IssueResult struct {
	Issue Issue `json:"issue"`
}

// ListIssues возвращает список задач по фильтрам и информацию о пагинации.
func (c *Client) ListIssues(filter IssueFilter) ([]Issue, *PageInfo, error) {
	query := `
query ListIssues($first: Int, $after: String, $filter: IssueFilter) {
  issues(first: $first, after: $after, filter: $filter) {
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
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}`

	variables := map[string]any{
		"first": filter.Limit,
	}
	if filter.After != "" {
		variables["after"] = filter.After
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
		return nil, nil, err
	}
	pi := result.Issues.PageInfo
	return result.Issues.Nodes, &pi, nil
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
	if result.Issue.ID == "" {
		return nil, fmt.Errorf("задача не найдена: %s", id)
	}
	return &result.Issue, nil
}

// ListComments возвращает список комментариев к задаче по её ID.
func (c *Client) ListComments(issueID string) ([]Comment, error) {
	query := `
query ListComments($issueId: String!) {
  issue(id: $issueId) {
    comments(first: 250) {
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
		Issue *struct {
			Comments struct {
				Nodes []Comment `json:"nodes"`
			} `json:"comments"`
		} `json:"issue"`
	}
	if err := c.Do(query, map[string]any{"issueId": issueID}, &result); err != nil {
		return nil, err
	}
	if result.Issue == nil {
		return nil, fmt.Errorf("задача не найдена: %s", issueID)
	}
	return result.Issue.Comments.Nodes, nil
}

// Project представляет проект Linear.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	State       string `json:"state"`
	StartDate   string `json:"startDate"`
	TargetDate  string `json:"targetDate"`
	URL         string `json:"url"`
	Lead        *User  `json:"lead"`
}

// ProjectConnection — ответ на запрос списка проектов.
type ProjectConnection struct {
	Projects struct {
		Nodes []Project `json:"nodes"`
	} `json:"projects"`
}

// TeamConnection — ответ на запрос списка команд.
type TeamConnection struct {
	Teams struct {
		Nodes []Team `json:"nodes"`
	} `json:"teams"`
}

// ListProjects возвращает список проектов.
func (c *Client) ListProjects() ([]Project, error) {
	query := `
query ListProjects {
  projects(first: 250) {
    nodes {
      id
      name
      description
      state
      startDate
      targetDate
      url
      lead { id name displayName email }
    }
  }
}`

	var result ProjectConnection
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Projects.Nodes, nil
}

// GetProject возвращает проект по ID.
func (c *Client) GetProject(id string) (*Project, error) {
	query := `
query GetProject($id: String!) {
  project(id: $id) {
    id
    name
    description
    state
    startDate
    targetDate
    url
    lead { id name displayName email }
  }
}`
	var result struct {
		Project *Project `json:"project"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Project == nil {
		return nil, fmt.Errorf("проект не найден: %s", id)
	}
	return result.Project, nil
}

// ListTeams возвращает список команд.
func (c *Client) ListTeams() ([]Team, error) {
	query := `
query ListTeams {
  teams(first: 250) {
    nodes {
      id
      key
      name
    }
  }
}`

	var result TeamConnection
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Teams.Nodes, nil
}

// IssueLabel представляет метку задачи Linear.
type IssueLabel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// ListIssueLabels возвращает все метки организации.
func (c *Client) ListIssueLabels() ([]IssueLabel, error) {
	query := `
query ListIssueLabels {
  issueLabels(first: 250) {
    nodes {
      id
      name
      color
      description
    }
  }
}`
	var result struct {
		IssueLabels struct {
			Nodes []IssueLabel `json:"nodes"`
		} `json:"issueLabels"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.IssueLabels.Nodes, nil
}

// GetIssueLabel возвращает метку по ID.
func (c *Client) GetIssueLabel(id string) (*IssueLabel, error) {
	query := `
query GetIssueLabel($id: String!) {
  issueLabel(id: $id) {
    id
    name
    color
    description
  }
}`
	var result struct {
		IssueLabel *IssueLabel `json:"issueLabel"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.IssueLabel == nil {
		return nil, fmt.Errorf("метка не найдена: %s", id)
	}
	return result.IssueLabel, nil
}

// IssueRelation представляет связь между задачами Linear.
type IssueRelation struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Issue        Issue  `json:"issue"`
	RelatedIssue Issue  `json:"relatedIssue"`
}

// IssueRelationIssue — краткое представление задачи в связи.
type IssueRelationIssue struct {
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
}

// ListIssueRelations возвращает список связей задачи по её ID.
func (c *Client) ListIssueRelations(issueID string) ([]IssueRelation, error) {
	query := `
query ListIssueRelations($id: String!) {
  issue(id: $id) {
    relations(first: 250) {
      nodes {
        id
        type
        relatedIssue { id identifier title }
      }
    }
  }
}`
	var result struct {
		Issue *struct {
			Relations struct {
				Nodes []IssueRelation `json:"nodes"`
			} `json:"relations"`
		} `json:"issue"`
	}
	if err := c.Do(query, map[string]any{"id": issueID}, &result); err != nil {
		return nil, err
	}
	if result.Issue == nil {
		return nil, fmt.Errorf("задача не найдена: %s", issueID)
	}
	return result.Issue.Relations.Nodes, nil
}

// SearchIssues выполняет поиск задач по строке запроса.
func (c *Client) SearchIssues(query string, limit int) ([]Issue, *PageInfo, error) {
	gql := `
query SearchIssues($query: String!, $first: Int) {
  issueSearch(query: $query, first: $first) {
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
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}`

	variables := map[string]any{
		"query": query,
		"first": limit,
	}

	var result struct {
		IssueSearch struct {
			Nodes    []Issue  `json:"nodes"`
			PageInfo PageInfo `json:"pageInfo"`
		} `json:"issueSearch"`
	}
	if err := c.Do(gql, variables, &result); err != nil {
		return nil, nil, err
	}
	pi := result.IssueSearch.PageInfo
	return result.IssueSearch.Nodes, &pi, nil
}

// ProjectMilestone представляет веху (milestone) проекта Linear.
type ProjectMilestone struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	TargetDate  string `json:"targetDate"`
	Description string `json:"description"`
}

// ListProjectMilestones возвращает список вех проекта по его ID.
func (c *Client) ListProjectMilestones(projectID string) ([]ProjectMilestone, error) {
	query := `
query ListProjectMilestones($id: String!) {
  project(id: $id) {
    projectMilestones(first: 250) {
      nodes {
        id
        name
        targetDate
        description
      }
    }
  }
}`
	var result struct {
		Project *struct {
			ProjectMilestones struct {
				Nodes []ProjectMilestone `json:"nodes"`
			} `json:"projectMilestones"`
		} `json:"project"`
	}
	if err := c.Do(query, map[string]any{"id": projectID}, &result); err != nil {
		return nil, err
	}
	if result.Project == nil {
		return nil, fmt.Errorf("проект не найден: %s", projectID)
	}
	return result.Project.ProjectMilestones.Nodes, nil
}

// ProjectUpdate представляет обновление (запись журнала) проекта Linear.
type ProjectUpdate struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	Health    string    `json:"health"`
	CreatedAt time.Time `json:"createdAt"`
	User      *User     `json:"user"`
}

// ListProjectUpdates возвращает список обновлений проекта по его ID.
func (c *Client) ListProjectUpdates(projectID string) ([]ProjectUpdate, error) {
	query := `
query ListProjectUpdates($id: String!) {
  project(id: $id) {
    projectUpdates(first: 250) {
      nodes {
        id
        body
        health
        createdAt
        user { id name displayName email }
      }
    }
  }
}`
	var result struct {
		Project *struct {
			ProjectUpdates struct {
				Nodes []ProjectUpdate `json:"nodes"`
			} `json:"projectUpdates"`
		} `json:"project"`
	}
	if err := c.Do(query, map[string]any{"id": projectID}, &result); err != nil {
		return nil, err
	}
	if result.Project == nil {
		return nil, fmt.Errorf("проект не найден: %s", projectID)
	}
	return result.Project.ProjectUpdates.Nodes, nil
}

// ProjectLabel представляет метку проекта Linear.
type ProjectLabel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// ListProjectLabels возвращает все метки проектов организации.
func (c *Client) ListProjectLabels() ([]ProjectLabel, error) {
	query := `
query ListProjectLabels {
  projectLabels(first: 250) {
    nodes {
      id
      name
      color
      description
    }
  }
}`
	var result struct {
		ProjectLabels struct {
			Nodes []ProjectLabel `json:"nodes"`
		} `json:"projectLabels"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.ProjectLabels.Nodes, nil
}

// ProjectStatus представляет статус проекта Linear.
type ProjectStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Position    float64 `json:"position"`
}

// ListProjectStatuses возвращает все статусы проектов организации.
func (c *Client) ListProjectStatuses() ([]ProjectStatus, error) {
	query := `
query ListProjectStatuses {
  projectStatuses(first: 250) {
    nodes {
      id
      name
      type
      color
      description
      position
    }
  }
}`
	var result struct {
		ProjectStatuses struct {
			Nodes []ProjectStatus `json:"nodes"`
		} `json:"projectStatuses"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.ProjectStatuses.Nodes, nil
}

// SearchProjects выполняет поиск проектов по строке запроса.
func (c *Client) SearchProjects(query string) ([]Project, error) {
	gql := `
query SearchProjects($filter: ProjectFilter) {
  projects(first: 50, filter: $filter) {
    nodes {
      id
      name
      description
      state
      startDate
      targetDate
      url
      lead { id name displayName email }
    }
  }
}`

	variables := map[string]any{
		"filter": map[string]any{
			"name": map[string]any{
				"containsIgnoreCase": query,
			},
		},
	}

	var result ProjectConnection
	if err := c.Do(gql, variables, &result); err != nil {
		return nil, err
	}
	return result.Projects.Nodes, nil
}

// Cycle представляет цикл (спринт) команды Linear.
type Cycle struct {
	ID          string     `json:"id"`
	Number      int        `json:"number"`
	Name        string     `json:"name"`
	StartsAt    time.Time  `json:"startsAt"`
	EndsAt      time.Time  `json:"endsAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Team        Team       `json:"team"`
}

// ListCycles возвращает список циклов команды по её ID.
func (c *Client) ListCycles(teamID string) ([]Cycle, error) {
	query := `
query ListCycles($teamId: String!) {
  cycles(filter: { team: { id: { eq: $teamId } } }, first: 250) {
    nodes {
      id
      number
      name
      startsAt
      endsAt
      completedAt
      team { id key name }
    }
  }
}`
	var result struct {
		Cycles struct {
			Nodes []Cycle `json:"nodes"`
		} `json:"cycles"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return nil, err
	}
	return result.Cycles.Nodes, nil
}

// GetCycle возвращает цикл по ID.
func (c *Client) GetCycle(id string) (*Cycle, error) {
	query := `
query GetCycle($id: String!) {
  cycle(id: $id) {
    id
    number
    name
    startsAt
    endsAt
    completedAt
    team { id key name }
  }
}`
	var result struct {
		Cycle *Cycle `json:"cycle"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Cycle == nil {
		return nil, fmt.Errorf("цикл не найден: %s", id)
	}
	return result.Cycle, nil
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

// WorkflowState представляет статус задачи (workflow state) в Linear.
type WorkflowState struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Team  Team   `json:"team"`
}

// ListWorkflowStates возвращает список статусов задач команды.
func (c *Client) ListWorkflowStates(teamID string) ([]WorkflowState, error) {
	query := `
query ListWorkflowStates($teamId: String!) {
  workflowStates(filter: { team: { id: { eq: $teamId } } }, first: 250) {
    nodes {
      id
      name
      type
      color
      team { id key name }
    }
  }
}`
	var result struct {
		WorkflowStates struct {
			Nodes []WorkflowState `json:"nodes"`
		} `json:"workflowStates"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return nil, err
	}
	return result.WorkflowStates.Nodes, nil
}

// TeamMembership представляет членство пользователя в команде.
type TeamMembership struct {
	ID   string `json:"id"`
	User User   `json:"user"`
	Team Team   `json:"team"`
	Role string `json:"role"`
}

// ListTeamMembers возвращает список членов команды.
func (c *Client) ListTeamMembers(teamID string) ([]TeamMembership, error) {
	query := `
query ListTeamMembers($teamId: String!) {
  team(id: $teamId) {
    members(first: 250) {
      nodes {
        id
        user { id name displayName email }
        team { id key name }
        role
      }
    }
  }
}`
	var result struct {
		Team struct {
			Members struct {
				Nodes []TeamMembership `json:"nodes"`
			} `json:"members"`
		} `json:"team"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return nil, err
	}
	return result.Team.Members.Nodes, nil
}
