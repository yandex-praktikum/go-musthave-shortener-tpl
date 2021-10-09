# go-url-shortener

Шаблон репозитория для практического трек "Веб-разработка на Go"

## Начало работы

1. Склонируйте репозиторий в любую подходящую директорию на вашем компьютере
2. В корне репозитория выполните команду `go mod init <name>` (где `<name>` - адрес вашего репозитория на Github без префикса `https://`) для создания модуля

## Обновление шаблона

Чтобы иметь возможность получать обновления автотестов и других частей шаблона выполните следующую команды:

```bash
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

Для обновления кода автотестов выполните команду:

```bash
git fetch template && git checkout template/main .github
```

затем добавьте полученные изменения в свой репозиторий.

## Запуск тестов

```go test -v ./...```

## Запуск сервиса

```go run cmd/shortener/main.go```
