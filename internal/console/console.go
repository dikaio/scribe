package console

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
	"github.com/dikaio/scribe/internal/templates"
)

// Console represents the web-based management console
type Console struct {
	config              config.Config
	sitePath            string
	port                int
	tmpl                *template.Template
	loadSiteStats       func(string) ([]content.Page, []content.Page, error)
	originalLoadSiteStats func(string) ([]content.Page, []content.Page, error)
}

// NewConsole creates a new console instance
func NewConsole(cfg config.Config, sitePath string, port int) *Console {
	c := &Console{
		config:   cfg,
		sitePath: sitePath,
		port:     port,
	}
	c.loadSiteStats = c.loadContentStats
	return c
}

// Start starts the console server
func (c *Console) Start() error {
	// Initialize templates
	if err := c.initTemplates(); err != nil {
		return err
	}

	// Register handlers
	http.HandleFunc("/", c.handleDashboard)
	http.HandleFunc("/content", c.handleContent)
	http.HandleFunc("/content/new", c.handleNewContent)
	http.HandleFunc("/settings", c.handleSettings)
	http.HandleFunc("/build", c.handleBuild)

	// Start server
	addr := fmt.Sprintf(":%d", c.port)
	fmt.Printf("Console running at http://localhost:%d/\n", c.port)
	return http.ListenAndServe(addr, nil)
}

// initTemplates initializes the HTML templates for the console
func (c *Console) initTemplates() error {
	// Create a new template
	c.tmpl = template.New("console")

	// Add helper functions
	c.tmpl.Funcs(template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("January 2, 2006 15:04")
		},
	})

	// Parse templates from the templates package
	var err error
	c.tmpl, err = c.tmpl.Parse(templates.ConsoleBaseTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse base template: %w", err)
	}

	for name, tmplContent := range templates.ConsoleTemplates {
		_, err = c.tmpl.New(name).Parse(tmplContent)
		if err != nil {
			return fmt.Errorf("failed to parse %s template: %w", name, err)
		}
	}

	return nil
}

// handleDashboard handles the dashboard page
func (c *Console) handleDashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Load site stats
	contentPath := filepath.Join(c.sitePath, c.config.ContentDir)
	posts, pages, err := c.loadSiteStats(contentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"Site":        c.config,
		"Title":       "Dashboard",
		"PostCount":   len(posts),
		"PageCount":   len(pages),
		"RecentPosts": getRecentItems(posts, 5),
		"RecentPages": getRecentItems(pages, 5),
	}

	// Render template
	w.Header().Set("Content-Type", "text/html")
	
	// Debug template availability
	if c.tmpl == nil {
		errMsg := "Template is nil in handleDashboard"
		fmt.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	
	if err := c.tmpl.ExecuteTemplate(w, "dashboard", data); err != nil {
		errMsg := fmt.Sprintf("Template execution error: %v", err)
		fmt.Println(errMsg) // Log error to console for debugging
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

// handleContent handles the content listing page
func (c *Console) handleContent(w http.ResponseWriter, r *http.Request) {
	contentType := r.URL.Query().Get("type")
	if contentType == "" {
		contentType = "all"
	}

	// Load content
	contentPath := filepath.Join(c.sitePath, c.config.ContentDir)
	posts, pages, err := c.loadSiteStats(contentPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []content.Page
	switch contentType {
	case "posts":
		items = posts
	case "pages":
		items = pages
	default:
		items = append(posts, pages...)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Site":        c.config,
		"Title":       "Content",
		"ContentType": contentType,
		"Items":       items,
	}

	// Render template
	w.Header().Set("Content-Type", "text/html")
	if err := c.tmpl.ExecuteTemplate(w, "content", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleNewContent handles the new content page
func (c *Console) handleNewContent(w http.ResponseWriter, r *http.Request) {
	contentType := r.URL.Query().Get("type")
	if contentType == "" {
		contentType = "post"
	}

	if r.Method == "POST" {
		// Process form submission (not implemented)
		http.Redirect(w, r, "/content", http.StatusSeeOther)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"Site":        c.config,
		"Title":       "New " + contentType,
		"ContentType": contentType,
	}

	// Render template
	w.Header().Set("Content-Type", "text/html")
	if err := c.tmpl.ExecuteTemplate(w, "new_content", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleSettings handles the settings page
func (c *Console) handleSettings(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Process form submission (not implemented)
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"Site":  c.config,
		"Title": "Settings",
	}

	// Render template
	w.Header().Set("Content-Type", "text/html")
	if err := c.tmpl.ExecuteTemplate(w, "settings", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleBuild handles the build action
func (c *Console) handleBuild(w http.ResponseWriter, r *http.Request) {
	// In a real implementation, this would trigger the build process
	// For now, just redirect back to the dashboard
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// loadContentStats loads posts and pages from the content directory
func (c *Console) loadContentStats(contentPath string) ([]content.Page, []content.Page, error) {
	var posts, pages []content.Page

	// Walk content directory
	err := filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// Load page
		page, err := content.LoadPage(path, c.config.BaseURL)
		if err != nil {
			return err
		}

		// Add to appropriate list
		if page.IsPost {
			posts = append(posts, page)
		} else {
			pages = append(pages, page)
		}

		return nil
	})

	return posts, pages, err
}

// getRecentItems returns the n most recent items
func getRecentItems(items []content.Page, n int) []content.Page {
	// Sort by date (newest first)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.After(items[j].Date)
	})

	// Return at most n items
	if len(items) > n {
		return items[:n]
	}
	return items
}

