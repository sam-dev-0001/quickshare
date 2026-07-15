package handlers

import (
	"encoding/json"
	"net/http"
	"quickshare/internal/models"
	"time"
)

// ViewSnippetData packages the snippet model with formatted display strings.
type ViewSnippetData struct {
	Snippet            *models.Snippet
	FormattedCreatedAt string
	FormattedExpiresAt string
	ExpiresAt          time.Time
}

// CreateSnippet handles posting a new snippet.
func (h *Handler) CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Maximum snippet payload protection (e.g. 10MB limit)
	r.Body = http.MaxBytesReader(w, r.Body, 10*1024*1024)

	err := r.ParseForm()
	if err != nil {
		h.RenderError(w, r, http.StatusBadRequest, "Payload Too Large", "The snippet content you are attempting to upload exceeds our storage limits.")
		return
	}

	text := r.FormValue("text")
	language := r.FormValue("language")
	expiry := r.FormValue("expiry")

	// Input Validation
	if len(text) == 0 {
		h.RenderError(w, r, http.StatusBadRequest, "Validation Error", "Snippet content cannot be completely empty.")
		return
	}
	if len(text) > 500000 {
		h.RenderError(w, r, http.StatusBadRequest, "Validation Error", "Snippet content exceeds the maximum limit of 500,000 characters.")
		return
	}

	// Clean language input
	if language == "" {
		language = "plaintext"
	}

	// Clean expiration input
	switch expiry {
	case models.ExpiryNever, models.Expiry10Min, models.Expiry1Hour, models.Expiry1Day, models.Expiry1Week:
		// valid
	default:
		expiry = models.ExpiryNever
	}

	// Generate secure ID
	id, err := h.Storage.GenerateID()
	if err != nil {
		h.RenderError(w, r, http.StatusInternalServerError, "Storage Error", "Failed to generate a secure sharing link. Please try again.")
		return
	}

	snippet := &models.Snippet{
		ID:        id,
		Text:      text,
		Language:  language,
		Expiry:    expiry,
		CreatedAt: time.Now(),
	}

	err = h.Storage.Save(snippet)
	if err != nil {
		h.RenderError(w, r, http.StatusInternalServerError, "Storage Error", "We failed to write the snippet to server files. Please retry.")
		return
	}

	// If XMLHttpRequest (HTMX), return full view template with pushed URL header
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("HX-Push-Url", "/snippet/"+id)
	}

	// Pre-format times for display
	ctx := ViewSnippetData{
		Snippet:            snippet,
		FormattedCreatedAt: snippet.CreatedAt.Format(time.RFC3339),
		FormattedExpiresAt: snippet.ExpiresAt().Format(time.RFC3339),
		ExpiresAt:          snippet.ExpiresAt(),
	}

	h.RenderTemplate(w, r, "view.html", ctx, http.StatusCreated)
}

// ViewSnippet renders the single snippet template page.
func (h *Handler) ViewSnippet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.RenderError(w, r, http.StatusBadRequest, "Bad Request", "Missing snippet identification key.")
		return
	}

	snippet, err := h.Storage.Get(id)
	if err != nil {
		h.RenderError(w, r, http.StatusNotFound, "Snippet Not Found", "This snippet does not exist, has expired, or was removed by administrators.")
		return
	}

	ctx := ViewSnippetData{
		Snippet:            snippet,
		FormattedCreatedAt: snippet.CreatedAt.Format(time.RFC3339),
		FormattedExpiresAt: snippet.ExpiresAt().Format(time.RFC3339),
		ExpiresAt:          snippet.ExpiresAt(),
	}

	h.RenderTemplate(w, r, "view.html", ctx, http.StatusOK)
}

// RawSnippet renders the clean raw plain text code.
func (h *Handler) RawSnippet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Bad Request: Missing ID", http.StatusBadRequest)
		return
	}

	snippet, err := h.Storage.Get(id)
	if err != nil {
		http.Error(w, "Snippet not found or has expired", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(snippet.Text))
}

// DownloadSnippet forces a file download browser behavior.
func (h *Handler) DownloadSnippet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		h.RenderError(w, r, http.StatusBadRequest, "Bad Request", "Missing snippet identification key.")
		return
	}

	snippet, err := h.Storage.Get(id)
	if err != nil {
		h.RenderError(w, r, http.StatusNotFound, "Snippet Not Found", "This snippet does not exist, has expired, or was removed.")
		return
	}

	// Choose appropriate filename extension based on language type
	ext := ".txt"
	switch snippet.Language {
	case "javascript":
		ext = ".js"
	case "go":
		ext = ".go"
	case "python":
		ext = ".py"
	case "html":
		ext = ".html"
	case "css":
		ext = ".css"
	case "json":
		ext = ".json"
	case "rust":
		ext = ".rs"
	case "cpp":
		ext = ".cpp"
	case "java":
		ext = ".java"
	case "sql":
		ext = ".sql"
	case "markdown":
		ext = ".md"
	}

	filename := "quickshare_" + id + ext

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(snippet.Text))
}

// Health is a simple uptime diagnostics endpoint.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "timestamp": time.Now().UTC().Format(time.RFC3339)})
}
