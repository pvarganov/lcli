package client

import (
	"fmt"
	"time"
)

// CreateComment добавляет комментарий к задаче.
func (c *Client) CreateComment(issueID, body string) (*Comment, error) {
	mutation := `
mutation CreateComment($input: CommentCreateInput!) {
  commentCreate(input: $input) {
    success
    comment {
      id
      body
      createdAt
      user { id name displayName email }
    }
  }
}`

	var result struct {
		CommentCreate struct {
			Comment Comment `json:"comment"`
			Success bool    `json:"success"`
		} `json:"commentCreate"`
	}
	if err := c.Do(mutation, map[string]any{
		"input": map[string]any{
			"issueId": issueID,
			"body":    body,
		},
	}, &result); err != nil {
		return nil, err
	}
	if !result.CommentCreate.Success {
		return nil, fmt.Errorf("commentCreate вернул success=false")
	}
	return &result.CommentCreate.Comment, nil
}

// CreateIssueInput — входные данные для создания задачи.
type CreateIssueInput struct {
	TeamID      string
	Title       string
	Description string
	AssigneeID  string
	Priority    int
}

// UpdateIssueInput — входные данные для обновления задачи.
type UpdateIssueInput struct {
	Title      string
	StateID    string
	AssigneeID string
	Priority   *int
}

type createIssueResult struct {
	IssueCreate struct {
		Issue   Issue `json:"issue"`
		Success bool  `json:"success"`
	} `json:"issueCreate"`
}

type updateIssueResult struct {
	IssueUpdate struct {
		Issue   Issue `json:"issue"`
		Success bool  `json:"success"`
	} `json:"issueUpdate"`
}

