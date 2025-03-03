# Database Migrations

This component handles database migrations using Liquibase. It uses the official Liquibase Docker image and applies
migrations to the PostgreSQL database.

## Structure

```
migrations/
├── Dockerfile                # Docker image definition for migrations
├── master-changelog.yaml     # Master changelog that includes all migrations
├── migrations/               # Migration files
│   ├── v1_initial_schema.yaml 
│   ├── v2_create_audit_tables.yaml
│   ├── v3_create_templates_tables.yaml
│   ├── v4_add_configuration_tables.yaml
│   └── dev/                  # Environment-specific migrations
│       └── v20250228_add_test_data.yaml
└── README.md                 # This file
```

## How It Works

The migration service:

1. Waits for the database to be ready
2. Generates a Liquibase properties file based on environment variables
3. Runs Liquibase update to apply all pending migrations
4. Terminates when migrations are complete

## Running Migrations Manually

You can run migrations manually using the Docker Compose:

```bash
docker-compose up migrations
```

Or using the Liquibase Docker image directly:

```bash
docker run --rm \
  -v "$(pwd)/migrations:/liquibase/changelog" \
  -e "LIQUIBASE_COMMAND_USERNAME=template_user" \
  -e "LIQUIBASE_COMMAND_PASSWORD=template_pass" \
  -e "LIQUIBASE_COMMAND_URL=jdbc:postgresql://postgres:5432/template_db" \
  --network=template-system_template-network \
  liquibase/liquibase:alpine \
  --changelog-file=master-changelog.yaml \
  update
```

## Adding New Migrations

To add a new migration:

1. Create a new YAML file in the `migrations/` directory with an incremental version number
2. Add the migration to `master-changelog.yaml`
3. Run the migrations service

Example migration file:

```yaml
databaseChangeLog:
  - changeSet:
      id: [ unique_id ]
      author: [ your_name ]
      comment: [ description ]
      preConditions:
        - onFail: MARK_RAN
          not:
            - tableExists:
                schemaName: template_service
                tableName: [ table_name ]
      changes:
      # Your schema changes here
      rollback:
      # Rollback instructions here
```

## Environment-Specific Migrations

For environment-specific migrations, place them in the appropriate subdirectory:

- `migrations/dev/` - Development environment
- `migrations/test/` - Test environment
- `migrations/prod/` - Production environment

These migrations will be applied only when the corresponding environment is specified in the `ENVIRONMENT` variable.