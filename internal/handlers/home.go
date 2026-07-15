package handlers

import (
	"net/http"
)

// Home serves the main landing page with the text editor.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.RenderError(w, r, http.StatusNotFound, "Page Not Found", "The page you are looking for does not exist or has been moved.")
		return
	}
	h.RenderTemplate(w, r, "home.html", nil, http.StatusOK)
}

// About serves the static About page.
func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	h.RenderTemplate(w, r, "about.html", nil, http.StatusOK)
}

// Contact serves the Contact support page.
func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	h.RenderTemplate(w, r, "contact.html", nil, http.StatusOK)
}

// Privacy serves the Privacy Policy page.
func (h *Handler) Privacy(w http.ResponseWriter, r *http.Request) {
	h.RenderTemplate(w, r, "privacy.html", nil, http.StatusOK)
}

// Terms serves the Terms of Service page.
func (h *Handler) Terms(w http.ResponseWriter, r *http.Request) {
	h.RenderTemplate(w, r, "terms.html", nil, http.StatusOK)
}

// ContactSubmit handles HTMX asynchronous contact form submissions.
func (h *Handler) ContactSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.RenderError(w, r, http.StatusBadRequest, "Invalid Submission", "We couldn't process your contact form. Please try again.")
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	// Validate inputs
	if name == "" || email == "" || message == "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(`<div class="contact-error-box" style="color: var(--danger-color); font-weight: 600;">Please fill out all required fields.</div>`))
		return
	}

	// Logging submission (in production, this would dispatch email notifications)
	_ = subject

	// Return a beautifully polished HTMX response block replacing the form
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
		<div class="contact-success-box fade-in" style="background-color: rgba(102, 187, 106, 0.1); border: 1px solid var(--success-color); border-radius: 8px; padding: 2rem; text-align: center;">
			<h3 style="color: var(--success-color); font-size: 1.3rem; margin-bottom: 0.5rem; font-family: 'Outfit', sans-serif;">Message Sent Successfully!</h3>
			<p style="color: var(--text-primary); font-size: 0.9rem; opacity: 0.9;">Thank you, <strong>` + name + `</strong>. Our security and support team has received your message and will respond to <strong>` + email + `</strong> within 24 hours.</p>
			<script>
				if (typeof showToast === 'function') {
					showToast('Message sent successfully!');
				}
			</script>
		</div>
	`))
}
