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
