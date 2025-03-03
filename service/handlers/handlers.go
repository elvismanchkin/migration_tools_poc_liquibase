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
	"github.com/elvismanchkin/migration_tools_poc_liquibase/models"
	"github.com/gorilla/mux"
)

var FS embed.FS

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

func HandleListTemplates(w http.ResponseWriter, r *http.Request) {
	_ = r
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

	htmlTemplate, err := template.ParseFS(FS, "templates/layout.html", "templates/template-list.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

func HandleNewTemplateForm(w http.ResponseWriter, r *http.Request) {
	_ = r
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

	_, err = strconv.Atoi(categoryID)
	if err != nil {
		http.Error(w, "Invalid category ID: "+err.Error(), http.StatusBadRequest)
		return
	}

	templateID, err := models.CreateTemplate(name, categoryID, content, format, "web_user")
	if err != nil {
		http.Error(w, "Error creating template: "+err.Error(), http.StatusInternalServerError)
		return
	}

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

	http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

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

	variables, err := models.GetTemplateVariables(id)
	if err != nil {
		http.Error(w, "Error fetching template variables: "+err.Error(), http.StatusInternalServerError)
		return
	}

	varMap := make(map[string]interface{})
	for _, v := range variables {
		value := r.FormValue(v.VariableName)
		if value == "" {
			value = v.DefaultValue
		}
		varMap[v.VariableName] = value
	}

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

	var templatePath string
	if tmpl.Format == "html" {
		templatePath = "templates/template-pdf-preview.html"
	} else {
		templatePath = "templates/template-rendered.html"
	}

	htmlTemplate, err := template.ParseFS(FS, "templates/layout.html", templatePath)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = htmlTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}

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

	variables, err := models.GetTemplateVariables(id)
	if err != nil {
		http.Error(w, "Error fetching template variables: "+err.Error(), http.StatusInternalServerError)
		return
	}

	varMap := make(map[string]interface{})
	for _, v := range variables {
		value := r.FormValue(v.VariableName)
		if value == "" {
			value = v.DefaultValue
		}
		varMap[v.VariableName] = value
	}

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

	pdfGen, err := wkhtmltopdf.NewPDFGenerator()

	pdfGen.Dpi.Set(300)
	pdfGen.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfGen.MarginTop.Set(20)
	pdfGen.MarginBottom.Set(20)
	pdfGen.MarginLeft.Set(20)
	pdfGen.MarginRight.Set(20)

	if err != nil {
		http.Error(w, "Error creating PDF generator: "+err.Error(), http.StatusInternalServerError)
		return
	}

	page := wkhtmltopdf.NewPageReader(strings.NewReader(renderedBuffer.String()))
	pdfGen.AddPage(page)

	err = pdfGen.Create()
	if err != nil {
		http.Error(w, "Error generating PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.pdf\"", tmpl.Name))

	_, err = w.Write(pdfGen.Bytes())
	if err != nil {
		http.Error(w, "Error sending PDF: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