// CreateIssue создаёт новую задачу в Linear.
func (c *Client) CreateIssue(input CreateIssueInput) (*Issue, error) {
	mutation := `
mutation CreateIssue($input: IssueCreateInput!) {
  issueCreate(input: $input) {
    success
    issue {
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
  }
}`

	gqlInput := map[string]any{
		"teamId": input.TeamID,
		"title":  input.Title,
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	if input.AssigneeID != "" {
		gqlInput["assigneeId"] = input.AssigneeID
	}
	if input.Priority != 0 {
		gqlInput["priority"] = input.Priority
	}

	var result createIssueResult
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.IssueCreate.Success {
		return nil, fmt.Errorf("issueCreate вернул success=false")
	}
	return &result.IssueCreate.Issue, nil
}

// UpdateIssue обновляет существующую задачу в Linear.
func (c *Client) UpdateIssue(id string, input UpdateIssueInput) (*Issue, error) {
	mutation := `
mutation UpdateIssue($id: String!, $input: IssueUpdateInput!) {
  issueUpdate(id: $id, input: $input) {
    success
    issue {
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
  }
}`

	gqlInput := map[string]any{}
	if input.Title != "" {
		gqlInput["title"] = input.Title
	}
	if input.StateID != "" {
		gqlInput["stateId"] = input.StateID
	}
	if input.AssigneeID != "" {
		gqlInput["assigneeId"] = input.AssigneeID
	}
	if input.Priority != nil {
		gqlInput["priority"] = *input.Priority
	}

	var result updateIssueResult
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.IssueUpdate.Success {
		return nil, fmt.Errorf("issueUpdate вернул success=false")
	}
	return &result.IssueUpdate.Issue, nil
}

// GetTeamByKey возвращает команду по её ключу (например, "ENG").
func (c *Client) GetTeamByKey(key string) (*Team, error) {
	query := `
query GetTeams {
  teams(first: 250) {
    nodes {
      id
      key
      name
    }
  }
}`

	var result struct {
		Teams struct {
			Nodes []Team `json:"nodes"`
		} `json:"teams"`
	}
	if err := c.Do(query, nil, &result); err != nil {
		return nil, err
	}
	for _, t := range result.Teams.Nodes {
		if t.Key == key {
			return &t, nil
		}
	}
	return nil, nil
}

// FindWorkflowStateByName ищет состояние задачи по имени в рамках команды.
func (c *Client) FindWorkflowStateByName(teamID, stateName string) (string, error) {
	query := `
query GetWorkflowStates($teamId: ID!) {
  workflowStates(filter: { team: { id: { eq: $teamId } } }) {
    nodes {
      id
      name
    }
  }
}`

	var result struct {
		WorkflowStates struct {
			Nodes []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"nodes"`
		} `json:"workflowStates"`
	}
	if err := c.Do(query, map[string]any{"teamId": teamID}, &result); err != nil {
		return "", err
	}
	for _, s := range result.WorkflowStates.Nodes {
		if s.Name == stateName {
			return s.ID, nil
		}
	}
	return "", nil
}

// CreateIssueLabelInput — входные данные для создания метки.
type CreateIssueLabelInput struct {
	Name        string
	Color       string
	Description string
	TeamID      string
}

// UpdateIssueLabelInput — входные данные для обновления метки.
type UpdateIssueLabelInput struct {
	Name        string
	Color       string
	Description string
}

// CreateIssueLabel создаёт новую метку задачи.
func (c *Client) CreateIssueLabel(input CreateIssueLabelInput) (*IssueLabel, error) {
	mutation := `
mutation CreateIssueLabel($input: IssueLabelCreateInput!) {
  issueLabelCreate(input: $input) {
    success
    issueLabel {
      id
      name
      color
      description
    }
  }
}`
	gqlInput := map[string]any{
		"name":  input.Name,
		"color": input.Color,
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	if input.TeamID != "" {
		gqlInput["teamId"] = input.TeamID
	}
	var result struct {
		IssueLabelCreate struct {
			IssueLabel IssueLabel `json:"issueLabel"`
			Success    bool       `json:"success"`
		} `json:"issueLabelCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.IssueLabelCreate.Success {
		return nil, fmt.Errorf("issueLabelCreate вернул success=false")
	}
	return &result.IssueLabelCreate.IssueLabel, nil
}

// UpdateIssueLabel обновляет существующую метку задачи.
func (c *Client) UpdateIssueLabel(id string, input UpdateIssueLabelInput) (*IssueLabel, error) {
	mutation := `
mutation UpdateIssueLabel($id: String!, $input: IssueLabelUpdateInput!) {
  issueLabelUpdate(id: $id, input: $input) {
    success
    issueLabel {
      id
      name
      color
      description
    }
  }
}`
	gqlInput := map[string]any{}
	if input.Name != "" {
		gqlInput["name"] = input.Name
	}
	if input.Color != "" {
		gqlInput["color"] = input.Color
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		IssueLabelUpdate struct {
			IssueLabel IssueLabel `json:"issueLabel"`
			Success    bool       `json:"success"`
		} `json:"issueLabelUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.IssueLabelUpdate.Success {
		return nil, fmt.Errorf("issueLabelUpdate вернул success=false")
	}
	return &result.IssueLabelUpdate.IssueLabel, nil
}

// DeleteIssueLabel удаляет метку задачи по ID.
func (c *Client) DeleteIssueLabel(id string) error {
	mutation := `
mutation DeleteIssueLabel($id: String!) {
  issueLabelDelete(id: $id) {
    success
  }
}`
	var result struct {
		IssueLabelDelete struct {
			Success bool `json:"success"`
		} `json:"issueLabelDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueLabelDelete.Success {
		return fmt.Errorf("issueLabelDelete вернул success=false")
	}
	return nil
}

// getIssueLabelIDs возвращает текущие ID меток задачи.
func (c *Client) getIssueLabelIDs(issueID string) ([]string, error) {
	query := `
query GetIssueLabelIDs($id: String!) {
  issue(id: $id) {
    labels(first: 250) {
      nodes {
        id
      }
    }
  }
}`
	var result struct {
		Issue *struct {
			Labels struct {
				Nodes []struct {
					ID string `json:"id"`
				} `json:"nodes"`
			} `json:"labels"`
		} `json:"issue"`
	}
	if err := c.Do(query, map[string]any{"id": issueID}, &result); err != nil {
		return nil, err
	}
	if result.Issue == nil {
		return nil, fmt.Errorf("задача не найдена: %s", issueID)
	}
	ids := make([]string, 0, len(result.Issue.Labels.Nodes))
	for _, n := range result.Issue.Labels.Nodes {
		ids = append(ids, n.ID)
	}
	return ids, nil
}

// updateIssueLabelIDs обновляет метки задачи через issueUpdate.
func (c *Client) updateIssueLabelIDs(issueID string, labelIDs []string) error {
	mutation := `
mutation UpdateIssueLabels($id: String!, $input: IssueUpdateInput!) {
  issueUpdate(id: $id, input: $input) {
    success
  }
}`
	var result struct {
		IssueUpdate struct {
			Success bool `json:"success"`
		} `json:"issueUpdate"`
	}
	if err := c.Do(mutation, map[string]any{
		"id":    issueID,
		"input": map[string]any{"labelIds": labelIDs},
	}, &result); err != nil {
		return err
	}
	if !result.IssueUpdate.Success {
		return fmt.Errorf("issueUpdate вернул success=false")
	}
	return nil
}

// AddLabelToIssue добавляет метку к задаче по ID метки.
func (c *Client) AddLabelToIssue(issueID, labelID string) error {
	currentIDs, err := c.getIssueLabelIDs(issueID)
	if err != nil {
		return err
	}
	for _, id := range currentIDs {
		if id == labelID {
			return nil
		}
	}
	return c.updateIssueLabelIDs(issueID, append(currentIDs, labelID))
}

// RemoveLabelFromIssue удаляет метку из задачи по ID метки.
func (c *Client) RemoveLabelFromIssue(issueID, labelID string) error {
	currentIDs, err := c.getIssueLabelIDs(issueID)
	if err != nil {
		return err
	}
	newIDs := make([]string, 0, len(currentIDs))
	for _, id := range currentIDs {
		if id != labelID {
			newIDs = append(newIDs, id)
		}
	}
	return c.updateIssueLabelIDs(issueID, newIDs)
}

// CreateIssueRelationInput — входные данные для создания связи.
type CreateIssueRelationInput struct {
	IssueID        string
	RelatedIssueID string
	Type           string // blocks, blocked_by, related, duplicate, duplicate_of
}

// CreateIssueRelation создаёт связь между двумя задачами.
func (c *Client) CreateIssueRelation(input CreateIssueRelationInput) (*IssueRelation, error) {
	mutation := `
mutation CreateIssueRelation($input: IssueRelationCreateInput!) {
  issueRelationCreate(input: $input) {
    success
    issueRelation {
      id
      type
      relatedIssue { id identifier title }
    }
  }
}`
	var result struct {
		IssueRelationCreate struct {
			IssueRelation IssueRelation `json:"issueRelation"`
			Success       bool          `json:"success"`
		} `json:"issueRelationCreate"`
	}
	if err := c.Do(mutation, map[string]any{
		"input": map[string]any{
			"issueId":        input.IssueID,
			"relatedIssueId": input.RelatedIssueID,
			"type":           input.Type,
		},
	}, &result); err != nil {
		return nil, err
	}
	if !result.IssueRelationCreate.Success {
		return nil, fmt.Errorf("issueRelationCreate вернул success=false")
	}
	return &result.IssueRelationCreate.IssueRelation, nil
}

// DeleteIssueRelation удаляет связь между задачами по ID связи.
func (c *Client) DeleteIssueRelation(id string) error {
	mutation := `
mutation DeleteIssueRelation($id: String!) {
  issueRelationDelete(id: $id) {
    success
  }
}`
	var result struct {
		IssueRelationDelete struct {
			Success bool `json:"success"`
		} `json:"issueRelationDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueRelationDelete.Success {
		return fmt.Errorf("issueRelationDelete вернул success=false")
	}
	return nil
}

// ArchiveIssue архивирует задачу по ID.
func (c *Client) ArchiveIssue(id string) error {
	mutation := `
mutation ArchiveIssue($id: String!) {
  issueArchive(id: $id) {
    success
  }
}`
	var result struct {
		IssueArchive struct {
			Success bool `json:"success"`
		} `json:"issueArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueArchive.Success {
		return fmt.Errorf("issueArchive вернул success=false")
	}
	return nil
}

// UnarchiveIssue разархивирует задачу по ID.
func (c *Client) UnarchiveIssue(id string) error {
	mutation := `
mutation UnarchiveIssue($id: String!) {
  issueUnarchive(id: $id) {
    success
  }
}`
	var result struct {
		IssueUnarchive struct {
			Success bool `json:"success"`
		} `json:"issueUnarchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueUnarchive.Success {
		return fmt.Errorf("issueUnarchive вернул success=false")
	}
	return nil
}

// DeleteIssue удаляет задачу по ID.
func (c *Client) DeleteIssue(id string) error {
	mutation := `
mutation DeleteIssue($id: String!) {
  issueDelete(id: $id) {
    success
  }
}`
	var result struct {
		IssueDelete struct {
			Success bool `json:"success"`
		} `json:"issueDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueDelete.Success {
		return fmt.Errorf("issueDelete вернул success=false")
	}
	return nil
}

// SubscribeToIssue подписывает текущего пользователя на задачу.
func (c *Client) SubscribeToIssue(id string) error {
	mutation := `
mutation SubscribeToIssue($id: String!) {
  issueSubscribe(id: $id) {
    success
  }
}`
	var result struct {
		IssueSubscribe struct {
			Success bool `json:"success"`
		} `json:"issueSubscribe"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueSubscribe.Success {
		return fmt.Errorf("issueSubscribe вернул success=false")
	}
	return nil
}

// UnsubscribeFromIssue отписывает текущего пользователя от задачи.
func (c *Client) UnsubscribeFromIssue(id string) error {
	mutation := `
mutation UnsubscribeFromIssue($id: String!) {
  issueUnsubscribe(id: $id) {
    success
  }
}`
	var result struct {
		IssueUnsubscribe struct {
			Success bool `json:"success"`
		} `json:"issueUnsubscribe"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.IssueUnsubscribe.Success {
		return fmt.Errorf("issueUnsubscribe вернул success=false")
	}
	return nil
}

// BatchCreateIssuesInput — входные данные для пакетного создания задач.
type BatchCreateIssuesInput struct {
	Issues []CreateIssueInput
}

// BatchCreateIssues создаёт несколько задач за одну транзакцию.
func (c *Client) BatchCreateIssues(input BatchCreateIssuesInput) ([]Issue, error) {
	mutation := `
mutation IssueBatchCreate($input: IssueBatchCreateInput!) {
  issueBatchCreate(input: $input) {
    success
    issues {
      id
      identifier
      title
      priority
      state { name type }
      team { id key name }
    }
  }
}`
	issues := make([]map[string]any, 0, len(input.Issues))
	for _, i := range input.Issues {
		gi := map[string]any{
			"teamId": i.TeamID,
			"title":  i.Title,
		}
		if i.Description != "" {
			gi["description"] = i.Description
		}
		if i.AssigneeID != "" {
			gi["assigneeId"] = i.AssigneeID
		}
		if i.Priority != 0 {
			gi["priority"] = i.Priority
		}
		issues = append(issues, gi)
	}
	var result struct {
		IssueBatchCreate struct {
			Issues  []Issue `json:"issues"`
			Success bool    `json:"success"`
		} `json:"issueBatchCreate"`
	}
	if err := c.Do(mutation, map[string]any{
		"input": map[string]any{"issues": issues},
	}, &result); err != nil {
		return nil, err
	}
	if !result.IssueBatchCreate.Success {
		return nil, fmt.Errorf("issueBatchCreate вернул success=false")
	}
	return result.IssueBatchCreate.Issues, nil
}

// BatchUpdateIssuesInput — входные данные для пакетного обновления задач.
type BatchUpdateIssuesInput struct {
	IDs    []string
	Update UpdateIssueInput
}

// BatchUpdateIssues обновляет несколько задач за один запрос (до 50 за раз).
func (c *Client) BatchUpdateIssues(input BatchUpdateIssuesInput) ([]Issue, error) {
	mutation := `
mutation IssueBatchUpdate($ids: [UUID!]!, $input: IssueUpdateInput!) {
  issueBatchUpdate(ids: $ids, input: $input) {
    success
    issues {
      id
      identifier
      title
      priority
      state { name type }
      team { id key name }
    }
  }
}`
	gqlInput := map[string]any{}
	if input.Update.Title != "" {
		gqlInput["title"] = input.Update.Title
	}
	if input.Update.StateID != "" {
		gqlInput["stateId"] = input.Update.StateID
	}
	if input.Update.AssigneeID != "" {
		gqlInput["assigneeId"] = input.Update.AssigneeID
	}
	if input.Update.Priority != nil {
		gqlInput["priority"] = *input.Update.Priority
	}
	var result struct {
		IssueBatchUpdate struct {
			Issues  []Issue `json:"issues"`
			Success bool    `json:"success"`
		} `json:"issueBatchUpdate"`
	}
	if err := c.Do(mutation, map[string]any{
		"ids":   input.IDs,
		"input": gqlInput,
	}, &result); err != nil {
		return nil, err
	}
	if !result.IssueBatchUpdate.Success {
		return nil, fmt.Errorf("issueBatchUpdate вернул success=false")
	}
	return result.IssueBatchUpdate.Issues, nil
}

// UpdateComment обновляет тело комментария по ID.
func (c *Client) UpdateComment(id, body string) (*Comment, error) {
	mutation := `
mutation UpdateComment($id: String!, $input: CommentUpdateInput!) {
  commentUpdate(id: $id, input: $input) {
    success
    comment {
      id
      body
      createdAt
      user { id name displayName email }
    }
  }
}`
	var result struct {
		CommentUpdate struct {
			Comment Comment `json:"comment"`
			Success bool    `json:"success"`
		} `json:"commentUpdate"`
	}
	if err := c.Do(mutation, map[string]any{
		"id":    id,
		"input": map[string]any{"body": body},
	}, &result); err != nil {
		return nil, err
	}
	if !result.CommentUpdate.Success {
		return nil, fmt.Errorf("commentUpdate вернул success=false")
	}
	return &result.CommentUpdate.Comment, nil
}

// DeleteComment удаляет комментарий по ID.
func (c *Client) DeleteComment(id string) error {
	mutation := `
mutation DeleteComment($id: String!) {
  commentDelete(id: $id) {
    success
  }
}`
	var result struct {
		CommentDelete struct {
			Success bool `json:"success"`
		} `json:"commentDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CommentDelete.Success {
		return fmt.Errorf("commentDelete вернул success=false")
	}
	return nil
}

// ResolveComment помечает комментарий как разрешённый.
func (c *Client) ResolveComment(id string) (*Comment, error) {
	mutation := `
mutation ResolveComment($id: String!) {
  commentResolve(id: $id) {
    success
    comment {
      id
      body
      createdAt
      user { id name displayName email }
    }
  }
}`
	var result struct {
		CommentResolve struct {
			Comment Comment `json:"comment"`
			Success bool    `json:"success"`
		} `json:"commentResolve"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if !result.CommentResolve.Success {
		return nil, fmt.Errorf("commentResolve вернул success=false")
	}
	return &result.CommentResolve.Comment, nil
}

// UnresolveComment снимает пометку "разрешён" с комментария.
func (c *Client) UnresolveComment(id string) (*Comment, error) {
	mutation := `
mutation UnresolveComment($id: String!) {
  commentUnresolve(id: $id) {
    success
    comment {
      id
      body
      createdAt
      user { id name displayName email }
    }
  }
}`
	var result struct {
		CommentUnresolve struct {
			Comment Comment `json:"comment"`
			Success bool    `json:"success"`
		} `json:"commentUnresolve"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if !result.CommentUnresolve.Success {
		return nil, fmt.Errorf("commentUnresolve вернул success=false")
	}
	return &result.CommentUnresolve.Comment, nil
}

// FindUserByName ищет пользователя по displayName или email.
func (c *Client) FindUserByName(name string) (*User, error) {
	query := `
query GetUsers {
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
	for _, u := range result.Users.Nodes {
		if u.DisplayName == name || u.Email == name || u.Name == name {
			return &u, nil
		}
	}
	return nil, nil
}

// CreateProjectInput — входные данные для создания проекта.
type CreateProjectInput struct {
	Name        string
	Description string
	TeamIDs     []string
	StartDate   string
	TargetDate  string
	LeadID      string
}

// UpdateProjectInput — входные данные для обновления проекта.
type UpdateProjectInput struct {
	Name        string
	Description string
	StartDate   string
	TargetDate  string
	LeadID      string
	State       string
}

// CreateProject создаёт новый проект в Linear.
func (c *Client) CreateProject(input CreateProjectInput) (*Project, error) {
	mutation := `
mutation CreateProject($input: ProjectCreateInput!) {
  projectCreate(input: $input) {
    success
    project {
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
	gqlInput := map[string]any{
		"name":    input.Name,
		"teamIds": input.TeamIDs,
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	if input.StartDate != "" {
		gqlInput["startDate"] = input.StartDate
	}
	if input.TargetDate != "" {
		gqlInput["targetDate"] = input.TargetDate
	}
	if input.LeadID != "" {
		gqlInput["leadId"] = input.LeadID
	}

	var result struct {
		ProjectCreate struct {
			Project Project `json:"project"`
			Success bool    `json:"success"`
		} `json:"projectCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectCreate.Success {
		return nil, fmt.Errorf("projectCreate вернул success=false")
	}
	return &result.ProjectCreate.Project, nil
}

// UpdateProject обновляет существующий проект в Linear.
func (c *Client) UpdateProject(id string, input UpdateProjectInput) (*Project, error) {
	mutation := `
mutation UpdateProject($id: String!, $input: ProjectUpdateInput!) {
  projectUpdate(id: $id, input: $input) {
    success
    project {
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
	gqlInput := map[string]any{}
	if input.Name != "" {
		gqlInput["name"] = input.Name
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	if input.StartDate != "" {
		gqlInput["startDate"] = input.StartDate
	}
	if input.TargetDate != "" {
		gqlInput["targetDate"] = input.TargetDate
	}
	if input.LeadID != "" {
		gqlInput["leadId"] = input.LeadID
	}
	if input.State != "" {
		gqlInput["state"] = input.State
	}

	var result struct {
		ProjectUpdate struct {
			Project Project `json:"project"`
			Success bool    `json:"success"`
		} `json:"projectUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectUpdate.Success {
		return nil, fmt.Errorf("projectUpdate вернул success=false")
	}
	return &result.ProjectUpdate.Project, nil
}

// DeleteProject удаляет проект из Linear.
func (c *Client) DeleteProject(id string) error {
	mutation := `
mutation DeleteProject($id: String!) {
  projectDelete(id: $id) {
    success
  }
}`
	var result struct {
		ProjectDelete struct {
			Success bool `json:"success"`
		} `json:"projectDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectDelete.Success {
		return fmt.Errorf("projectDelete вернул success=false")
	}
	return nil
}

// ArchiveProject архивирует проект в Linear.
func (c *Client) ArchiveProject(id string) error {
	mutation := `
mutation ArchiveProject($id: String!) {
  projectArchive(id: $id) {
    success
  }
}`
	var result struct {
		ProjectArchive struct {
			Success bool `json:"success"`
		} `json:"projectArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectArchive.Success {
		return fmt.Errorf("projectArchive вернул success=false")
	}
	return nil
}

// CreateProjectMilestoneInput — входные данные для создания вехи проекта.
type CreateProjectMilestoneInput struct {
	ProjectID   string
	Name        string
	TargetDate  string
	Description string
}

// UpdateProjectMilestoneInput — входные данные для обновления вехи проекта.
type UpdateProjectMilestoneInput struct {
	Name        string
	TargetDate  string
	Description string
}

// CreateProjectMilestone создаёт новую веху проекта.
func (c *Client) CreateProjectMilestone(input CreateProjectMilestoneInput) (*ProjectMilestone, error) {
	mutation := `
mutation CreateProjectMilestone($input: ProjectMilestoneCreateInput!) {
  projectMilestoneCreate(input: $input) {
    success
    projectMilestone {
      id
      name
      targetDate
      description
    }
  }
}`
	gqlInput := map[string]any{
		"projectId": input.ProjectID,
		"name":      input.Name,
	}
	if input.TargetDate != "" {
		gqlInput["targetDate"] = input.TargetDate
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		ProjectMilestoneCreate struct {
			ProjectMilestone ProjectMilestone `json:"projectMilestone"`
			Success          bool             `json:"success"`
		} `json:"projectMilestoneCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectMilestoneCreate.Success {
		return nil, fmt.Errorf("projectMilestoneCreate вернул success=false")
	}
	return &result.ProjectMilestoneCreate.ProjectMilestone, nil
}

// UpdateProjectMilestone обновляет существующую веху проекта.
func (c *Client) UpdateProjectMilestone(id string, input UpdateProjectMilestoneInput) (*ProjectMilestone, error) {
	mutation := `
mutation UpdateProjectMilestone($id: String!, $input: ProjectMilestoneUpdateInput!) {
  projectMilestoneUpdate(id: $id, input: $input) {
    success
    projectMilestone {
      id
      name
      targetDate
      description
    }
  }
}`
	gqlInput := map[string]any{}
	if input.Name != "" {
		gqlInput["name"] = input.Name
	}
	if input.TargetDate != "" {
		gqlInput["targetDate"] = input.TargetDate
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		ProjectMilestoneUpdate struct {
			ProjectMilestone ProjectMilestone `json:"projectMilestone"`
			Success          bool             `json:"success"`
		} `json:"projectMilestoneUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectMilestoneUpdate.Success {
		return nil, fmt.Errorf("projectMilestoneUpdate вернул success=false")
	}
	return &result.ProjectMilestoneUpdate.ProjectMilestone, nil
}

// DeleteProjectMilestone удаляет веху проекта по ID.
func (c *Client) DeleteProjectMilestone(id string) error {
	mutation := `
mutation DeleteProjectMilestone($id: String!) {
  projectMilestoneDelete(id: $id) {
    success
  }
}`
	var result struct {
		ProjectMilestoneDelete struct {
			Success bool `json:"success"`
		} `json:"projectMilestoneDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectMilestoneDelete.Success {
		return fmt.Errorf("projectMilestoneDelete вернул success=false")
	}
	return nil
}

// CreateProjectUpdateInput — входные данные для создания обновления проекта.
type CreateProjectUpdateInput struct {
	ProjectID string
	Body      string
	Health    string
}

// UpdateProjectUpdateInput — входные данные для обновления записи журнала проекта.
type UpdateProjectUpdateInput struct {
	Body   string
	Health string
}

// CreateProjectUpdate создаёт новую запись журнала проекта.
func (c *Client) CreateProjectUpdate(input CreateProjectUpdateInput) (*ProjectUpdate, error) {
	mutation := `
mutation CreateProjectUpdate($input: ProjectUpdateCreateInput!) {
  projectUpdateCreate(input: $input) {
    success
    projectUpdate {
      id
      body
      health
      createdAt
      user { id name displayName email }
    }
  }
}`
	gqlInput := map[string]any{
		"projectId": input.ProjectID,
		"body":      input.Body,
	}
	if input.Health != "" {
		gqlInput["health"] = input.Health
	}
	var result struct {
		ProjectUpdateCreate struct {
			ProjectUpdate ProjectUpdate `json:"projectUpdate"`
			Success       bool          `json:"success"`
		} `json:"projectUpdateCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectUpdateCreate.Success {
		return nil, fmt.Errorf("projectUpdateCreate вернул success=false")
	}
	return &result.ProjectUpdateCreate.ProjectUpdate, nil
}

// UpdateProjectUpdate обновляет запись журнала проекта.
func (c *Client) UpdateProjectUpdate(id string, input UpdateProjectUpdateInput) (*ProjectUpdate, error) {
	mutation := `
mutation UpdateProjectUpdate($id: String!, $input: ProjectUpdateUpdateInput!) {
  projectUpdateUpdate(id: $id, input: $input) {
    success
    projectUpdate {
      id
      body
      health
      createdAt
      user { id name displayName email }
    }
  }
}`
	gqlInput := map[string]any{}
	if input.Body != "" {
		gqlInput["body"] = input.Body
	}
	if input.Health != "" {
		gqlInput["health"] = input.Health
	}
	var result struct {
		ProjectUpdateUpdate struct {
			ProjectUpdate ProjectUpdate `json:"projectUpdate"`
			Success       bool          `json:"success"`
		} `json:"projectUpdateUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectUpdateUpdate.Success {
		return nil, fmt.Errorf("projectUpdateUpdate вернул success=false")
	}
	return &result.ProjectUpdateUpdate.ProjectUpdate, nil
}

// ArchiveProjectUpdate архивирует запись журнала проекта.
func (c *Client) ArchiveProjectUpdate(id string) error {
	mutation := `
mutation ArchiveProjectUpdate($id: String!) {
  projectUpdateArchive(id: $id) {
    success
  }
}`
	var result struct {
		ProjectUpdateArchive struct {
			Success bool `json:"success"`
		} `json:"projectUpdateArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectUpdateArchive.Success {
		return fmt.Errorf("projectUpdateArchive вернул success=false")
	}
	return nil
}

// CreateProjectLabelInput — входные данные для создания метки проекта.
type CreateProjectLabelInput struct {
	Name        string
	Color       string
	Description string
}

// UpdateProjectLabelInput — входные данные для обновления метки проекта.
type UpdateProjectLabelInput struct {
	Name        string
	Color       string
	Description string
}

// CreateProjectLabel создаёт новую метку проекта.
func (c *Client) CreateProjectLabel(input CreateProjectLabelInput) (*ProjectLabel, error) {
	mutation := `
mutation CreateProjectLabel($input: ProjectLabelCreateInput!) {
  projectLabelCreate(input: $input) {
    success
    projectLabel {
      id
      name
      color
      description
    }
  }
}`
	gqlInput := map[string]any{
		"name": input.Name,
	}
	if input.Color != "" {
		gqlInput["color"] = input.Color
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		ProjectLabelCreate struct {
			ProjectLabel ProjectLabel `json:"projectLabel"`
			Success      bool         `json:"success"`
		} `json:"projectLabelCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectLabelCreate.Success {
		return nil, fmt.Errorf("projectLabelCreate вернул success=false")
	}
	return &result.ProjectLabelCreate.ProjectLabel, nil
}

// UpdateProjectLabel обновляет метку проекта.
func (c *Client) UpdateProjectLabel(id string, input UpdateProjectLabelInput) (*ProjectLabel, error) {
	mutation := `
mutation UpdateProjectLabel($id: String!, $input: ProjectLabelUpdateInput!) {
  projectLabelUpdate(id: $id, input: $input) {
    success
    projectLabel {
      id
      name
      color
      description
    }
  }
}`
	gqlInput := map[string]any{}
	if input.Name != "" {
		gqlInput["name"] = input.Name
	}
	if input.Color != "" {
		gqlInput["color"] = input.Color
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		ProjectLabelUpdate struct {
			ProjectLabel ProjectLabel `json:"projectLabel"`
			Success      bool         `json:"success"`
		} `json:"projectLabelUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectLabelUpdate.Success {
		return nil, fmt.Errorf("projectLabelUpdate вернул success=false")
	}
	return &result.ProjectLabelUpdate.ProjectLabel, nil
}

// DeleteProjectLabel удаляет метку проекта.
func (c *Client) DeleteProjectLabel(id string) error {
	mutation := `
mutation DeleteProjectLabel($id: String!) {
  projectLabelDelete(id: $id) {
    success
  }
}`
	var result struct {
		ProjectLabelDelete struct {
			Success bool `json:"success"`
		} `json:"projectLabelDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectLabelDelete.Success {
		return fmt.Errorf("projectLabelDelete вернул success=false")
	}
	return nil
}

// UnarchiveProject разархивирует проект в Linear.
func (c *Client) UnarchiveProject(id string) error {
	mutation := `
mutation UnarchiveProject($id: String!) {
  projectUnarchive(id: $id) {
    success
  }
}`
	var result struct {
		ProjectUnarchive struct {
			Success bool `json:"success"`
		} `json:"projectUnarchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectUnarchive.Success {
		return fmt.Errorf("projectUnarchive вернул success=false")
	}
	return nil
}

// CreateProjectStatusInput — входные данные для создания статуса проекта.
type CreateProjectStatusInput struct {
	Name        string
	Type        string
	Color       string
	Description string
	Position    float64
}

// UpdateProjectStatusInput — входные данные для обновления статуса проекта.
type UpdateProjectStatusInput struct {
	Name        string
	Color       string
	Description string
}

// CreateProjectStatus создаёт новый статус проекта.
func (c *Client) CreateProjectStatus(input CreateProjectStatusInput) (*ProjectStatus, error) {
	mutation := `
mutation CreateProjectStatus($input: ProjectStatusCreateInput!) {
  projectStatusCreate(input: $input) {
    success
    projectStatus {
      id
      name
      type
      color
      description
      position
    }
  }
}`
	gqlInput := map[string]any{
		"name":     input.Name,
		"type":     input.Type,
		"color":    input.Color,
		"position": input.Position,
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		ProjectStatusCreate struct {
			ProjectStatus ProjectStatus `json:"projectStatus"`
			Success       bool          `json:"success"`
		} `json:"projectStatusCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectStatusCreate.Success {
		return nil, fmt.Errorf("projectStatusCreate вернул success=false")
	}
	return &result.ProjectStatusCreate.ProjectStatus, nil
}

// UpdateProjectStatus обновляет статус проекта.
func (c *Client) UpdateProjectStatus(id string, input UpdateProjectStatusInput) (*ProjectStatus, error) {
	mutation := `
mutation UpdateProjectStatus($id: String!, $input: ProjectStatusUpdateInput!) {
  projectStatusUpdate(id: $id, input: $input) {
    success
    projectStatus {
      id
      name
      type
      color
      description
      position
    }
  }
}`
	gqlInput := map[string]any{}
	if input.Name != "" {
		gqlInput["name"] = input.Name
	}
	if input.Color != "" {
		gqlInput["color"] = input.Color
	}
	if input.Description != "" {
		gqlInput["description"] = input.Description
	}
	var result struct {
		ProjectStatusUpdate struct {
			ProjectStatus ProjectStatus `json:"projectStatus"`
			Success       bool          `json:"success"`
		} `json:"projectStatusUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectStatusUpdate.Success {
		return nil, fmt.Errorf("projectStatusUpdate вернул success=false")
	}
	return &result.ProjectStatusUpdate.ProjectStatus, nil
}

// ArchiveProjectStatus архивирует статус проекта.
func (c *Client) ArchiveProjectStatus(id string) error {
	mutation := `
mutation ArchiveProjectStatus($id: String!) {
  projectStatusArchive(id: $id) {
    success
  }
}`
	var result struct {
		ProjectStatusArchive struct {
			Success bool `json:"success"`
		} `json:"projectStatusArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectStatusArchive.Success {
		return fmt.Errorf("projectStatusArchive вернул success=false")
	}
	return nil
}

// ProjectRelation представляет связь между проектами.
type ProjectRelation struct {
	ID             string  `json:"id"`
	Type           string  `json:"type"`
	Project        Project `json:"project"`
	RelatedProject Project `json:"relatedProject"`
}

// CreateProjectRelationInput — входные данные для создания связи между проектами.
type CreateProjectRelationInput struct {
	ProjectID        string
	RelatedProjectID string
	Type             string
}

// CreateProjectRelation создаёт связь между двумя проектами.
func (c *Client) CreateProjectRelation(input CreateProjectRelationInput) (*ProjectRelation, error) {
	mutation := `
mutation CreateProjectRelation($input: ProjectRelationCreateInput!) {
  projectRelationCreate(input: $input) {
    success
    projectRelation {
      id
      type
      project { id name }
      relatedProject { id name }
    }
  }
}`
	var result struct {
		ProjectRelationCreate struct {
			ProjectRelation ProjectRelation `json:"projectRelation"`
			Success         bool            `json:"success"`
		} `json:"projectRelationCreate"`
	}
	if err := c.Do(mutation, map[string]any{
		"input": map[string]any{
			"projectId":        input.ProjectID,
			"relatedProjectId": input.RelatedProjectID,
			"type":             input.Type,
		},
	}, &result); err != nil {
		return nil, err
	}
	if !result.ProjectRelationCreate.Success {
		return nil, fmt.Errorf("projectRelationCreate вернул success=false")
	}
	return &result.ProjectRelationCreate.ProjectRelation, nil
}

// DeleteProjectRelation удаляет связь между проектами по ID связи.
func (c *Client) DeleteProjectRelation(id string) error {
	mutation := `
mutation DeleteProjectRelation($id: String!) {
  projectRelationDelete(id: $id) {
    success
  }
}`
	var result struct {
		ProjectRelationDelete struct {
			Success bool `json:"success"`
		} `json:"projectRelationDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ProjectRelationDelete.Success {
		return fmt.Errorf("projectRelationDelete вернул success=false")
	}
	return nil
}

// CreateCycle создаёт новый цикл (спринт) для команды.
func (c *Client) CreateCycle(teamID, name, startsAt, endsAt string) (*Cycle, error) {
	mutation := `
mutation CreateCycle($input: CycleCreateInput!) {
  cycleCreate(input: $input) {
    success
    cycle {
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
	input := map[string]any{
		"teamId":   teamID,
		"startsAt": startsAt,
		"endsAt":   endsAt,
	}
	if name != "" {
		input["name"] = name
	}
	var result struct {
		CycleCreate struct {
			Cycle   Cycle `json:"cycle"`
			Success bool  `json:"success"`
		} `json:"cycleCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CycleCreate.Success {
		return nil, fmt.Errorf("cycleCreate вернул success=false")
	}
	return &result.CycleCreate.Cycle, nil
}

// UpdateCycle обновляет цикл по ID.
func (c *Client) UpdateCycle(id, name, startsAt, endsAt string) (*Cycle, error) {
	mutation := `
mutation UpdateCycle($id: String!, $input: CycleUpdateInput!) {
  cycleUpdate(id: $id, input: $input) {
    success
    cycle {
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
	input := map[string]any{}
	if name != "" {
		input["name"] = name
	}
	if startsAt != "" {
		input["startsAt"] = startsAt
	}
	if endsAt != "" {
		input["endsAt"] = endsAt
	}
	var result struct {
		CycleUpdate struct {
			Cycle   Cycle `json:"cycle"`
			Success bool  `json:"success"`
		} `json:"cycleUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CycleUpdate.Success {
		return nil, fmt.Errorf("cycleUpdate вернул success=false")
	}
	return &result.CycleUpdate.Cycle, nil
}

// ArchiveCycle архивирует цикл по ID.
func (c *Client) ArchiveCycle(id string) error {
	mutation := `
mutation ArchiveCycle($id: String!) {
  cycleArchive(id: $id) {
    success
  }
}`
	var result struct {
		CycleArchive struct {
			Success bool `json:"success"`
		} `json:"cycleArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CycleArchive.Success {
		return fmt.Errorf("cycleArchive вернул success=false")
	}
	return nil
}

// CreateWorkflowState создаёт новый статус задачи.
func (c *Client) CreateWorkflowState(teamID, name, stateType, color string) (*WorkflowState, error) {
	mutation := `
mutation CreateWorkflowState($input: WorkflowStateCreateInput!) {
  workflowStateCreate(input: $input) {
    success
    workflowState {
      id
      name
      type
      color
      team { id key name }
    }
  }
}`
	input := map[string]any{
		"teamId": teamID,
		"name":   name,
		"type":   stateType,
	}
	if color != "" {
		input["color"] = color
	}
	var result struct {
		WorkflowStateCreate struct {
			WorkflowState WorkflowState `json:"workflowState"`
			Success       bool          `json:"success"`
		} `json:"workflowStateCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.WorkflowStateCreate.Success {
		return nil, fmt.Errorf("workflowStateCreate вернул success=false")
	}
	return &result.WorkflowStateCreate.WorkflowState, nil
}

// UpdateWorkflowState обновляет статус задачи по ID.
func (c *Client) UpdateWorkflowState(id, name, color string) (*WorkflowState, error) {
	mutation := `
mutation UpdateWorkflowState($id: String!, $input: WorkflowStateUpdateInput!) {
  workflowStateUpdate(id: $id, input: $input) {
    success
    workflowState {
      id
      name
      type
      color
      team { id key name }
    }
  }
}`
	input := map[string]any{}
	if name != "" {
		input["name"] = name
	}
	if color != "" {
		input["color"] = color
	}
	var result struct {
		WorkflowStateUpdate struct {
			WorkflowState WorkflowState `json:"workflowState"`
			Success       bool          `json:"success"`
		} `json:"workflowStateUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.WorkflowStateUpdate.Success {
		return nil, fmt.Errorf("workflowStateUpdate вернул success=false")
	}
	return &result.WorkflowStateUpdate.WorkflowState, nil
}

// ArchiveWorkflowState архивирует статус задачи по ID.
func (c *Client) ArchiveWorkflowState(id string) error {
	mutation := `
mutation ArchiveWorkflowState($id: String!) {
  workflowStateArchive(id: $id) {
    success
  }
}`
	var result struct {
		WorkflowStateArchive struct {
			Success bool `json:"success"`
		} `json:"workflowStateArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.WorkflowStateArchive.Success {
		return fmt.Errorf("workflowStateArchive вернул success=false")
	}
	return nil
}

// CreateTeam создаёт новую команду.
func (c *Client) CreateTeam(name, key string) (*Team, error) {
	mutation := `
mutation CreateTeam($input: TeamCreateInput!) {
  teamCreate(input: $input) {
    success
    team {
      id
      key
      name
    }
  }
}`
	var result struct {
		TeamCreate struct {
			Team    Team `json:"team"`
			Success bool `json:"success"`
		} `json:"teamCreate"`
	}
	input := map[string]any{"name": name, "key": key}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.TeamCreate.Success {
		return nil, fmt.Errorf("teamCreate вернул success=false")
	}
	return &result.TeamCreate.Team, nil
}

// UpdateTeamInput — входные данные для обновления команды.
type UpdateTeamInput struct {
	Name string
	Key  string
}

// UpdateTeam обновляет команду.
func (c *Client) UpdateTeam(id string, input UpdateTeamInput) (*Team, error) {
	mutation := `
mutation UpdateTeam($id: String!, $input: TeamUpdateInput!) {
  teamUpdate(id: $id, input: $input) {
    success
    team {
      id
      key
      name
    }
  }
}`
	var result struct {
		TeamUpdate struct {
			Team    Team `json:"team"`
			Success bool `json:"success"`
		} `json:"teamUpdate"`
	}
	gqlInput := map[string]any{}
	if input.Name != "" {
		gqlInput["name"] = input.Name
	}
	if input.Key != "" {
		gqlInput["key"] = input.Key
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.TeamUpdate.Success {
		return nil, fmt.Errorf("teamUpdate вернул success=false")
	}
	return &result.TeamUpdate.Team, nil
}

// DeleteTeam удаляет команду по ID.
func (c *Client) DeleteTeam(id string) error {
	mutation := `
mutation DeleteTeam($id: String!) {
  teamDelete(id: $id) {
    success
  }
}`
	var result struct {
		TeamDelete struct {
			Success bool `json:"success"`
		} `json:"teamDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.TeamDelete.Success {
		return fmt.Errorf("teamDelete вернул success=false")
	}
	return nil
}

// CreateTeamMembership добавляет пользователя в команду.
func (c *Client) CreateTeamMembership(teamID, userID string) (*TeamMembership, error) {
	mutation := `
mutation CreateTeamMembership($input: TeamMembershipCreateInput!) {
  teamMembershipCreate(input: $input) {
    success
    teamMembership {
      id
      user { id name displayName email }
      team { id key name }
      role
    }
  }
}`
	var result struct {
		TeamMembershipCreate struct {
			TeamMembership TeamMembership `json:"teamMembership"`
			Success        bool           `json:"success"`
		} `json:"teamMembershipCreate"`
	}
	input := map[string]any{"teamId": teamID, "userId": userID}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.TeamMembershipCreate.Success {
		return nil, fmt.Errorf("teamMembershipCreate вернул success=false")
	}
	return &result.TeamMembershipCreate.TeamMembership, nil
}

// UpdateTeamMembership обновляет роль члена команды.
func (c *Client) UpdateTeamMembership(id, role string) (*TeamMembership, error) {
	mutation := `
mutation UpdateTeamMembership($id: String!, $input: TeamMembershipUpdateInput!) {
  teamMembershipUpdate(id: $id, input: $input) {
    success
    teamMembership {
      id
      user { id name displayName email }
      team { id key name }
      role
    }
  }
}`
	var result struct {
		TeamMembershipUpdate struct {
			TeamMembership TeamMembership `json:"teamMembership"`
			Success        bool           `json:"success"`
		} `json:"teamMembershipUpdate"`
	}
	gqlInput := map[string]any{}
	if role != "" {
		gqlInput["role"] = role
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": gqlInput}, &result); err != nil {
		return nil, err
	}
	if !result.TeamMembershipUpdate.Success {
		return nil, fmt.Errorf("teamMembershipUpdate вернул success=false")
	}
	return &result.TeamMembershipUpdate.TeamMembership, nil
}

// DeleteTeamMembership удаляет членство в команде по ID.
func (c *Client) DeleteTeamMembership(id string) error {
	mutation := `
mutation DeleteTeamMembership($id: String!) {
  teamMembershipDelete(id: $id) {
    success
  }
}`
	var result struct {
		TeamMembershipDelete struct {
			Success bool `json:"success"`
		} `json:"teamMembershipDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.TeamMembershipDelete.Success {
		return fmt.Errorf("teamMembershipDelete вернул success=false")
	}
	return nil
}

// MarkNotificationRead отмечает уведомление как прочитанное.
func (c *Client) MarkNotificationRead(id string) error {
	mutation := `
mutation MarkNotificationRead($id: String!, $input: NotificationUpdateInput!) {
  notificationUpdate(id: $id, input: $input) {
    success
  }
}`
	var result struct {
		NotificationUpdate struct {
			Success bool `json:"success"`
		} `json:"notificationUpdate"`
	}
	input := map[string]any{"readAt": time.Now().UTC().Format(time.RFC3339)}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return err
	}
	if !result.NotificationUpdate.Success {
		return fmt.Errorf("notificationUpdate вернул success=false")
	}
	return nil
}

// MarkAllNotificationsRead отмечает все уведомления как прочитанные.
func (c *Client) MarkAllNotificationsRead() error {
	mutation := `
mutation MarkAllNotificationsRead {
  notificationMarkReadAll {
    success
  }
}`
	var result struct {
		NotificationMarkReadAll struct {
			Success bool `json:"success"`
		} `json:"notificationMarkReadAll"`
	}
	if err := c.Do(mutation, nil, &result); err != nil {
		return err
	}
	if !result.NotificationMarkReadAll.Success {
		return fmt.Errorf("notificationMarkReadAll вернул success=false")
	}
	return nil
}

// ArchiveNotification архивирует уведомление по ID.
func (c *Client) ArchiveNotification(id string) error {
	mutation := `
mutation ArchiveNotification($id: String!) {
  notificationArchive(id: $id) {
    success
  }
}`
	var result struct {
		NotificationArchive struct {
			Success bool `json:"success"`
		} `json:"notificationArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.NotificationArchive.Success {
		return fmt.Errorf("notificationArchive вернул success=false")
	}
	return nil
}

// WebhookCreateInput — входные данные для создания вебхука.
type WebhookCreateInput struct {
	URL           string
	TeamID        string
	Enabled       bool
	Secret        string
	ResourceTypes []string
}

// CreateWebhook создаёт вебхук.
func (c *Client) CreateWebhook(input WebhookCreateInput) (*Webhook, error) {
	mutation := `
mutation CreateWebhook($input: WebhookCreateInput!) {
  webhookCreate(input: $input) {
    success
    webhook {
      id
      url
      enabled
      secret
      resourceTypes
      team { id key name }
    }
  }
}`
	inp := map[string]any{
		"url":    input.URL,
		"teamId": input.TeamID,
	}
	if input.Secret != "" {
		inp["secret"] = input.Secret
	}
	if len(input.ResourceTypes) > 0 {
		inp["resourceTypes"] = input.ResourceTypes
	}
	inp["enabled"] = input.Enabled

	var result struct {
		WebhookCreate struct {
			Webhook Webhook `json:"webhook"`
			Success bool    `json:"success"`
		} `json:"webhookCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": inp}, &result); err != nil {
		return nil, err
	}
	if !result.WebhookCreate.Success {
		return nil, fmt.Errorf("webhookCreate вернул success=false")
	}
	return &result.WebhookCreate.Webhook, nil
}

// UpdateWebhook обновляет вебхук по ID.
func (c *Client) UpdateWebhook(id string, input map[string]any) (*Webhook, error) {
	mutation := `
mutation UpdateWebhook($id: String!, $input: WebhookUpdateInput!) {
  webhookUpdate(id: $id, input: $input) {
    success
    webhook {
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
		WebhookUpdate struct {
			Webhook Webhook `json:"webhook"`
			Success bool    `json:"success"`
		} `json:"webhookUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.WebhookUpdate.Success {
		return nil, fmt.Errorf("webhookUpdate вернул success=false")
	}
	return &result.WebhookUpdate.Webhook, nil
}

