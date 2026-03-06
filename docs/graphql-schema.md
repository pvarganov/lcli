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
    }
  }
}
```

Переменные (`input`):
- `title` (String!) — заголовок
- `teamId` (String!) — ID команды
- `description` (String) — описание
- `assigneeId` (String) — ID исполнителя
- `priority` (Int) — приоритет (0-4)

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
    }
  }
}
```

Переменные (`input`):
- `stateId` (String) — ID статуса
- `assigneeId` (String) — ID исполнителя
- `priority` (Int) — приоритет (0-4)
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
    }
  }
}
```

Переменные (`input`):
- `issueId` (String!) — ID задачи
- `body` (String!) — текст комментария
