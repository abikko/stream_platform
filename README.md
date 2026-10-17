# stream_platform
Пет-проект который практический боевой. ИИшка не пишет код - а советует как реализовать правильно, так же он коммитит изменения и обновляет README.

Пет-проект стриминговой платформы на Go. Сейчас есть auth-сервис с регистрацией и входом пользователей.

## Требования

- Go 1.27 или новее.
- Docker (для PostgreSQL).
- [goose](https://github.com/pressly/goose) для миграций.
- [golangci-lint](https://golangci-lint.run) v2: `brew install golangci-lint`.

## Настройка

Скопировать `.env.example` в `.env` и задать `PASSWORD`:

```sh
cp .env.example .env
```

## База данных

Первый запуск (создаёт контейнер `stream-postgres` с данными в томе `stream-postgres-data`):

```sh
set -a; source .env; set +a
docker run -d --name stream-postgres \
  -e POSTGRES_USER=stream \
  -e POSTGRES_PASSWORD="$PASSWORD" \
  -e POSTGRES_DB=stream_platform \
  -p 127.0.0.1:5432:5432 \
  -v stream-postgres-data:/var/lib/postgresql/data \
  postgres:17
```

Дальше контейнер нужно только запускать (после перезапуска Docker он остаётся остановленным):

```sh
docker start stream-postgres
```

Миграции:

```sh
set -a; source .env; set +a
goose -dir auth/migrations postgres "$DATABASE_URL" up
```

Миграции запускать после каждого `git pull`: схема базы не обновляется сама.

## Запуск auth-сервиса

В VS Code: конфигурация «Auth service» (F5), переменные берутся из `.env`.

Из терминала:

```sh
set -a; source .env; set +a
go run ./auth/cmd
```

Сервер слушает `http://localhost:8080`.

```sh
curl -i -X POST localhost:8080/register -d '{"email":"user@example.com","password":"secret123"}'
```

Ответы: `201` с `{"id":"..."}`, `400` на некорректный запрос, `409` если email уже занят.

```sh
curl -i -X POST localhost:8080/login -d '{"email":"user@example.com","password":"secret123"}'
```

Ответы: `200` при верном пароле, `400` на некорректный запрос, `401` если email или пароль неверны.

## Проверка

```sh
make check
```

Форматирует код, запускает линтеры (`.golangci.yml`) и тесты. Остальные команды: `make help`.

Правила для AI-агента находятся в `AGENTS.md`.
