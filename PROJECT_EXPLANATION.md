# Template System Architecture

This document explains the architecture of the Template System, a 3-component system designed to demonstrate database migrations with Liquibase.

## Overview

The system is split into three separate components, each with its own responsibility:

1. **PostgreSQL Database** - Stores all template data
2. **Migrations Service** - Handles database schema changes using Liquibase
3. **Template Web Service** - Provides the user interface and business logic

Each component lives in its own directory, has its own Dockerfile, but they are orchestrated together using a single Docker Compose file.

## Component Details

### 1. PostgreSQL Database (postgres/)

A standard PostgreSQL 15 database with:
- Custom configuration in `postgres/config/postgresql.conf`
- Initialization scripts in `postgres/init-scripts/`
- Data persistence through a Docker volume

The database initializes with the `template_service` and `audit` schemas that our application uses.

### 2. Migrations Service (migrations/)

This is a specialized service that:
- Uses the official Liquibase Alpine Docker image
- Waits for the database to be available before running migrations
- Applies all pending migrations defined in YAML files
- Exits successfully once migrations are complete

The migrations are defined in YAML format and stored in `migrations/migrations/`. The service follows a dependency pattern in Docker Compose where the template service waits for migrations to complete before starting.

### 3. Template Web Service (service/)

A Go web application that:
- Provides a user interface for managing templates
- Implements HTMX for interactive frontend functionality
- Generates PDFs from templates using wkhtmltopdf
- Connects to the PostgreSQL database

The service follows a modular structure:
- `db/` - Database connection and management
- `models/` - Data structures and database operations
- `handlers/` - HTTP request handlers
- `templates/` - HTML templates for the frontend
- `static/` - Static assets (CSS, JS, images)
- `main.go` - Application entry point

The service waits for both the database to be ready and migrations to be complete before starting.

## Docker Compose Integration

The `docker-compose.yml` file orchestrates the three components:

1. **Dependency Chain**: 
   - PostgreSQL starts first
   - Migrations service starts after PostgreSQL is healthy
   - Template service starts after migrations service completes successfully

2. **Network Configuration**:
   - All services are on the same Docker network
   - Only the web service exposes a port to the host

3. **Environment Variables**:
   - Environment variables are defined in the Docker Compose file
   - Each service gets only the variables it needs

## Running the System

You can run the entire system with a single command:

```bash
docker-compose up -d
```

This will:
1. Start PostgreSQL database
2. Run database migrations
3. Start the template web service

## Database Migration Design

The database migration approach has these key characteristics:

1. **Independence**: Migrations run as a separate service, not embedded in the application
2. **Versioning**: Each migration is numbered and tracked in the database
3. **Idempotence**: Migrations use preConditions to ensure they can be run multiple times safely
4. **Rollback Support**: Most migrations include rollback instructions
5. **Environment Awareness**: Support for environment-specific migrations (dev/test/prod)

## Integration with CI/CD

This architecture is ideal for integration with CI/CD pipelines:

1. Build and push the Docker images to a registry
2. Deploy the system with Docker Compose or Kubernetes
3. Run the migrations service to update the database schema
4. Restart or redeploy the service to use the updated schema

Example deployment script:
```bash
#!/bin/bash
# Pull the latest images
docker-compose pull

# Run migrations only
docker-compose up -d --no-deps migrations

# Wait for migrations to complete
docker wait template-system_migrations_1

# Restart the template service
docker-compose restart service
```

## Benefits of This Architecture

1. **Separation of Concerns**: 
   - Database, migrations, and application logic are cleanly separated
   - Each component can be developed and tested independently

2. **Operational Flexibility**:
   - Run migrations separately from the application
   - Easily roll back to previous versions if needed
   - Test migrations in lower environments before production

3. **Scalability**:
   - Database can be scaled independently
   - Multiple instances of the web service can be deployed
   - Migrations only need to run once, no matter how many service instances