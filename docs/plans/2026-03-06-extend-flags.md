# Расширение флагов основных операций lcli

## Overview

Добавить все полезные флаги для основных операций, поддерживаемых Linear GraphQL API, но не реализованных в текущем CLI. Цель — покрыть IssueCreateInput, IssueUpdateInput, IssueFilter, CommentCreateInput, ProjectCreateInput, ProjectUpdateInput.

**Что меняется:**
- `issues create` — добавить: `--due-date`, `--estimate`, `--labels`, `--parent`, `--state`, `--cycle-id`, `--project-id`, `--milestone-id`
- `issues update` — добавить: `--description`, `--due-date`, `--estimate`, `--parent`, `--cycle-id`, `--project-id`, `--milestone-id`, `--add-labels`, `--remove-labels`, `--snooze-until`
- `issues list` — добавить фильтры: `--priority`, `--label`, `--project-id`, `--cycle-id`, `--creator`, `--order-by`
- `issue comment` — добавить: `--parent-id` (вложенные комментарии)
- `projects create/update` — добавить: `--color`, `--icon`, `--priority`, `--member-ids`, `--content`

## Context (from discovery)

- Команды задач: `cmd/issue_create.go`, `cmd/issue_update.go`, `cmd/issues.go`, `cmd/issue_comment.go`
- Команды проектов: `cmd/project_mutations.go`, `cmd/projects.go`
- Клиентский слой: `internal/client/mutations.go`, `internal/client/queries.go`
- Тесты: `cmd/issue_mutations_test.go`, `cmd/issues_test.go`, `cmd/issue_comment_test.go`, `cmd/project_mutations_test.go`
- API схема: `docs/Linear-API@current--#@!api!@#.json`

## Development Approach

- **Testing approach**: TDD — тесты пишутся до реализации
- Каждый таск завершается работающими тестами перед переходом к следующему
- Изменения в client-слое: struct + метод; изменения в cmd-слое: флаги + логика разрешения имён
- Для полей типа `[String!]` (labelIds и т.д.) принимаем запятые-разделённые значения в CLI

## Testing Strategy

- **Unit tests**: table-driven тесты через mock-клиент (паттерн уже есть в проекте)
- Проверяем: флаг передаётся в GraphQL input, флаг не передаётся если не задан

## Progress Tracking

- Отмечать `[x]` сразу после завершения
- `➕` — новая обнаруженная задача
- `⚠️` — блокер

## What Goes Where

- **Implementation Steps**: изменения кода, тесты, lint
- **Post-Completion**: ручное тестирование с реальным Linear-аккаунтом

## Implementation Steps

### Task 1: Расширить CreateIssueInput и клиент issues create

- [x] добавить поля в `CreateIssueInput` в `internal/client/mutations.go`: `DueDate string`, `Estimate *int`, `LabelIDs []string`, `ParentID string`, `CycleID string`, `ProjectID string`, `MilestoneID string`, `StateID string`
- [x] расширить `CreateIssue` — передавать новые поля в `gqlInput` если не пусты
- [x] добавить метод `FindLabelsByNames(teamID string, names []string) ([]string, error)` в `internal/client/queries.go`
- [x] написать тесты в `cmd/issue_mutations_test.go` для каждого нового поля (due-date, estimate, labels, parent, state, cycle, project, milestone)
- [x] проверить, что тесты падают (TDD red)
- [x] запустить `go test ./...` — тесты должны пройти после реализации

### Task 2: Добавить флаги в `issues create` (cmd-слой)

- [x] добавить флаги в `cmd/issue_create.go`: `--due-date`, `--estimate`, `--labels` (comma-sep), `--parent`, `--state`, `--cycle-id`, `--project-id`, `--milestone-id`
- [x] добавить логику разрешения: `--labels` → `FindLabelsByNames`, `--state` → `FindWorkflowStateByName`, `--parent` → передавать как identifier напрямую
- [x] написать тесты в `cmd/issue_mutations_test.go` на новые флаги (success + отсутствие поля при незаданном флаге)
- [x] запустить `go test ./cmd/...` — все тесты зелёные

### Task 3: Расширить UpdateIssueInput и клиент issues update

- [x] добавить поля в `UpdateIssueInput` в `internal/client/mutations.go`: `Description string`, `DueDate string`, `Estimate *int`, `ParentID string`, `CycleID string`, `ProjectID string`, `MilestoneID string`, `AddedLabelIDs []string`, `RemovedLabelIDs []string`, `SnoozedUntilAt string`
- [x] расширить `UpdateIssue` — передавать новые поля в `gqlInput` если не пусты
- [x] написать тесты в `cmd/issue_mutations_test.go` для каждого нового update-поля
- [x] запустить `go test ./...` — все зелёные

