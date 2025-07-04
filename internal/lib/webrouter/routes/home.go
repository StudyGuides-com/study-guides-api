package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/studyguides-com/study-guides-api/internal/lib/webrouter/utils"
)

// HomeHandler handles home page requests
type HomeHandler struct {
	templates *template.Template
}

// NewHomeHandler creates a new home handler
func NewHomeHandler(templates *template.Template) *HomeHandler {
	return &HomeHandler{
		templates: templates,
	}
}

// Handle responds to home page requests
func (h *HomeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Study Guides API",
		"Message": "Welcome to the Study Guides API",
	}

	// Add environment data
	data = utils.MergeWithEnvData(data)

	if err := h.templates.ExecuteTemplate(w, "home.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
} 