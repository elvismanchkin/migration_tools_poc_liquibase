package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/template-service/models"
)

// APIResponse is a standard structure for all API responses
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// TemplateRequest is the structure for template creation/update requests
type TemplateRequest struct {
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	Content    string `json:"content"`
	Format     string `json:"format"`
}

// TemplateVariableRequest is the structure for template variable requests
type TemplateVariableRequest struct {
	VariableName string `json:"variable_name"`
	Description  string `json:"description"`
	DefaultValue string `json:"default_value"`
	IsRequired   bool   `json:"is_required"`
	VariableType string `json:"variable_type,omitempty"`
}

// RenderRequest is the structure for template rendering requests
type RenderRequest struct {
	Variables map[string]string `json:"variables"`
}

// API helper functions

// respondWithJSON sends a JSON response with the given status code
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"success":false,"error":"Error marshalling JSON response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// respondWithError sends an error response with the given status code
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		Success: false,
		Error:   message,
	})
}

// API handler functions

// APIGetTemplates returns all templates as JSON
func APIGetTemplates(w http.ResponseWriter, r *http.Request) {
	templates, err := models.GetTemplates()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching templates: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    templates,
	})
}

// APIGetTemplate returns a specific template as JSON
func APIGetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	template, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    template,
	})
}

// APICreateTemplate creates a new template
func APICreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req TemplateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Validate request
	if req.Name == "" || req.CategoryID == "" || req.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Create template
	templateID, err := models.CreateTemplate(req.Name, req.CategoryID, req.Content, req.Format, "api_user")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating template: "+err.Error())
		return
	}

	// Get the created template to return
	template, err := models.GetTemplateByID(templateID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Template created but could not be retrieved: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    template,
	})
}

// APIUpdateTemplate updates an existing template
func APIUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if template exists
	_, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

	var req TemplateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Validate request
	if req.Name == "" || req.CategoryID == "" || req.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Update template
	err = models.UpdateTemplate(id, req.Name, req.CategoryID, req.Content, req.Format, "api_user")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating template: "+err.Error())
		return
	}

	// Get the updated template to return
	template, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Template updated but could not be retrieved: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    template,
	})
}

// APIDeleteTemplate deletes a template
func APIDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if template exists
	_, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

	// Delete template
	err = models.DeleteTemplate(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting template: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    "Template deleted successfully",
	})
}

// APIGetTemplateVariables returns all variables for a template
func APIGetTemplateVariables(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if template exists
	_, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

	variables, err := models.GetTemplateVariables(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching template variables: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    variables,
	})
}

// APIAddTemplateVariable adds a variable to a template
func APIAddTemplateVariable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if template exists
	_, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

	var req TemplateVariableRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Validate request
	if req.VariableName == "" {
		respondWithError(w, http.StatusBadRequest, "Variable name is required")
		return
	}

	// Add variable
	err = models.AddTemplateVariable(id, req.VariableName, req.Description, req.DefaultValue, req.IsRequired)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error adding template variable: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    "Variable added successfully",
	})
}

// APIGetCategories returns all template categories
func APIGetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := models.GetTemplateCategories()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching categories: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    categories,
	})
}

// APIRenderTemplate renders a template with provided variables
func APIRenderTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Check if template exists
	tmpl, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

	// Parse the request body to get variables
	var renderReq RenderRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&renderReq); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Get template variables to check required fields and set defaults
	templateVars, err := models.GetTemplateVariables(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error fetching template variables: "+err.Error())
		return
	}

	// Build variable map for template rendering
	varMap := make(map[string]interface{})
	for _, v := range templateVars {
		value, exists := renderReq.Variables[v.VariableName]

		// Check if required variables are provided
		if v.IsRequired && (!exists || value == "") {
			respondWithError(w, http.StatusBadRequest, "Required variable missing: "+v.VariableName)
			return
		}

		// Use provided value or default
		if exists && value != "" {
			varMap[v.VariableName] = value
		} else {
			varMap[v.VariableName] = v.DefaultValue
		}
	}

	// Render the template to a string
	var renderedBuffer bytes.Buffer
	htmlTmpl, err := template.New("render").Parse(tmpl.Content)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing template: "+err.Error())
		return
	}

	err = htmlTmpl.Execute(&renderedBuffer, varMap)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error rendering template: "+err.Error())
		return
	}

	// Return the rendered template
	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    renderedBuffer.String(),
	})
}

// APIHealthCheck provides a simple health check endpoint
func APIHealthCheck(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":    "up",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    status,
	})
}
