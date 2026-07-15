package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"quickshare/internal/handlers"
	"quickshare/internal/middleware"
	"quickshare/internal/services"
	"time"
)

func main() {
	log.Println("Initializing QuickShare application...")

	// 1. Resolve Environment Configurations
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		// Fallback for local development
		appURL = "http://localhost:3000"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 2. Instantiate services
	storage, err := services.NewStorageService("data")
	if err != nil {
		log.Fatalf("CRITICAL: Failed to initialize storage service: %v", err)
	}

	// 3. Compile Go templates
	templates := make(map[string]*template.Template)
	pages := []string{
		"home.html",
		"view.html",
		"about.html",
		"contact.html",
		"privacy.html",
		"terms.html",
		"error.html",
	}
	for _, page := range pages {
		t, err := template.ParseFiles("templates/base.html", "templates/"+page)
		if err != nil {
			log.Fatalf("CRITICAL: Failed to compile HTML template %s: %v", page, err)
		}
		templates[page] = t
	}

	// 4. Create centralized route handler context
	h := handlers.NewHandler(storage, templates, appURL)

	// 5. Setup route multiplexer using Go 1.22+ routing features
	mux := http.NewServeMux()

	// Web pages routes
	mux.HandleFunc("GET /", h.Home)
	mux.HandleFunc("GET /about", h.About)
	mux.HandleFunc("GET /contact", h.Contact)
	mux.HandleFunc("POST /contact/submit", h.ContactSubmit)
	mux.HandleFunc("GET /privacy", h.Privacy)
	mux.HandleFunc("GET /terms", h.Terms)

	// Snippets operations
	mux.HandleFunc("POST /snippet", h.CreateSnippet)
	mux.HandleFunc("GET /snippet/{id}", h.ViewSnippet)
	mux.HandleFunc("GET /raw/{id}", h.RawSnippet)
	mux.HandleFunc("GET /snippet/{id}/download", h.DownloadSnippet)

	// SEO, Uptime and indexing assets
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /sitemap.xml", h.Sitemap)
	
	// Serve static files (CSS, JS, favicon)
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	
	// Also explicitly handle robots.txt and favicon.ico from static directory
	mux.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/robots.txt")
	})
	mux.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/favicon.ico")
	})

	// 6. Setup Rate limiting
	// Allows 2 requests per client per second with a burst threshold of 20
	rl := middleware.NewRateLimiter(500*time.Millisecond, 20)

	// 7. Chain global middlewares
	var handler http.Handler = mux
	handler = rl.Limit(handler)
	handler = middleware.Gzip(handler)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.Logger(handler)
	handler = middleware.Recovery(handler)

	// 8. Launch Server
	addr := "0.0.0.0:" + port
	log.Printf("Server listening successfully on http://%s (Public App URL: %s)", addr, appURL)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("CRITICAL: Server crashed: %v", err)
	}
}
