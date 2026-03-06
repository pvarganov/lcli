# Релиз и установка через Homebrew

## Требования

- [`gh`](https://cli.github.com/) — GitHub CLI (`brew install gh && gh auth login`)
- Публичный репозиторий `pvarganov/lcli` на GitHub
- Публичный tap-репозиторий `pvarganov/homebrew-tap` на GitHub

## Структура tap-репозитория

```
homebrew-tap/
└── Formula/
    └── lcli.rb
```

Создать вручную:
```bash
mkdir -p ../homebrew-tap/Formula
cd ../homebrew-tap && git init && git remote add origin git@github.com:pvarganov/homebrew-tap.git
```

## Создать релиз

```bash
# Из корня репозитория lcli
./scripts/release.sh v0.1.0

# Или через make
make release VERSION=v0.1.0
```

Скрипт:
1. Собирает бинарники для darwin/amd64, darwin/arm64, linux/amd64, linux/arm64
2. Упаковывает в `.tar.gz`
3. Создаёт GitHub Release через `gh`
4. Генерирует `Formula/lcli.rb` и пушит в `../homebrew-tap`

## Установка пользователями

```bash
brew tap pvarganov/tap
brew install lcli
```
