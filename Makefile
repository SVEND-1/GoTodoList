include .env
export

PROJECT_ROOT := $(shell pwd)

env-up:
	docker compose up todoapp-postgres -d

env-down:
	docker compose down todoapp-postgres

env-cleanup:
	@read -p "Уверен что хочешь все очистить?[y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down todoapp-postgres && \
	  rm -rf out/pgdata && \
	  echo "Очистка завершена"; \
	else \
	  echo "Очистка отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Ошибка: параметр seq не указан"; \
		echo "Использование: make migrate-create seq=название_миграции"; \
		exit 1; \
	else \
		echo "Создание миграции: $(seq)"; \
		docker compose run --rm todoapp-postgres-migrate create -ext sql -dir /migrations -seq $(seq); \
	fi

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Ошибка: параметр action не указан"; \
		echo "Использование: migrate-action action=up"; \
		exit 1; \
	else \
		docker compose run --rm todoapp-postgres-migrate \
			-path /migrations \
			-database postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5433/$(POSTGRES_DB)?sslmode=disable \
			$(action); \
	fi