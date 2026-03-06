# Полное покрытие Linear GraphQL API в lcli

## Overview

Расширение CLI-утилиты lcli для покрытия всего публичного Linear GraphQL API.
Текущее состояние: реализовано ~4% API (8 команд).
Цель: покрыть все значимые Query и Mutation Linear API.

Схема API: `docs/Linear-API@current--#@!api!@#.json`

## Context (from discovery)

- Файлы клиента: `internal/client/queries.go`, `internal/client/mutations.go`, `internal/client/client.go`
- Команды: `cmd/issues.go`, `cmd/issue_create.go`, `cmd/issue_update.go`, `cmd/issue_view.go`, `cmd/issue_comment.go`, `cmd/projects.go`, `cmd/teams.go`, `cmd/auth.go`
- Формат: таблица (default) или JSON (`-o json`)
- Паттерн команды: cobra + client.Do(graphql) + format.TableWriter
- Тесты: table-driven, моки через замену `newLinearClient`

## Development Approach

- **Testing approach**: TDD (тесты первыми)
- Каждый task = один логический блок (одна группа команд)
- Сначала тесты, потом реализация
- Все тесты должны пройти перед переходом к следующему task
- Обратная совместимость с существующими командами

## Testing Strategy

- **Unit tests**: table-driven, мок клиента через замену фабрики `newLinearClient`
- Паттерн из существующих тестов: `cmd/issues_test.go`, `cmd/issue_mutations_test.go`
- Тесты для client-методов: `internal/client/queries_test.go`, `internal/client/mutations_test.go`

## Progress Tracking

- Отмечать выполненные пункты `[x]` сразу после завершения
- Новые задачи добавлять с префиксом ➕
- Блокеры отмечать с префиксом ⚠️

---

## Implementation Steps

### Phase 1: Issues — расширенная функциональность

#### Task 1: Issue Labels (метки задач) — Query

- [x] добавить типы `IssueLabel` в `internal/client/queries.go`
- [x] добавить метод `ListIssueLabels()` — все метки организации
- [x] написать тесты для `ListIssueLabels` (success + empty)
- [x] добавить метод `GetIssueLabel(id)` — одна метка по ID
- [x] написать тесты для `GetIssueLabel` (success + not found)
- [x] запустить тесты — должны пройти

#### Task 2: Issue Labels (метки задач) — Mutations + CMD

- [x] написать тесты для cmd `issues labels list` (table + json output)
- [x] добавить команды `issues labels list` в `cmd/issue_labels.go`
- [x] написать тесты для мутаций `issueLabelCreate`, `issueLabelUpdate`, `issueLabelDelete`
- [x] добавить методы `CreateIssueLabel`, `UpdateIssueLabel`, `DeleteIssueLabel` в `internal/client/mutations.go`
- [x] написать тесты для cmd `issues labels create/update/delete`
- [x] добавить команды `issues labels create`, `issues labels update`, `issues labels delete`
- [x] запустить тесты — должны пройти

#### Task 3: Issue Labels — добавление/удаление меток на задачах

- [x] написать тесты для мутаций `issueAddLabel`, `issueRemoveLabel`
- [x] добавить методы `AddLabelToIssue(issueID, labelID)`, `RemoveLabelFromIssue(issueID, labelID)` в mutations.go
- [x] написать тесты для cmd `issues label add <ISSUE-ID> <label-name>` и `issues label remove`
- [x] добавить команды `issues label add` и `issues label remove` (поиск метки по имени)
- [x] запустить тесты — должны пройти

#### Task 4: Issue Relations (связи задач)

- [x] добавить тип `IssueRelation` в queries.go (id, type, issue, relatedIssue)
- [x] написать тесты для `ListIssueRelations(issueID)`
- [x] добавить метод `ListIssueRelations(issueID)` в queries.go
- [x] написать тесты для мутаций `issueRelationCreate`, `issueRelationDelete`
- [x] добавить методы `CreateIssueRelation`, `DeleteIssueRelation` в mutations.go
- [x] написать тесты для cmd `issues relations list <ISSUE-ID>`, `issues relations add`, `issues relations remove`
- [x] добавить команды в `cmd/issue_relations.go`
- [x] запустить тесты — должны пройти

#### Task 5: Issue Archive/Delete/Unarchive

