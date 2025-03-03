# Template Service with Database Migrations

A minimal proof-of-concept system with three separate components:

1. PostgreSQL database
2. Database migration service
3. Template web service

## Project Structure

```
template-system/
├── docker-compose.yml        # Main Docker Compose file for all services
├── postgres/                 # PostgreSQL configuration
│   ├── Dockerfile
│   ├── init-scripts/
│   └── config/
├── migrations/               # Database migration service
│   ├── Dockerfile
│   ├── migrations/           # Migration files
│   │   ├── v1_initial_schema.yaml
│   │   ├── v2_create_audit_tables.yaml
│   │   ├── v3_create_templates_tables.yaml
│   │   ├── v4_add_configuration_tables.yaml
│   │   └── dev/             # Environment-specific migrations
│   │       └── v20250228_add_test_data.yaml
│   ├── scripts/              # Liquibase execution scripts
│   │   ├── liquibase-migrate.sh
│   │   ├── liquibase-validate.sh
│   │   └── liquibase-api.sh
│   └── master-changelog.yaml # Liquibase master changelog
├── service/                  # Go template service
│   ├── Dockerfile
│   ├── templates/            # HTML templates for the frontend
│   │   ├── layout.html
│   │   ├── templates-list.html
│   │   ├── template-form.html
│   │   ├── template-view.html
│   │   └── template-rendered.html
│   ├── static/               # Static assets (CSS, JS, images)
│   ├── main.go               # Go service entry point
│   ├── go.mod
│   └── go.sum
└── README.md                 # This file
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/template-system.git
cd template-system

# Start all services
docker-compose up -d

# Access the web application
open http://localhost:8080
```

## Components

### 1. PostgreSQL Database (postgres/)

Standalone PostgreSQL database with custom configuration.

### 2. Database Migrations (migrations/)

Separate service using the official Liquibase Docker image for database migrations. This service:

- Waits for the database to be ready
- Applies migrations using Liquibase
- Exits after completing migrations

This service uses the official Liquibase Docker image (`liquibase/liquibase:4.20-alpine`) for simplicity and
maintainability.

### 3. Template Service (service/)

Go web service that:

- Waits for both the database and migrations to be complete
- Provides a web interface for managing templates
- Allows creating, viewing, and rendering templates
- Generates PDFs from templates

## Running with Docker Compose

The `docker-compose.yml` file orchestrates all three services:

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Restart a specific service
docker-compose restart service
```

## Development

For development purposes, you can:

1. Run just the database:
   ```bash
   docker-compose up -d postgres
   ```

2. Run migrations manually:
   ```bash
   cd migrations
   ./scripts/liquibase-api.sh migrate dev
   ```

3. Run the service locally:
   ```bash
   cd service
   go run main.go
   ```

## Integration with Rundeck

To use this with Rundeck, create a job that:

1. Pulls the latest code
2. Runs the migrations service
3. Restarts the template service

Example Rundeck job:

```bash
#!/bin/bash
cd /path/to/template-system
git pull
docker-compose up -d --no-deps migrations
docker-compose restart service
```

## License

Internal use only