// DeleteWebhook удаляет вебхук по ID.
func (c *Client) DeleteWebhook(id string) error {
	mutation := `
mutation DeleteWebhook($id: String!) {
  webhookDelete(id: $id) {
    success
  }
}`
	var result struct {
		WebhookDelete struct {
			Success bool `json:"success"`
		} `json:"webhookDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.WebhookDelete.Success {
		return fmt.Errorf("webhookDelete вернул success=false")
	}
	return nil
}

// RotateWebhookSecret обновляет секрет вебхука по ID.
func (c *Client) RotateWebhookSecret(id string) (*Webhook, error) {
	mutation := `
mutation RotateWebhookSecret($id: String!) {
  webhookRotateSecret(id: $id) {
    success
    webhook {
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
		WebhookRotateSecret struct {
			Webhook Webhook `json:"webhook"`
			Success bool    `json:"success"`
		} `json:"webhookRotateSecret"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	if !result.WebhookRotateSecret.Success {
		return nil, fmt.Errorf("webhookRotateSecret вернул success=false")
	}
	return &result.WebhookRotateSecret.Webhook, nil
}

// AttachmentLinkURL привязывает URL к задаче как вложение.
func (c *Client) AttachmentLinkURL(issueID, url, title string) (*Attachment, error) {
	mutation := `
mutation AttachmentLinkURL($issueId: String!, $url: String!, $title: String) {
  attachmentLinkURL(issueId: $issueId, url: $url, title: $title) {
    success
    attachment {
      id
      title
      url
      sourceType
      subtitle
    }
  }
}`
	var result struct {
		AttachmentLinkURL struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentLinkURL"`
	}
	vars := map[string]any{"issueId": issueID, "url": url}
	if title != "" {
		vars["title"] = title
	}
	if err := c.Do(mutation, vars, &result); err != nil {
		return nil, err
	}
	if !result.AttachmentLinkURL.Success {
		return nil, fmt.Errorf("attachmentLinkURL вернул success=false")
	}
	return &result.AttachmentLinkURL.Attachment, nil
}

// AttachmentLinkGitHubPR привязывает GitHub PR к задаче как вложение.
func (c *Client) AttachmentLinkGitHubPR(issueID, url, title string) (*Attachment, error) {
	mutation := `
mutation AttachmentLinkGitHubPR($issueId: String!, $url: String!, $title: String) {
  attachmentLinkGitHubPR(issueId: $issueId, url: $url, title: $title) {
    success
    attachment {
      id
      title
      url
      sourceType
      subtitle
    }
  }
}`
	var result struct {
		AttachmentLinkGitHubPR struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentLinkGitHubPR"`
	}
	vars := map[string]any{"issueId": issueID, "url": url}
	if title != "" {
		vars["title"] = title
	}
	if err := c.Do(mutation, vars, &result); err != nil {
		return nil, err
	}
	if !result.AttachmentLinkGitHubPR.Success {
		return nil, fmt.Errorf("attachmentLinkGitHubPR вернул success=false")
	}
	return &result.AttachmentLinkGitHubPR.Attachment, nil
}

// AttachmentLinkGitHubIssue привязывает GitHub Issue к задаче как вложение.
func (c *Client) AttachmentLinkGitHubIssue(issueID, url, title string) (*Attachment, error) {
	mutation := `
mutation AttachmentLinkGitHubIssue($issueId: String!, $url: String!, $title: String) {
  attachmentLinkGitHubIssue(issueId: $issueId, url: $url, title: $title) {
    success
    attachment {
      id
      title
      url
      sourceType
      subtitle
    }
  }
}`
	var result struct {
		AttachmentLinkGitHubIssue struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentLinkGitHubIssue"`
	}
	vars := map[string]any{"issueId": issueID, "url": url}
	if title != "" {
		vars["title"] = title
	}
	if err := c.Do(mutation, vars, &result); err != nil {
		return nil, err
	}
	if !result.AttachmentLinkGitHubIssue.Success {
		return nil, fmt.Errorf("attachmentLinkGitHubIssue вернул success=false")
	}
	return &result.AttachmentLinkGitHubIssue.Attachment, nil
}

// AttachmentLinkGitLabMR привязывает GitLab MR к задаче как вложение.
func (c *Client) AttachmentLinkGitLabMR(issueID, url string, number float64, projectPath, title string) (*Attachment, error) {
	mutation := `
mutation AttachmentLinkGitLabMR($issueId: String!, $url: String!, $number: Float!, $projectPathWithNamespace: String!, $title: String) {
  attachmentLinkGitLabMR(issueId: $issueId, url: $url, number: $number, projectPathWithNamespace: $projectPathWithNamespace, title: $title) {
    success
    attachment {
      id
      title
      url
      sourceType
      subtitle
    }
  }
}`
	var result struct {
		AttachmentLinkGitLabMR struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentLinkGitLabMR"`
	}
	vars := map[string]any{
		"issueId":                  issueID,
		"url":                      url,
		"number":                   number,
		"projectPathWithNamespace": projectPath,
	}
	if title != "" {
		vars["title"] = title
	}
	if err := c.Do(mutation, vars, &result); err != nil {
		return nil, err
	}
	if !result.AttachmentLinkGitLabMR.Success {
		return nil, fmt.Errorf("attachmentLinkGitLabMR вернул success=false")
	}
	return &result.AttachmentLinkGitLabMR.Attachment, nil
}

// AttachmentDelete удаляет вложение по ID.
func (c *Client) AttachmentDelete(id string) error {
	mutation := `
mutation AttachmentDelete($id: String!) {
  attachmentDelete(id: $id) {
    success
  }
}`
	var result struct {
		AttachmentDelete struct {
			Success bool `json:"success"`
		} `json:"attachmentDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.AttachmentDelete.Success {
		return fmt.Errorf("attachmentDelete вернул success=false")
	}
	return nil
}

// AttachmentUpdateInput — входные данные для обновления вложения.
type AttachmentUpdateInput struct {
	Title    string `json:"title,omitempty"`
	Subtitle string `json:"subtitle,omitempty"`
}

// AttachmentUpdate обновляет вложение по ID.
func (c *Client) AttachmentUpdate(id string, input AttachmentUpdateInput) (*Attachment, error) {
	mutation := `
mutation AttachmentUpdate($id: String!, $input: AttachmentUpdateInput!) {
  attachmentUpdate(id: $id, input: $input) {
    success
    attachment {
      id
      title
      url
      sourceType
      subtitle
    }
  }
}`
	var result struct {
		AttachmentUpdate struct {
			Attachment Attachment `json:"attachment"`
			Success    bool       `json:"success"`
		} `json:"attachmentUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.AttachmentUpdate.Success {
		return nil, fmt.Errorf("attachmentUpdate вернул success=false")
	}
	return &result.AttachmentUpdate.Attachment, nil
}

// CreateDocument создаёт новый документ.
func (c *Client) CreateDocument(title, content, projectID string) (*Document, error) {
	mutation := `
mutation DocumentCreate($input: DocumentCreateInput!) {
  documentCreate(input: $input) {
    success
    document {
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
	input := map[string]any{
		"title":   title,
		"content": content,
	}
	if projectID != "" {
		input["projectId"] = projectID
	}
	var result struct {
		DocumentCreate struct {
			Document Document `json:"document"`
			Success  bool     `json:"success"`
		} `json:"documentCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.DocumentCreate.Success {
		return nil, fmt.Errorf("documentCreate вернул success=false")
	}
	return &result.DocumentCreate.Document, nil
}

// UpdateDocument обновляет документ по ID.
func (c *Client) UpdateDocument(id string, input map[string]any) (*Document, error) {
	mutation := `
mutation DocumentUpdate($id: String!, $input: DocumentUpdateInput!) {
  documentUpdate(id: $id, input: $input) {
    success
    document {
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
		DocumentUpdate struct {
			Document Document `json:"document"`
			Success  bool     `json:"success"`
		} `json:"documentUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.DocumentUpdate.Success {
		return nil, fmt.Errorf("documentUpdate вернул success=false")
	}
	return &result.DocumentUpdate.Document, nil
}

// DeleteDocument удаляет документ по ID.
func (c *Client) DeleteDocument(id string) error {
	mutation := `
mutation DocumentDelete($id: String!) {
  documentDelete(id: $id) {
    success
  }
}`
	var result struct {
		DocumentDelete struct {
			Success bool `json:"success"`
		} `json:"documentDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.DocumentDelete.Success {
		return fmt.Errorf("documentDelete вернул success=false")
	}
	return nil
}

// CreateInitiative создаёт новую инициативу.
func (c *Client) CreateInitiative(name, description, ownerID string) (*Initiative, error) {
	mutation := `
mutation InitiativeCreate($input: InitiativeCreateInput!) {
  initiativeCreate(input: $input) {
    success
    initiative {
      id
      name
      description
      status
      owner { id name displayName email }
      archivedAt
    }
  }
}`
	input := map[string]any{"name": name}
	if description != "" {
		input["description"] = description
	}
	if ownerID != "" {
		input["ownerId"] = ownerID
	}
	var result struct {
		InitiativeCreate struct {
			Initiative Initiative `json:"initiative"`
			Success    bool       `json:"success"`
		} `json:"initiativeCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.InitiativeCreate.Success {
		return nil, fmt.Errorf("initiativeCreate вернул success=false")
	}
	return &result.InitiativeCreate.Initiative, nil
}

// UpdateInitiative обновляет инициативу.
func (c *Client) UpdateInitiative(id string, input map[string]any) (*Initiative, error) {
	mutation := `
mutation InitiativeUpdate($id: String!, $input: InitiativeUpdateInput!) {
  initiativeUpdate(id: $id, input: $input) {
    success
    initiative {
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
		InitiativeUpdate struct {
			Initiative Initiative `json:"initiative"`
			Success    bool       `json:"success"`
		} `json:"initiativeUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.InitiativeUpdate.Success {
		return nil, fmt.Errorf("initiativeUpdate вернул success=false")
	}
	return &result.InitiativeUpdate.Initiative, nil
}