- [x] написать тесты для мутаций `issueArchive`, `issueUnarchive`, `issueDelete`
- [x] добавить методы `ArchiveIssue(id)`, `UnarchiveIssue(id)`, `DeleteIssue(id)` в mutations.go
- [x] написать тесты для cmd `issues archive <ID>`, `issues unarchive <ID>`, `issues delete <ID>`
- [x] добавить команды в `cmd/issue_archive.go`
- [x] запустить тесты — должны пройти

#### Task 6: Issue Search

- [x] добавить тип `IssueSearchResult` (nodes + pageInfo) в queries.go
- [x] написать тесты для `SearchIssues(query, limit)`
- [x] добавить метод `SearchIssues(query string, limit int)` в queries.go (использовать `searchIssues` mutation)
- [x] написать тесты для cmd `issues search <query>`
- [x] добавить команду `issues search` в `cmd/issue_search.go` (флаги: --limit, --team, -o json)
- [x] запустить тесты — должны пройти

#### Task 7: Issue Subscribe/Unsubscribe

- [x] написать тесты для мутаций `issueSubscribe`, `issueUnsubscribe`
- [x] добавить методы `SubscribeToIssue(id)`, `UnsubscribeFromIssue(id)` в mutations.go
- [x] написать тесты для cmd `issues subscribe <ID>`, `issues unsubscribe <ID>`
- [x] добавить команды в `cmd/issue_subscribe.go`
- [x] запустить тесты — должны пройти

#### Task 8: Issue Batch Operations

- [x] написать тесты для мутаций `issueBatchCreate`, `issueBatchUpdate`
- [x] добавить типы и методы `BatchCreateIssues`, `BatchUpdateIssues` в mutations.go
- [x] написать тесты для cmd `issues batch-update --status <S> --ids <id1,id2,...>`
- [x] добавить команду `issues batch-update` в `cmd/issue_batch.go`
- [x] запустить тесты — должны пройти

#### Task 9: Comments — Update/Delete/Resolve

- [x] написать тесты для мутаций `commentUpdate`, `commentDelete`, `commentResolve`, `commentUnresolve`
- [x] добавить методы `UpdateComment(id, body)`, `DeleteComment(id)`, `ResolveComment(id)`, `UnresolveComment(id)` в mutations.go
- [x] написать тесты для cmd `issues comment update <ID>`, `issues comment delete <ID>`, `issues comment resolve <ID>`
- [x] добавить команды в `cmd/issue_comment.go` (расширить существующий файл)
- [x] запустить тесты — должны пройти

---

### Phase 2: Projects — расширенная функциональность

#### Task 10: Projects — Create/Update/Delete

- [x] расширить тип `Project` в queries.go (добавить startDate, targetDate, lead, members, teams, url)
- [x] написать тесты для `GetProject(id)` query
- [x] добавить метод `GetProject(id)` в queries.go
- [x] написать тесты для мутаций `projectCreate`, `projectUpdate`, `projectDelete`, `projectArchive`, `projectUnarchive`
- [x] добавить методы `CreateProject`, `UpdateProject`, `DeleteProject`, `ArchiveProject`, `UnarchiveProject` в mutations.go
- [x] написать тесты для cmd `projects create`, `projects update <ID>`, `projects delete <ID>`, `projects view <ID>`
- [x] добавить команды в `cmd/projects.go` (расширить) и `cmd/project_mutations.go`
- [x] запустить тесты — должны пройти

#### Task 11: Project Milestones

- [x] добавить тип `ProjectMilestone` в queries.go (id, name, targetDate, description)
- [x] написать тесты для `ListProjectMilestones(projectID)`
- [x] добавить метод `ListProjectMilestones(projectID)` в queries.go
- [x] написать тесты для мутаций `projectMilestoneCreate`, `projectMilestoneUpdate`, `projectMilestoneDelete`
- [x] добавить методы `CreateProjectMilestone`, `UpdateProjectMilestone`, `DeleteProjectMilestone` в mutations.go
- [x] написать тесты для cmd `projects milestones list <PROJECT-ID>`, `projects milestones create`, `projects milestones update`, `projects milestones delete`
- [x] добавить команды в `cmd/project_milestones.go`
- [x] запустить тесты — должны пройти

#### Task 12: Project Updates (журнал обновлений проекта)

