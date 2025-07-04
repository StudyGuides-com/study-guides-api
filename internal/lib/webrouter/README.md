# Web Router

This package provides a modular HTTP router for handling web routes with templates, separate from the gRPC services.

## Structure

```
internal/lib/webrouter/
├── webrouter.go          # Main router implementation
├── routes/
│   ├── routes.go         # Route registration and management
│   ├── home.go           # Home page handler
│   └── notfound.go       # 404 error handler
└── README.md             # This file
```

## Adding New Routes

To add a new route, follow these steps:

### 1. Create a Route Handler

Create a new file in `internal/lib/webrouter/routes/` (e.g., `about.go`):

```go
package routes

import (
    "html/template"
    "log"
    "net/http"
)

// AboutHandler handles about page requests
type AboutHandler struct {
    templates *template.Template
}

// NewAboutHandler creates a new about handler
func NewAboutHandler(templates *template.Template) *AboutHandler {
    return &AboutHandler{
        templates: templates,
    }
}

// Handle responds to about page requests
func (h *AboutHandler) Handle(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "Title":   "About",
        "Message": "About Study Guides API",
    }

    if err := h.templates.ExecuteTemplate(w, "about.html", data); err != nil {
        log.Printf("Error executing template: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
}
```

### 2. Create a Template

Create a corresponding HTML template in `templates/` (e.g., `about.html`):

```html
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>{{.Title}}</title>
    <link rel="stylesheet" href="/static/css/global.css" />
  </head>
  <body>
    <div class="container">
      <h1>{{.Title}}</h1>
      <div class="message">{{.Message}}</div>
    </div>
  </body>
</html>
```

### 3. Register the Route

Add the route registration in `internal/lib/webrouter/webrouter.go`:

```go
// registerRoutes registers all available routes
func (wr *WebRouter) registerRoutes() {
    // Register home page
    wr.routes.Register("/", routes.NewHomeHandler(wr.templates))

    // Register about page
    wr.routes.Register("/about", routes.NewAboutHandler(wr.templates))

    // Register 404 handler (this will be used as fallback)
    wr.routes.Register("*", routes.NewNotFoundHandler(wr.templates))
}
```

## Route Types

### Static Routes

For simple static pages like `/about`, `/contact`, etc.

### Web Pages

For HTML pages with templates like `/docs`, `/about`, etc.

### Static Files

For serving static assets like images, CSS, and JavaScript files under `/static/`.

### Dynamic Routes

For routes with parameters, you can extend the routing system to support path parameters.

## Template Data

Each handler can pass different data to templates:

```go
data := map[string]interface{}{
    "Title":   "Page Title",
    "Message": "Page message",
    "Items":   []string{"item1", "item2"},
    "User":    userObject,
}
```

## Error Handling

- Template errors are logged and return a 500 error
- Missing routes automatically fall back to the 404 handler

## Testing

To test your new route:

1. Start the server: `go run cmd/server/main.go`
2. Visit `http://localhost:8080/your-route`
3. Check the logs for any errors

## Static Files

Static files are automatically served from the `static/` directory:

- Images: `http://localhost:8080/static/images/logo.png`
- CSS: `http://localhost:8080/static/css/global.css`
- JavaScript: `http://localhost:8080/static/js/app.js`

### Global CSS

The application uses a global CSS file (`/static/css/global.css`) that provides:

- Consistent styling across all pages
- Responsive design for mobile devices
- Modern typography and spacing
- Pre-defined classes for common elements

All templates should include the global CSS link:

```html
<link rel="stylesheet" href="/static/css/global.css" />
```

## Best Practices

1. Keep handlers focused on a single responsibility
2. Use consistent naming conventions
3. Include proper error handling
4. Add logging for debugging
5. Use semantic HTML and accessible design in templates
