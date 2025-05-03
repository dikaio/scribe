package render

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/dikaio/scribe/internal/config"
)

func TestNewTemplateManager(t *testing.T) {
	// Create a simple config
	cfg := config.Config{
		Theme: "default",
	}

	// Create a new template manager
	tm := NewTemplateManager(cfg)

	// Check that the template manager was created successfully
	if tm == nil {
		t.Fatal("Expected template manager to be created, got nil")
	}

	// Check that the templates map was initialized
	if tm.templates == nil {
		t.Error("Expected templates map to be initialized")
	}

	// Check that the config was set correctly
	if tm.config.Theme != "default" {
		t.Errorf("Expected theme to be 'default', got '%s'", tm.config.Theme)
	}
}

func TestGetTemplate(t *testing.T) {
	// Create a template manager with a template
	tm := &TemplateManager{
		templates: make(map[string]*template.Template),
		config:    config.Config{},
	}

	// Add a template
	tmpl, err := template.New("test").Parse("Hello {{.Name}}")
	if err != nil {
		t.Fatalf("Failed to parse template: %v", err)
	}
	tm.templates["test"] = tmpl

	// Test getting an existing template
	got, err := tm.GetTemplate("test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if got == nil {
		t.Error("Expected template, got nil")
	}

	// Test getting a non-existent template
	_, err = tm.GetTemplate("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

func TestLoadTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	themeDir := filepath.Join(tempDir, "themes", "default", "layouts")
	siteLayoutDir := filepath.Join(tempDir, "layouts")

	// Create theme directories
	err = os.MkdirAll(themeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create theme dir: %v", err)
	}

	// Create site layout directory
	err = os.MkdirAll(siteLayoutDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create site layout dir: %v", err)
	}

	// Create base template in theme
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

	// Create single template in theme
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

	// Create home template in site (override)
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
	err = os.WriteFile(filepath.Join(siteLayoutDir, "home.html"), []byte(homeContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write home template: %v", err)
	}

	// Create config
	cfg := config.Config{
		Theme:     "default",
		LayoutDir: "layouts",
	}

	// Create template manager
	tm := NewTemplateManager(cfg)

	// Load templates
	err = tm.LoadTemplates(tempDir)
	if err != nil {
		t.Fatalf("Failed to load templates: %v", err)
	}

	// Check that templates were loaded
	if len(tm.templates) != 3 {
		t.Errorf("Expected 3 templates, got %d", len(tm.templates))
	}

	// Check if base template exists
	_, err = tm.GetTemplate("base")
	if err != nil {
		t.Errorf("Expected base template to be loaded, got error: %v", err)
	}

	// Check if single template exists
	_, err = tm.GetTemplate("single")
	if err != nil {
		t.Errorf("Expected single template to be loaded, got error: %v", err)
	}

	// Check if home template exists
	homeTemplate, err := tm.GetTemplate("home")
	if err != nil {
		t.Errorf("Expected home template to be loaded, got error: %v", err)
	}

	// Test that the home template renders correctly
	if homeTemplate != nil {
		var buf bytes.Buffer
		data := map[string]interface{}{
			"Site": struct {
				Title string
			}{
				Title: "Test Site",
			},
			"Pages": []struct {
				Title string
				URL   string
			}{
				{Title: "Test Post", URL: "/posts/test-post"},
			},
		}

		err = homeTemplate.Execute(&buf, data)
		if err != nil {
			t.Errorf("Failed to execute template: %v", err)
		}

		if !bytes.Contains(buf.Bytes(), []byte("Welcome to Test Site")) {
			t.Error("Template did not render correctly")
		}
	}
}

func TestLoadTemplatesNoBaseTemplate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-test-no-base")
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

	// Create single template without base
	singleContent := `<article>
  <h1>{{.Page.Title}}</h1>
  <div>{{.Content}}</div>
</article>`
	err = os.WriteFile(filepath.Join(themeDir, "single.html"), []byte(singleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write single template: %v", err)
	}

	// Create config
	cfg := config.Config{
		Theme:     "default",
		LayoutDir: "layouts",
	}

	// Create template manager
	tm := NewTemplateManager(cfg)

	// Load templates should fail without base template
	err = tm.LoadTemplates(tempDir)
	if err == nil {
		t.Error("Expected error for missing base template, got nil")
	}
}