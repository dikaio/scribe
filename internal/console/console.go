package console

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/dikaio/scribes/internal/config"
	"github.com/dikaio/scribes/internal/content"
)

// Console represents the web-based management console
type Console struct {
	config   config.Config
	sitePath string
	port     int
	tmpl     *template.Template
}

// NewConsole creates a new console instance
func NewConsole(cfg config.Config, sitePath string, port int) *Console {
	return &Console{
		config:   cfg,
		sitePath: sitePath,
		port:     port,
	}
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
	// In a real implementation, templates would be embedded in the binary
	// For now, we'll create them dynamically
	c.tmpl = template.New("console")

	// Add helper functions
	c.tmpl.Funcs(template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("January 2, 2006 15:04")
		},
	})

	// Parse templates
	var err error
	c.tmpl, err = c.tmpl.Parse(baseTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse base template: %w", err)
	}

	for name, tmplContent := range templates {
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
	if err := c.tmpl.ExecuteTemplate(w, "dashboard", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// loadSiteStats loads posts and pages from the content directory
func (c *Console) loadSiteStats(contentPath string) ([]content.Page, []content.Page, error) {
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

// HTML templates
var baseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Scribes Console</title>
    <style>
        :root {
            --primary-color: #0077cc;
            --text-color: #333;
            --background-color: #fff;
            --light-gray: #f5f5f5;
            --border-color: #ddd;
        }
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
            line-height: 1.6;
            color: var(--text-color);
            background-color: var(--background-color);
            display: flex;
            min-height: 100vh;
        }
        
        .sidebar {
            width: 200px;
            background-color: var(--light-gray);
            padding: 20px;
            border-right: 1px solid var(--border-color);
        }
        
        .sidebar h1 {
            font-size: 1.5rem;
            margin-bottom: 20px;
            color: var(--primary-color);
        }
        
        .sidebar nav ul {
            list-style: none;
        }
        
        .sidebar nav li {
            margin-bottom: 10px;
        }
        
        .sidebar nav a {
            color: var(--text-color);
            text-decoration: none;
            display: block;
            padding: 5px 0;
        }
        
        .sidebar nav a:hover {
            color: var(--primary-color);
        }
        
        .main-content {
            flex: 1;
            padding: 20px;
            max-width: 1000px;
        }
        
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--border-color);
        }
        
        .header h2 {
            font-size: 1.8rem;
        }
        
        .btn {
            display: inline-block;
            background-color: var(--primary-color);
            color: white;
            padding: 8px 16px;
            border-radius: 4px;
            text-decoration: none;
            border: none;
            cursor: pointer;
            font-size: 14px;
        }
        
        .btn:hover {
            opacity: 0.9;
        }
        
        .card {
            background-color: white;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            padding: 20px;
            margin-bottom: 20px;
        }
        
        .card h3 {
            margin-bottom: 15px;
            border-bottom: 1px solid var(--border-color);
            padding-bottom: 10px;
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
        }
        
        table th,
        table td {
            text-align: left;
            padding: 10px;
            border-bottom: 1px solid var(--border-color);
        }
        
        table th {
            background-color: var(--light-gray);
        }
        
        .form-group {
            margin-bottom: 15px;
        }
        
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        
        input[type="text"],
        input[type="url"],
        textarea,
        select {
            width: 100%;
            padding: 8px;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            font-family: inherit;
            font-size: inherit;
        }
        
        textarea {
            min-height: 200px;
        }
        
        .dashboard-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
        }
    </style>
</head>
<body>
    <div class="sidebar">
        <h1>Scribes Console</h1>
        <nav>
            <ul>
                <li><a href="/">Dashboard</a></li>
                <li><a href="/content">Content</a></li>
                <li><a href="/settings">Settings</a></li>
                <li><a href="/build">Build Site</a></li>
                <li><a href="/" target="_blank">View Site</a></li>
            </ul>
        </nav>
    </div>
    <div class="main-content">
        {{template "content" .}}
    </div>
</body>
</html>`

var templates = map[string]string{
	"dashboard": `{{define "content"}}
<div class="header">
    <h2>Dashboard</h2>
    <a href="/content/new" class="btn">New Content</a>
</div>

