package models

import (
	"strconv"
	"time"

	"github.com/elvismanchkin/migration_tools_poc_liquibase/db"
	"github.com/elvismanchkin/migration_tools_poc_liquibase/ent"
	"github.com/elvismanchkin/migration_tools_poc_liquibase/ent/template"
	"github.com/elvismanchkin/migration_tools_poc_liquibase/ent/templatecategory"
	"github.com/elvismanchkin/migration_tools_poc_liquibase/ent/templatevariable"
	"github.com/google/uuid"
)

// For compatibility with existing code
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
	CategoryName string
}

type TemplateVariable struct {
	ID           int
	TemplateID   string
	VariableName string
	Description  string
	DefaultValue string
	IsRequired   bool
	VariableType string
}

type TemplateCategory struct {
	ID          int
	Name        string
	Description string
}

// Convert Ent template to our model
func toTemplateModel(t *ent.Template) Template {
	template := Template{
		ID:        t.ID.String(),
		Name:      t.Name,
		Content:   t.Content,
		Format:    t.Format,
		Version:   t.Version,
		IsActive:  t.IsActive,
		CreatedBy: t.CreatedBy,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}

	// Handle nullable fields
	if t.UpdatedBy != nil {
		template.UpdatedBy = *t.UpdatedBy
	}

	// Set category information if available
	if t.Edges.Category != nil {
		template.CategoryID = t.Edges.Category.ID
		template.CategoryName = t.Edges.Category.Name
	}

	return template
}

// GetTemplates returns all templates
func GetTemplates() ([]Template, error) {
	templates, err := db.EntClient.Template.
		Query().
		Where(template.IsActiveEQ(true)).
		Order(ent.Desc(template.FieldCreatedAt)).
		WithCategory().
		All(db.Ctx)

	if err != nil {
		return nil, err
	}

	result := make([]Template, len(templates))
	for i, t := range templates {
		result[i] = toTemplateModel(t)
	}

	return result, nil
}

// GetTemplateByID returns a template by ID
func GetTemplateByID(id string) (Template, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return Template{}, err
	}

	t, err := db.EntClient.Template.
		Query().
		Where(template.ID(uid)).
		WithCategory().
		Only(db.Ctx)

	if err != nil {
		return Template{}, err
	}

	return toTemplateModel(t), nil
}

// GetTemplateVariables returns variables for a template
func GetTemplateVariables(templateID string) ([]TemplateVariable, error) {
	uid, err := uuid.Parse(templateID)
	if err != nil {
		return nil, err
	}

	vars, err := db.EntClient.TemplateVariable.
		Query().
		Where(templatevariable.HasTemplateWith(template.ID(uid))).
		All(db.Ctx)

	if err != nil {
		return nil, err
	}

	result := make([]TemplateVariable, len(vars))
	for i, v := range vars {
		result[i] = TemplateVariable{
			ID:           v.ID,
			TemplateID:   templateID,
			VariableName: v.VariableName,
			IsRequired:   v.IsRequired,
			VariableType: v.VariableType,
		}

		// Handle nullable fields
		if v.Description != nil {
			result[i].Description = *v.Description
		}

		if v.DefaultValue != nil {
			result[i].DefaultValue = *v.DefaultValue
		}
	}

	return result, nil
}

// GetTemplateCategories returns all template categories
func GetTemplateCategories() ([]TemplateCategory, error) {
	categories, err := db.EntClient.TemplateCategory.
		Query().
		Order(ent.Asc(templatecategory.FieldName)).
		All(db.Ctx)

	if err != nil {
		return nil, err
	}

	result := make([]TemplateCategory, len(categories))
	for i, c := range categories {
		result[i] = TemplateCategory{
			ID:   c.ID,
			Name: c.Name,
		}

		// Handle nullable field
		if c.Description != nil {
			result[i].Description = *c.Description
		}
	}

	return result, nil
}

// CreateTemplate adds a new template to the database
func CreateTemplate(name, categoryIDStr, content, format, createdBy string) (string, error) {
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return "", err
	}

	// Create template
	t, err := db.EntClient.Template.
		Create().
		SetName(name).
		SetCategoryID(categoryID).
		SetContent(content).
		SetFormat(format).
		SetVersion(1).
		SetIsActive(true).
		SetCreatedBy(createdBy).
		Save(db.Ctx)

	if err != nil {
		return "", err
	}

	return t.ID.String(), nil
}

// AddTemplateVariable adds a variable to a template
func AddTemplateVariable(templateID, variableName, description, defaultValue string, isRequired bool) error {
	tid, err := uuid.Parse(templateID)
	if err != nil {
		return err
	}

	create := db.EntClient.TemplateVariable.
		Create().
		SetVariableName(variableName).
		SetIsRequired(isRequired).
		SetVariableType("string"). // Default
		SetTemplateID(tid)

	// Set nullable fields only if they have values
	if description != "" {
		create.SetDescription(description)
	}

	if defaultValue != "" {
		create.SetDefaultValue(defaultValue)
	}

	_, err = create.Save(db.Ctx)
	return err
}

// UpdateTemplate updates an existing template
func UpdateTemplate(id, name, categoryIDStr, content, format, updatedBy string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return err
	}

	// Get current version
	t, err := db.EntClient.Template.Get(db.Ctx, uid)
	if err != nil {
		return err
	}

	// Update template
	_, err = db.EntClient.Template.
		UpdateOne(t).
		SetName(name).
		SetCategoryID(categoryID).
		SetContent(content).
		SetFormat(format).
		SetUpdatedBy(updatedBy).
		SetVersion(t.Version + 1).
		Save(db.Ctx)

	return err
}

// DeleteTemplate marks a template as inactive
func DeleteTemplate(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	_, err = db.EntClient.Template.
		UpdateOneID(uid).
		SetIsActive(false).
		Save(db.Ctx)

	return err
}
