# lcli — Linear CLI

CLI-утилита для работы с [Linear](https://linear.app) через GraphQL API.
Позволяет просматривать, создавать, обновлять задачи и оставлять комментарии прямо из терминала.

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

## Команды

### Issues

```bash
# список задач
lcli issues list
lcli issues list --team ENG --status "In Progress" --assignee "John Doe" --limit 50

# детали задачи
lcli issue view ENG-123

# создать задачу
lcli issue create --title "Новая задача" --team ENG
lcli issue create --title "Баг в авторизации" --team ENG --description "Описание" --priority 1

# обновить задачу
lcli issue update ENG-123 --status "Done"
lcli issue update ENG-123 --assignee "Jane Doe" --priority 2

# добавить комментарий
lcli issue comment ENG-123 --body "Комментарий к задаче"

# список комментариев
lcli issue comments ENG-123
```

### Projects

```bash
lcli projects list
```

### Teams

```bash
lcli teams list
```

### Флаги вывода

```bash
# JSON-вывод для machine-readable обработки
lcli issues list --output json
lcli projects list --output json
```

## Приоритеты

| Значение | Описание    |
|----------|-------------|
| 0        | No priority |
| 1        | Urgent      |
| 2        | High        |
| 3        | Medium      |
| 4        | Low         |

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

- Go 1.21+
- Linear API-токен
