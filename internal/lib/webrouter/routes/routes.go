package routes

import (
	"net/http"
)

// RouteHandler defines the interface for route handlers
type RouteHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// Route represents a single route with its handler
type Route struct {
	Path    string
	Handler RouteHandler
}

// Routes holds all registered routes
type Routes struct {
	routes []Route
}

// NewRoutes creates a new Routes instance
func NewRoutes() *Routes {
	return &Routes{
		routes: make([]Route, 0),
	}
}

// Register adds a new route
func (r *Routes) Register(path string, handler RouteHandler) {
	r.routes = append(r.routes, Route{
		Path:    path,
		Handler: handler,
	})
}

// GetRoutes returns all registered routes
func (r *Routes) GetRoutes() []Route {
	return r.routes
}

// FindRoute finds a route that matches the given path
func (r *Routes) FindRoute(path string) (RouteHandler, bool) {
	for _, route := range r.routes {
		if route.Path == path {
			return route.Handler, true
		}
	}
	return nil, false
} 