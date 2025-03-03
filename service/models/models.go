package models

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/template-service/db"
)

// Template struct with GORM tags
type Template struct {
	ID           string           `gorm:"type:uuid;primaryKey" json:"id"`
	Name         string           `gorm:"size:200;not null" json:"name"`
	CategoryID   int              `gorm:"not null" json:"category_id"`
	Content      string           `gorm:"type:text;not null" json:"content"`
	Format       string           `gorm:"size:50;not null;default:html" json:"format"`
	Version      int              `gorm:"not null;default:1" json:"version"`
	IsActive     bool             `gorm:"not null;default:true" json:"is_active"`
	CreatedBy    string           `gorm:"size:100;not null" json:"created_by"`
	CreatedAt    time.Time        `gorm:"autoCreateTime" json:"created_at"`
	UpdatedBy    string           `gorm:"size:100" json:"updated_by"`
	UpdatedAt    time.Time        `gorm:"autoUpdateTime" json:"updated_at"`
	Category     TemplateCategory `gorm:"foreignKey:CategoryID" json:"-"`
	CategoryName string           `gorm:"-" json:"category_name"`
}

// TemplateVariable for template parameters
type TemplateVariable struct {
	ID           int    `gorm:"primaryKey;autoIncrement" json:"id"`
	TemplateID   string `gorm:"type:uuid;not null" json:"template_id"`
	VariableName string `gorm:"size:100;not null" json:"variable_name"`
	Description  string `gorm:"type:text" json:"description"`
	DefaultValue string `gorm:"type:text" json:"default_value"`
	IsRequired   bool   `gorm:"not null;default:false" json:"is_required"`
	VariableType string `gorm:"size:50;not null;default:string" json:"variable_type"`
}

// TemplateCategory represents a template category
type TemplateCategory struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName sets the table name for Template
func (Template) TableName() string {
	return "template_service.template"
}

// TableName sets the table name for TemplateVariable
func (TemplateVariable) TableName() string {
	return "template_service.template_variable"
}

// TableName sets the table name for TemplateCategory
func (TemplateCategory) TableName() string {
	return "template_service.template_category"
}

// GetTemplates returns all templates
func GetTemplates() ([]Template, error) {
	var templates []Template
	result := db.GORMDB.Preload("Category").Where("is_active = ?", true).Order("created_at DESC").Find(&templates)

	// Set CategoryName from the preloaded Category
	for i := range templates {
		templates[i].CategoryName = templates[i].Category.Name
	}

	return templates, result.Error
}

func GetTemplateByID(id string) (Template, error) {
	var template Template
	result := db.GORMDB.Preload("Category").Where("id = ?", id).First(&template)

	// Set CategoryName from the preloaded Category
	if result.Error == nil {
		template.CategoryName = template.Category.Name
	}

	return template, result.Error
}

func GetTemplateVariables(templateID string) ([]TemplateVariable, error) {
	var variables []TemplateVariable
	result := db.GORMDB.Where("template_id = ?", templateID).Order("id").Find(&variables)
	return variables, result.Error
}

func GetTemplateCategories() ([]TemplateCategory, error) {
	var categories []TemplateCategory
	result := db.GORMDB.Order("name").Find(&categories)
	return categories, result.Error
}

// CreateTemplate adds a new template to the database
func CreateTemplate(name, categoryIDStr, content, format, createdBy string) (string, error) {
	templateID := uuid.New().String()

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return "", err
	}

	template := Template{
		ID:         templateID,
		Name:       name,
		CategoryID: categoryID,
		Content:    content,
		Format:     format,
		Version:    1,
		IsActive:   true,
		CreatedBy:  createdBy,
	}

	result := db.GORMDB.Create(&template)
	if result.Error != nil {
		return "", result.Error
	}

	return templateID, nil
}

// AddTemplateVariable adds a variable to a template
func AddTemplateVariable(templateID, variableName, description, defaultValue string, isRequired bool) error {
	variable := TemplateVariable{
		TemplateID:   templateID,
		VariableName: variableName,
		Description:  description,
		DefaultValue: defaultValue,
		IsRequired:   isRequired,
		VariableType: "string", // Default
	}

	result := db.GORMDB.Create(&variable)
	return result.Error
}

// UpdateTemplate updates an existing template
func UpdateTemplate(id, name, categoryIDStr, content, format, updatedBy string) error {
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		return err
	}

	// First get the template to increment version
	var template Template
	if err := db.GORMDB.First(&template, "id = ?", id).Error; err != nil {
		return err
	}

	// Update the template
	result := db.GORMDB.Model(&Template{ID: id}).Updates(map[string]interface{}{
		"name":        name,
		"category_id": categoryID,
		"content":     content,
		"format":      format,
		"updated_by":  updatedBy,
		"version":     template.Version + 1,
	})

	return result.Error
}

// DeleteTemplate marks a template as inactive
func DeleteTemplate(id string) error {
	result := db.GORMDB.Model(&Template{ID: id}).Updates(map[string]interface{}{
		"is_active": false,
	})

	return result.Error
}
