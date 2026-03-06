# lcli — Linear CLI

CLI-утилита для работы с [Linear](https://linear.app) через GraphQL API.
Позволяет управлять задачами, проектами, командами и всеми сущностями Linear прямо из терминала.

## Установка

### Из исходников

```bash
git clone https://github.com/pavelvarganov/lcli
cd lcli
make install
```

### Сборка бинарника

```bash
make build
# бинарник: ./bin/lcli
```

## Конфигурация

Токен авторизации читается из следующих источников (в порядке приоритета):

1. Флаг `--token`
2. Переменная окружения `LINEAR_API_KEY`
3. Файл `~/.config/lcli/config.yaml`

### Сохранить токен

```bash
lcli auth login
# введите токен при запросе
```

### Проверить текущий токен

```bash
lcli auth status
```

Получить API-токен можно в настройках Linear: Settings > API > Personal API keys.

## Полный список команд

### Issues (задачи)

```bash
# список задач
lcli issues list
lcli issues list --team ENG --status "In Progress" --assignee "John Doe" --limit 50
lcli issues list --after <cursor>     # следующая страница

# просмотр задачи
lcli issue view ENG-123
lcli issue view ENG-123 -o json

# создать задачу
lcli issue create --title "Новая задача" --team ENG
lcli issue create --title "Баг" --team ENG --description "Описание" --priority 1
lcli issue create --title "Задача" --team ENG --assignee "Jane Doe"

# обновить задачу
lcli issue update ENG-123 --status "Done"
lcli issue update ENG-123 --assignee "Jane Doe" --priority 2
lcli issue update ENG-123 --title "Новый заголовок"

# поиск задач
lcli issues search "текст поиска"
lcli issues search "bug" --limit 20 -o json

# архивирование и удаление
lcli issues archive ENG-123
lcli issues unarchive ENG-123
lcli issues delete ENG-123

# подписка
lcli issues subscribe ENG-123
lcli issues unsubscribe ENG-123

# пакетное обновление
lcli issues batch-update --status "Done" --ids "ENG-1,ENG-2,ENG-3"
```

### Comments (комментарии)

```bash
# список комментариев
lcli issue comments ENG-123
lcli issue comments ENG-123 -o json

# добавить комментарий
lcli issue comment ENG-123 --body "Комментарий"

# обновить комментарий
lcli issue comment update <COMMENT-ID> --body "Новый текст"

# удалить комментарий
lcli issue comment delete <COMMENT-ID>

# resolve/unresolve комментарий
lcli issue comment resolve <COMMENT-ID>
lcli issue comment unresolve <COMMENT-ID>

# добавить реакцию на комментарий
lcli issue comment react <COMMENT-ID> 👍
```

### Issue Labels (метки задач)

```bash
# список меток
lcli issues labels list
lcli issues labels list -o json

# создать/обновить/удалить метку
lcli issues labels create --name "bug" --color "#ff0000" --team ENG
lcli issues labels update <LABEL-ID> --name "critical" --color "#cc0000"
lcli issues labels delete <LABEL-ID>

# добавить/убрать метку с задачи
lcli issues label add ENG-123 "bug"
lcli issues label remove ENG-123 "bug"
```

### Issue Relations (связи задач)

```bash
# список связей
lcli issues relations list ENG-123

# добавить/убрать связь
lcli issues relations add ENG-123 --related ENG-456 --type "blocks"
lcli issues relations remove <RELATION-ID>
```

### Issue Attachments (вложения)

```bash
# список вложений
lcli issues attachments list ENG-123

# прикрепить URL
lcli issues attachments link-url ENG-123 --url "https://example.com" --title "Ссылка"

# прикрепить GitHub PR
lcli issues attachments link-github-pr ENG-123 --url "https://github.com/org/repo/pull/1"

# удалить вложение
lcli issues attachments delete <ATTACHMENT-ID>
```

### Projects (проекты)

```bash
# список проектов
lcli projects list
lcli projects list -o json

# просмотр проекта
lcli projects view <PROJECT-ID>

# создать проект
lcli projects create --name "Новый проект" --team ENG
lcli projects create --name "Проект" --team ENG --description "Описание" --state "started"

# обновить проект
lcli projects update <PROJECT-ID> --name "Новое название" --state "completed"

# архивирование и удаление
lcli projects archive <PROJECT-ID>
lcli projects unarchive <PROJECT-ID>
lcli projects delete <PROJECT-ID>

# поиск проектов
lcli projects search "запрос"
```

### Project Milestones (вехи проекта)

```bash
lcli projects milestones list <PROJECT-ID>
lcli projects milestones create --project <PROJECT-ID> --name "MVP" --target-date "2025-06-01"
lcli projects milestones update <MILESTONE-ID> --name "v1.0"
lcli projects milestones delete <MILESTONE-ID>
```

### Project Updates (обновления проекта)

```bash
lcli projects updates list <PROJECT-ID>
lcli projects updates create --project <PROJECT-ID> --body "Прогресс за неделю" --health "onTrack"
lcli projects updates delete <UPDATE-ID>
```

### Project Labels (метки проекта)

```bash
lcli projects labels list
lcli projects labels create --name "feature" --color "#00ff00"
lcli projects labels update <LABEL-ID> --name "enhancement"
lcli projects labels delete <LABEL-ID>
```

### Project Statuses (статусы проекта)

```bash
lcli projects statuses list
lcli projects statuses create --name "Planning" --color "#0000ff" --type "backlog"
lcli projects statuses update <STATUS-ID> --name "In Flight"
```

### Project Relations (связи проектов)

```bash
lcli projects relations add --project <PROJECT-ID> --related <RELATED-ID> --type "blocks"
lcli projects relations remove <RELATION-ID>
```

### Cycles (циклы)

```bash
# список циклов команды
lcli cycles list --team ENG
lcli cycles list --team ENG -o json

# просмотр цикла
lcli cycles view <CYCLE-ID>

# создать цикл
lcli cycles create --team ENG --starts-at "2025-01-01" --ends-at "2025-01-14"
lcli cycles create --team ENG --name "Sprint 1" --starts-at "2025-01-01" --ends-at "2025-01-14"

# обновить цикл
lcli cycles update <CYCLE-ID> --name "Sprint 2"

# архивировать цикл
lcli cycles archive <CYCLE-ID>
```

### Workflow States (статусы задач)

```bash
# список статусов команды
lcli workflow-states list --team ENG
lcli workflow-states list --team ENG -o json

# создать статус
lcli workflow-states create --team ENG --name "Review" --type "started" --color "#ffaa00"

# обновить статус
lcli workflow-states update <STATE-ID> --name "Code Review" --color "#ff8800"

# архивировать статус
lcli workflow-states archive <STATE-ID>
```

### Teams (команды)

```bash
# список команд
lcli teams list
lcli teams list -o json

# создать команду
lcli teams create --name "Backend" --key "BE"
lcli teams create --name "Frontend" --key "FE" --description "Frontend team"

# обновить команду
lcli teams update <TEAM-ID> --name "Backend Team"

# удалить команду
lcli teams delete <TEAM-ID>
```

### Team Members (участники команды)

```bash
lcli teams members list <TEAM-ID>
lcli teams members add --team <TEAM-ID> --user <USER-ID>
lcli teams members remove <MEMBERSHIP-ID>
```

### Users (пользователи)

```bash
# список пользователей
lcli users list
lcli users list -o json

# просмотр пользователя
lcli users view <USER-ID>

# информация о текущем пользователе
lcli users me
```

### Notifications (уведомления)

```bash
# список уведомлений
lcli notifications list
lcli notifications list --limit 50 -o json

# количество непрочитанных
lcli notifications unread-count

# отметить все как прочитанные
lcli notifications mark-read

# архивировать уведомление
lcli notifications archive <NOTIFICATION-ID>
```

### Webhooks

```bash
# список вебхуков
lcli webhooks list
lcli webhooks list -o json

# просмотр вебхука
lcli webhooks view <WEBHOOK-ID>

# создать вебхук
lcli webhooks create --url "https://example.com/hook" --team ENG --resource-types "Issue,Comment"

# обновить вебхук
lcli webhooks update <WEBHOOK-ID> --url "https://example.com/new-hook" --enabled

# удалить вебхук
lcli webhooks delete <WEBHOOK-ID>

# ротация секрета
lcli webhooks rotate-secret <WEBHOOK-ID>
```

### Documents (документы)

```bash
# список документов
lcli documents list
lcli documents list -o json

# просмотр документа
lcli documents view <DOCUMENT-ID>

# поиск документов
lcli documents search "запрос"

# создать документ
lcli documents create --title "Дизайн API" --content "Содержимое"
lcli documents create --title "Дизайн" --project <PROJECT-ID> --content "Текст"

# обновить документ
lcli documents update <DOCUMENT-ID> --title "Новое название" --content "Новый текст"

# удалить документ
lcli documents delete <DOCUMENT-ID>
```

### Initiatives (инициативы)

```bash
# список инициатив
lcli initiatives list
lcli initiatives list -o json

# просмотр инициативы
lcli initiatives view <INITIATIVE-ID>

# создать инициативу
lcli initiatives create --name "Q1 Goals"
lcli initiatives create --name "Q1 Goals" --description "Цели на первый квартал"

# обновить инициативу
lcli initiatives update <INITIATIVE-ID> --name "Q1 2025 Goals"

# архивировать инициативу
lcli initiatives archive <INITIATIVE-ID>

# обновления инициативы
lcli initiatives updates list <INITIATIVE-ID>
lcli initiatives updates create --initiative <INITIATIVE-ID> --body "Прогресс"

# связать/отвязать проект
lcli initiatives link-project --initiative <INITIATIVE-ID> --project <PROJECT-ID>
lcli initiatives unlink-project --initiative <INITIATIVE-ID> --project <PROJECT-ID>
```

### Roadmaps (дорожные карты)

```bash
# список roadmap
lcli roadmaps list
lcli roadmaps list -o json

# просмотр roadmap
lcli roadmaps view <ROADMAP-ID>

# создать roadmap
lcli roadmaps create --name "2025 Roadmap"
lcli roadmaps create --name "2025 Roadmap" --description "Планы на 2025 год"

# обновить roadmap
lcli roadmaps update <ROADMAP-ID> --name "2025 Product Roadmap"

# удалить/архивировать roadmap
lcli roadmaps delete <ROADMAP-ID>

# добавить/убрать проект из roadmap
lcli roadmaps add-project --roadmap <ROADMAP-ID> --project <PROJECT-ID>
lcli roadmaps remove-project --roadmap <ROADMAP-ID> --project <PROJECT-ID>
```

### Customers (клиенты / CRM)

```bash
# список клиентов
lcli customers list
lcli customers list -o json

# просмотр клиента
lcli customers view <CUSTOMER-ID>

# создать клиента
lcli customers create --name "Acme Corp" --website "https://acme.com"

# обновить клиента
lcli customers update <CUSTOMER-ID> --name "Acme Corporation"

# удалить клиента
lcli customers delete <CUSTOMER-ID>

# потребности клиентов
lcli customers needs list
lcli customers needs create --customer <CUSTOMER-ID> --body "Нужна интеграция с Slack"

# статусы клиентов
lcli customers statuses list

# уровни клиентов
lcli customers tiers list
```

### Templates (шаблоны)

```bash
# список шаблонов
lcli templates list
lcli templates list -o json

# просмотр шаблона
lcli templates view <TEMPLATE-ID>

# создать шаблон
lcli templates create --name "Bug Report" --type issue

# обновить шаблон
lcli templates update <TEMPLATE-ID> --name "Bug Report v2"

# удалить шаблон
lcli templates delete <TEMPLATE-ID>
```

### Organization (организация)

```bash
# просмотр информации об организации
lcli org view

# список приглашений
lcli org invites list

# создать приглашение
lcli org invites create --email user@example.com

# удалить приглашение
lcli org invites delete <INVITE-ID>

# повторно отправить приглашение
lcli org invites resend <INVITE-ID>
```

### Custom Views (пользовательские представления)

```bash
# список представлений
lcli views list
lcli views list -o json

# просмотр представления
lcli views view <VIEW-ID>

# создать представление
lcli views create --name "My Issues" --description "Мои задачи"

# обновить представление
lcli views update <VIEW-ID> --name "My Active Issues"

# удалить представление
lcli views delete <VIEW-ID>
```

### Favorites (избранное)

```bash
# список избранного
lcli favorites list
lcli favorites list -o json

# добавить в избранное
lcli favorites add --issue ENG-123
lcli favorites add --project <PROJECT-ID>

# удалить из избранного
lcli favorites remove <FAVORITE-ID>
```

### Emojis (эмодзи)

```bash
# список эмодзи
lcli emojis list
lcli emojis list -o json

# создать кастомный эмодзи
lcli emojis create --name "custom_emoji" --url "https://example.com/emoji.png"

# удалить эмодзи
lcli emojis delete <EMOJI-ID>
```

### Releases (релизы) [ALPHA]

```bash
# список релизов
lcli releases list
lcli releases list -o json

# просмотр релиза
lcli releases view <RELEASE-ID>

# создать релиз
lcli releases create --pipeline <PIPELINE-ID> --title "v1.0.0"

# завершить релиз
lcli releases complete <RELEASE-ID>

# список пайплайнов релизов
lcli releases pipelines list

# создать пайплайн
lcli releases pipelines create --team ENG --name "Production"
```

### Integrations (интеграции)

```bash
# список интеграций
lcli integrations list
lcli integrations list -o json

# удалить интеграцию
lcli integrations delete <INTEGRATION-ID>
```

### Git Automation (git автоматизация)

```bash
# список состояний git автоматизации
lcli git-automation states list --team ENG

# создать состояние
lcli git-automation states create --team ENG --state <WORKFLOW-STATE-ID> --branch-pattern "feature/*"

# удалить состояние
lcli git-automation states delete <STATE-ID>
```

### Audit Log (журнал аудита)

```bash
# список записей аудита
lcli audit list
lcli audit list --limit 50 -o json

# список типов событий аудита
lcli audit types
```

### Misc (разное)

```bash
# проверить статус rate limit
lcli rate-limit

# список расписаний дежурств
lcli time-schedules list
lcli time-schedules list --team ENG

# создать расписание
lcli time-schedules create --team ENG --name "Дежурство"

# обновить расписание
lcli time-schedules update <SCHEDULE-ID> --name "Дежурство Backend"

# удалить расписание
lcli time-schedules delete <SCHEDULE-ID>

# ответственные за triage
lcli triage-responsibilities list --team ENG
```

## Флаги вывода

Флаг `--output json` (или `-o json`) поддерживается всеми командами:

```bash
lcli issues list -o json
lcli issue view ENG-123 -o json
lcli projects list -o json
lcli teams list -o json
lcli users list -o json
# ...и другие команды
```

## Приоритеты

| Значение | Описание    |
|----------|-------------|
| 0        | No priority |
| 1        | Urgent      |
| 2        | High        |
| 3        | Medium      |
| 4        | Low         |

При `issue create --priority 0` приоритет не задаётся (No priority нельзя указать явно при создании).
При `issue update --priority 0` приоритет сбрасывается в No priority.

## Разработка

```bash
# запуск тестов
make test

# линтер
make lint

# сборка
make build
```

## Требования

- Go 1.25+
- Linear API-токен
