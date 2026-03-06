# GraphQL-запросы lcli

Все запросы отправляются методом POST на `https://api.linear.app/graphql`.

Заголовки:
- `Content-Type: application/json`
- `Authorization: Bearer <token>`

Тело запроса:
```json
{
  "query": "...",
  "variables": { "key": "value" }
}
```

## Queries

### ListIssues

Список задач с фильтрацией и пагинацией.

```graphql
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
}
```

Переменные:
- `first` (Int) — количество записей (по умолчанию 25)
- `after` (String) — курсор пагинации
- `filter` (IssueFilter) — объект фильтрации:
  - `team.key.eq` — ключ команды
  - `state.name.eq` — имя статуса
  - `assignee.displayName.eq` — имя исполнителя

### GetIssue

Детали одной задачи по ID.

```graphql
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
}
```

### ListComments

Список комментариев к задаче.

```graphql
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
}
```

### ListProjects

Список проектов.

```graphql
query ListProjects {
  projects {
    nodes {
      id
      name
      description
      state
    }
  }
}
```

### ListTeams

Список команд.

```graphql
query ListTeams {
  teams {
    nodes {
      id
      key
      name
    }
  }
}
```

### GetTeams (используется в issue create --team)

Получение команды по ключу. Используется при создании задачи для преобразования ключа команды (например, `ENG`) в ID.

```graphql
query GetTeams {
  teams {
    nodes {
      id
      key
      name
    }
  }
}
```

Матчинг по полю `key` выполняется на стороне клиента.

### GetWorkflowStates (используется в issue update --status)

Список состояний задач в рамках команды. Используется при обновлении статуса задачи для поиска `stateId` по имени.

```graphql
query GetWorkflowStates($teamId: ID!) {
  workflowStates(filter: { team: { id: { eq: $teamId } } }) {
    nodes {
      id
      name
    }
  }
}
```

Переменные:
- `teamId` (ID!) — ID команды

### GetUsers (используется в --assignee)

Список пользователей для поиска по `displayName`, `email` или `name`. Используется при указании исполнителя в `issue create` и `issue update`.

```graphql
query GetUsers {
  users {
    nodes {
      id
      name
      displayName
      email
    }
  }
}
```

Матчинг по `displayName`, `email` или `name` выполняется на стороне клиента.

## Mutations

### CreateIssue

Создать задачу.

```graphql
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
}
```

Переменные (`input`):
- `title` (String!) — заголовок
- `teamId` (String!) — ID команды
- `description` (String) — описание
- `assigneeId` (String) — ID исполнителя
- `priority` (Int) — приоритет (1-4; 0 не передаётся)

### UpdateIssue

Обновить задачу.

```graphql
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
}
```

Переменные (`input`):
- `stateId` (String) — ID статуса
- `assigneeId` (String) — ID исполнителя
- `priority` (Int) — приоритет (0-4; 0 сбрасывает приоритет)
- `title` (String) — заголовок

### CreateComment

Добавить комментарий к задаче.

```graphql
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
}
```

Переменные (`input`):
- `issueId` (String!) — ID задачи
- `body` (String!) — текст комментария
