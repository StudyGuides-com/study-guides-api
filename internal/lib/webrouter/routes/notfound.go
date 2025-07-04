package routes

import (
	"html/template"
	"log"
	"net/http"

	"github.com/studyguides-com/study-guides-api/internal/lib/webrouter/utils"
)

// NotFoundHandler handles 404 errors
type NotFoundHandler struct {
	templates *template.Template
}

// NewNotFoundHandler creates a new not found handler
func NewNotFoundHandler(templates *template.Template) *NotFoundHandler {
	return &NotFoundHandler{
		templates: templates,
	}
}

// Handle responds to 404 requests
func (h *NotFoundHandler) Handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	data := map[string]interface{}{
		"Title":   "Page Not Found",
		"Message": "The requested page could not be found",
		"Path":    r.URL.Path,
	}

	// Add environment data
	data = utils.MergeWithEnvData(data, r)

	if err := h.templates.ExecuteTemplate(w, "404.html", data); err != nil {
		log.Printf("Error executing 404 template: %v", err)
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}
} 