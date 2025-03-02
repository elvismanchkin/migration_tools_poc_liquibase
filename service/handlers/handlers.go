package handlers

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/gorilla/mux"
	"github.com/yourusername/template-service/models"
)

// FS is the embedded filesystem for HTML templates
var FS embed.FS

// HandleIndex redirects to templates list
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

// HandleListTemplates lists all templates
func HandleListTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := models.GetTemplates()
	if err != nil {
		http.Error(w, "Error fetching templates: "+err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := models.GetTemplateCategories()
	if err != nil {
		http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Templates  []models.Template
		Categories []models.TemplateCategory
	}{
		Templates:  templates,
		Categories: categories,
	}

	// Load templates and parse them
	htmlTemplate, err := template.ParseFS(FS, "templates/layout.html", "templates/templates-list.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// HandleNewTemplateForm displays form to create a new template
func HandleNewTemplateForm(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetTemplateCategories()
	if err != nil {
		http.Error(w, "Error fetching categories: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []models.TemplateCategory
	}{
		Categories: categories,
	}

	// Load templates and parse them
	htmlTemplate, err := template.ParseFS(FS, "templates/layout.html", "templates/template-form.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// HandleCreateTemplate handles new template creation
func HandleCreateTemplate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	categoryID := r.FormValue("category_id")
	content := r.FormValue("content")
	format := r.FormValue("format")

	if name == "" || categoryID == "" || content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Validate category ID is a number
	_, err = strconv.Atoi(categoryID)
	if err != nil {
		http.Error(w, "Invalid category ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create template in database
	templateID, err := models.CreateTemplate(name, categoryID, content, format, "web_user")
	if err != nil {
		http.Error(w, "Error creating template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Add variables if they were included
	varName := r.FormValue("var_name")
	varDesc := r.FormValue("var_description")
	varDefault := r.FormValue("var_default")
	varRequired := r.FormValue("var_required") == "on"

	if varName != "" {
		err = models.AddTemplateVariable(templateID, varName, varDesc, varDefault, varRequired)
		if err != nil {
			log.Printf("Warning: Failed to add variable to template: %v", err)
		}
	}

	// Redirect to template list
	http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

// HandleViewTemplate displays a single template
func HandleViewTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tmpl, err := models.GetTemplateByID(id)
	if err != nil {
		http.Error(w, "Error fetching template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	variables, err := models.GetTemplateVariables(id)
	if err != nil {
		http.Error(w, "Error fetching template variables: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Template  models.Template
		Variables []models.TemplateVariable
	}{
		Template:  tmpl,
		Variables: variables,
	}

	// Load templates and parse them
	htmlTemplate, err := template.ParseFS(FS, "templates/layout.html", "templates/template-view.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// HandleRenderTemplate renders a template with variables
func HandleRenderTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := models.GetTemplateByID(id)
	if err != nil {
		http.Error(w, "Error fetching template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all variables for this template
	variables, err := models.GetTemplateVariables(id)
	if err != nil {
		http.Error(w, "Error fetching template variables: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build variable map for template rendering
	varMap := make(map[string]interface{})
	for _, v := range variables {
		value := r.FormValue(v.VariableName)
		if value == "" {
			value = v.DefaultValue
		}
		varMap[v.VariableName] = value
	}

	// Render the template content to a string
	var renderedBuffer bytes.Buffer
	htmlTmpl, err := template.New("render").Parse(tmpl.Content)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTmpl.Execute(&renderedBuffer, varMap)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create data for the HTML template
	data := struct {
		Template        models.Template
		Variables       []models.TemplateVariable
		RenderedContent template.HTML
		FormValues      map[string]interface{}
	}{
		Template:        tmpl,
		Variables:       variables,
		RenderedContent: template.HTML(renderedBuffer.String()),
		FormValues:      varMap,
	}

	// Load and parse the page template
	htmlTemplate, err := template.ParseFS(FS, "templates/layout.html", "templates/template-rendered.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the page template
	err = htmlTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

// HandleGeneratePDF generates a PDF from a rendered template
func HandleGeneratePDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, err := models.GetTemplateByID(id)
	if err != nil {
		http.Error(w, "Error fetching template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get all variables for this template
	variables, err := models.GetTemplateVariables(id)
	if err != nil {
		http.Error(w, "Error fetching template variables: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build variable map for template rendering
	varMap := make(map[string]interface{})
	for _, v := range variables {
		value := r.FormValue(v.VariableName)
		if value == "" {
			value = v.DefaultValue
		}
		varMap[v.VariableName] = value
	}

	// Render the template to a string
	var renderedBuffer bytes.Buffer
	htmlTmpl, err := template.New("render").Parse(tmpl.Content)
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTmpl.Execute(&renderedBuffer, varMap)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate PDF using wkhtmltopdf
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		http.Error(w, "Error creating PDF generator: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new page from string
	page := wkhtmltopdf.NewPageReader(strings.NewReader(renderedBuffer.String()))
	pdfg.AddPage(page)

	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		http.Error(w, "Error generating PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set response headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", tmpl.Name))

	// Write PDF to response
	_, err = w.Write(pdfg.Bytes())
	if err != nil {
		http.Error(w, "Error sending PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
