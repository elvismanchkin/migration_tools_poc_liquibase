package models

import (
	"database/sql"
	"log"
	"time"

	"github.com/elvismanchkin/migration_tools_poc_liquibase/db"
	"github.com/google/uuid"
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
	rows, err := db.DB.Query(`
		SELECT 
			t.id, t.name, t.category_id, t.content, t.format, 
			t.version, t.is_active, t.created_by, t.created_at, 
			t.updated_by, t.updated_at, c.name as category_name
		FROM template_service.template t
		JOIN template_service.template_category c ON t.category_id = c.id
		WHERE t.is_active = true
		ORDER BY t.created_at DESC
	`)
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

	err := db.DB.QueryRow(`
		SELECT 
			t.id, t.name, t.category_id, t.content, t.format, 
			t.version, t.is_active, t.created_by, t.created_at, 
			t.updated_by, t.updated_at, c.name as category_name
		FROM template_service.template t
		JOIN template_service.template_category c ON t.category_id = c.id
		WHERE t.id = $1
	`, id).Scan(
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
	rows, err := db.DB.Query(`
		SELECT 
			id, template_id, variable_name, description, 
			default_value, is_required, variable_type
		FROM template_service.template_variable
		WHERE template_id = $1
		ORDER BY id
	`, templateID)
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
	rows, err := db.DB.Query(`
		SELECT id, name, description
		FROM template_service.template_category
		ORDER BY name
	`)
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
	// Create new template in database
	templateID := uuid.New().String()
	_, err := db.DB.Exec(`
		INSERT INTO template_service.template 
		(id, name, category_id, content, format, version, is_active, created_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		templateID, name, categoryID, content, format, 1, true, createdBy)

	if err != nil {
		return "", err
	}

	return templateID, nil
}

// AddTemplateVariable adds a variable to a template
func AddTemplateVariable(templateID, variableName, description, defaultValue string, isRequired bool) error {
	_, err := db.DB.Exec(`
		INSERT INTO template_service.template_variable 
		(template_id, variable_name, description, default_value, is_required) 
		VALUES ($1, $2, $3, $4, $5)`,
		templateID, variableName, description, defaultValue, isRequired)

	return err
}

// UpdateTemplate updates an existing template
func UpdateTemplate(id, name, categoryID, content, format, updatedBy string) error {
	_, err := db.DB.Exec(`
		UPDATE template_service.template 
		SET name = $1, category_id = $2, content = $3, format = $4, 
		    updated_by = $5, updated_at = CURRENT_TIMESTAMP, version = version + 1
		WHERE id = $6`,
		name, categoryID, content, format, updatedBy, id)

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
