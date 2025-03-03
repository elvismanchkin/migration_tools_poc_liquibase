# Database Migrations with Flyway

This component handles database migrations using Flyway. It applies migrations to the PostgreSQL database in a
version-controlled manner.

## Structure

```
flyway/
├── Dockerfile              # Docker image definition for Flyway migrations
├── entrypoint.sh           # Script to run Flyway migrations
├── sql/                    # SQL migration files
│   ├── V1__Initial_Schema.sql  
│   ├── V2__Create_Audit_Tables.sql
│   ├── V3__Create_Templates_Tables.sql
│   ├── V4__Add_Configuration_Tables.sql
│   └── R__Dev_Data.sql     # Repeatable migration for test data
└── README.md               # This file
```

## Migration Naming Convention

Flyway uses a specific naming convention for migration files:

- `V<VERSION>__<DESCRIPTION>.sql` - Versioned migrations that run once in order
- `R__<DESCRIPTION>.sql` - Repeatable migrations that run when their content changes
- `U<VERSION>__<DESCRIPTION>.sql` - Undo migrations (not used in this project)

The description part of the filename uses underscores instead of spaces.

## How It Works

The migration service:

1. Waits for the database to be ready
2. Runs Flyway migrate to apply all pending migrations
3. Terminates when migrations are complete

## Running Migrations Manually

You can run migrations manually using the Docker Compose:

```bash
# Using the Flyway setup
docker-compose -f docker-compose-flyway.yml up flyway
```

Or using the Flyway Docker image directly:

```bash
docker run --rm \
  -v "$(pwd)/flyway/sql:/flyway/sql" \
  -e "FLYWAY_URL=jdbc:postgresql://postgres:5432/template_db" \
  -e "FLYWAY_USER=template_user" \
  -e "FLYWAY_PASSWORD=template_pass" \
  -e "FLYWAY_DEFAULT_SCHEMA=template_service" \
  --network=template-system_template-network \
  flyway/flyway:9-alpine \
  migrate
```

## Adding New Migrations

To add a new migration:

1. Create a new SQL file in the `sql/` directory with an incremental version number
2. Follow the naming convention `V<VERSION>__<DESCRIPTION>.sql`
3. Run the Flyway migrations service

Example migration file:

```sql
-- V5__Add_User_Table.sql
CREATE TABLE template_service.users
(
    id         SERIAL PRIMARY KEY,
    username   VARCHAR(100) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Environment-Specific Data

For environment-specific data, use repeatable migrations with conditional logic:

- `R__Dev_Data.sql` - Development environment data
- `R__Test_Data.sql` - Test environment data
