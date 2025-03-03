# Template Service

This is a Go web service that provides a simple UI for managing templates and generating PDFs.

## Project Structure

```
service/
├── db/                   # Database connection and utilities
│   └── db.go
├── handlers/             # HTTP request handlers
│   └── handlers.go
├── models/               # Data models and database access
│   └── models.go
├── templates/            # HTML templates for the UI
│   ├── layout.html
│   ├── templates-list.html
│   ├── template-form.html
│   ├── template-view.html
│   └── template-rendered.html
├── static/               # Static assets (CSS, JS, images)
├── main.go               # Application entry point
├── Dockerfile            # Docker configuration
└── README.md             # This file
```

## Dependencies

- [gorilla/mux](https://github.com/gorilla/mux) - HTTP router
- [lib/pq](https://github.com/lib/pq) - PostgreSQL driver
- [go-wkhtmltopdf](https://github.com/SebastiaanKlippert/go-wkhtmltopdf) - PDF generation
- [godotenv](https://github.com/joho/godotenv) - Environment variable loading
- [google/uuid](https://github.com/google/uuid) - UUID generation

## Building and Running

### Locally

```bash
# Build the service
go build -o template-service

# Run the service
./template-service
```

### Docker

```bash
# Build the Docker image
docker build -t template-service .

# Run the Docker container
docker run -p 8080:8080 template-service
```

## Configuration

The service can be configured using environment variables or a `.env` file:

| Variable    | Description       | Default          |
|-------------|-------------------|------------------|
| DB_HOST     | Database host     | localhost        |
| DB_PORT     | Database port     | 5432             |
| DB_NAME     | Database name     | template_db      |
| DB_USER     | Database user     | template_user    |
| DB_PASSWORD | Database password | template_pass    |
| DB_SCHEMA   | Database schema   | template_service |
| SERVER_PORT | Web server port   | 8080             |
| ENVIRONMENT | Environment name  | dev              |

## API Endpoints

- `GET /` - Redirect to templates list
- `GET /templates` - List all templates
- `GET /templates/new` - Show new template form
- `POST /templates` - Create a new template
- `GET /templates/{id}` - View a specific template
- `POST /templates/{id}/render` - Render a template with variables
- `POST /templates/{id}/pdf` - Generate a PDF from a template
- `GET /health` - Health check endpoint

## REST API Endpoints

The service provides a REST API for programmatic access:

- `GET /api/health` - API health check
- `GET /api/templates` - List all templates
- `POST /api/templates` - Create a new template
- `GET /api/templates/{id}` - Get a specific template
- `PUT /api/templates/{id}` - Update a template
- `DELETE /api/templates/{id}` - Delete a template
- `GET /api/templates/{id}/variables` - Get template variables
- `POST /api/templates/{id}/variables` - Add a variable to a template
- `POST /api/templates/{id}/render` - Render a template with variables
- `GET /api/categories` - List all template categories

All API endpoints return JSON responses with a standard format:

```json
{
    "success": true,
    "data": {
        "some": "data"
    }
}
```

Or in case of error:

```json
{
    "success": false,
    "error": "Error message"
}
```

## ORM Integration

This project demonstrates integration with two popular Go ORMs: GORM and Ent.

### Using Ent ORM

[Ent](https://entgo.io/) is a powerful entity framework for Go that provides a simple API for modeling, querying, and
maintaining your database schema.

#### Setup

1. Make sure you have the Ent dependencies installed:

```bash
go get entgo.io/ent
```

2. The schema definitions are in the `service/ent/schema` directory:
    - `template.go`: Defines the Template entity
    - `template_category.go`: Defines the TemplateCategory entity
    - `template_variable.go`: Defines the TemplateVariable entity

3. Generate the Ent code:

```bash
cd service
go generate ./ent
```

#### Key Features

- **Type-Safe API**: Ent provides a fully type-safe API generated from your schema definitions.
- **Rich Query Interface**: Create complex queries with a simple and fluent API.
- **Schema Migrations**: Ent can automatically create database schema based on your entity definitions.
- **Relationships**: Automatically handles entity relationships with edges.

#### Sample Usage

```go
// Query active templates with their categories
templates, err := db.EntClient.Template.
Query().
Where(template.IsActiveEQ(true)).
WithCategory().
All(ctx)

// Create a new template
template, err := db.EntClient.Template.
Create().
SetName("Example Template").
SetCategoryID(1).
SetContent("<h1>Hello, {{.name}}!</h1>").
SetFormat("html").
SetCreatedBy("admin").
Save(ctx)
```

### Using GORM

[GORM](https://gorm.io/) is a developer-friendly ORM library for Go that provides a more traditional ORM approach.

#### Setup

1. Install GORM dependencies:

```bash
go get gorm.io/gorm
go get gorm.io/driver/postgres
```

2. Models are defined in `service/models/models.go` with GORM tags

3. The database connection is set up in `service/db/db.go`

For more information, see the [GORM Documentation](https://gorm.io/docs/)
and [Ent Documentation](https://entgo.io/docs/getting-started).