package handlers

import (
	"html/template"
	"log"
	"net/http"
	"quickshare/internal/services"
)

// Handler contains all dependencies for HTTP routing.
type Handler struct {
	Storage   *services.StorageService
	Templates map[string]*template.Template
	AppURL    string
}

// GlobalData represents parameters passed to almost every page template.
type GlobalData struct {
	AppURL string
	Data   interface{}
}

// ErrorData represents the parameters required to render custom errors.
type ErrorData struct {
	AppURL  string
	Status  int
	Title   string
	Message string
}

// NewHandler constructs and configures a Handler.
func NewHandler(storage *services.StorageService, tmpls map[string]*template.Template, appURL string) *Handler {
	return &Handler{
		Storage:   storage,
		Templates: tmpls,
		AppURL:    appURL,
	}
}

// RenderTemplate safely merges base and page blocks and writes to the client.
func (h *Handler) RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	// Wrap in GlobalData to guarantee AppURL is always available to layouts
	gd := GlobalData{
		AppURL: h.AppURL,
		Data:   data,
	}

	tmpl, ok := h.Templates[name]
	if !ok {
		log.Printf("Template not found: %s", name)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", gd)
	if err != nil {
		log.Printf("Template execution error for %s: %v", name, err)
		// If template failed, fall back to standard HTTP error
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// RenderError serves a beautifully styled error page in our dark theme.
func (h *Handler) RenderError(w http.ResponseWriter, r *http.Request, statusCode int, title, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)

	data := ErrorData{
		AppURL:  h.AppURL,
		Status:  statusCode,
		Title:   title,
		Message: message,
	}

	// If XMLHttpRequest, we can send a minimal HTML snippet instead of the full layout
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Write([]byte("<div class='error-container fade-in'><h1 class='error-code'>Error " + string(rune(statusCode)) + "</h1><h2 class='error-title'>" + title + "</h2><p class='error-message'>" + message + "</p></div>"))
		return
	}

	tmpl, ok := h.Templates["error.html"]
	if !ok {
		log.Printf("Template not found: error.html")
		http.Error(w, message, statusCode)
		return
	}

	err := tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Printf("Failed to render custom error page: %v", err)
		http.Error(w, message, statusCode)
	}
}
