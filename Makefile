include .env.example

COMPOSE=docker compose --env-file .env.dev -f compose.yml -p ${PROJECT_NAME}
PROD=docker compose --env-file .env.prod -f compose.yml -p ${PROJECT_NAME}


.PHONY: build start stop front vps deploy

build:
	${COMPOSE} build
start: build
	${COMPOSE} up
stop:
	${COMPOSE} down

front:
	docker exec -it ${PROJECT_NAME}-app /bin/bash

generate: sqlc
sqlc:
	docker exec -it ${PROJECT_NAME}-api /bin/bash -c "sqlc generate"

migration: # usage: make migration name="create_users_table". Creates a new migration file
	@if [ -z "$(name)" ]; then \
		echo "Error: name argument is required. Usage: make migration name=\"create_users_table\""; \
		exit 1; \
	fi
	docker exec -it ${PROJECT_NAME}-api /bin/bash -c "migrate create -ext sql -dir db/migrations -seq $(name)"
migrate-up:
	docker exec -it ${PROJECT_NAME}-api /bin/bash -c 'migrate -path db/migrations -database "$$DATABASE_URL" up'
migrate-down:
	docker exec -it ${PROJECT_NAME}-api /bin/bash -c 'migrate -path db/migrations -database "$$DATABASE_URL" down'
migrate-drop:
	docker exec -it ${PROJECT_NAME}-api /bin/bash -c 'migrate -path db/migrations -database "$$DATABASE_URL" drop'
seed:
	docker exec -it ${PROJECT_NAME}-api /bin/bash -c 'psql "$$DATABASE_URL" -f db/seed/seed.sql'

db-url:
	@docker exec -it ${PROJECT_NAME}-api /bin/bash -c 'echo "$$DATABASE_URL"'