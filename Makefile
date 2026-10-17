.PHONY: help fmt lint test check run db-start migrate unmigrate

# Переменные из .env подхватываются только в целях, которым нужна база.
ENV = set -a; . ./.env; set +a;

help: ## Показать список команд
	@grep -E '^[a-z-]+:.*## ' $(MAKEFILE_LIST) | awk -F ':.*## ' '{printf "  %-10s %s\n", $$1, $$2}'

fmt: ## Отформатировать код (gofumpt + goimports)
	golangci-lint fmt ./...

lint: ## Запустить линтеры
	golangci-lint run ./...

test: ## Запустить тесты
	go test ./...

check: fmt lint test ## Всё перед коммитом: формат, линтеры, тесты

run: ## Запустить auth-сервис
	$(ENV) go run ./auth/cmd

db-start: ## Поднять контейнер с Postgres
	docker start stream-postgres

migrate: ## Применить миграции
	$(ENV) goose -dir auth/migrations postgres "$$DATABASE_URL" up

unmigrate: ## Откатить последнюю миграцию
	$(ENV) goose -dir auth/migrations postgres "$$DATABASE_URL" down
