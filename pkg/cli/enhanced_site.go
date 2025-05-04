package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/ui"
)

// createSiteEnhanced is an enhanced version of createNewSite that uses
// the improved UI components for a better interactive experience
func (a *App) createSiteEnhanced(initialName string) error {
	ui.Header("Create New Site")

	// Prompt for site name
	ui.Title("Site Setup", "Let's configure your new Scribe site")
	sitePath := "."
	siteName := ui.Prompt("Site Name (leave empty to use current directory)", initialName)
	
	if siteName != "" {
		sitePath = siteName
		// Create site directory
		if err := os.MkdirAll(sitePath, 0755); err != nil {
			return fmt.Errorf("failed to create site directory: %w", err)
		}
	}

	// Template selection
	templateOptions := []ui.Option{
		{Label: "None (minimal)", Value: "none"},
		{Label: "Blog", Value: "blog"},
		{Label: "Documentation", Value: "docs"},
		{Label: "Kitchen sink (all features)", Value: "kitchen-sink"},
	}

	template := ui.SelectOption(
		"Site Template",
		"Select a template for your new site:",
		templateOptions,
		0, // Default to first option (None)
	)

	// Display banner
	ui.Info(fmt.Sprintf("Creating new Scribe site in '%s' with '%s' template...", sitePath, template))
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
	
	// Customize config based on template
	switch template {
	case "blog":
		cfg.Title = "My Blog"
		cfg.Description = "A blog created with Scribe"
	case "docs":
		cfg.Title = "Documentation"
		cfg.Description = "Documentation site built with Scribe"
	case "kitchen-sink":
		cfg.Title = "Scribe Demo Site"
		cfg.Description = "Showcasing all Scribe features"
	}
	
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
	
	// Initialize git repository if requested
	initGit := ui.ConfirmYesNo("Initialize git repository?", false)
	
	if initGit {
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
	}

	// Success message
	ui.Divider()
	ui.Success("Site created successfully!")
	ui.Info("Run 'scribe serve' to start the development server.")
	
	return nil
}