// DeleteInitiative удаляет инициативу по ID.
func (c *Client) DeleteInitiative(id string) error {
	mutation := `
mutation InitiativeDelete($id: String!) {
  initiativeDelete(id: $id) {
    success
  }
}`
	var result struct {
		InitiativeDelete struct {
			Success bool `json:"success"`
		} `json:"initiativeDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.InitiativeDelete.Success {
		return fmt.Errorf("initiativeDelete вернул success=false")
	}
	return nil
}

// ArchiveInitiative архивирует инициативу по ID.
func (c *Client) ArchiveInitiative(id string) error {
	mutation := `
mutation InitiativeArchive($id: String!) {
  initiativeArchive(id: $id) {
    success
  }
}`
	var result struct {
		InitiativeArchive struct {
			Success bool `json:"success"`
		} `json:"initiativeArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.InitiativeArchive.Success {
		return fmt.Errorf("initiativeArchive вернул success=false")
	}
	return nil
}

// CreateInitiativeUpdateInput — входные данные для создания обновления инициативы.
type CreateInitiativeUpdateInput struct {
	InitiativeID string `json:"initiativeId"`
	Body         string `json:"body,omitempty"`
	Health       string `json:"health,omitempty"`
}

// UpdateInitiativeUpdateInput — входные данные для изменения обновления инициативы.
type UpdateInitiativeUpdateInput struct {
	Body   string `json:"body,omitempty"`
	Health string `json:"health,omitempty"`
}

