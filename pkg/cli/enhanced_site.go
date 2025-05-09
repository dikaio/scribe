package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dikaio/scribe/internal/config"
)

// createSiteEnhanced creates a new site with the provided name
// using the default theme
func (a *App) createSiteEnhanced(initialName string) error {
	Header("Create New Site")

	// Use the provided site name or error if empty
	siteName := initialName
	if siteName == "" {
		return fmt.Errorf("site name is required")
	}
	
	// Create site directory
	sitePath := siteName
	if err := os.MkdirAll(sitePath, 0755); err != nil {
		return fmt.Errorf("failed to create site directory: %w", err)
	}

	// Display banner
	Info(fmt.Sprintf("Creating new Scribe site in '%s'...", sitePath))
	Divider()

	// Create directories
	dirPaths := []string{
		"content",
		"content/posts",
		"layouts",
		"static",
		"themes/default",
		"themes/default/layouts",
		"themes/default/static",
	}

	for _, dir := range dirPaths {
		path := filepath.Join(sitePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}
	}

	// Create default config file
	cfg := config.DefaultConfig()
	
	// Set default values
	cfg.Title = siteName
	cfg.Description = fmt.Sprintf("A %s site created with Scribe", siteName)
	
	if err := cfg.Save(sitePath); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Create sample content
	if err := a.createSampleContent(sitePath); err != nil {
		return fmt.Errorf("failed to create sample content: %w", err)
	}

	// Create templates
	if err := a.createDefaultTemplates(sitePath, false); err != nil {
		return fmt.Errorf("failed to create default templates: %w", err)
	}
	
	// Create .gitignore and initialize git repository
	a.createGitignore(sitePath, false)
	a.initGitRepo(sitePath)
	
	// Success message with specific instruction
	Divider()
	fmt.Printf("Your new site \"%s\" is ready, to start: cd \"%s\" && scribe serve\n", siteName, siteName)
	
	return nil
}

// createGitignore creates a .gitignore file
func (a *App) createGitignore(sitePath string, _ bool) {
	gitignorePath := filepath.Join(sitePath, ".gitignore")
	
	// Standard gitignore content
	gitignoreContent := "# Output directory\npublic/\n\n# IDE files\n.idea/\n.vscode/\n\n# System files\n.DS_Store\nThumbs.db\n"
	
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		Warning(fmt.Sprintf("Failed to create .gitignore file: %v", err))
	}
}

// initGitRepo initializes a git repository
func (a *App) initGitRepo(sitePath string) {
	Info("Initializing git repository...")
	gitCmd := exec.Command("git", "init", sitePath)
	if err := gitCmd.Run(); err != nil {
		Warning(fmt.Sprintf("Failed to initialize git repository: %v", err))
	}
}

// showNextSteps is removed since we now display a simple success message directly
// in the createSiteEnhanced function