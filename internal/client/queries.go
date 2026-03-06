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
	Assignee  string
	Status    string
	Team      string
	Limit     int
	After     string // курсор для пагинации
	Priority  int    // -1 = не задан; 0=нет, 1=срочно, 2=высокий, 3=средний, 4=низкий
	Label     string
	ProjectID string
	CycleID   string
	Creator   string
	OrderBy   string
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
query ListIssues($first: Int, $after: String, $filter: IssueFilter, $orderBy: PaginationOrderBy) {
  issues(first: $first, after: $after, filter: $filter, orderBy: $orderBy) {
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
	if filter.OrderBy != "" {
		switch filter.OrderBy {
		case "createdAt", "updatedAt", "priority", "manualOrder":
			variables["orderBy"] = filter.OrderBy
		default:
			return nil, nil, fmt.Errorf("недопустимое значение --order-by %q: допустимые значения: createdAt, updatedAt, priority, manualOrder", filter.OrderBy)
		}
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
	if filter.Priority >= 0 {
		issueFilter["priority"] = map[string]any{
			"eq": filter.Priority,
		}
	}
	if filter.Label != "" {
		issueFilter["labels"] = map[string]any{
			"some": map[string]any{
				"name": map[string]any{"eqIgnoreCase": filter.Label},
			},
		}
	}
	if filter.ProjectID != "" {
		issueFilter["project"] = map[string]any{
			"id": map[string]any{"eq": filter.ProjectID},
		}
	}
	if filter.CycleID != "" {
		issueFilter["cycle"] = map[string]any{
			"id": map[string]any{"eq": filter.CycleID},
		}
	}
	if filter.Creator != "" {
		issueFilter["creator"] = map[string]any{
			"displayName": map[string]any{"eq": filter.Creator},
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

// FindLabelsByNames ищет метки по именам в рамках команды и возвращает их ID.
func (c *Client) FindLabelsByNames(teamID string, names []string) ([]string, error) {
	query := `
query GetIssueLabels($teamId: ID!) {
  issueLabels(first: 250, filter: { team: { id: { eq: $teamId } } }) {
    nodes {
      id
      name
    }
  }
}`
	var result struct {
		IssueLabels struct {
			Nodes []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"nodes"`
		} `json:"issueLabels"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return nil, err
	}
	nameToID := make(map[string]string, len(result.IssueLabels.Nodes))
	for _, lbl := range result.IssueLabels.Nodes {
		nameToID[lbl.Name] = lbl.ID
	}
	ids := make([]string, 0, len(names))
	for _, name := range names {
		id, ok := nameToID[name]
		if !ok {
			return nil, fmt.Errorf("метка не найдена: %s", name)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// IssueRelation представляет связь между задачами Linear.
type IssueRelation struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Issue        Issue  `json:"issue"`
	RelatedIssue Issue  `json:"relatedIssue"`
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

// ListUsers возвращает всех пользователей организации.
func (c *Client) ListUsers() ([]User, error) {
	query := `
query ListUsers {
  users(first: 250) {
    nodes {
      id
      name
      displayName
      email
    }
  }
}`
	var result struct {
		Users struct {
			Nodes []User `json:"nodes"`
		} `json:"users"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Users.Nodes, nil
}

// GetUser возвращает пользователя по ID.
func (c *Client) GetUser(id string) (*User, error) {
	query := `
query GetUser($id: String!) {
  user(id: $id) {
    id
    name
    displayName
    email
  }
}`
	var result struct {
		User *User `json:"user"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.User == nil {
		return nil, fmt.Errorf("пользователь не найден: %s", id)
	}
	return result.User, nil
}

// GetViewer возвращает текущего аутентифицированного пользователя.
func (c *Client) GetViewer() (*User, error) {
	query := `
query GetViewer {
  viewer {
    id
    name
    displayName
    email
  }
}`
	var result struct {
		Viewer *User `json:"viewer"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	if result.Viewer == nil {
		return nil, fmt.Errorf("не удалось получить данные текущего пользователя")
	}
	return result.Viewer, nil
}

// Notification представляет уведомление Linear.
type Notification struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	ReadAt    *time.Time `json:"readAt"`
	CreatedAt time.Time `json:"createdAt"`
	Issue     *Issue    `json:"issue"`
	Comment   *Comment  `json:"comment"`
	Project   *Project  `json:"project"`
}

// ListNotifications возвращает список уведомлений с пагинацией.
func (c *Client) ListNotifications(limit int, after string) ([]Notification, *PageInfo, error) {
	if limit <= 0 {
		limit = 50
	}
	query := `
query ListNotifications($first: Int, $after: String) {
  notifications(first: $first, after: $after) {
    nodes {
      id
      type
      readAt
      createdAt
      issue {
        id
        identifier
        title
      }
      comment {
        id
        body
      }
      project {
        id
        name
      }
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}`
	vars := map[string]any{"first": limit}
	if after != "" {
		vars["after"] = after
	}
	var result struct {
		Notifications struct {
			Nodes    []Notification `json:"nodes"`
			PageInfo PageInfo       `json:"pageInfo"`
		} `json:"notifications"`
	}
	if err := c.Do(query, vars, &result); err != nil {
		return nil, nil, err
	}
	pi := result.Notifications.PageInfo
	return result.Notifications.Nodes, &pi, nil
}

// GetNotificationsUnreadCount возвращает количество непрочитанных уведомлений.
func (c *Client) GetNotificationsUnreadCount() (int, error) {
	query := `
query GetNotificationsUnreadCount($first: Int, $after: String) {
  notifications(first: $first, filter: { readAt: { null: true } }, after: $after) {
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      id
    }
  }
}`
	type resultType struct {
		Notifications struct {
			PageInfo struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
			Nodes []struct {
				ID string `json:"id"`
			} `json:"nodes"`
		} `json:"notifications"`
	}
	count := 0
	var cursor *string
	for {
		vars := map[string]any{"first": 250}
		if cursor != nil {
			vars["after"] = *cursor
		}
		var result resultType
		if err := c.Do(query, vars, &result); err != nil {
			return 0, err
		}
		count += len(result.Notifications.Nodes)
		if !result.Notifications.PageInfo.HasNextPage {
			break
		}
		end := result.Notifications.PageInfo.EndCursor
		cursor = &end
	}
	return count, nil
}

// Webhook представляет вебхук Linear.
type Webhook struct {
	ID            string   `json:"id"`
	URL           string   `json:"url"`
	Enabled       bool     `json:"enabled"`
	Secret        string   `json:"secret"`
	ResourceTypes []string `json:"resourceTypes"`
	Team          *Team    `json:"team"`
}

// ListWebhooks возвращает список всех вебхуков организации.
func (c *Client) ListWebhooks() ([]Webhook, error) {
	query := `
query ListWebhooks {
  webhooks {
    nodes {
      id
      url
      enabled
      secret
      resourceTypes
      team { id key name }
    }
  }
}`
	var result struct {
		Webhooks struct {
			Nodes []Webhook `json:"nodes"`
		} `json:"webhooks"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Webhooks.Nodes, nil
}

// GetWebhook возвращает вебхук по ID.
func (c *Client) GetWebhook(id string) (*Webhook, error) {
	query := `
query GetWebhook($id: String!) {
  webhook(id: $id) {
    id
    url
    enabled
    secret
    resourceTypes
    team { id key name }
  }
}`
	var result struct {
		Webhook *Webhook `json:"webhook"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Webhook == nil {
		return nil, fmt.Errorf("вебхук не найден: %s", id)
	}
	return result.Webhook, nil
}

// Attachment представляет вложение к задаче Linear.
type Attachment struct {
	ID         string  `json:"id"`
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	SourceType string  `json:"sourceType"`
	Subtitle   string  `json:"subtitle"`
	Issue      *Issue  `json:"issue"`
}

// ListAttachments возвращает список вложений задачи.
func (c *Client) ListAttachments(issueID string) ([]Attachment, error) {
	query := `
query ListAttachments($id: String!) {
  issue(id: $id) {
    attachments {
      nodes {
        id
        title
        url
        sourceType
        subtitle
      }
    }
  }
}`
	var result struct {
		Issue *struct {
			Attachments struct {
				Nodes []Attachment `json:"nodes"`
			} `json:"attachments"`
		} `json:"issue"`
	}
	if err := c.Do(query, map[string]any{"id": issueID}, &result); err != nil {
		return nil, err
	}
	if result.Issue == nil {
		return nil, fmt.Errorf("задача не найдена: %s", issueID)
	}
	return result.Issue.Attachments.Nodes, nil
}

// Document представляет документ Linear.
type Document struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Project   *Project  `json:"project"`
	Creator   *User     `json:"creator"`
}

// ListDocuments возвращает список документов организации.
func (c *Client) ListDocuments() ([]Document, error) {
	query := `
query {
  documents {
    nodes {
      id
      title
      content
      createdAt
      updatedAt
      project { id name }
      creator { id name displayName }
    }
  }
}`
	var result struct {
		Documents struct {
			Nodes []Document `json:"nodes"`
		} `json:"documents"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Documents.Nodes, nil
}

// GetDocument возвращает документ по ID.
func (c *Client) GetDocument(id string) (*Document, error) {
	query := `
query GetDocument($id: String!) {
  document(id: $id) {
    id
    title
    content
    createdAt
    updatedAt
    project { id name }
    creator { id name displayName }
  }
}`
	var result struct {
		Document *Document `json:"document"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Document == nil {
		return nil, fmt.Errorf("документ %q не найден", id)
	}
	return result.Document, nil
}

// SearchDocuments выполняет поиск документов по запросу.
func (c *Client) SearchDocuments(query string) ([]Document, error) {
	gqlQuery := `
query SearchDocuments($term: String!) {
  searchDocuments(term: $term) {
    nodes {
      id
      title
      content
      createdAt
      updatedAt
      project { id name }
      creator { id name displayName }
    }
  }
}`
	var result struct {
		SearchDocuments struct {
			Nodes []Document `json:"nodes"`
		} `json:"searchDocuments"`
	}
	if err := c.Do(gqlQuery, map[string]any{"term": query}, &result); err != nil {
		return nil, err
	}
	return result.SearchDocuments.Nodes, nil
}

// Initiative представляет инициативу в Linear.
type Initiative struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
	Owner       *User   `json:"owner"`
	ArchivedAt  *string `json:"archivedAt"`
}

// ListInitiatives возвращает список инициатив организации.
func (c *Client) ListInitiatives() ([]Initiative, error) {
	query := `
query {
  initiatives {
    nodes {
      id
      name
      description
      status
      owner { id name displayName email }
      archivedAt
    }
  }
}`
	var result struct {
		Initiatives struct {
			Nodes []Initiative `json:"nodes"`
		} `json:"initiatives"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Initiatives.Nodes, nil
}

// GetInitiative возвращает инициативу по ID.
func (c *Client) GetInitiative(id string) (*Initiative, error) {
	query := `
query GetInitiative($id: String!) {
  initiative(id: $id) {
    id
    name
    description
    status
    owner { id name displayName email }
    archivedAt
  }
}`
	var result struct {
		Initiative *Initiative `json:"initiative"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Initiative == nil {
		return nil, fmt.Errorf("инициатива %q не найдена", id)
	}
	return result.Initiative, nil
}

// InitiativeUpdate представляет обновление (запись журнала) инициативы Linear.
type InitiativeUpdate struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	Health    string    `json:"health"`
	CreatedAt time.Time `json:"createdAt"`
	User      *User     `json:"user"`
}

// ListInitiativeUpdates возвращает список обновлений инициативы по её ID.
func (c *Client) ListInitiativeUpdates(initiativeID string) ([]InitiativeUpdate, error) {
	query := `
query ListInitiativeUpdates($id: String!) {
  initiative(id: $id) {
    initiativeUpdates(first: 250) {
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
		Initiative *struct {
			InitiativeUpdates struct {
				Nodes []InitiativeUpdate `json:"nodes"`
			} `json:"initiativeUpdates"`
		} `json:"initiative"`
	}
	if err := c.Do(query, map[string]any{"id": initiativeID}, &result); err != nil {
		return nil, err
	}
	if result.Initiative == nil {
		return nil, fmt.Errorf("инициатива не найдена: %s", initiativeID)
	}
	return result.Initiative.InitiativeUpdates.Nodes, nil
}

// InitiativeToProject представляет связь инициативы с проектом.
type InitiativeToProject struct {
	ID         string     `json:"id"`
	Initiative Initiative `json:"initiative"`
	Project    Project    `json:"project"`
}

// CustomerStatus представляет статус клиента.
type CustomerStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// CustomerTier представляет уровень клиента.
type CustomerTier struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// CustomerNeed представляет потребность клиента, привязанную к задаче или проекту.
type CustomerNeed struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	Priority  float64   `json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
	Creator   *User     `json:"creator"`
	Customer  *Customer `json:"customer"`
	Issue     *Issue    `json:"issue"`
}

// Customer представляет клиента Linear (CRM).
type Customer struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	LogoURL   string          `json:"logoUrl"`
	SlugID    string          `json:"slugId"`
	Revenue   int             `json:"revenue"`
	Size      float64         `json:"size"`
	CreatedAt time.Time       `json:"createdAt"`
	Owner     *User           `json:"owner"`
	Status    *CustomerStatus `json:"status"`
	Tier      *CustomerTier   `json:"tier"`
}

// ListCustomers возвращает список всех клиентов организации.
func (c *Client) ListCustomers() ([]Customer, error) {
	query := `
query ListCustomers {
  customers(first: 250) {
    nodes {
      id
      name
      logoUrl
      slugId
      revenue
      size
      createdAt
      owner { id name displayName email }
      status { id name displayName color }
      tier { id name displayName color }
    }
  }
}`
	var result struct {
		Customers struct {
			Nodes []Customer `json:"nodes"`
		} `json:"customers"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Customers.Nodes, nil
}

// GetCustomer возвращает клиента по ID.
func (c *Client) GetCustomer(id string) (*Customer, error) {
	query := `
query GetCustomer($id: String!) {
  customer(id: $id) {
    id
    name
    logoUrl
    slugId
    revenue
    size
    createdAt
    owner { id name displayName email }
    status { id name displayName color }
    tier { id name displayName color }
  }
}`
	var result struct {
		Customer *Customer `json:"customer"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Customer == nil {
		return nil, fmt.Errorf("клиент не найден: %s", id)
	}
	return result.Customer, nil
}

// ListCustomerNeeds возвращает список потребностей клиентов.
func (c *Client) ListCustomerNeeds() ([]CustomerNeed, error) {
	query := `
query ListCustomerNeeds {
  customerNeeds(first: 250) {
    nodes {
      id
      body
      priority
      createdAt
      creator { id name displayName email }
      customer { id name }
      issue { id identifier title }
    }
  }
}`
	var result struct {
		CustomerNeeds struct {
			Nodes []CustomerNeed `json:"nodes"`
		} `json:"customerNeeds"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.CustomerNeeds.Nodes, nil
}

// ListCustomerStatuses возвращает список статусов клиентов.
func (c *Client) ListCustomerStatuses() ([]CustomerStatus, error) {
	query := `
query ListCustomerStatuses {
  customerStatuses(first: 250) {
    nodes {
      id
      name
      displayName
      color
      description
    }
  }
}`
	var result struct {
		CustomerStatuses struct {
			Nodes []CustomerStatus `json:"nodes"`
		} `json:"customerStatuses"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.CustomerStatuses.Nodes, nil
}

// ListCustomerTiers возвращает список уровней клиентов.
func (c *Client) ListCustomerTiers() ([]CustomerTier, error) {
	query := `
query ListCustomerTiers {
  customerTiers(first: 250) {
    nodes {
      id
      name
      displayName
      color
      description
    }
  }
}`
	var result struct {
		CustomerTiers struct {
			Nodes []CustomerTier `json:"nodes"`
		} `json:"customerTiers"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.CustomerTiers.Nodes, nil
}

// Roadmap представляет дорожную карту Linear.
type Roadmap struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Owner       *User     `json:"owner"`
}

// RoadmapToProject представляет связь дорожной карты с проектом.
type RoadmapToProject struct {
	ID      string  `json:"id"`
	Roadmap Roadmap `json:"roadmap"`
	Project Project `json:"project"`
}

// ListRoadmaps возвращает список всех дорожных карт организации.
func (c *Client) ListRoadmaps() ([]Roadmap, error) {
	query := `
query ListRoadmaps {
  roadmaps(first: 250) {
    nodes {
      id
      name
      description
      createdAt
      owner { id name displayName email }
    }
  }
}`
	var result struct {
		Roadmaps struct {
			Nodes []Roadmap `json:"nodes"`
		} `json:"roadmaps"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Roadmaps.Nodes, nil
}

// GetRoadmap возвращает дорожную карту по ID.
func (c *Client) GetRoadmap(id string) (*Roadmap, error) {
	query := `
query GetRoadmap($id: String!) {
  roadmap(id: $id) {
    id
    name
    description
    createdAt
    owner { id name displayName email }
  }
}`
	var result struct {
		Roadmap *Roadmap `json:"roadmap"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Roadmap == nil {
		return nil, fmt.Errorf("дорожная карта не найдена: %s", id)
	}
	return result.Roadmap, nil
}

// Template представляет шаблон Linear.
type Template struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Type         string    `json:"type"`
	TemplateData string    `json:"templateData"`
	CreatedAt    time.Time `json:"createdAt"`
	Creator      *User     `json:"creator"`
	Team         *Team     `json:"team"`
}

// ListTemplates возвращает список шаблонов организации.
func (c *Client) ListTemplates() ([]Template, error) {
	query := `
query ListTemplates {
  templates {
    id
    name
    description
    type
    templateData
    createdAt
    creator { id name displayName email }
    team { id name key }
  }
}`
	var result struct {
		Templates []Template `json:"templates"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Templates, nil
}

// GetTemplate возвращает шаблон по ID.
func (c *Client) GetTemplate(id string) (*Template, error) {
	query := `
query GetTemplate($id: String!) {
  template(id: $id) {
    id
    name
    description
    type
    templateData
    createdAt
    creator { id name displayName email }
    team { id name key }
  }
}`
	var result struct {
		Template *Template `json:"template"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Template == nil {
		return nil, fmt.Errorf("шаблон не найден: %s", id)
	}
	return result.Template, nil
}

// Organization представляет организацию Linear.
type Organization struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	URLKey             string    `json:"urlKey"`
	LogoURL            string    `json:"logoUrl"`
	CreatedAt          time.Time `json:"createdAt"`
	PeriodUploadVolume float64   `json:"periodUploadVolume"`
}

// OrganizationInvite представляет приглашение в организацию.
type OrganizationInvite struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	Inviter   *User     `json:"inviter"`
}

// GetOrganization возвращает информацию о текущей организации.
func (c *Client) GetOrganization() (*Organization, error) {
	query := `
query GetOrganization {
  organization {
    id
    name
    urlKey
    logoUrl
    createdAt
    periodUploadVolume
  }
}`
	var result struct {
		Organization *Organization `json:"organization"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	if result.Organization == nil {
		return nil, fmt.Errorf("организация не найдена")
	}
	return result.Organization, nil
}

// ListOrganizationInvites возвращает список приглашений в организацию.
func (c *Client) ListOrganizationInvites() ([]OrganizationInvite, error) {
	query := `
query ListOrganizationInvites {
  organizationInvites {
    nodes {
      id
      email
      role
      createdAt
      expiresAt
      inviter { id name displayName email }
    }
  }
}`
	var result struct {
		OrganizationInvites struct {
			Nodes []OrganizationInvite `json:"nodes"`
		} `json:"organizationInvites"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.OrganizationInvites.Nodes, nil
}

// CustomView представляет пользовательское представление в Linear.
type CustomView struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Owner       *User  `json:"owner"`
}

// ListCustomViews возвращает список пользовательских представлений.
func (c *Client) ListCustomViews() ([]CustomView, error) {
	query := `
query ListCustomViews {
  customViews {
    nodes {
      id
      name
      description
      icon
      color
      owner { id name displayName email }
    }
  }
}`
	var result struct {
		CustomViews struct {
			Nodes []CustomView `json:"nodes"`
		} `json:"customViews"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.CustomViews.Nodes, nil
}

// Favorite представляет избранный элемент в Linear.
type Favorite struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Issue      *Issue      `json:"issue"`
	Project    *Project    `json:"project"`
	Label      *IssueLabel `json:"label"`
	CustomView *CustomView `json:"customView"`
}

// ListFavorites возвращает список всех избранных элементов текущего пользователя.
func (c *Client) ListFavorites() ([]Favorite, error) {
	query := `
query ListFavorites {
  favorites {
    nodes {
      id
      type
      issue { id identifier title }
      project { id name }
      label { id name color }
      customView { id name }
    }
  }
}`
	var result struct {
		Favorites struct {
			Nodes []Favorite `json:"nodes"`
		} `json:"favorites"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Favorites.Nodes, nil
}

// GetCustomView возвращает пользовательское представление по ID.
func (c *Client) GetCustomView(id string) (*CustomView, error) {
	query := `
query GetCustomView($id: String!) {
  customView(id: $id) {
    id
    name
    description
    icon
    color
    owner { id name displayName email }
  }
}`
	var result struct {
		CustomView *CustomView `json:"customView"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.CustomView == nil {
		return nil, fmt.Errorf("представление не найдено: %s", id)
	}
	return result.CustomView, nil
}

// Emoji представляет эмодзи в организации Linear.
type Emoji struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Creator *User  `json:"creator"`
}

// ListEmojis возвращает список всех эмодзи организации.
func (c *Client) ListEmojis() ([]Emoji, error) {
	query := `
query ListEmojis {
  emojis {
    nodes {
      id
      name
      url
      creator { id name displayName }
    }
  }
}`
	var result struct {
		Emojis struct {
			Nodes []Emoji `json:"nodes"`
		} `json:"emojis"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Emojis.Nodes, nil
}

// Release представляет выпуск (релиз) в Linear.
type Release struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Version     string          `json:"version"`
	CommitSha   string          `json:"commitSha"`
	CreatedAt   string          `json:"createdAt"`
	CompletedAt string          `json:"completedAt"`
	CanceledAt  string          `json:"canceledAt"`
	Pipeline    *ReleasePipeline `json:"pipeline"`
	Stage       *ReleaseStage   `json:"stage"`
}

// ReleasePipeline представляет пайплайн релизов в Linear.
type ReleasePipeline struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	SlugID    string          `json:"slugId"`
	Type      string          `json:"type"`
	CreatedAt string          `json:"createdAt"`
	Stages    []ReleaseStage  `json:"stages"`
}

// ReleaseStage представляет стадию в пайплайне релизов.
type ReleaseStage struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	Position float64 `json:"position"`
	Frozen   bool   `json:"frozen"`
	Type     string `json:"type"`
}

// ListReleases возвращает список релизов.
func (c *Client) ListReleases() ([]Release, error) {
	query := `
query ListReleases {
  releases {
    nodes {
      id
      name
      description
      createdAt
      completedAt
      canceledAt
      pipeline { id name }
      stage { id name color }
    }
  }
}`
	var result struct {
		Releases struct {
			Nodes []Release `json:"nodes"`
		} `json:"releases"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Releases.Nodes, nil
}

// GetRelease возвращает релиз по ID.
func (c *Client) GetRelease(id string) (*Release, error) {
	query := `
query GetRelease($id: String!) {
  release(id: $id) {
    id
    name
    description
    createdAt
    completedAt
    canceledAt
    commitSha
    pipeline { id name }
    stage { id name color }
  }
}`
	var result struct {
		Release *Release `json:"release"`
	}
	if err := c.Do(query, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if result.Release == nil {
		return nil, fmt.Errorf("релиз не найден: %s", id)
	}
	return result.Release, nil
}

// SearchReleases ищет релизы по строке запроса.
func (c *Client) SearchReleases(q string) ([]Release, error) {
	query := `
query SearchReleases($term: String!) {
  releaseSearch(term: $term) {
    nodes {
      id
      name
      description
      createdAt
      completedAt
      pipeline { id name }
      stage { id name color }
    }
  }
}`
	var result struct {
		ReleaseSearch struct {
			Nodes []Release `json:"nodes"`
		} `json:"releaseSearch"`
	}
	if err := c.Do(query, map[string]any{"term": q}, &result); err != nil {
		return nil, err
	}
	return result.ReleaseSearch.Nodes, nil
}

// ListReleasePipelines возвращает список пайплайнов релизов.
func (c *Client) ListReleasePipelines() ([]ReleasePipeline, error) {
	query := `
query ListReleasePipelines {
  releasePipelines {
    nodes {
      id
      name
      slugId
      type
      createdAt
    }
  }
}`
	var result struct {
		ReleasePipelines struct {
			Nodes []ReleasePipeline `json:"nodes"`
		} `json:"releasePipelines"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.ReleasePipelines.Nodes, nil
}

// Integration представляет интеграцию Linear с внешним сервисом.
type Integration struct {
	ID           string  `json:"id"`
	Service      string  `json:"service"`
	CreatedAt    string  `json:"createdAt"`
	Team         *Team   `json:"team"`
	Organization *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"organization"`
}

// ListIntegrations возвращает список интеграций организации.
func (c *Client) ListIntegrations() ([]Integration, error) {
	query := `
query ListIntegrations {
  integrations {
    nodes {
      id
      service
      createdAt
      team { id key name }
    }
  }
}`
	var result struct {
		Integrations struct {
			Nodes []Integration `json:"nodes"`
		} `json:"integrations"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.Integrations.Nodes, nil
}

// GitAutomationTargetBranch представляет целевую ветку для git автоматизации.
type GitAutomationTargetBranch struct {
	ID            string `json:"id"`
	BranchPattern string `json:"branchPattern"`
	IsRegex       bool   `json:"isRegex"`
	Team          *Team  `json:"team"`
}

// GitAutomationState представляет правило git автоматизации.
type GitAutomationState struct {
	ID            string                     `json:"id"`
	Event         string                     `json:"event"`
	BranchPattern string                     `json:"branchPattern"`
	State         *WorkflowState             `json:"state"`
	TargetBranch  *GitAutomationTargetBranch `json:"targetBranch"`
	Team          *Team                      `json:"team"`
}

// ListGitAutomationStates возвращает список правил git автоматизации для команды.
func (c *Client) ListGitAutomationStates(teamID string) ([]GitAutomationState, error) {
	query := `
query ListGitAutomationStates($teamId: String!) {
  team(id: $teamId) {
    gitAutomationStates {
      nodes {
        id
        event
        branchPattern
        state { id name type color }
        targetBranch { id branchPattern isRegex }
        team { id key name }
      }
    }
  }
}`
	var result struct {
		Team struct {
			GitAutomationStates struct {
				Nodes []GitAutomationState `json:"nodes"`
			} `json:"gitAutomationStates"`
		} `json:"team"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return nil, err
	}
	return result.Team.GitAutomationStates.Nodes, nil
}

// AuditEntry представляет запись аудит-лога Linear.
type AuditEntry struct {
	ID          string         `json:"id"`
	Type        string         `json:"type"`
	ActorID     string         `json:"actorId"`
	CreatedAt   string         `json:"createdAt"`
	IP          string         `json:"ip"`
	CountryCode string         `json:"countryCode"`
	Metadata    map[string]any `json:"metadata"`
}

// AuditEntryType представляет тип записи аудит-лога.
type AuditEntryType struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// AuditEntryFilter — параметры фильтрации аудит-лога.
type AuditEntryFilter struct {
	Type  string
	Limit int
	After string
}

// ListAuditEntries возвращает список записей аудит-лога.
func (c *Client) ListAuditEntries(filter AuditEntryFilter) ([]AuditEntry, *PageInfo, error) {
	query := `
query ListAuditEntries($first: Int, $after: String, $filter: AuditEntryFilter) {
  auditEntries(first: $first, after: $after, filter: $filter) {
    nodes {
      id
      type
      actorId
      createdAt
      ip
      countryCode
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
	if filter.Type != "" {
		variables["filter"] = map[string]any{"type": map[string]any{"eq": filter.Type}}
	}
	var result struct {
		AuditEntries struct {
			Nodes    []AuditEntry `json:"nodes"`
			PageInfo PageInfo     `json:"pageInfo"`
		} `json:"auditEntries"`
	}
	if err := c.Do(query, variables, &result); err != nil {
		return nil, nil, err
	}
	return result.AuditEntries.Nodes, &result.AuditEntries.PageInfo, nil
}

// RateLimitResultPayload — лимит одного типа запросов.
type RateLimitResultPayload struct {
	Type            string  `json:"type"`
	AllowedAmount   float64 `json:"allowedAmount"`
	RequestedAmount float64 `json:"requestedAmount"`
	RemainingAmount float64 `json:"remainingAmount"`
	Period          float64 `json:"period"`
	Reset           string  `json:"reset"`
}

// RateLimitPayload — информация о лимите запросов API.
type RateLimitPayload struct {
	Identifier string                   `json:"identifier"`
	Kind       string                   `json:"kind"`
	Limits     []RateLimitResultPayload `json:"limits"`
}

// GetRateLimitStatus возвращает текущий статус rate limit.
func (c *Client) GetRateLimitStatus() (*RateLimitPayload, error) {
	query := `
query {
  rateLimitStatus {
    identifier
    kind
    limits {
      type
      allowedAmount
      requestedAmount
      remainingAmount
      period
      reset
    }
  }
}`
	var result struct {
		RateLimitStatus RateLimitPayload `json:"rateLimitStatus"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return &result.RateLimitStatus, nil
}

// TimeSchedule представляет расписание для команды.
type TimeSchedule struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ListTimeSchedules возвращает список расписаний.
func (c *Client) ListTimeSchedules() ([]TimeSchedule, error) {
	query := `
query {
  timeSchedules {
    nodes {
      id
      name
      createdAt
      updatedAt
    }
  }
}`
	var result struct {
		TimeSchedules struct {
			Nodes []TimeSchedule `json:"nodes"`
		} `json:"timeSchedules"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.TimeSchedules.Nodes, nil
}

// TriageResponsibility представляет ответственность за триаж задач.
type TriageResponsibility struct {
	ID          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	Action      string `json:"action"`
	Team        Team   `json:"team"`
	CurrentUser *User  `json:"currentUser"`
}

// ListTriageResponsibilities возвращает ответственность за триаж для команды.
func (c *Client) ListTriageResponsibilities(teamID string) ([]TriageResponsibility, error) {
	query := `
query ListTriageResponsibilities($teamId: String!) {
  team(id: $teamId) {
    triageResponsibility {
      id
      createdAt
      action
      team { id key name }
      currentUser { id name displayName email }
    }
  }
}`
	var result struct {
		Team struct {
			TriageResponsibility *TriageResponsibility `json:"triageResponsibility"`
		} `json:"team"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return nil, err
	}
	if result.Team.TriageResponsibility == nil {
		return []TriageResponsibility{}, nil
	}
	return []TriageResponsibility{*result.Team.TriageResponsibility}, nil
}

// ListAuditEntryTypes возвращает список типов аудит-лога.
func (c *Client) ListAuditEntryTypes() ([]AuditEntryType, error) {
	query := `
query {
  auditEntryTypes {
    type
    description
  }
}`
	var result struct {
		AuditEntryTypes []AuditEntryType `json:"auditEntryTypes"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	return result.AuditEntryTypes, nil
}