// CreateInitiativeUpdate создаёт новую запись журнала инициативы.
func (c *Client) CreateInitiativeUpdate(input CreateInitiativeUpdateInput) (*InitiativeUpdate, error) {
	mutation := `
mutation CreateInitiativeUpdate($input: InitiativeUpdateCreateInput!) {
  initiativeUpdateCreate(input: $input) {
    success
    initiativeUpdate {
      id
      body
      health
      createdAt
      user { id name displayName email }
    }
  }
}`
	var result struct {
		InitiativeUpdateCreate struct {
			InitiativeUpdate InitiativeUpdate `json:"initiativeUpdate"`
			Success          bool             `json:"success"`
		} `json:"initiativeUpdateCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.InitiativeUpdateCreate.Success {
		return nil, fmt.Errorf("initiativeUpdateCreate вернул success=false")
	}
	return &result.InitiativeUpdateCreate.InitiativeUpdate, nil
}

// UpdateInitiativeUpdate обновляет запись журнала инициативы.
func (c *Client) UpdateInitiativeUpdate(id string, input UpdateInitiativeUpdateInput) (*InitiativeUpdate, error) {
	mutation := `
mutation UpdateInitiativeUpdate($id: String!, $input: InitiativeUpdateUpdateInput!) {
  initiativeUpdateUpdate(id: $id, input: $input) {
    success
    initiativeUpdate {
      id
      body
      health
      createdAt
      user { id name displayName email }
    }
  }
}`
	var result struct {
		InitiativeUpdateUpdate struct {
			InitiativeUpdate InitiativeUpdate `json:"initiativeUpdate"`
			Success          bool             `json:"success"`
		} `json:"initiativeUpdateUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.InitiativeUpdateUpdate.Success {
		return nil, fmt.Errorf("initiativeUpdateUpdate вернул success=false")
	}
	return &result.InitiativeUpdateUpdate.InitiativeUpdate, nil
}

