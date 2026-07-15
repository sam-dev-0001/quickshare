package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// Sitemap generates and streams a compliant, dynamic XML sitemap.
func (h *Handler) Sitemap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	now := time.Now().UTC().Format("2006-01-02")

	sitemap := `<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <url>
        <loc>` + h.AppURL + `/</loc>
        <lastmod>` + now + `</lastmod>
        <changefreq>daily</changefreq>
        <priority>1.0</priority>
    </url>
    <url>
        <loc>` + h.AppURL + `/about</loc>
        <lastmod>` + now + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.8</priority>
    </url>
    <url>
        <loc>` + h.AppURL + `/contact</loc>
        <lastmod>` + now + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.7</priority>
    </url>
    <url>
        <loc>` + h.AppURL + `/privacy</loc>
        <lastmod>` + now + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.5</priority>
    </url>
    <url>
        <loc>` + h.AppURL + `/terms</loc>
        <lastmod>` + now + `</lastmod>
        <changefreq>monthly</changefreq>
        <priority>0.5</priority>
    </url>
</urlset>`

	_, _ = fmt.Fprint(w, sitemap)
}
