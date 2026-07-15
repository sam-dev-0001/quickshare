# QuickShare

QuickShare is a secure, lightweight, and modern server-rendered web application for sharing text and code snippets with unique, secure links. Built as an original product with an aesthetic dark design and a highly performant architecture, QuickShare utilizes **Go** for backend logic, **HTMX** for progress-enhanced frontend interactivity, and standard file-based JSON storage.

---

## 🚀 Key Features

* **Instant Creation:** Create text or code shares with custom styling, high-speed delivery, and low overhead.
* **Auto-Expiration Protection:** Automatically self-cleans snippets. Select duration intervals of **10 Minutes**, **1 Hour**, **1 Day**, **1 Week**, or **Never**.
* **Beautiful Syntax Highlighting:** Integrated seamlessly for 12 popular computer, query, and markup languages using the **Highlight.js** engine.
* **Modern Tab Indentation Support:** High-efficiency code editing with standard tab overrides (converts Tab key inputs into 4 spacing blocks without breaking field focus).
* **One-Click Clipboard copy:** Instant snippet copying via the native Browser Clipboard API.
* **Robust Security Suite:** Hardened against XSS injections, clickjacking, mime sniffing, and heavy automated bot requests using custom token-bucket rate limiters.
* **Dynamic SEO Integration:** Dynamic standard XML sitemap generator, `robots.txt` crawler policies, and OpenGraph/Twitter card headers for search engines.
* **AdSense Compliant Layout:** Features pre-configured, fluid Google AdSense placeholders along with essential Legal Policy pages (About, Contact, Privacy, and Terms of Service) for direct advertising compliance.

---

## 🛠️ Architecture and Tech Stack

* **Backend Engine:** Go (Standard Library `net/http`) — zero router dependencies.
* **Interactions:** HTMX — for high-performance partial layout swaps and fast interactive state transitions.
* **Styling Framework:** Custom Vanilla CSS3 — engineered with cohesive, eye-safe custom variables, fluid grid structures, and touch-target safe dimensions.
* **Syntax Highlighting:** Highlight.js (github-dark theme CDN).
* **Storage Layer:** File-system JSON database storage inside the `data/` directory.

---

## 📂 Project Structure

```text
├── cmd/
│   └── server/
│       └── main.go         # App Entrypoint, Router compilation & server listener
├── internal/
│   ├── handlers/
│   │   ├── handlers.go     # Core HTTP handlers, error rendering, global parameters
│   │   ├── home.go         # Homepage, static legal pages, HTMX form responders
│   │   ├── seo.go          # Dynamic XML sitemap generator
│   │   └── snippet.go      # Snippet routing controllers (Create, View, Raw, Download, Health)
│   ├── middleware/
│   │   └── middleware.go   # Global stack (Logging, Recovery, Gzip, CSRF/XSS, Rate limiter)
│   ├── models/
│   │   └── snippet.go      # Snippet structures and active expiration computations
│   └── services/
│       └── storage.go      # File-based JSON CRUD & periodic background self-cleaner
├── static/
│   ├── css/
│   │   └── style.css       # Unified modern CSS3 design sheet
│   ├── js/
│   │   └── main.js         # Textarea Tab overrides, clipboard actions & Toast UI engine
│   ├── robots.txt          # Crawler instructions
│   └── favicon.ico         # Placeholder icon
├── templates/
│   ├── base.html           # HTML5 structural layout shell
│   ├── home.html           # Text snippet creation canvas
│   ├── view.html           # Snippet viewer stage with custom syntax runner
│   ├── error.html          # Custom error views (400, 404, 429, 500)
│   ├── about.html          # Informative About Page
│   ├── contact.html        # Interactive HTMX Support Page
│   ├── privacy.html        # AdSense-ready Privacy Policy
│   └── terms.html          # AdSense-ready Terms of Service
├── go.mod                  # Go dependency manifest
├── package.json            # Node wrapper dev helper for sandboxed environments
├── Dockerfile              # Multi-stage optimized builder & production runner container
├── .dockerignore           # Container builder ignore listings
└── render.yaml             # Render deployment configurations
```

---

## 💻 Local Development Setup

To boot QuickShare locally on your computer, ensure you have **Go 1.22+** installed:

1. **Clone the repository** and navigate to the project directory.
2. **Start the server:**
   ```bash
   go run cmd/server/main.go
   ```
3. **Open in Browser:** Visit [http://localhost:3000](http://localhost:3000)

---

## 🧪 Running Unit Tests

Execute Go's internal testing suite:
```bash
go test -v ./...
```

---

## 🐳 Docker Deployment

To build and run the self-contained production Docker container locally:

1. **Build the image:**
   ```bash
   docker build -t quickshare:latest .
   ```
2. **Launch the container:**
   ```bash
   docker run -p 3000:3000 -v $(pwd)/data:/app/data quickshare:latest
   ```

---

## 🌐 Cloud Deployment (Render)

QuickShare is ready for immediate deployment to **Render** via our pre-configured `render.yaml` infrastructure-as-code template:

1. Push this codebase to your **GitHub** account.
2. Go to **Render Dashboard** and select **Blueprints**.
3. Link your GitHub repository and click **Deploy**.
4. Configure the following Environment Variable in your Render dashboard:
   - `APP_URL`: Set this to your live Render public URL (e.g. `https://quickshare.onrender.com`) to enable dynamic Sitemap generation and OpenGraph canonical links.