// ArchiveInitiativeUpdate архивирует запись журнала инициативы.
func (c *Client) ArchiveInitiativeUpdate(id string) error {
	mutation := `
mutation ArchiveInitiativeUpdate($id: String!) {
  initiativeUpdateArchive(id: $id) {
    success
  }
}`
	var result struct {
		InitiativeUpdateArchive struct {
			Success bool `json:"success"`
		} `json:"initiativeUpdateArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.InitiativeUpdateArchive.Success {
		return fmt.Errorf("initiativeUpdateArchive вернул success=false")
	}
	return nil
}

// CreateInitiativeToProject создаёт связь инициативы с проектом.
func (c *Client) CreateInitiativeToProject(initiativeID, projectID string) (*InitiativeToProject, error) {
	mutation := `
mutation CreateInitiativeToProject($input: InitiativeToProjectCreateInput!) {
  initiativeToProjectCreate(input: $input) {
    success
    initiativeToProject {
      id
      initiative { id name }
      project { id name }
    }
  }
}`
	var result struct {
		InitiativeToProjectCreate struct {
			InitiativeToProject InitiativeToProject `json:"initiativeToProject"`
			Success             bool                `json:"success"`
		} `json:"initiativeToProjectCreate"`
	}
	input := map[string]any{
		"initiativeId": initiativeID,
		"projectId":    projectID,
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.InitiativeToProjectCreate.Success {
		return nil, fmt.Errorf("initiativeToProjectCreate вернул success=false")
	}
	return &result.InitiativeToProjectCreate.InitiativeToProject, nil
}

// DeleteInitiativeToProject удаляет связь инициативы с проектом.
func (c *Client) DeleteInitiativeToProject(id string) error {
	mutation := `
mutation DeleteInitiativeToProject($id: String!) {
  initiativeToProjectDelete(id: $id) {
    success
  }
}`
	var result struct {
		InitiativeToProjectDelete struct {
			Success bool `json:"success"`
		} `json:"initiativeToProjectDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.InitiativeToProjectDelete.Success {
		return fmt.Errorf("initiativeToProjectDelete вернул success=false")
	}
	return nil
}

// CreateRoadmap создаёт новую дорожную карту.
func (c *Client) CreateRoadmap(name, description, ownerID string) (*Roadmap, error) {
	mutation := `
mutation CreateRoadmap($input: RoadmapCreateInput!) {
  roadmapCreate(input: $input) {
    success
    roadmap {
      id
      name
      description
      createdAt
      owner { id name displayName email }
    }
  }
}`
	input := map[string]any{"name": name}
	if description != "" {
		input["description"] = description
	}
	if ownerID != "" {
		input["ownerId"] = ownerID
	}
	var result struct {
		RoadmapCreate struct {
			Roadmap Roadmap `json:"roadmap"`
			Success bool    `json:"success"`
		} `json:"roadmapCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.RoadmapCreate.Success {
		return nil, fmt.Errorf("roadmapCreate вернул success=false")
	}
	return &result.RoadmapCreate.Roadmap, nil
}

// UpdateRoadmap обновляет дорожную карту по ID.
func (c *Client) UpdateRoadmap(id string, input map[string]any) (*Roadmap, error) {
	mutation := `
mutation UpdateRoadmap($id: String!, $input: RoadmapUpdateInput!) {
  roadmapUpdate(id: $id, input: $input) {
    success
    roadmap {
      id
      name
      description
      createdAt
      owner { id name displayName email }
    }
  }
}`
	var result struct {
		RoadmapUpdate struct {
			Roadmap Roadmap `json:"roadmap"`
			Success bool    `json:"success"`
		} `json:"roadmapUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.RoadmapUpdate.Success {
		return nil, fmt.Errorf("roadmapUpdate вернул success=false")
	}
	return &result.RoadmapUpdate.Roadmap, nil
}

