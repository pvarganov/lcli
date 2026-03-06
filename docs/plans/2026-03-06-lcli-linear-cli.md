# lcli — CLI-утилита для Linear

## Overview

Полнофункциональная CLI-утилита `lcli` для работы с Linear через GraphQL API.
Позволяет просматривать, создавать, обновлять задачи и оставлять комментарии прямо из терминала.
Аутентификация через Personal API Token, вывод в виде таблицы (аналогично `gh` CLI).

## Context (from discovery)

- Проект: новый, директория пустая
- Язык: Go
- CLI-фреймворк: cobra
- GraphQL: простой HTTP-клиент (без сторонних GraphQL-библиотек)
- API: https://api.linear.app/graphql (Bearer-токен в заголовке)
- Формат вывода: таблица по умолчанию

## Development Approach

- **Testing approach**: Regular (code first, then tests)
- Завершать каждую задачу полностью перед переходом к следующей
- Небольшие, сфокусированные изменения
- **CRITICAL: каждая задача ДОЛЖНА включать тесты** для изменённого кода
- **CRITICAL: все тесты должны проходить перед началом следующей задачи**
- Запускать `go test ./...` после каждого изменения

## Testing Strategy

- **Unit tests**: обязательны для каждой задачи
- Моки для HTTP-клиента (интерфейс + тестовый сервер `httptest`)
- Тесты покрывают успешные сценарии и обработку ошибок

## Progress Tracking

- Отмечать выполненные пункты `[x]` сразу при завершении
- Добавлять новые задачи с префиксом ➕
- Документировать блокеры с префиксом ⚠️

## What Goes Where

- **Implementation Steps** (`[ ]`): изменения кода, тесты, документация
- **Post-Completion**: ручное тестирование, публикация бинарника

## Implementation Steps

### Task 1: Инициализация Go-модуля и структуры проекта

- [x] Выполнить `go mod init github.com/pavelvarganov/lcli`
- [x] Создать `main.go` с точкой входа
- [x] Создать структуру директорий:
  - `cmd/` — cobra-команды
  - `internal/client/` — GraphQL HTTP-клиент
  - `internal/config/` — конфиг (токен, org)
  - `internal/format/` — форматирование вывода (таблицы)
- [x] Добавить зависимости: `cobra`, `charmbracelet/lipgloss` (опционально для цветов), `olekukonko/tablewriter`
- [x] Создать корневую cobra-команду в `cmd/root.go` с флагом `--token`
- [x] Написать тест для `cmd/root.go` (проверка инициализации команды)
- [x] Запустить `go test ./...` — должен пройти

### Task 2: GraphQL HTTP-клиент

- [x] Создать `internal/client/client.go` с интерфейсом `LinearClient`
- [x] Реализовать метод `Do(query string, variables map[string]any, result any) error`
- [x] Добавить Bearer-токен в заголовок `Authorization`
- [x] Обрабатывать ошибки GraphQL (поле `errors` в ответе)
- [x] Написать тест с `httptest.NewServer` для успешного запроса
- [x] Написать тест для ошибки GraphQL и ошибки сети
- [x] Запустить `go test ./...` — должен пройти

### Task 3: Конфигурация (токен и workspace)

- [x] Создать `internal/config/config.go` для чтения токена из:
  1. Флаг `--token`
  2. Переменная окружения `LINEAR_API_KEY`
  3. Файл `~/.config/lcli/config.yaml`
- [x] Команда `lcli auth login` — сохраняет токен в конфиг-файл
- [x] Команда `lcli auth status` — показывает текущий токен (замаскированный)
- [x] Написать тесты для чтения конфига из env и файла
- [x] Запустить `go test ./...` — должен пройти

### Task 4: Команды для Issues — просмотр

- [x] Создать `cmd/issues.go` с субкомандой `lcli issues list`
  - Флаги: `--assignee`, `--status`, `--team`, `--limit` (default: 25)
  - Вывод: таблица (ID, Title, Status, Assignee, Priority, Updated)