- [x] добавить тип `ProjectUpdate` в queries.go (id, body, createdAt, user, health)
- [x] написать тесты для `ListProjectUpdates(projectID)`
- [x] добавить метод `ListProjectUpdates(projectID)` в queries.go
- [x] написать тесты для мутаций `projectUpdateCreate`, `projectUpdateUpdate`, `projectUpdateArchive`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `projects updates list <PROJECT-ID>`, `projects updates create`, `projects updates delete`
- [x] добавить команды в `cmd/project_updates.go`
- [x] запустить тесты — должны пройти

#### Task 13: Project Labels (метки проекта)

- [x] добавить тип `ProjectLabel` в queries.go
- [x] написать тесты для `ListProjectLabels()`
- [x] добавить метод `ListProjectLabels()` в queries.go
- [x] написать тесты для мутаций `projectLabelCreate`, `projectLabelUpdate`, `projectLabelDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `projects labels list`, `projects labels create`, `projects labels update`, `projects labels delete`
- [x] добавить команды в `cmd/project_labels.go`
- [x] запустить тесты — должны пройти

#### Task 14: Project Statuses

- [x] добавить тип `ProjectStatus` в queries.go
- [x] написать тесты для `ListProjectStatuses()`
- [x] добавить метод `ListProjectStatuses()` в queries.go
- [x] написать тесты для мутаций `projectStatusCreate`, `projectStatusUpdate`, `projectStatusArchive`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `projects statuses list`, `projects statuses create`, `projects statuses update`
- [x] добавить команды в `cmd/project_statuses.go`
- [x] запустить тесты — должны пройти

#### Task 15: Project Search + Relations

- [x] написать тесты для `SearchProjects(query)`
- [x] добавить метод `SearchProjects(query string)` в queries.go
- [x] написать тесты для cmd `projects search <query>`
- [x] добавить команду `projects search` в cmd/projects.go
- [x] написать тесты для мутаций `projectRelationCreate`, `projectRelationDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `projects relations add`, `projects relations remove`
- [x] добавить команды в `cmd/project_relations.go`
- [x] запустить тесты — должны пройти

---

### Phase 3: Cycles (циклы)

#### Task 16: Cycles — Query

- [x] добавить тип `Cycle` в queries.go (id, number, name, startsAt, endsAt, team, issues)
- [x] написать тесты для `ListCycles(teamKey)`, `GetCycle(id)`
- [x] добавить методы `ListCycles(teamID)`, `GetCycle(id)` в queries.go
- [x] написать тесты для cmd `cycles list --team <KEY>`, `cycles view <ID>`
- [x] добавить команды в `cmd/cycles.go`
- [x] запустить тесты — должны пройти

#### Task 17: Cycles — Mutations

- [x] написать тесты для мутаций `cycleCreate`, `cycleUpdate`, `cycleArchive`
- [x] добавить методы `CreateCycle`, `UpdateCycle`, `ArchiveCycle` в mutations.go
- [x] написать тесты для cmd `cycles create`, `cycles update <ID>`, `cycles archive <ID>`
- [x] добавить команды в `cmd/cycles.go`
- [x] запустить тесты — должны пройти

---

### Phase 4: Workflow States (статусы задач)

#### Task 18: Workflow States — Query + Mutations

- [x] добавить тип `WorkflowState` в queries.go (id, name, type, color, team)
- [x] написать тесты для `ListWorkflowStates(teamID)`
- [x] добавить метод `ListWorkflowStates(teamID)` в queries.go (расширить существующий)
- [x] написать тесты для мутаций `workflowStateCreate`, `workflowStateUpdate`, `workflowStateArchive`
- [x] добавить методы `CreateWorkflowState`, `UpdateWorkflowState`, `ArchiveWorkflowState` в mutations.go
- [x] написать тесты для cmd `workflow-states list --team <KEY>`, `workflow-states create`, `workflow-states update`, `workflow-states archive`
- [x] добавить команды в `cmd/workflow_states.go`
- [x] запустить тесты — должны пройти

---

### Phase 5: Teams — расширенная функциональность

#### Task 19: Teams — Create/Update/Delete + Memberships