<div class="dashboard-grid">
    <div class="card">
        <h3>Site Overview</h3>
        <ul>
            <li><strong>Title:</strong> {{.Site.Title}}</li>
            <li><strong>URL:</strong> {{.Site.BaseURL}}</li>
            <li><strong>Posts:</strong> {{.PostCount}}</li>
            <li><strong>Pages:</strong> {{.PageCount}}</li>
            <li><strong>Theme:</strong> {{.Site.Theme}}</li>
        </ul>
    </div>
    
    <div class="card">
        <h3>Recent Posts</h3>
        {{if .RecentPosts}}
        <table>
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Date</th>
                </tr>
            </thead>
            <tbody>
                {{range .RecentPosts}}
                <tr>
                    <td><a href="/content/edit?path={{.Path}}">{{.Title}}</a></td>
                    <td>{{formatDate .Date}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <p>No posts found.</p>
        {{end}}
    </div>
    
    <div class="card">
        <h3>Recent Pages</h3>
        {{if .RecentPages}}
        <table>
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Date</th>
                </tr>
            </thead>
            <tbody>
                {{range .RecentPages}}
                <tr>
                    <td><a href="/content/edit?path={{.Path}}">{{.Title}}</a></td>
                    <td>{{formatDate .Date}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <p>No pages found.</p>
        {{end}}
    </div>
</div>
{{end}}`,

	"content": `{{define "content"}}
<div class="header">
    <h2>Content</h2>
    <a href="/content/new" class="btn">New Content</a>
</div>

<div class="card">
    <div style="margin-bottom: 20px;">
        <a href="/content?type=all" {{if eq .ContentType "all"}}style="font-weight: bold;"{{end}}>All</a> |
        <a href="/content?type=posts" {{if eq .ContentType "posts"}}style="font-weight: bold;"{{end}}>Posts</a> |
        <a href="/content?type=pages" {{if eq .ContentType "pages"}}style="font-weight: bold;"{{end}}>Pages</a>
    </div>
    
    {{if .Items}}
    <table>
        <thead>
            <tr>
                <th>Title</th>
                <th>Type</th>
                <th>Date</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td><a href="/content/edit?path={{.Path}}">{{.Title}}</a></td>
                <td>{{if .IsPost}}Post{{else}}Page{{end}}</td>
                <td>{{formatDate .Date}}</td>
                <td>{{if .Draft}}Draft{{else}}Published{{end}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    {{else}}
    <p>No content found.</p>
    {{end}}
</div>
{{end}}`,

	"new_content": `{{define "content"}}
<div class="header">
    <h2>New {{if eq .ContentType "post"}}Post{{else}}Page{{end}}</h2>
</div>

<div class="card">
    <form action="/content/new?type={{.ContentType}}" method="post">
        <div class="form-group">
            <label for="title">Title</label>
            <input type="text" id="title" name="title" required>
        </div>
        
        <div class="form-group">
            <label for="description">Description</label>
            <input type="text" id="description" name="description">
        </div>
        
        {{if eq .ContentType "post"}}
        <div class="form-group">
            <label for="tags">Tags (comma separated)</label>
            <input type="text" id="tags" name="tags">
        </div>
        {{end}}
        
        <div class="form-group">
            <label for="content">Content</label>
            <textarea id="content" name="content" required></textarea>
        </div>
        
        <div class="form-group">
            <label>
                <input type="checkbox" name="draft"> Draft
            </label>
        </div>
        
        <button type="submit" class="btn">Create {{if eq .ContentType "post"}}Post{{else}}Page{{end}}</button>
    </form>
</div>
{{end}}`,

	"settings": `{{define "content"}}
<div class="header">
    <h2>Settings</h2>
</div>

<div class="card">
    <form action="/settings" method="post">
        <div class="form-group">
            <label for="title">Site Title</label>
            <input type="text" id="title" name="title" value="{{.Site.Title}}" required>
        </div>
        
        <div class="form-group">
            <label for="baseURL">Base URL</label>
            <input type="url" id="baseURL" name="baseURL" value="{{.Site.BaseURL}}" required>
        </div>
        
        <div class="form-group">
            <label for="description">Site Description</label>
            <input type="text" id="description" name="description" value="{{.Site.Description}}">
        </div>
        
        <div class="form-group">
            <label for="author">Author</label>
            <input type="text" id="author" name="author" value="{{.Site.Author}}">
        </div>
        
        <div class="form-group">
            <label for="theme">Theme</label>
            <select id="theme" name="theme">
                <option value="default" {{if eq .Site.Theme "default"}}selected{{end}}>Default</option>
                <!-- Add other themes here as they become available -->
            </select>
        </div>
        
        <button type="submit" class="btn">Save Settings</button>
    </form>
</div>
{{end}}`,
}