// DeleteRoadmap удаляет дорожную карту по ID.
func (c *Client) DeleteRoadmap(id string) error {
	mutation := `
mutation DeleteRoadmap($id: String!) {
  roadmapDelete(id: $id) {
    success
  }
}`
	var result struct {
		RoadmapDelete struct {
			Success bool `json:"success"`
		} `json:"roadmapDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.RoadmapDelete.Success {
		return fmt.Errorf("roadmapDelete вернул success=false")
	}
	return nil
}

// ArchiveRoadmap архивирует дорожную карту по ID.
func (c *Client) ArchiveRoadmap(id string) error {
	mutation := `
mutation ArchiveRoadmap($id: String!) {
  roadmapArchive(id: $id) {
    success
  }
}`
	var result struct {
		RoadmapArchive struct {
			Success bool `json:"success"`
		} `json:"roadmapArchive"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.RoadmapArchive.Success {
		return fmt.Errorf("roadmapArchive вернул success=false")
	}
	return nil
}

// CreateRoadmapToProject создаёт связь дорожной карты с проектом.
func (c *Client) CreateRoadmapToProject(roadmapID, projectID string) (*RoadmapToProject, error) {
	mutation := `
mutation CreateRoadmapToProject($input: RoadmapToProjectCreateInput!) {
  roadmapToProjectCreate(input: $input) {
    success
    roadmapToProject {
      id
      roadmap { id name }
      project { id name }
    }
  }
}`
	input := map[string]any{
		"roadmapId": roadmapID,
		"projectId": projectID,
	}
	var result struct {
		RoadmapToProjectCreate struct {
			RoadmapToProject RoadmapToProject `json:"roadmapToProject"`
			Success          bool             `json:"success"`
		} `json:"roadmapToProjectCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.RoadmapToProjectCreate.Success {
		return nil, fmt.Errorf("roadmapToProjectCreate вернул success=false")
	}
	return &result.RoadmapToProjectCreate.RoadmapToProject, nil
}

// CreateCustomer создаёт нового клиента.
func (c *Client) CreateCustomer(name string, input map[string]any) (*Customer, error) {
	mutation := `
mutation CreateCustomer($input: CustomerCreateInput!) {
  customerCreate(input: $input) {
    success
    customer {
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
	if input == nil {
		input = map[string]any{}
	}
	input["name"] = name
	var result struct {
		CustomerCreate struct {
			Customer Customer `json:"customer"`
			Success  bool     `json:"success"`
		} `json:"customerCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerCreate.Success {
		return nil, fmt.Errorf("customerCreate вернул success=false")
	}
	return &result.CustomerCreate.Customer, nil
}

// UpdateCustomer обновляет клиента по ID.
func (c *Client) UpdateCustomer(id string, input map[string]any) (*Customer, error) {
	mutation := `
mutation UpdateCustomer($id: String!, $input: CustomerUpdateInput!) {
  customerUpdate(id: $id, input: $input) {
    success
    customer {
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
		CustomerUpdate struct {
			Customer Customer `json:"customer"`
			Success  bool     `json:"success"`
		} `json:"customerUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerUpdate.Success {
		return nil, fmt.Errorf("customerUpdate вернул success=false")
	}
	return &result.CustomerUpdate.Customer, nil
}

// DeleteCustomer удаляет клиента по ID.
func (c *Client) DeleteCustomer(id string) error {
	mutation := `
mutation DeleteCustomer($id: String!) {
  customerDelete(id: $id) {
    success
  }
}`
	var result struct {
		CustomerDelete struct {
			Success bool `json:"success"`
		} `json:"customerDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CustomerDelete.Success {
		return fmt.Errorf("customerDelete вернул success=false")
	}
	return nil
}

// UpsertCustomer создаёт или обновляет клиента.
func (c *Client) UpsertCustomer(input map[string]any) (*Customer, error) {
	mutation := `
mutation UpsertCustomer($input: CustomerUpsertInput!) {
  customerUpsert(input: $input) {
    success
    customer {
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
		CustomerUpsert struct {
			Customer Customer `json:"customer"`
			Success  bool     `json:"success"`
		} `json:"customerUpsert"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerUpsert.Success {
		return nil, fmt.Errorf("customerUpsert вернул success=false")
	}
	return &result.CustomerUpsert.Customer, nil
}

// CreateCustomerNeed создаёт потребность клиента.
func (c *Client) CreateCustomerNeed(input map[string]any) (*CustomerNeed, error) {
	mutation := `
mutation CreateCustomerNeed($input: CustomerNeedCreateInput!) {
  customerNeedCreate(input: $input) {
    success
    customerNeed {
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
		CustomerNeedCreate struct {
			CustomerNeed CustomerNeed `json:"customerNeed"`
			Success      bool         `json:"success"`
		} `json:"customerNeedCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerNeedCreate.Success {
		return nil, fmt.Errorf("customerNeedCreate вернул success=false")
	}
	return &result.CustomerNeedCreate.CustomerNeed, nil
}

// UpdateCustomerNeed обновляет потребность клиента по ID.
func (c *Client) UpdateCustomerNeed(id string, input map[string]any) (*CustomerNeed, error) {
	mutation := `
mutation UpdateCustomerNeed($id: String!, $input: CustomerNeedUpdateInput!) {
  customerNeedUpdate(id: $id, input: $input) {
    success
    customerNeed {
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
		CustomerNeedUpdate struct {
			CustomerNeed CustomerNeed `json:"customerNeed"`
			Success      bool         `json:"success"`
		} `json:"customerNeedUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerNeedUpdate.Success {
		return nil, fmt.Errorf("customerNeedUpdate вернул success=false")
	}
	return &result.CustomerNeedUpdate.CustomerNeed, nil
}

// DeleteCustomerNeed удаляет потребность клиента по ID.
func (c *Client) DeleteCustomerNeed(id string) error {
	mutation := `
mutation DeleteCustomerNeed($id: String!) {
  customerNeedDelete(id: $id) {
    success
  }
}`
	var result struct {
		CustomerNeedDelete struct {
			Success bool `json:"success"`
		} `json:"customerNeedDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CustomerNeedDelete.Success {
		return fmt.Errorf("customerNeedDelete вернул success=false")
	}
	return nil
}

// CreateCustomerStatus создаёт статус клиента.
func (c *Client) CreateCustomerStatus(input map[string]any) (*CustomerStatus, error) {
	mutation := `
mutation CreateCustomerStatus($input: CustomerStatusCreateInput!) {
  customerStatusCreate(input: $input) {
    success
    customerStatus {
      id
      name
      displayName
      color
      description
    }
  }
}`
	var result struct {
		CustomerStatusCreate struct {
			CustomerStatus CustomerStatus `json:"customerStatus"`
			Success        bool           `json:"success"`
		} `json:"customerStatusCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerStatusCreate.Success {
		return nil, fmt.Errorf("customerStatusCreate вернул success=false")
	}
	return &result.CustomerStatusCreate.CustomerStatus, nil
}

// UpdateCustomerStatus обновляет статус клиента по ID.
func (c *Client) UpdateCustomerStatus(id string, input map[string]any) (*CustomerStatus, error) {
	mutation := `
mutation UpdateCustomerStatus($id: String!, $input: CustomerStatusUpdateInput!) {
  customerStatusUpdate(id: $id, input: $input) {
    success
    customerStatus {
      id
      name
      displayName
      color
      description
    }
  }
}`
	var result struct {
		CustomerStatusUpdate struct {
			CustomerStatus CustomerStatus `json:"customerStatus"`
			Success        bool           `json:"success"`
		} `json:"customerStatusUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerStatusUpdate.Success {
		return nil, fmt.Errorf("customerStatusUpdate вернул success=false")
	}
	return &result.CustomerStatusUpdate.CustomerStatus, nil
}

// DeleteCustomerStatus удаляет статус клиента по ID.
func (c *Client) DeleteCustomerStatus(id string) error {
	mutation := `
mutation DeleteCustomerStatus($id: String!) {
  customerStatusDelete(id: $id) {
    success
  }
}`
	var result struct {
		CustomerStatusDelete struct {
			Success bool `json:"success"`
		} `json:"customerStatusDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CustomerStatusDelete.Success {
		return fmt.Errorf("customerStatusDelete вернул success=false")
	}
	return nil
}

// CreateCustomerTier создаёт уровень клиента.
func (c *Client) CreateCustomerTier(input map[string]any) (*CustomerTier, error) {
	mutation := `
mutation CreateCustomerTier($input: CustomerTierCreateInput!) {
  customerTierCreate(input: $input) {
    success
    customerTier {
      id
      name
      displayName
      color
      description
    }
  }
}`
	var result struct {
		CustomerTierCreate struct {
			CustomerTier CustomerTier `json:"customerTier"`
			Success      bool         `json:"success"`
		} `json:"customerTierCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerTierCreate.Success {
		return nil, fmt.Errorf("customerTierCreate вернул success=false")
	}
	return &result.CustomerTierCreate.CustomerTier, nil
}

// UpdateCustomerTier обновляет уровень клиента по ID.
func (c *Client) UpdateCustomerTier(id string, input map[string]any) (*CustomerTier, error) {
	mutation := `
mutation UpdateCustomerTier($id: String!, $input: CustomerTierUpdateInput!) {
  customerTierUpdate(id: $id, input: $input) {
    success
    customerTier {
      id
      name
      displayName
      color
      description
    }
  }
}`
	var result struct {
		CustomerTierUpdate struct {
			CustomerTier CustomerTier `json:"customerTier"`
			Success      bool         `json:"success"`
		} `json:"customerTierUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomerTierUpdate.Success {
		return nil, fmt.Errorf("customerTierUpdate вернул success=false")
	}
	return &result.CustomerTierUpdate.CustomerTier, nil
}

// DeleteCustomerTier удаляет уровень клиента по ID.
func (c *Client) DeleteCustomerTier(id string) error {
	mutation := `
mutation DeleteCustomerTier($id: String!) {
  customerTierDelete(id: $id) {
    success
  }
}`
	var result struct {
		CustomerTierDelete struct {
			Success bool `json:"success"`
		} `json:"customerTierDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CustomerTierDelete.Success {
		return fmt.Errorf("customerTierDelete вернул success=false")
	}
	return nil
}

// DeleteRoadmapToProject удаляет связь дорожной карты с проектом.
func (c *Client) DeleteRoadmapToProject(id string) error {
	mutation := `
mutation DeleteRoadmapToProject($id: String!) {
  roadmapToProjectDelete(id: $id) {
    success
  }
}`
	var result struct {
		RoadmapToProjectDelete struct {
			Success bool `json:"success"`
		} `json:"roadmapToProjectDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.RoadmapToProjectDelete.Success {
		return fmt.Errorf("roadmapToProjectDelete вернул success=false")
	}
	return nil
}

// CreateTemplate создаёт новый шаблон.
func (c *Client) CreateTemplate(name, templateType string, input map[string]any) (*Template, error) {
	mutation := `
mutation CreateTemplate($input: TemplateCreateInput!) {
  templateCreate(input: $input) {
    success
    template {
      id
      name
      description
      type
      templateData
      createdAt
      creator { id name displayName email }
      team { id name key }
    }
  }
}`
	inp := map[string]any{
		"name": name,
		"type": templateType,
	}
	for k, v := range input {
		inp[k] = v
	}
	var result struct {
		TemplateCreate struct {
			Template Template `json:"template"`
			Success  bool     `json:"success"`
		} `json:"templateCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": inp}, &result); err != nil {
		return nil, err
	}
	if !result.TemplateCreate.Success {
		return nil, fmt.Errorf("templateCreate вернул success=false")
	}
	return &result.TemplateCreate.Template, nil
}

// UpdateTemplate обновляет шаблон по ID.
func (c *Client) UpdateTemplate(id string, input map[string]any) (*Template, error) {
	mutation := `
mutation UpdateTemplate($id: String!, $input: TemplateUpdateInput!) {
  templateUpdate(id: $id, input: $input) {
    success
    template {
      id
      name
      description
      type
      templateData
      createdAt
      creator { id name displayName email }
      team { id name key }
    }
  }
}`
	var result struct {
		TemplateUpdate struct {
			Template Template `json:"template"`
			Success  bool     `json:"success"`
		} `json:"templateUpdate"`
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.TemplateUpdate.Success {
		return nil, fmt.Errorf("templateUpdate вернул success=false")
	}
	return &result.TemplateUpdate.Template, nil
}

// DeleteTemplate удаляет шаблон по ID.
func (c *Client) DeleteTemplate(id string) error {
	mutation := `
mutation DeleteTemplate($id: String!) {
  templateDelete(id: $id) {
    success
  }
}`
	var result struct {
		TemplateDelete struct {
			Success bool `json:"success"`
		} `json:"templateDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.TemplateDelete.Success {
		return fmt.Errorf("templateDelete вернул success=false")
	}
	return nil
}

// CreateOrganizationInvite создаёт приглашение в организацию.
func (c *Client) CreateOrganizationInvite(email, role string) (*OrganizationInvite, error) {
	mutation := `
mutation OrganizationInviteCreate($input: OrganizationInviteCreateInput!) {
  organizationInviteCreate(input: $input) {
    success
    organizationInvite {
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
		OrganizationInviteCreate struct {
			OrganizationInvite OrganizationInvite `json:"organizationInvite"`
			Success            bool               `json:"success"`
		} `json:"organizationInviteCreate"`
	}
	input := map[string]any{"email": email}
	if role != "" {
		input["role"] = role
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.OrganizationInviteCreate.Success {
		return nil, fmt.Errorf("organizationInviteCreate вернул success=false")
	}
	return &result.OrganizationInviteCreate.OrganizationInvite, nil
}

// DeleteOrganizationInvite удаляет приглашение в организацию по ID.
func (c *Client) DeleteOrganizationInvite(id string) error {
	mutation := `
mutation OrganizationInviteDelete($id: String!) {
  organizationInviteDelete(id: $id) {
    success
  }
}`
	var result struct {
		OrganizationInviteDelete struct {
			Success bool `json:"success"`
		} `json:"organizationInviteDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.OrganizationInviteDelete.Success {
		return fmt.Errorf("organizationInviteDelete вернул success=false")
	}
	return nil
}

// ResendOrganizationInvite повторно отправляет приглашение по ID.
func (c *Client) ResendOrganizationInvite(id string) error {
	mutation := `
mutation ResendOrganizationInvite($id: String!) {
  resendOrganizationInvite(id: $id) {
    success
  }
}`
	var result struct {
		ResendOrganizationInvite struct {
			Success bool `json:"success"`
		} `json:"resendOrganizationInvite"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ResendOrganizationInvite.Success {
		return fmt.Errorf("resendOrganizationInvite вернул success=false")
	}
	return nil
}

// CreateCustomView создаёт новое пользовательское представление.
func (c *Client) CreateCustomView(name, description, icon, color string) (*CustomView, error) {
	mutation := `
mutation CustomViewCreate($input: CustomViewCreateInput!) {
  customViewCreate(input: $input) {
    success
    customView {
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
		CustomViewCreate struct {
			CustomView CustomView `json:"customView"`
			Success    bool       `json:"success"`
		} `json:"customViewCreate"`
	}
	input := map[string]any{"name": name}
	if description != "" {
		input["description"] = description
	}
	if icon != "" {
		input["icon"] = icon
	}
	if color != "" {
		input["color"] = color
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomViewCreate.Success {
		return nil, fmt.Errorf("customViewCreate вернул success=false")
	}
	return &result.CustomViewCreate.CustomView, nil
}

// UpdateCustomView обновляет пользовательское представление по ID.
func (c *Client) UpdateCustomView(id, name, description string) (*CustomView, error) {
	mutation := `
mutation CustomViewUpdate($id: String!, $input: CustomViewUpdateInput!) {
  customViewUpdate(id: $id, input: $input) {
    success
    customView {
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
		CustomViewUpdate struct {
			CustomView CustomView `json:"customView"`
			Success    bool       `json:"success"`
		} `json:"customViewUpdate"`
	}
	input := map[string]any{}
	if name != "" {
		input["name"] = name
	}
	if description != "" {
		input["description"] = description
	}
	if err := c.Do(mutation, map[string]any{"id": id, "input": input}, &result); err != nil {
		return nil, err
	}
	if !result.CustomViewUpdate.Success {
		return nil, fmt.Errorf("customViewUpdate вернул success=false")
	}
	return &result.CustomViewUpdate.CustomView, nil
}

// CreateFavorite создаёт избранный элемент.
// entityType: "issue", "project", "label", "customView"
// entityID: ID соответствующего элемента
func (c *Client) CreateFavorite(entityType, entityID string) (*Favorite, error) {
	mutation := `
mutation FavoriteCreate($input: FavoriteCreateInput!) {
  favoriteCreate(input: $input) {
    success
    favorite {
      id
      type
      issue { id identifier title }
      project { id name }
      label { id name color }
      customView { id name }
    }
  }
}`
	input := map[string]any{}
	switch entityType {
	case "issue":
		input["issueId"] = entityID
	case "project":
		input["projectId"] = entityID
	case "label":
		input["issueLabelId"] = entityID
	case "customView":
		input["customViewId"] = entityID
	default:
		return nil, fmt.Errorf("неизвестный тип избранного: %s", entityType)
	}

	var result struct {
		FavoriteCreate struct {
			Favorite Favorite `json:"favorite"`
			Success  bool     `json:"success"`
		} `json:"favoriteCreate"`
	}
	if err := c.Do(mutation, map[string]any{"input": input}, &result); err != nil {
		return nil, err
	}
	if !result.FavoriteCreate.Success {
		return nil, fmt.Errorf("favoriteCreate вернул success=false")
	}
	return &result.FavoriteCreate.Favorite, nil
}

// DeleteFavorite удаляет избранный элемент по ID.
func (c *Client) DeleteFavorite(id string) error {
	mutation := `
mutation FavoriteDelete($id: String!) {
  favoriteDelete(id: $id) {
    success
  }
}`
	var result struct {
		FavoriteDelete struct {
			Success bool `json:"success"`
		} `json:"favoriteDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.FavoriteDelete.Success {
		return fmt.Errorf("favoriteDelete вернул success=false")
	}
	return nil
}

// UpdateFavorite обновляет избранный элемент по ID (например, порядок сортировки).
func (c *Client) UpdateFavorite(id string, sortOrder float64) (*Favorite, error) {
	mutation := `
mutation FavoriteUpdate($id: String!, $input: FavoriteUpdateInput!) {
  favoriteUpdate(id: $id, input: $input) {
    success
    favorite {
      id
      type
    }
  }
}`
	var result struct {
		FavoriteUpdate struct {
			Favorite Favorite `json:"favorite"`
			Success  bool     `json:"success"`
		} `json:"favoriteUpdate"`
	}
	if err := c.Do(mutation, map[string]any{
		"id":    id,
		"input": map[string]any{"sortOrder": sortOrder},
	}, &result); err != nil {
		return nil, err
	}
	if !result.FavoriteUpdate.Success {
		return nil, fmt.Errorf("favoriteUpdate вернул success=false")
	}
	return &result.FavoriteUpdate.Favorite, nil
}

// DeleteCustomView удаляет пользовательское представление по ID.
func (c *Client) DeleteCustomView(id string) error {
	mutation := `
mutation CustomViewDelete($id: String!) {
  customViewDelete(id: $id) {
    success
  }
}`
	var result struct {
		CustomViewDelete struct {
			Success bool `json:"success"`
		} `json:"customViewDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.CustomViewDelete.Success {
		return fmt.Errorf("customViewDelete вернул success=false")
	}
	return nil
}

// CreateReaction добавляет реакцию на комментарий.
func (c *Client) CreateReaction(commentID, emoji string) (string, error) {
	mutation := `
mutation ReactionCreate($input: ReactionCreateInput!) {
  reactionCreate(input: $input) {
    success
    reaction {
      id
    }
  }
}`
	var result struct {
		ReactionCreate struct {
			Reaction struct {
				ID string `json:"id"`
			} `json:"reaction"`
			Success bool `json:"success"`
		} `json:"reactionCreate"`
	}
	if err := c.Do(mutation, map[string]any{
		"input": map[string]any{
			"commentId": commentID,
			"emoji":     emoji,
		},
	}, &result); err != nil {
		return "", err
	}
	if !result.ReactionCreate.Success {
		return "", fmt.Errorf("reactionCreate вернул success=false")
	}
	return result.ReactionCreate.Reaction.ID, nil
}

// DeleteReaction удаляет реакцию по ID.
func (c *Client) DeleteReaction(id string) error {
	mutation := `
mutation ReactionDelete($id: String!) {
  reactionDelete(id: $id) {
    success
  }
}`
	var result struct {
		ReactionDelete struct {
			Success bool `json:"success"`
		} `json:"reactionDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.ReactionDelete.Success {
		return fmt.Errorf("reactionDelete вернул success=false")
	}
	return nil
}

// CreateEmoji создаёт кастомный эмодзи.
func (c *Client) CreateEmoji(name, url string) (*Emoji, error) {
	mutation := `
mutation EmojiCreate($input: EmojiCreateInput!) {
  emojiCreate(input: $input) {
    success
    emoji {
      id
      name
      url
    }
  }
}`
	var result struct {
		EmojiCreate struct {
			Emoji   Emoji `json:"emoji"`
			Success bool  `json:"success"`
		} `json:"emojiCreate"`
	}
	if err := c.Do(mutation, map[string]any{
		"input": map[string]any{
			"name": name,
			"url":  url,
		},
	}, &result); err != nil {
		return nil, err
	}
	if !result.EmojiCreate.Success {
		return nil, fmt.Errorf("emojiCreate вернул success=false")
	}
	return &result.EmojiCreate.Emoji, nil
}

// DeleteEmoji удаляет кастомный эмодзи по ID.
func (c *Client) DeleteEmoji(id string) error {
	mutation := `
mutation EmojiDelete($id: String!) {
  emojiDelete(id: $id) {
    success
  }
}`
	var result struct {
		EmojiDelete struct {
			Success bool `json:"success"`
		} `json:"emojiDelete"`
	}
	if err := c.Do(mutation, map[string]any{"id": id}, &result); err != nil {
		return err
	}
	if !result.EmojiDelete.Success {
		return fmt.Errorf("emojiDelete вернул success=false")
	}
	return nil
}
