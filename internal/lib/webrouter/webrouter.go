package webrouter

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/studyguides-com/study-guides-api/internal/lib/webrouter/routes"
)

// WebRouter handles HTTP routes for web pages and API endpoints
type WebRouter struct {
	templates *template.Template
	store     interface{} // You can type this more specifically if needed
	routes    *routes.Routes
}

// NewWebRouter creates a new web router instance
func NewWebRouter(store interface{}) *WebRouter {
	router := &WebRouter{
		store:  store,
		routes: routes.NewRoutes(),
	}
	
	// Load templates
	router.loadTemplates()
	
	// Register routes
	router.registerRoutes()
	
	return router
}

// loadTemplates loads HTML templates from the templates directory
func (wr *WebRouter) loadTemplates() {
	// You can customize the template path as needed
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Printf("Warning: Could not load templates: %v", err)
		// Create an empty template set to avoid nil pointer errors
		wr.templates = template.New("")
	} else {
		wr.templates = tmpl
	}
}

// registerRoutes registers all available routes
func (wr *WebRouter) registerRoutes() {
	// Register health endpoint
	wr.routes.Register("/health", routes.NewHealthHandler())
	
	// Register home page
	wr.routes.Register("/", routes.NewHomeHandler(wr.templates))
	
	// Register 404 handler (this will be used as fallback)
	wr.routes.Register("*", routes.NewNotFoundHandler(wr.templates))
}

// ServeHTTP handles all HTTP requests and routes them appropriately
func (wr *WebRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("=== WEB REQUEST START ===")
	log.Printf("Method: %s", r.Method)
	log.Printf("Path: %s", r.URL.Path)
	log.Printf("Protocol: %s", r.Proto)
	log.Printf("Content-Type: %s", r.Header.Get("Content-Type"))
	log.Printf("User-Agent: %s", r.Header.Get("User-Agent"))

	// Handle static files first
	if strings.HasPrefix(r.URL.Path, "/static/") {
		wr.serveStatic(w, r)
		return
	}

	// Try to find a matching route
	handler, found := wr.routes.FindRoute(r.URL.Path)
	if found {
		handler.Handle(w, r)
	} else {
		// Use 404 handler
		notFoundHandler, _ := wr.routes.FindRoute("*")
		notFoundHandler.Handle(w, r)
	}

	log.Printf("=== WEB REQUEST END ===")
}

// serveStatic serves static files from the static directory
func (wr *WebRouter) serveStatic(w http.ResponseWriter, r *http.Request) {
	// Remove the /static/ prefix to get the file path
	filePath := strings.TrimPrefix(r.URL.Path, "/static/")
	
	// Serve files from the static directory
	http.ServeFile(w, r, "static/"+filePath)
} 