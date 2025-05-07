package render

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"

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

func TestTemplateCaching(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-cache-test")
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
	baseTemplatePath := filepath.Join(themeDir, "base.html")
	err = os.WriteFile(baseTemplatePath, []byte(baseContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write base template: %v", err)
	}

	// Create single template
	singleContent := `{{define "content"}}<h1>Original Content</h1>{{end}}`
	singleTemplatePath := filepath.Join(themeDir, "single.html")
	err = os.WriteFile(singleTemplatePath, []byte(singleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write single template: %v", err)
	}

	// Create config
	cfg := config.Config{
		Theme:     "default",
		LayoutDir: "layouts",
	}

	// Test with caching enabled
	t.Run("CachingEnabled", func(t *testing.T) {
		// Create template manager with caching enabled
		tm := NewTemplateManager(cfg)
		tm.EnableCaching()

		// Load templates first time
		err = tm.LoadTemplates(tempDir)
		if err != nil {
			t.Fatalf("Failed to load templates: %v", err)
		}

		// Make sure the template was cached
		if len(tm.cache) == 0 {
			t.Error("Expected cache to contain templates, but it was empty")
		}

		// Get reference to the original template
		originalTemplate, err := tm.GetTemplate("single")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		// Simulate template not changing (reload with same files)
		err = tm.LoadTemplates(tempDir)
		if err != nil {
			t.Fatalf("Failed to reload templates: %v", err)
		}

		// Get reference to the template again
		unchangedTemplate, err := tm.GetTemplate("single")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		// Verify the template was NOT reloaded
		if originalTemplate != unchangedTemplate {
			t.Error("Expected template to be reused from cache when file hasn't changed")
		}

		// Modify the template file
		updatedContent := `{{define "content"}}<h1>Updated Content</h1>{{end}}`
		err = os.WriteFile(singleTemplatePath, []byte(updatedContent), 0644)
		if err != nil {
			t.Fatalf("Failed to update template file: %v", err)
		}
		
		// Wait a moment to ensure filesystem has time to update
		time.Sleep(10 * time.Millisecond)

		// Reload templates
		err = tm.LoadTemplates(tempDir)
		if err != nil {
			t.Fatalf("Failed to reload templates: %v", err)
		}

		// Get reference to the reloaded template
		reloadedTemplate, err := tm.GetTemplate("single")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		// Verify the template was updated in cache
		if originalTemplate == reloadedTemplate {
			t.Error("Expected template to be reloaded when file changes, but got same template reference")
		}
	})

	// Test with caching disabled
	t.Run("CachingDisabled", func(t *testing.T) {
		// Create template manager with caching disabled
		tm := NewTemplateManager(cfg)
		tm.DisableCaching()

		// Load templates first time
		err = tm.LoadTemplates(tempDir)
		if err != nil {
			t.Fatalf("Failed to load templates: %v", err)
		}

		// Make sure the template cache is empty
		if len(tm.cache) != 0 {
			t.Error("Expected cache to be empty when caching is disabled")
		}

		// Get reference to the original template
		originalTemplate, err := tm.GetTemplate("single")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		// Reload templates
		err = tm.LoadTemplates(tempDir)
		if err != nil {
			t.Fatalf("Failed to reload templates: %v", err)
		}

		// Get reference to the reloaded template
		reloadedTemplate, err := tm.GetTemplate("single")
		if err != nil {
			t.Fatalf("Failed to get template: %v", err)
		}

		// Verify the template was reloaded even without changes
		if originalTemplate == reloadedTemplate {
			t.Error("Expected template to be reloaded when caching is disabled, but got same template reference")
		}
	})
}

func TestGetFileModTime(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "modtime-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files with different timestamps
	file1Path := filepath.Join(tempDir, "file1.txt")
	file2Path := filepath.Join(tempDir, "file2.txt")
	
	// Create first file
	err = os.WriteFile(file1Path, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file1: %v", err)
	}
	
	// Wait a moment to ensure different timestamp
	time.Sleep(100 * time.Millisecond)
	
	// Create second file (newer)
	err = os.WriteFile(file2Path, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file2: %v", err)
	}
	
	// Get file info for verification
	file1Info, err := os.Stat(file1Path)
	if err != nil {
		t.Fatalf("Failed to stat file1: %v", err)
	}
	
	file2Info, err := os.Stat(file2Path)
	if err != nil {
		t.Fatalf("Failed to stat file2: %v", err)
	}
	
	// Test with single file
	t.Run("SingleFile", func(t *testing.T) {
		modTime, err := getFileModTime(file1Path)
		if err != nil {
			t.Fatalf("Failed to get mod time: %v", err)
		}
		
		if !modTime.Equal(file1Info.ModTime()) {
			t.Errorf("Expected mod time %v, got %v", file1Info.ModTime(), modTime)
		}
	})
	
	// Test with multiple files
	t.Run("MultipleFiles", func(t *testing.T) {
		modTime, err := getFileModTime(file1Path, file2Path)
		if err != nil {
			t.Fatalf("Failed to get mod time: %v", err)
		}
		
		// Should return the latest mod time
		if !modTime.Equal(file2Info.ModTime()) {
			t.Errorf("Expected latest mod time %v, got %v", file2Info.ModTime(), modTime)
		}
	})
	
	// Test with non-existent file
	t.Run("NonExistentFile", func(t *testing.T) {
		_, err := getFileModTime(filepath.Join(tempDir, "nonexistent.txt"))
		if err == nil {
			t.Error("Expected error for non-existent file, got nil")
		}
	})
}

func TestTemplateNeedsUpdate(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "template-update-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test file
	filePath := filepath.Join(tempDir, "test.html")
	err = os.WriteFile(filePath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Get file mod time
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	fileModTime := fileInfo.ModTime()

	cfg := config.Config{}
	tm := NewTemplateManager(cfg)

	// Test when cache is empty
	t.Run("EmptyCache", func(t *testing.T) {
		tm.cache = make(map[string]TemplateCache)
		needsUpdate, _, err := tm.templateNeedsUpdate("test", []string{filePath})
		if err != nil {
			t.Fatalf("Error checking if template needs update: %v", err)
		}
		if !needsUpdate {
			t.Error("Expected template to need update when cache is empty")
		}
	})

	// Test when cache entry exists but is older
	t.Run("OlderCacheEntry", func(t *testing.T) {
		tm.cache = make(map[string]TemplateCache)
		tm.cache["test"] = TemplateCache{
			ModTime: fileModTime.Add(-1 * time.Hour), // 1 hour older
			Files:   []string{filePath},
		}
		
		needsUpdate, _, err := tm.templateNeedsUpdate("test", []string{filePath})
		if err != nil {
			t.Fatalf("Error checking if template needs update: %v", err)
		}
		if !needsUpdate {
			t.Error("Expected template to need update when cache entry is older")
		}
	})

	// Test when cache entry exists and is newer
	t.Run("NewerCacheEntry", func(t *testing.T) {
		tm.cache = make(map[string]TemplateCache)
		tm.cache["test"] = TemplateCache{
			ModTime: fileModTime.Add(1 * time.Hour), // 1 hour newer (shouldn't happen in practice)
			Files:   []string{filePath},
		}
		
		needsUpdate, _, err := tm.templateNeedsUpdate("test", []string{filePath})
		if err != nil {
			t.Fatalf("Error checking if template needs update: %v", err)
		}
		if needsUpdate {
			t.Error("Expected template to not need update when cache entry is newer")
		}
	})

	// Test when file list has changed
	t.Run("DifferentFiles", func(t *testing.T) {
		tm.cache = make(map[string]TemplateCache)
		tm.cache["test"] = TemplateCache{
			ModTime: fileModTime,
			Files:   []string{filePath, filepath.Join(tempDir, "other.html")}, // Different list
		}
		
		needsUpdate, _, err := tm.templateNeedsUpdate("test", []string{filePath})
		if err != nil {
			t.Fatalf("Error checking if template needs update: %v", err)
		}
		if !needsUpdate {
			t.Error("Expected template to need update when file list has changed")
		}
	})
}