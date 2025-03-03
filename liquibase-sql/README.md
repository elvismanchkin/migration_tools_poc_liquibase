# Database Migrations with Liquibase SQL

This component handles database migrations using Liquibase with SQL-formatted migrations. It applies migrations to the
PostgreSQL database in a version-controlled manner.

## Structure

```
liquibase-sql/
├── Dockerfile                # Docker image definition for Liquibase migrations
├── entrypoint.sh             # Script to run Liquibase migrations
├── master-changelog.xml      # Master changelog that includes all migrations
├── sql/                      # SQL migration files
│   ├── v1_initial_schema.sql
│   ├── v2_create_audit_tables.sql
│   ├── v3_create_templates_tables.sql
│   ├── v4_add_configuration_tables.sql
│   └── dev/                  # Environment-specific migrations
│       └── v20250228_add_test_data.sql
└── README.md                 # This file
```

## Migration File Format

Liquibase SQL format uses special comment markers to define changesets:

```sql
--liquibase formatted sql

--changeset author:id
--comment A description of the changes
SQL statements...

--changeset author:id2
SQL statements...

--rollback SQL statements to undo the change
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
# Using the Liquibase SQL setup
docker-compose -f docker-compose-liquibase-sql.yml up liquibase-sql
```

Or using the Liquibase Docker image directly:

```bash
docker run --rm \
  -v "$(pwd)/liquibase-sql:/liquibase/changelog" \
  -e "LIQUIBASE_COMMAND_USERNAME=template_user" \
  -e "LIQUIBASE_COMMAND_PASSWORD=template_pass" \
  -e "LIQUIBASE_COMMAND_URL=jdbc:postgresql://postgres:5432/template_db" \
  --network=template-system_template-network \
  liquibase/liquibase:4.20-alpine \
  --changelog-file=master-changelog.xml \
  update
```

## Adding New Migrations

To add a new migration:

1. Create a new SQL file in the `sql/` directory with an incremental version number
2. Follow the Liquibase SQL format with changeset markers
3. Add the migration to `master-changelog.xml`
4. Run the migrations service

Example migration file:

```sql
--liquibase formatted sql

--changeset authornamehere:5
--comment Add User Table
CREATE TABLE template_service.users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(100) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--rollback DROP TABLE template_service.users;
```

## Environment-Specific Data

For environment-specific data, use the `context` attribute in your changeset:

```sql
--changeset authornamehere:20250301001
--comment Add Test Data
--context dev
--runAlways true
```

This will only run when the `ENVIRONMENT` variable is set to `dev`.