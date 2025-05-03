package render

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
)

// setupTestEnvironment creates a temporary test environment with templates
func setupTestEnvironment(t *testing.T) (string, *Renderer, func()) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "renderer-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test directory structure
	themeDir := filepath.Join(tempDir, "themes", "default", "layouts")
	outputDir := filepath.Join(tempDir, "public")

	// Create directories
	err = os.MkdirAll(themeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create theme dir: %v", err)
	}
	err = os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create output dir: %v", err)
	}

	// Create base template
	baseContent := `<!DOCTYPE html>
<html>
<head>
  <title>{{.Site.Title}} - {{.Page.Title}}</title>
</head>
<body>
  {{block "content" .}}Default content{{end}}
</body>
</html>`
	err = os.WriteFile(filepath.Join(themeDir, "base.html"), []byte(baseContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write base template: %v", err)
	}

	// Create single template
	singleContent := `{{define "content"}}
<article>
  <h1>{{.Page.Title}}</h1>
  <div>{{.Content}}</div>
</article>
{{end}}`
	err = os.WriteFile(filepath.Join(themeDir, "single.html"), []byte(singleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write single template: %v", err)
	}

	// Create page template
	pageContent := `{{define "content"}}
<div class="page">
  <h1>{{.Page.Title}}</h1>
  <div>{{.Content}}</div>
</div>
{{end}}`
	err = os.WriteFile(filepath.Join(themeDir, "page.html"), []byte(pageContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write page template: %v", err)
	}

	// Create list template
	listContent := `{{define "content"}}
<div class="list">
  <h1>{{.Title}}</h1>
  <ul>
    {{range .Pages}}
    <li><a href="{{.URL}}">{{.Title}}</a></li>
    {{end}}
  </ul>
</div>
{{end}}`
	err = os.WriteFile(filepath.Join(themeDir, "list.html"), []byte(listContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write list template: %v", err)
	}

	// Create home template
	homeContent := `{{define "content"}}
<div class="home">
  <h1>Welcome to {{.Site.Title}}</h1>
  <div class="posts">
    {{range .Pages}}
    <div class="post">
      <h2><a href="{{.URL}}">{{.Title}}</a></h2>
    </div>
    {{end}}
  </div>
</div>
{{end}}`
	err = os.WriteFile(filepath.Join(themeDir, "home.html"), []byte(homeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write home template: %v", err)
	}

	// Create config
	cfg := config.Config{
		Title:     "Test Site",
		Theme:     "default",
		LayoutDir: "layouts",
		OutputDir: "public",
	}

	// Create renderer
	renderer := NewRenderer(cfg)

	// Initialize renderer with templates
	err = renderer.Init(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize renderer: %v", err)
	}

	// Return temp dir, renderer, and cleanup function
	return tempDir, renderer, func() {
		os.RemoveAll(tempDir)
	}
}

func TestNewRenderer(t *testing.T) {
	// Create a simple config
	cfg := config.Config{
		Theme: "default",
	}

	// Create a new renderer
	r := NewRenderer(cfg)

	// Check that the renderer was created successfully
	if r == nil {
		t.Fatal("Expected renderer to be created, got nil")
	}

	// Check that the template manager was initialized
	if r.templateManager == nil {
		t.Error("Expected template manager to be initialized")
	}

	// Check that the config was set correctly
	if r.config.Theme != "default" {
		t.Errorf("Expected theme to be 'default', got '%s'", r.config.Theme)
	}
}

func TestInit(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "renderer-init-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	themeDir := filepath.Join(tempDir, "themes", "default", "layouts")
	err = os.MkdirAll(themeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create theme dir: %v", err)
	}

	// Create base template
	baseContent := `<!DOCTYPE html><html><body>{{block "content" .}}{{end}}</body></html>`
	err = os.WriteFile(filepath.Join(themeDir, "base.html"), []byte(baseContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write base template: %v", err)
	}

	// Create config
	cfg := config.Config{
		Theme:     "default",
		LayoutDir: "layouts",
	}

	// Create renderer
	r := NewRenderer(cfg)

	// Initialize renderer
	err = r.Init(tempDir)
	if err != nil {
		t.Fatalf("Failed to initialize renderer: %v", err)
	}

	// Initialize with invalid path
	err = r.Init("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestRenderPage(t *testing.T) {
	// Setup test environment
	tempDir, renderer, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create a test page
	page := content.Page{
		Title:   "Test Page",
		URL:     "/test-page/",
		Content: "# Test\nThis is a test.",
		HTML:    "<h1>Test</h1>\n<p>This is a test.</p>",
		Date:    time.Now(),
		IsPost:  true,
	}

	// Create output path
	outputPath := filepath.Join(tempDir, "public", "test-page", "index.html")

	// Render the page
	err := renderer.RenderPage(page, outputPath)
	if err != nil {
		t.Fatalf("Failed to render page: %v", err)
	}

	// Check that the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to exist at %s", outputPath)
	}

	// Check file content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify that the content contains expected elements
	expectedElements := []string{
		"<title>Test Site - Test Page</title>",
		"<h1>Test Page</h1>",
		"<h1>Test</h1>",
		"<p>This is a test.</p>",
	}

	for _, element := range expectedElements {
		if !bytes.Contains(content, []byte(element)) {
			t.Errorf("Expected output to contain %q", element)
		}
	}

	// Test with non-post page
	page.IsPost = false
	page.Title = "About Page"
	page.URL = "/about/"

	outputPath = filepath.Join(tempDir, "public", "about", "index.html")

	err = renderer.RenderPage(page, outputPath)
	if err != nil {
		t.Fatalf("Failed to render page: %v", err)
	}

	// Check that the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to exist at %s", outputPath)
	}
}

func TestRenderList(t *testing.T) {
	// Setup test environment
	tempDir, renderer, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test pages
	pages := []content.Page{
		{
			Title:  "Test Post 1",
			URL:    "/post1/",
			IsPost: true,
		},
		{
			Title:  "Test Post 2",
			URL:    "/post2/",
			IsPost: true,
		},
	}

	// Create output path
	outputPath := filepath.Join(tempDir, "public", "tags", "test", "index.html")

	// Render the list
	err := renderer.RenderList("Test Tag", pages, outputPath)
	if err != nil {
		t.Fatalf("Failed to render list: %v", err)
	}

	// Check that the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to exist at %s", outputPath)
	}

	// Check file content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify that the content contains expected elements
	expectedElements := []string{
		"<h1>Test Tag</h1>",
		"<li><a href=\"/post1/\">Test Post 1</a></li>",
		"<li><a href=\"/post2/\">Test Post 2</a></li>",
	}

	for _, element := range expectedElements {
		if !bytes.Contains(content, []byte(element)) {
			t.Errorf("Expected output to contain %q", element)
		}
	}
}

func TestRenderHome(t *testing.T) {
	// Setup test environment
	tempDir, renderer, cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Create test pages
	pages := []content.Page{
		{
			Title:  "Test Post 1",
			URL:    "/post1/",
			IsPost: true,
		},
		{
			Title:  "Test Post 2",
			URL:    "/post2/",
			IsPost: true,
		},
	}

	// Create output path
	outputPath := filepath.Join(tempDir, "public", "index.html")

	// Render the home page
	err := renderer.RenderHome(pages, outputPath)
	if err != nil {
		t.Fatalf("Failed to render home: %v", err)
	}

	// Check that the file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Expected output file to exist at %s", outputPath)
	}

	// Check file content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Verify that the content contains expected elements
	expectedElements := []string{
		"<h1>Welcome to Test Site</h1>",
		"<h2><a href=\"/post1/\">Test Post 1</a></h2>",
		"<h2><a href=\"/post2/\">Test Post 2</a></h2>",
	}

	for _, element := range expectedElements {
		if !bytes.Contains(content, []byte(element)) {
			t.Errorf("Expected output to contain %q", element)
		}
	}
}