- [x] Создать `cmd/issue_view.go` с командой `lcli issue view <ID>`
  - Вывод: детальная информация + описание
- [x] GraphQL-запросы в `internal/client/queries.go`
- [x] Написать тесты для форматирования таблицы
- [x] Написать тест с моком клиента для `issues list`
- [x] Запустить `go test ./...` — должен пройти

### Task 5: Команды для Issues — создание и обновление

- [ ] Команда `lcli issue create` с флагами:
  - `--title` (required), `--description`, `--team` (required), `--assignee`, `--priority`
- [ ] Команда `lcli issue update <ID>` с флагами:
  - `--status`, `--assignee`, `--priority`, `--title`
- [ ] GraphQL-мутации в `internal/client/mutations.go`
- [ ] Написать тест создания issue (мок клиент)
- [ ] Написать тест обновления issue (мок клиент)
- [ ] Запустить `go test ./...` — должен пройти

### Task 6: Команды для комментариев

- [ ] Команда `lcli issue comment <ID> --body "текст"` — добавить комментарий
- [ ] Команда `lcli issue comments <ID>` — список комментариев к задаче
- [ ] GraphQL-запрос и мутация для комментариев
- [ ] Написать тесты для обеих команд
- [ ] Запустить `go test ./...` — должен пройти

### Task 7: Команды для Projects и Teams

- [ ] Команда `lcli projects list` — список проектов
- [ ] Команда `lcli teams list` — список команд (для использования в `--team`)
- [ ] Написать тесты
- [ ] Запустить `go test ./...` — должен пройти

### Task 8: Форматирование вывода и UX

- [ ] Добавить флаг `--output json` для machine-readable вывода
- [ ] Цветовое выделение статусов (In Progress — синий, Done — зелёный, Cancelled — серый)
- [ ] Пагинация для больших списков (`--limit` + `after`-cursor)
- [ ] Написать тесты для JSON-вывода
- [ ] Запустить `go test ./...` — должен пройти

### Task 9: Verify acceptance criteria

- [ ] Проверить все требования из Overview реализованы
- [ ] Проверить обработку edge cases (нет токена, нет сети, неверный ID)
- [ ] Запустить `go test ./...` — все тесты проходят
- [ ] Запустить `go vet ./...` — нет предупреждений
- [ ] Запустить `golangci-lint run` если установлен
- [ ] Проверить покрытие тестами: `go test -cover ./...`

### Task 10: [Final] Документация

- [ ] Обновить/создать `README.md`: установка, конфигурация, примеры команд
- [ ] Добавить `Makefile` с целями: `build`, `test`, `lint`, `install`
- [ ] Описать схему GraphQL-запросов в `docs/`

*Note: ralphex automatically moves completed plans to `docs/plans/completed/`*

## Technical Details

### Структура команд

```
lcli
├── auth
│   ├── login       # сохранить токен
│   └── status      # показать текущий токен
├── issues
│   └── list        # список задач
├── issue
│   ├── view <ID>   # детали задачи
│   ├── create      # создать задачу
│   ├── update <ID> # обновить задачу
│   ├── comment <ID> # добавить комментарий
│   └── comments <ID> # список комментариев
├── projects
│   └── list
└── teams
    └── list
```

### GraphQL Endpoint

- URL: `https://api.linear.app/graphql`
- Auth: `Authorization: Bearer <token>`
- Метод: POST, Content-Type: application/json

### Структура запроса

```json
{
  "query": "query { ... }",
  "variables": { "key": "value" }
}
```

### Приоритеты Linear

- 0: No priority
- 1: Urgent
- 2: High
- 3: Medium
- 4: Low

## Post-Completion

**Ручное тестирование:**
- Протестировать все команды с реальным Linear API-токеном
- Проверить вывод таблиц в разных терминалах
- Проверить поведение при отсутствии сети

**Публикация (опционально):**
- Добавить GoReleaser для сборки бинарников под macOS/Linux/Windows
- Опубликовать в Homebrew tap
