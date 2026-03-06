package client

import "fmt"

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