- [x] написать тесты для мутаций `teamCreate`, `teamUpdate`, `teamDelete`
- [x] добавить методы `CreateTeam`, `UpdateTeam`, `DeleteTeam` в mutations.go
- [x] написать тесты для cmd `teams create`, `teams update <ID>`, `teams delete <ID>`
- [x] добавить тип `TeamMembership` в queries.go (id, user, team, role)
- [x] написать тесты для `ListTeamMembers(teamID)`
- [x] добавить метод `ListTeamMembers(teamID)` в queries.go
- [x] написать тесты для мутаций `teamMembershipCreate`, `teamMembershipDelete`, `teamMembershipUpdate`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `teams members list <TEAM-ID>`, `teams members add`, `teams members remove`
- [x] добавить команды в `cmd/teams.go` (расширить) и `cmd/team_members.go`
- [x] запустить тесты — должны пройти

---

### Phase 6: Users

#### Task 20: Users — List + View + Viewer

- [x] написать тесты для `ListUsers()`, `GetUser(id)`, `GetViewer()` в queries.go
- [x] добавить методы в queries.go (расширить существующий FindUserByName)
- [x] написать тесты для cmd `users list`, `users view <ID>`, `users me`
- [x] добавить команды в `cmd/users.go`
- [x] запустить тесты — должны пройти

---

### Phase 7: Notifications

#### Task 21: Notifications — List + управление

- [x] добавить тип `Notification` в queries.go (id, type, readAt, createdAt, issue, comment, project)
- [x] написать тесты для `ListNotifications(limit, after)`, `GetNotificationsUnreadCount()`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `notificationMarkReadAll`, `notificationArchive`, `notificationUpdate`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `notifications list`, `notifications unread-count`, `notifications mark-read`, `notifications archive`
- [x] добавить команды в `cmd/notifications.go`
- [x] запустить тесты — должны пройти

---

### Phase 8: Webhooks

#### Task 22: Webhooks — CRUD

- [x] добавить тип `Webhook` в queries.go (id, url, enabled, secret, resourceTypes, team)
- [x] написать тесты для `ListWebhooks()`, `GetWebhook(id)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `webhookCreate`, `webhookUpdate`, `webhookDelete`, `webhookRotateSecret`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `webhooks list`, `webhooks view <ID>`, `webhooks create`, `webhooks update <ID>`, `webhooks delete <ID>`, `webhooks rotate-secret <ID>`
- [x] добавить команды в `cmd/webhooks.go`
- [x] запустить тесты — должны пройти

---

### Phase 9: Attachments

#### Task 23: Attachments — Query + основные мутации

- [x] добавить тип `Attachment` в queries.go (id, title, url, sourceType, issue)
- [x] написать тесты для `ListAttachments(issueID)`
- [x] добавить метод `ListAttachments(issueID)` в queries.go
- [x] написать тесты для мутаций `attachmentLinkURL`, `attachmentLinkGitHubPR`, `attachmentLinkGitHubIssue`, `attachmentLinkGitLabMR`, `attachmentDelete`, `attachmentUpdate`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `issues attachments list <ISSUE-ID>`, `issues attachments link-url`, `issues attachments link-github-pr`, `issues attachments delete <ID>`
- [x] добавить команды в `cmd/issue_attachments.go`
- [x] запустить тесты — должны пройти

---

### Phase 10: Documents

#### Task 24: Documents — CRUD + Search

- [x] добавить тип `Document` в queries.go (id, title, content, createdAt, updatedAt, project, creator)
- [x] написать тесты для `ListDocuments()`, `GetDocument(id)`, `SearchDocuments(query)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `documentCreate`, `documentUpdate`, `documentDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `documents list`, `documents view <ID>`, `documents create`, `documents update <ID>`, `documents delete <ID>`, `documents search <query>`
- [x] добавить команды в `cmd/documents.go`
- [x] запустить тесты — должны пройти

---

### Phase 11: Initiatives (инициативы)

#### Task 25: Initiatives — Query + Mutations

- [x] добавить тип `Initiative` в queries.go (id, name, description, status, owner, projects)
- [x] написать тесты для `ListInitiatives()`, `GetInitiative(id)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `initiativeCreate`, `initiativeUpdate`, `initiativeDelete`, `initiativeArchive`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `initiatives list`, `initiatives view <ID>`, `initiatives create`, `initiatives update <ID>`, `initiatives archive <ID>`
- [x] добавить команды в `cmd/initiatives.go`
- [x] запустить тесты — должны пройти

#### Task 26: Initiative Updates + Relations с проектами

