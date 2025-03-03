package models

import (
	"database/sql"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/yourusername/template-service/db"
)

// Template struct to match the database schema
type Template struct {
	ID           string
	Name         string
	CategoryID   int
	Content      string
	Format       string
	Version      int
	IsActive     bool
	CreatedBy    string
	CreatedAt    time.Time
	UpdatedBy    string
	UpdatedAt    time.Time
	CategoryName string // Joined from category table
}

// TemplateVariable for template parameters
type TemplateVariable struct {
	ID           int
	TemplateID   string
	VariableName string
	Description  string
	DefaultValue string
	IsRequired   bool
	VariableType string
}

// TemplateCategory represents a template category
type TemplateCategory struct {
	ID          int
	Name        string
	Description string
}

// GetTemplates returns all templates
func GetTemplates() ([]Template, error) {
	query := db.StatementBuilder.
		Select("t.id", "t.name", "t.category_id", "t.content", "t.format",
			"t.version", "t.is_active", "t.created_by", "t.created_at",
			"t.updated_by", "t.updated_at", "c.name as category_name").
		From("template_service.template t").
		Join("template_service.template_category c ON t.category_id = c.id").
		Where(sq.Eq{"t.is_active": true}).
		OrderBy("t.created_at DESC")

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var templates []Template
	for rows.Next() {
		var t Template
		var updatedBy sql.NullString
		var updatedAt sql.NullTime

		err := rows.Scan(
			&t.ID, &t.Name, &t.CategoryID, &t.Content, &t.Format,
			&t.Version, &t.IsActive, &t.CreatedBy, &t.CreatedAt,
			&updatedBy, &updatedAt, &t.CategoryName,
		)
		if err != nil {
			return nil, err
		}

		if updatedBy.Valid {
			t.UpdatedBy = updatedBy.String
		}
		if updatedAt.Valid {
			t.UpdatedAt = updatedAt.Time
		}

		templates = append(templates, t)
	}

	return templates, nil
}

func GetTemplateByID(id string) (Template, error) {
	var t Template
	var updatedBy sql.NullString
	var updatedAt sql.NullTime

	query := db.StatementBuilder.
		Select("t.id", "t.name", "t.category_id", "t.content", "t.format",
			"t.version", "t.is_active", "t.created_by", "t.created_at",
			"t.updated_by", "t.updated_at", "c.name as category_name").
		From("template_service.template t").
		Join("template_service.template_category c ON t.category_id = c.id").
		Where(sq.Eq{"t.id": id})

	err := db.QueryRow(query).Scan(
		&t.ID, &t.Name, &t.CategoryID, &t.Content, &t.Format,
		&t.Version, &t.IsActive, &t.CreatedBy, &t.CreatedAt,
		&updatedBy, &updatedAt, &t.CategoryName,
	)

	if err != nil {
		return t, err
	}

	if updatedBy.Valid {
		t.UpdatedBy = updatedBy.String
	}
	if updatedAt.Valid {
		t.UpdatedAt = updatedAt.Time
	}

	return t, nil
}

func GetTemplateVariables(templateID string) ([]TemplateVariable, error) {
	query := db.StatementBuilder.
		Select("id", "template_id", "variable_name", "description",
			"default_value", "is_required", "variable_type").
		From("template_service.template_variable").
		Where(sq.Eq{"template_id": templateID}).
		OrderBy("id")

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var variables []TemplateVariable
	for rows.Next() {
		var v TemplateVariable
		if err := rows.Scan(&v.ID, &v.TemplateID, &v.VariableName, &v.Description,
			&v.DefaultValue, &v.IsRequired, &v.VariableType); err != nil {
			return nil, err
		}
		variables = append(variables, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return variables, nil
}

func GetTemplateCategories() ([]TemplateCategory, error) {
	query := db.StatementBuilder.
		Select("id", "name", "description").
		From("template_service.template_category").
		OrderBy("name")

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var categories []TemplateCategory
	for rows.Next() {
		var c TemplateCategory
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

// CreateTemplate adds a new template to the database
func CreateTemplate(name, categoryID, content, format, createdBy string) (string, error) {
	templateID := uuid.New().String()

	query := db.StatementBuilder.
		Insert("template_service.template").
		Columns("id", "name", "category_id", "content", "format", "version", "is_active", "created_by").
		Values(templateID, name, categoryID, content, format, 1, true, createdBy)

	_, err := db.Exec(query)
	if err != nil {
		return "", err
	}

	return templateID, nil
}

// AddTemplateVariable adds a variable to a template
func AddTemplateVariable(templateID, variableName, description, defaultValue string, isRequired bool) error {
	query := db.StatementBuilder.
		Insert("template_service.template_variable").
		Columns("template_id", "variable_name", "description", "default_value", "is_required").
		Values(templateID, variableName, description, defaultValue, isRequired)

	_, err := db.Exec(query)

	return err
}

// UpdateTemplate updates an existing template
func UpdateTemplate(id, name, categoryID, content, format, updatedBy string) error {
	query := db.StatementBuilder.
		Update("template_service.template").
		Set("name", name).
		Set("category_id", categoryID).
		Set("content", content).
		Set("format", format).
		Set("updated_by", updatedBy).
		Set("updated_at", sq.Expr("CURRENT_TIMESTAMP")).
		Set("version", sq.Expr("version + 1")).
		Where(sq.Eq{"id": id})

	_, err := db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}

// DeleteTemplate marks a template as inactive
func DeleteTemplate(id string) error {
	_, err := db.DB.Exec(`
		UPDATE template_service.template 
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`,
		id)

	return err
}
