package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp()

	// Check that the app was created successfully
	if app == nil {
		t.Fatal("Expected app to be created, got nil")
	}

	// Check that the app name is set
	if app.Name != "scribe" {
		t.Errorf("Expected app name to be 'scribe', got '%s'", app.Name)
	}

	// Check that the app version is set
	if app.Version == "" {
		t.Error("Expected app version to be set")
	}

	// Check that commands were registered
	expectedCommands := []string{"build", "serve", "new"}
	for _, cmd := range expectedCommands {
		if _, exists := app.Commands[cmd]; !exists {
			t.Errorf("Expected command '%s' to be registered", cmd)
		}
	}
}

func TestRegisterCommands(t *testing.T) {
	app := &App{
		Name:     "scribe-test",
		Version:  "0.0.1",
		Commands: make(map[string]Command),
	}

	// Register commands
	app.registerCommands()

	// Check that all expected commands are registered
	expectedCommands := []string{"build", "serve", "new"}
	for _, cmd := range expectedCommands {
		if _, exists := app.Commands[cmd]; !exists {
			t.Errorf("Expected command '%s' to be registered", cmd)
		}
	}

	// Check that each command has a name, description, and action
	for name, cmd := range app.Commands {
		if cmd.Name != name {
			t.Errorf("Expected command name to be '%s', got '%s'", name, cmd.Name)
		}
		if cmd.Description == "" {
			t.Errorf("Expected command '%s' to have a description", name)
		}
		if cmd.Action == nil {
			t.Errorf("Expected command '%s' to have an action", name)
		}
	}
}

func TestRun(t *testing.T) {
	app := NewApp()

	// Test with no args
	err := app.Run([]string{"scribe"})
	if err != nil {
		t.Errorf("Expected no error for help command, got: %v", err)
	}

	// Test with help command
	err = app.Run([]string{"scribe", "help"})
	if err != nil {
		t.Errorf("Expected no error for help command, got: %v", err)
	}

	// Test with version command
	err = app.Run([]string{"scribe", "version"})
	if err != nil {
		t.Errorf("Expected no error for version command, got: %v", err)
	}

	// Test with unknown command
	err = app.Run([]string{"scribe", "unknown-command"})
	if err == nil {
		t.Error("Expected error for unknown command, got nil")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("Expected error message to contain 'unknown command', got: %v", err)
	}
}

func TestShowHelp(t *testing.T) {
	// This is a visual test, so we just verify it doesn't panic
	app := NewApp()
	app.showHelp()
	// If we get here, it didn't panic
}

func TestCreateNewSite(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test site path
	sitePath := filepath.Join(tempDir, "test-site")

	// Create app
	app := NewApp()

	// Mock stdin to provide interactive input
	oldStdin := os.Stdin
	// Create a pipe to simulate user input
	r, w, _ := os.Pipe()
	os.Stdin = r
	
	// Write mock input for interactive prompts
	// Empty selection for template (default to "none")
	// "n" for git init prompt
	go func() {
		defer w.Close()
		w.Write([]byte("\nn\n"))
	}()
	
	// Create new site
	err = app.createNewSite(sitePath)
	
	// Restore stdin
	os.Stdin = oldStdin
	
	if err != nil {
		t.Fatalf("Failed to create new site: %v", err)
	}

	// Check that directories were created
	expectedDirs := []string{
		"content",
		"content/posts",
		"layouts",
		"static",
		"themes/default",
		"themes/default/layouts",
		"themes/default/static",
	}

	for _, dir := range expectedDirs {
		path := filepath.Join(sitePath, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected directory '%s' to be created", path)
		}
	}

	// Check that config file was created
	configPath := filepath.Join(sitePath, "config.jsonc")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Expected config file to be created at '%s'", configPath)
	}

	// Check that sample content was created
	samplePostPath := filepath.Join(sitePath, "content", "posts", "welcome.md")
	if _, err := os.Stat(samplePostPath); os.IsNotExist(err) {
		t.Errorf("Expected sample post to be created at '%s'", samplePostPath)
	}

	samplePagePath := filepath.Join(sitePath, "content", "about.md")
	if _, err := os.Stat(samplePagePath); os.IsNotExist(err) {
		t.Errorf("Expected sample page to be created at '%s'", samplePagePath)
	}

	// Check that templates were created
	expectedTemplates := []string{
		"base.html",
		"single.html",
		"list.html",
		"home.html",
		"page.html",
	}

	for _, template := range expectedTemplates {
		path := filepath.Join(sitePath, "themes", "default", "layouts", template)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected template '%s' to be created", path)
		}
	}

	// Check that CSS file was created
	cssPath := filepath.Join(sitePath, "themes", "default", "static", "css", "style.css")
	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("Expected CSS file to be created at '%s'", cssPath)
	}
}

func TestSampleContentAndTemplates(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "content-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create necessary subdirectories
	err = os.MkdirAll(filepath.Join(tempDir, "content", "posts"), 0755)
	if err != nil {
		t.Fatalf("Failed to create content/posts directory: %v", err)
	}

	// Create app
	app := NewApp()

	// Test createSampleContent
	err = app.createSampleContent(tempDir)
	if err != nil {
		t.Fatalf("Failed to create sample content: %v", err)
	}

	// Check that sample content was created
	postPath := filepath.Join(tempDir, "content", "posts", "welcome.md")
	if _, err := os.Stat(postPath); os.IsNotExist(err) {
		t.Errorf("Expected post file to be created at '%s'", postPath)
	}

	pagePath := filepath.Join(tempDir, "content", "about.md")
	if _, err := os.Stat(pagePath); os.IsNotExist(err) {
		t.Errorf("Expected page file to be created at '%s'", pagePath)
	}

	// Create necessary theme directories
	err = os.MkdirAll(filepath.Join(tempDir, "themes", "default", "layouts"), 0755)
	if err != nil {
		t.Fatalf("Failed to create themes directory: %v", err)
	}

	// Test createDefaultTemplates
	err = app.createDefaultTemplates(tempDir)
	if err != nil {
		t.Fatalf("Failed to create default templates: %v", err)
	}

	// Check that templates were created
	templateDir := filepath.Join(tempDir, "themes", "default", "layouts")
	templates := []string{"base.html", "single.html", "list.html", "home.html", "page.html"}
	for _, template := range templates {
		path := filepath.Join(templateDir, template)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected template file to be created at '%s'", path)
		}
	}

	// Check that CSS file was created
	cssPath := filepath.Join(tempDir, "themes", "default", "static", "css", "style.css")
	if _, err := os.Stat(cssPath); os.IsNotExist(err) {
		t.Errorf("Expected CSS file to be created at '%s'", cssPath)
	}
}