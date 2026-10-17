.PHONY: help fmt lint test check run db-start migrate unmigrate docs docs-lint

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

check: fmt lint test docs-lint ## Всё перед коммитом: формат, линтеры, тесты, спецификация API

run: ## Запустить auth-сервис
	$(ENV) go run ./auth/cmd

db-start: ## Поднять контейнер с Postgres
	docker start stream-postgres

migrate: ## Применить миграции
	$(ENV) goose -dir auth/migrations postgres "$$DATABASE_URL" up

unmigrate: ## Откатить последнюю миграцию
	$(ENV) goose -dir auth/migrations postgres "$$DATABASE_URL" down

docs: ## Открыть Swagger UI для auth/api/openapi.yaml на http://localhost:8081
	docker run --rm -p 127.0.0.1:8081:8080 -e SWAGGER_JSON=/spec/openapi.yaml -v "$(CURDIR)/auth/api:/spec:ro" swaggerapi/swagger-ui:v5.17.14

docs-lint: ## Проверить спецификацию API
	docker run --rm -v "$(CURDIR)/auth/api:/spec:ro" redocly/cli:1.34.3 lint --format=stylish /spec/openapi.yaml
