package handlers

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/yourusername/template-service/models"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type TemplateRequest struct {
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
	Content    string `json:"content"`
	Format     string `json:"format"`
}

type TemplateVariableRequest struct {
	VariableName string `json:"variable_name"`
	Description  string `json:"description"`
	DefaultValue string `json:"default_value"`
	IsRequired   bool   `json:"is_required"`
	VariableType string `json:"variable_type,omitempty"`
}

type RenderRequest struct {
	Variables map[string]string `json:"variables"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"success":false,"error":"Error marshalling JSON response"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(response); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, APIResponse{
		Success: false,
		Error:   message,
	})
}

func APIGetTemplates(w http.ResponseWriter, r *http.Request) {
	_ = r //explicitly ignored
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

func APIGetTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	retrievedTemplate, err := models.GetTemplateByID(id)
	if err != nil {
		log.Printf("Failed to retrieve template %s: %v", id, err)
		respondWithError(w, http.StatusNotFound, "Template not found")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    retrievedTemplate,
	})
}

func APICreateTemplate(w http.ResponseWriter, r *http.Request) {
	var req TemplateRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	if req.Name == "" || req.CategoryID == "" || req.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	templateID, err := models.CreateTemplate(req.Name, req.CategoryID, req.Content, req.Format, "api_user")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating template: "+err.Error())
		return
	}

	retrievedTemplate, err := models.GetTemplateByID(templateID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Template created but could not be retrieved: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    retrievedTemplate,
	})
}

func APIUpdateTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	if req.Name == "" || req.CategoryID == "" || req.Content == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	err = models.UpdateTemplate(id, req.Name, req.CategoryID, req.Content, req.Format, "api_user")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating template: "+err.Error())
		return
	}

	retrievedTemplate, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Template updated but could not be retrieved: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    retrievedTemplate,
	})
}

func APIDeleteTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := models.GetTemplateByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Template not found: "+err.Error())
		return
	}

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

func APIGetTemplateVariables(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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

func APIAddTemplateVariable(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	if req.VariableName == "" {
		respondWithError(w, http.StatusBadRequest, "Variable name is required")
		return
	}

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

func APIGetCategories(w http.ResponseWriter, r *http.Request) {
	_ = r //explicitly ignored
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

func APIRenderTemplate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	tmpl, err := models.GetTemplateByID(id)
	if err != nil {
		log.Printf("Failed to retrieve template %s: %v", id, err)
		respondWithError(w, http.StatusNotFound, "Template not found")
		return
	}

	var renderReq RenderRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&renderReq); err != nil {
		log.Printf("Invalid request payload: %v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("Error closing request body: %v", err)
		}
	}()

	templateVars, err := models.GetTemplateVariables(id)
	if err != nil {
		log.Printf("Error fetching template variables for %s: %v", id, err)
		respondWithError(w, http.StatusInternalServerError, "Error fetching template variables")
		return
	}

	varMap := make(map[string]interface{})
	for _, v := range templateVars {
		value, exists := renderReq.Variables[v.VariableName]
		if v.IsRequired && (!exists || value == "") {
			respondWithError(w, http.StatusBadRequest, "Required variable missing: "+v.VariableName)
			return
		}
		if exists && value != "" {
			varMap[v.VariableName] = value
		} else {
			varMap[v.VariableName] = v.DefaultValue
		}
	}

	var renderedBuffer bytes.Buffer
	htmlTmpl, err := template.New("render").Parse(tmpl.Content)
	if err != nil {
		log.Printf("Error parsing template %s: %v", id, err)
		respondWithError(w, http.StatusInternalServerError, "Error parsing template")
		return
	}

	if err = htmlTmpl.Execute(&renderedBuffer, varMap); err != nil {
		log.Printf("Error rendering template %s: %v", id, err)
		respondWithError(w, http.StatusInternalServerError, "Error rendering template")
		return
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    renderedBuffer.String(),
	})
}

func APIHealthCheck(w http.ResponseWriter, r *http.Request) {
	_ = r // explicitly ignored

	status := map[string]interface{}{
		"status":    "up",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	respondWithJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    status,
	})
}
