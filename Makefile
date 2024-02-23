.PHONY: postgres-up
postgres-up:
	docker run --name postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=password -e POSTGRES_DB=quotation_db -p 5433:5432 -d postgres:latest

.PHONY: build-app
build-app:
	docker pull alpine:3.19
	docker pull golang:1.21-alpine3.19

.PHONY: compose-up
compose-up:
	docker compose -p app -f docker-compose.yaml up --remove-orphans --abort-on-container-exit --exit-code-from=quotation-service

.PHONY: compose-down
compose-down:
	docker compose -p app down -v