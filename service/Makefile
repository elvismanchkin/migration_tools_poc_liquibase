.PHONY: up down restart build clean logs status migrate

# Default target
all: up

# Start all services
up:
	docker-compose up -d

# Start all services with logs
up-logs:
	docker-compose up

# Stop all services
down:
	docker-compose down

# Restart all services
restart:
	docker-compose restart

# Rebuild all services
build:
	docker-compose build

# Clean up volumes and containers
clean:
	docker-compose down -v

# View logs
logs:
	docker-compose logs -f

# Check service status
status:
	docker-compose ps

# Run database migrations manually
migrate:
	docker-compose up -d postgres
	docker-compose up migrations

# Just run the database
db:
	docker-compose up -d postgres

# Just run the service
service:
	docker-compose up -d service

# Test the API
test-api:
	curl -v http://localhost:8080/health

# Create directory structure
init:
	mkdir -p postgres/config postgres/init-scripts
	mkdir -p migrations/migrations migrations/scripts
	mkdir -p service/templates service/static
	@echo "Directory structure created successfully!"

# Help information
help:
	@echo "Template Service Management Commands:"
	@echo "  make up            - Start all services in detached mode"
	@echo "  make up-logs       - Start all services and show logs"
	@echo "  make down          - Stop all services"
	@echo "  make restart       - Restart all services"
	@echo "  make build         - Rebuild all services"
	@echo "  make clean         - Clean up volumes and containers"
	@echo "  make logs          - View logs from all services"
	@echo "  make status        - Check service status"
	@echo "  make migrate       - Run database migrations manually"
	@echo "  make db            - Start only the database service"
	@echo "  make service       - Start only the template service"
	@echo "  make test-api      - Test the API health endpoint"
	@echo "  make init          - Create initial directory structure"
	@echo "  make help          - Show this help message"