- [x] написать тесты для `ListInitiativeUpdates(initiativeID)`
- [x] добавить метод в queries.go
- [x] написать тесты для мутаций `initiativeUpdateCreate`, `initiativeUpdateUpdate`, `initiativeUpdateArchive`
- [x] добавить методы в mutations.go
- [x] написать тесты для мутаций `initiativeToProjectCreate`, `initiativeToProjectDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `initiatives updates list`, `initiatives updates create`, `initiatives link-project`, `initiatives unlink-project`
- [x] добавить команды в `cmd/initiatives.go`
- [x] запустить тесты — должны пройти

---

### Phase 12: Roadmaps (дорожные карты)

#### Task 27: Roadmaps — CRUD + Projects

- [x] добавить тип `Roadmap` в queries.go (id, name, description, owner)
- [x] написать тесты для `ListRoadmaps()`, `GetRoadmap(id)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `roadmapCreate`, `roadmapUpdate`, `roadmapDelete`, `roadmapArchive`
- [x] добавить методы в mutations.go
- [x] написать тесты для мутаций `roadmapToProjectCreate`, `roadmapToProjectDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `roadmaps list`, `roadmaps view <ID>`, `roadmaps create`, `roadmaps update <ID>`, `roadmaps delete <ID>`, `roadmaps add-project`, `roadmaps remove-project`
- [x] добавить команды в `cmd/roadmaps.go`
- [x] запустить тесты — должны пройти

---

### Phase 13: Customers (CRM)

#### Task 28: Customers — CRUD

- [x] добавить тип `Customer`, `CustomerNeed`, `CustomerStatus`, `CustomerTier` в queries.go
- [x] написать тесты для `ListCustomers()`, `GetCustomer(id)`, `ListCustomerNeeds()`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `customerCreate`, `customerUpdate`, `customerDelete`, `customerUpsert`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `customers list`, `customers view <ID>`, `customers create`, `customers update <ID>`, `customers delete <ID>`
- [x] добавить команды в `cmd/customers.go`
- [x] запустить тесты — должны пройти

#### Task 29: Customer Needs + Statuses + Tiers

- [x] написать тесты для мутаций `customerNeedCreate`, `customerNeedUpdate`, `customerNeedDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для мутаций `customerStatusCreate`, `customerStatusUpdate`, `customerStatusDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для мутаций `customerTierCreate`, `customerTierUpdate`, `customerTierDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `customers needs list`, `customers needs create`, `customers statuses list`, `customers tiers list`
- [x] добавить команды в `cmd/customers.go`
- [x] запустить тесты — должны пройти

---

### Phase 14: Templates

#### Task 30: Templates — CRUD

- [x] добавить тип `Template` в queries.go (id, name, description, type, templateData)
- [x] написать тесты для `ListTemplates()`, `GetTemplate(id)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `templateCreate`, `templateUpdate`, `templateDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `templates list`, `templates view <ID>`, `templates create`, `templates update <ID>`, `templates delete <ID>`
- [x] добавить команды в `cmd/templates.go`
- [x] запустить тесты — должны пройти

---

### Phase 15: Organization

#### Task 31: Organization — View + Invites

- [x] добавить тип `Organization` в queries.go (id, name, urlKey, logoUrl, createdAt, periodUploadVolume)
- [x] написать тесты для `GetOrganization()`
- [x] добавить метод `GetOrganization()` в queries.go
- [x] добавить тип `OrganizationInvite` в queries.go
- [x] написать тесты для `ListOrganizationInvites()`
- [x] добавить метод в queries.go
- [x] написать тесты для мутаций `organizationInviteCreate`, `organizationInviteDelete`, `resendOrganizationInvite`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `org view`, `org invites list`, `org invites create`, `org invites delete`, `org invites resend`
- [x] добавить команды в `cmd/org.go`
- [x] запустить тесты — должны пройти

---

### Phase 16: Custom Views

#### Task 32: Custom Views — CRUD

- [x] добавить тип `CustomView` в queries.go (id, name, description, filters, icon, color, owner)
- [x] написать тесты для `ListCustomViews()`, `GetCustomView(id)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `customViewCreate`, `customViewUpdate`, `customViewDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `views list`, `views view <ID>`, `views create`, `views update <ID>`, `views delete <ID>`
- [x] добавить команды в `cmd/custom_views.go`
- [x] запустить тесты — должны пройти

