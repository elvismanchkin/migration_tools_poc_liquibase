# Template Service with Database Migrations

A minimal proof-of-concept system with three separate components:

1. PostgreSQL database
2. Database migration service (using either Liquibase or Flyway)
3. Template web service (Go)

## Project Structure

```
template-system/
├── docker-compose-liquibase.yml
├── docker-compose-flyway.yml
├── postgres/                    # PostgreSQL configuration
│   ├── Dockerfile
│   ├── init-scripts/
│   └── config/
├── liquibase/                   # Liquibase database migrations
│   ├── Dockerfile
│   ├── migrations/              # Migration files
│   │   ├── v1_initial_schema.yaml
│   │   ├── v2_create_audit_tables.yaml
│   │   ├── v3_create_templates_tables.yaml
│   │   ├── v4_add_configuration_tables.yaml
│   │   └── dev/                 # Environment-specific migrations
│   │       └── v20250228_add_test_data.yaml
│   ├── scripts/              # Liquibase execution scripts
│   │   ├── liquibase-migrate.sh
│   │   ├── liquibase-validate.sh
│   │   └── liquibase-api.sh
│   └── master-changelog.yaml # Liquibase master changelog
├── flyway/                   # Flyway migration service
│   ├── Dockerfile
│   ├── entrypoint.sh
│   ├── sql/                  # Flyway SQL migrations
│   │   ├── V1__Initial_Schema.sql
│   │   ├── V2__Create_Audit_Tables.sql
│   │   ├── V3__Create_Templates_Tables.sql
│   │   ├── V4__Add_Configuration_Tables.sql
│   │   └── R__Dev_Data.sql   # Repeatable migration for dev data
│   └── README.md
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

## Start Here

```bash
# Start all services with Liquibase migration
docker-compose up -d

# Or use Flyway for migrations instead
docker-compose -f docker-compose-flyway.yml up -d

# Access the web application
open http://localhost:8080
```

## Components

### 1. PostgreSQL Database (postgres/)

Standalone PostgreSQL database with custom configuration.

### 2. Database Migrations

The project supports two migration tools:

#### Liquibase (migrations/)

Uses YAML format for migrations with a schema-driven approach.

- Waits for the database to be ready
- Applies migrations using Liquibase
- Exits after completing migrations

#### Flyway (flyway/)

Uses plain SQL migrations with a naming convention-based approach.

- Waits for the database to be ready
- Applies migrations using Flyway
- Exits after completing migrations

You can choose either tool based on your preference.

### 3. Template Service (service/)

Go web service that:

- Waits for both the database and migrations to be complete
- Provides a web interface for managing templates
- Allows creating, viewing, and rendering templates
- Generates PDFs from templates

## Running with Docker Compose

The `docker-compose.yml` file orchestrates all three services:

```bash
# Using Liquibase
docker-compose -f docker-compose-flyway.yml up --build

# Using Flyway
docker-compose -f docker-compose-flyway.yml up --build

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Restart a specific service
docker-compose restart service
```

```bash
# Shutdown and cleanup

# For Liquibase
docker-compose -f docker-compose-liquibase.yml down -v && docker volume prune -f

# For Flyway
docker-compose -f docker-compose-flyway.yml down -v && docker volume prune -f
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