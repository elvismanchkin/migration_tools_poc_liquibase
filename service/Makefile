.PHONY: all up up-logs down restart build clean logs status migrate migrate-sql migrate-flyway db service test-api init help

DC = docker-compose

all: up

up:
	$(DC) up -d

up-logs:
	$(DC) up

down:
	$(DC) down

restart:
	$(DC) restart

build:
	$(DC) build

clean:
	$(DC) down -v

logs:
	$(DC) logs -f

status:
	$(DC) ps

migrate: migrate-yaml

migrate-yaml:
	$(DC) up -d postgres
	$(DC) up migrations

migrate-sql:
	$(DC) -f docker-compose-liquibase-sql.yml up -d postgres
	$(DC) -f docker-compose-liquibase-sql.yml up liquibase-sql

migrate-flyway:
	$(DC) -f docker-compose-flyway.yml up -d postgres
	$(DC) up migrations

db:
	$(DC) up -d postgres

service:
	$(DC) up -d service

test-api:
	curl -v http://localhost:8080/health

init:
	mkdir -p postgres/config postgres/init-scripts \
	         migrations/migrations migrations/scripts \
	         service/templates service/static
	@echo "Directory structure created."

help:
	@echo "Commands:"; \
	echo "  make up           - Start all services (detached)"; \
	echo "  make up-logs      - Start services and display logs"; \
	echo "  make down         - Stop all services"; \
	echo "  make restart      - Restart services"; \
	echo "  make build        - Rebuild services"; \
	echo "  make clean        - Remove containers and volumes"; \
	echo "  make logs         - Tail service logs"; \
	echo "  make status       - Show service status"; \
	echo "  make migrate      - Run database migrations (YAML format)"; \
	echo "  make migrate-sql  - Run database migrations (SQL format)"; \
	echo "  make migrate-flyway - Run database migrations with Flyway"; \
	echo "  make db           - Start only the database"; \
	echo "  make service      - Start only the application service"; \
	echo "  make test-api     - Test API health endpoint"; \
	echo "  make init         - Create initial directory structure"