---

### Phase 17: Favorites + Reactions + Emoji

#### Task 33: Favorites

- [x] добавить тип `Favorite` в queries.go (id, type, issue, project, cycle, label, custom view)
- [x] написать тесты для `ListFavorites()`
- [x] добавить метод в queries.go
- [x] написать тесты для мутаций `favoriteCreate`, `favoriteDelete`, `favoriteUpdate`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `favorites list`, `favorites add`, `favorites remove`
- [x] добавить команды в `cmd/favorites.go`
- [x] запустить тесты — должны пройти

#### Task 34: Comment Reactions + Emojis

- [x] написать тесты для мутаций `reactionCreate`, `reactionDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для `ListEmojis()`
- [x] добавить метод в queries.go
- [x] написать тесты для мутаций `emojiCreate`, `emojiDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `issues comment react <COMMENT-ID> <emoji>`, `emojis list`, `emojis create`, `emojis delete`
- [x] добавить команды в `cmd/reactions.go`, `cmd/emojis.go`
- [x] запустить тесты — должны пройти

---

### Phase 18: Releases [ALPHA]

#### Task 35: Releases — CRUD + Pipelines

- [x] добавить типы `Release`, `ReleasePipeline`, `ReleaseStage` в queries.go
- [x] написать тесты для `ListReleases()`, `ListReleasePipelines()`, `SearchReleases(query)`
- [x] добавить методы в queries.go
- [x] написать тесты для мутаций `releaseCreate`, `releaseUpdate`, `releaseDelete`, `releaseComplete`
- [x] добавить методы в mutations.go
- [x] написать тесты для мутаций `releasePipelineCreate`, `releasePipelineUpdate`, `releasePipelineDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `releases list`, `releases view <ID>`, `releases create`, `releases complete <ID>`, `releases pipelines list`, `releases pipelines create`
- [x] добавить команды в `cmd/releases.go`
- [x] запустить тесты — должны пройти

---

### Phase 19: Integrations + Git Automation

#### Task 36: Integrations — List + View

- [x] добавить тип `Integration` в queries.go (id, service, createdAt, team, organization)
- [x] написать тесты для `ListIntegrations()`
- [x] добавить метод в queries.go
- [x] написать тесты для мутаций `integrationArchive`, `integrationDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `integrations list`, `integrations delete <ID>`
- [x] добавить команды в `cmd/integrations.go`
- [x] запустить тесты — должны пройти

#### Task 37: Git Automation States + Branch Automation

- [x] добавить типы `GitAutomationState`, `GitAutomationTargetBranch` в queries.go
- [x] написать тесты для мутаций `gitAutomationStateCreate`, `gitAutomationStateUpdate`, `gitAutomationStateDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для мутаций `gitAutomationTargetBranchCreate`, `gitAutomationTargetBranchUpdate`, `gitAutomationTargetBranchDelete`
- [x] добавить методы в mutations.go
- [x] написать тесты для cmd `git-automation states list --team <KEY>`, `git-automation states create`, `git-automation states delete <ID>`
- [x] добавить команды в `cmd/git_automation.go`
- [x] запустить тесты — должны пройти

---

### Phase 20: Audit Log + Rate Limit + Misc

#### Task 38: Audit Entries

- [x] добавить тип `AuditEntry` в queries.go (id, type, actorId, createdAt, ip, country, metadata)
- [x] написать тесты для `ListAuditEntries(filter, limit, after)`
- [x] добавить метод в queries.go
- [x] написать тесты для `ListAuditEntryTypes()`
- [x] добавить метод в queries.go
- [x] написать тесты для cmd `audit list`, `audit types`
- [x] добавить команды в `cmd/audit.go`
- [x] запустить тесты — должны пройти

#### Task 39: Rate Limit + Time Schedules + Triage

- [x] написать тесты для `GetRateLimitStatus()` в queries.go
- [x] добавить метод в queries.go
- [x] добавить тип `TimeSchedule` в queries.go
- [x] написать тесты для `ListTimeSchedules()`
- [x] добавить метод в queries.go
- [x] написать тесты для мутаций `timeScheduleCreate`, `timeScheduleUpdate`, `timeScheduleDelete`
- [x] добавить методы в mutations.go
- [x] добавить тип `TriageResponsibility` в queries.go
- [x] написать тесты для `ListTriageResponsibilities(teamID)`
- [x] добавить метод в queries.go
- [x] написать тесты для cmd `rate-limit`, `time-schedules list/create/update/delete`, `triage-responsibilities list`
- [x] добавить команды в `cmd/misc.go`
- [x] запустить тесты — должны пройти

---

### Task 40: Финальная проверка

- [ ] проверить, что все команды зарегистрированы в root command (rootCmd.AddCommand)
- [ ] проверить, что все команды имеют --output/-o json флаг
- [ ] запустить полный набор тестов (`go test ./...`)
- [ ] запустить линтер (`go vet ./...`)
- [ ] проверить покрытие тестов (`go test -cover ./...`)

### Task 41: Обновить документацию

- [ ] обновить README.md: добавить все новые команды с примерами
- [ ] добавить секцию с полным списком команд в README.md

---

## Technical Details

### Паттерн реализации команды (на примере существующего кода)

```go
// Структура типа в queries.go
type Foo struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// Query метод в queries.go
func (c *Client) ListFoos() ([]Foo, error) {
    query := `query { foos { nodes { id name } } }`
    var result struct { Foos struct { Nodes []Foo `json:"nodes"` } `json:"foos"` }
    if err := c.Do(query, nil, &result); err != nil { return nil, err }
    return result.Foos.Nodes, nil
}

