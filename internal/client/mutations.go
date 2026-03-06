package client

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
	return &result.IssueUpdate.Issue, nil
}

// GetTeamByKey возвращает команду по её ключу (например, "ENG").
func (c *Client) GetTeamByKey(key string) (*Team, error) {
	query := `
query GetTeams {
  teams {
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

// FindUserByName ищет пользователя по displayName или email.
func (c *Client) FindUserByName(name string) (*User, error) {
	query := `
query GetUsers {
  users {
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
