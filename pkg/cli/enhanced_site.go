package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/ui"
)

// createSiteEnhanced is a simplified version of createNewSite that just
// asks for the site name and creates a default site with git initialization
func (a *App) createSiteEnhanced(initialName string) error {
	ui.Header("Scribe")

	// Prompt for site name
	sitePath := "."
	siteName := initialName
	
	if siteName == "" {
		siteName = ui.Prompt("Site Name (leave empty to use current directory)", "")
	}
	
	if siteName != "" {
		sitePath = siteName
		// Create site directory
		if err := os.MkdirAll(sitePath, 0755); err != nil {
			return fmt.Errorf("failed to create site directory: %w", err)
		}
	}

	// Display banner
	ui.Info(fmt.Sprintf("Creating new Scribe site in '%s'...", sitePath))
	ui.Divider()

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
	cfg.Title = "My Scribe Site"
	cfg.Description = "A site created with Scribe"
	
	if err := cfg.Save(sitePath); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Create sample content
	if err := a.createSampleContent(sitePath); err != nil {
		return fmt.Errorf("failed to create sample content: %w", err)
	}

	// Create default templates
	if err := a.createDefaultTemplates(sitePath); err != nil {
		return fmt.Errorf("failed to create default templates: %w", err)
	}
	
	// Initialize git repository automatically
	ui.Info("Initializing git repository...")
	gitCmd := exec.Command("git", "init", sitePath)
	err := gitCmd.Run()
	if err != nil {
		ui.Warning(fmt.Sprintf("Failed to initialize git repository: %v", err))
	} else {
		// Create .gitignore
		gitignorePath := filepath.Join(sitePath, ".gitignore")
		gitignoreContent := "# Output directory\npublic/\n\n# IDE files\n.idea/\n.vscode/\n\n# System files\n.DS_Store\nThumbs.db\n"
		if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
			ui.Warning(fmt.Sprintf("Failed to create .gitignore file: %v", err))
		}
	}

	// Success message
	ui.Divider()
	ui.Success("Site created successfully!")
	
	// Navigation instructions
	if siteName != "" {
		ui.Info(fmt.Sprintf("Next steps:"))
		ui.Info(fmt.Sprintf("  cd %s", siteName))
		ui.Info(fmt.Sprintf("  scribe serve"))
		ui.Info(fmt.Sprintf("Then view your site at http://localhost:8080"))
	} else {
		ui.Info(fmt.Sprintf("Next steps:"))
		ui.Info(fmt.Sprintf("  scribe serve"))
		ui.Info(fmt.Sprintf("Then view your site at http://localhost:8080"))
	}
	
	return nil
}