// Mutation метод в mutations.go
func (c *Client) CreateFoo(name string) (*Foo, error) {
    mutation := `mutation CreateFoo($input: FooCreateInput!) { fooCreate(input: $input) { success foo { id name } } }`
    var result struct { FooCreate struct { Foo Foo `json:"foo"`; Success bool `json:"success"` } `json:"fooCreate"` }
    if err := c.Do(mutation, map[string]any{"input": map[string]any{"name": name}}, &result); err != nil { return nil, err }
    if !result.FooCreate.Success { return nil, fmt.Errorf("fooCreate вернул success=false") }
    return &result.FooCreate.Foo, nil
}

// Команда в cmd/foos.go — использует паттерн из cmd/issues.go
```

### Паттерн тестов

```go
// Тест для client метода
func TestListFoos(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"foos": map[string]any{"nodes": []map[string]any{{"id": "foo1", "name": "Test Foo"}}}}})
    }))
    defer srv.Close()
    c := client.New("test-token")
    // переопределить endpoint...
    foos, err := c.ListFoos()
    // assert...
}
```

### Организация файлов

```
cmd/
  issue_labels.go          # issues labels
  issue_relations.go       # issues relations
  issue_archive.go         # issues archive/delete/unarchive
  issue_search.go          # issues search
  issue_subscribe.go       # issues subscribe/unsubscribe
  issue_batch.go           # issues batch-update
  issue_attachments.go     # issues attachments
  project_mutations.go     # projects create/update/delete/archive
  project_milestones.go    # projects milestones
  project_updates.go       # projects updates
  project_labels.go        # projects labels
  project_statuses.go      # projects statuses
  project_relations.go     # projects relations
  cycles.go                # cycles CRUD
  workflow_states.go       # workflow-states
  team_members.go          # teams members
  users.go                 # users
  notifications.go         # notifications
  webhooks.go              # webhooks
  documents.go             # documents
  initiatives.go           # initiatives
  roadmaps.go              # roadmaps
  customers.go             # customers (CRM)
  templates.go             # templates
  org.go                   # org
  custom_views.go          # views
  favorites.go             # favorites
  reactions.go             # comment reactions
  emojis.go                # emojis
  releases.go              # releases [ALPHA]
  integrations.go          # integrations
  git_automation.go        # git automation
  audit.go                 # audit log
  misc.go                  # rate-limit, time-schedules, triage
```

## Post-Completion

**Ручное тестирование:**
- Тестирование каждой команды с реальным Linear API токеном
- Проверка форматов вывода (table и json) для каждой команды
- Тестирование edge cases: пустые результаты, несуществующие ID, ошибки авторизации

**Внешние зависимости:**
- Некоторые команды требуют платного тарифа Linear (releases, roadmaps)
- Customers (CRM) может быть недоступен для всех тарифов
- Команды помеченные [ALPHA] или [INTERNAL] в API могут быть нестабильны