### Task 4: Добавить флаги в `issues update` (cmd-слой)

- [x] добавить флаги в `cmd/issue_update.go`: `--description`, `--due-date`, `--estimate`, `--parent`, `--cycle-id`, `--project-id`, `--milestone-id`, `--add-labels`, `--remove-labels`, `--snooze-until`
- [x] обновить проверку "хотя бы один флаг задан" — включить новые флаги
- [x] добавить логику разрешения: `--add-labels`/`--remove-labels` → `FindLabelsByNames` (нужен teamID из issue)
- [x] написать тесты на новые флаги
- [x] запустить `go test ./cmd/...` — все зелёные

### Task 5: Расширить IssueFilter и ListIssues (issues list)

- [x] добавить поля в `IssueFilter` в `internal/client/queries.go`: `Priority int`, `Label string`, `ProjectID string`, `CycleID string`, `Creator string`, `OrderBy string`
- [x] расширить `ListIssues` — добавить новые поля в GraphQL filter и `orderBy` переменную
- [x] написать тесты в `cmd/issues_test.go` для каждого нового фильтра
- [x] добавить флаги в `cmd/issues.go`: `--priority`, `--label`, `--project-id`, `--cycle-id`, `--creator`, `--order-by`
- [x] запустить `go test ./...` — все зелёные

### Task 6: Добавить --parent-id к comment create

- [x] создать `CreateCommentInput` struct в `internal/client/mutations.go` (тело, issue ID, parent ID)
- [x] обновить `CreateComment` чтобы принимал `CreateCommentInput` вместо двух строк (или добавить новый метод)
- [x] обновить вызовы `CreateComment` в `cmd/issue_comment.go`
- [x] написать тесты в `cmd/issue_comment_test.go` для `--parent-id`
- [x] добавить флаг `--parent-id` в `issueCommentCmd`
- [x] запустить `go test ./...` — все зелёные

### Task 7: Расширить projects create/update

- [x] добавить поля в `CreateProjectInput` в `internal/client/mutations.go`: `Color string`, `Icon string`, `Priority *int`, `MemberIDs []string`, `Content string`
- [x] добавить поля в `UpdateProjectInput`: `Color string`, `Icon string`, `Priority *int`, `MemberIDs []string`, `Content string`
- [x] расширить `CreateProject` и `UpdateProject` — передавать новые поля
- [x] написать тесты в `cmd/project_mutations_test.go` для новых полей
- [x] добавить флаги в `cmd/project_mutations.go`: `--color`, `--icon`, `--priority`, `--member-ids`, `--content`
- [x] запустить `go test ./...` — все зелёные

### Task 8: Финальная верификация

- [ ] запустить полный `go test ./...` — все тесты зелёные
- [ ] запустить `go vet ./...` — нет ошибок
- [ ] запустить `golangci-lint run` если установлен (или `go build ./...`)
- [ ] проверить `lcli issue create --help` и `lcli issue update --help` — новые флаги видны
- [ ] обновить README.md если описаны команды с флагами

## Technical Details

**Разрешение labels по имени:**
```
FindLabelsByNames(teamID string, names []string) ([]string, error)
→ query labels by team, filter by name, return IDs
```

**Формат дат:**
- `dueDate`: TimelessDate → строка `YYYY-MM-DD` (без timezone)
- `snoozedUntilAt`: DateTime → строка RFC3339

**IssueFilter orderBy:**
- Допустимые значения из схемы: `updatedAt`, `createdAt`, `priority`, `manualOrder`

**Паттерн передачи необязательных полей:**
- Для `*int` (estimate, priority) — используем `cmd.Flags().Changed("flag-name")` чтобы отличить 0 от "не задан"
- Для `[]string` — пропускаем если len == 0

## Post-Completion

**Ручное тестирование** (требует реального Linear-токена):
- `lcli issue create --team ENG --title "test" --due-date 2026-04-01 --estimate 3 --labels "Bug"`
- `lcli issue update OVG-1 --add-labels "Bug" --due-date 2026-04-01`
- `lcli issues list --priority 1 --label "Bug" --order-by priority`
- `lcli issue comment OVG-1 --body "reply" --parent-id <comment-uuid>`

*Note: ralphex automatically moves completed plans to `docs/plans/completed/`*
