package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/templates"
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
		siteName = ui.Prompt("Enter site name:", "")
	}
	
	if siteName != "" {
		sitePath = siteName
		// Create site directory
		if err := os.MkdirAll(sitePath, 0755); err != nil {
			return fmt.Errorf("failed to create site directory: %w", err)
		}
	}

	// Prompt for CSS framework
	cssOptions := []ui.Option{
		{Label: "Default CSS (Simple, no dependencies)", Value: "default"},
		{Label: "Tailwind CSS (Utility-first CSS framework)", Value: "tailwind"},
	}
	cssChoice := ui.SelectOption("CSS Framework", "Choose a CSS framework for your site", cssOptions, 0)
	useTailwind := cssChoice == "tailwind"

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

	// Create templates based on CSS choice
	if err := a.createDefaultTemplates(sitePath, useTailwind); err != nil {
		return fmt.Errorf("failed to create default templates: %w", err)
	}
	
	// Initialize git repository automatically
	ui.Info("Initializing git repository...")
	gitCmd := exec.Command("git", "init", sitePath)
	err := gitCmd.Run()
	if err != nil {
		ui.Warning(fmt.Sprintf("Failed to initialize git repository: %v", err))
	} else {
		// Create .gitignore with appropriate content
		gitignorePath := filepath.Join(sitePath, ".gitignore")
		gitignoreContent := "# Output directory\npublic/\n\n# IDE files\n.idea/\n.vscode/\n\n# System files\n.DS_Store\nThumbs.db\n"
		
		// Add Node.js entries if using Tailwind
		if useTailwind {
			gitignoreContent = templates.TailwindGitignore
		}
		
		if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
			ui.Warning(fmt.Sprintf("Failed to create .gitignore file: %v", err))
		}
	}

	// Success message
	ui.Divider()
	ui.Success("Site created successfully!")
	
	// Navigation instructions with framework-specific guidance
	if siteName != "" {
		ui.Info("Next steps:")
		ui.Info(fmt.Sprintf("  cd %s", siteName))
		
		if useTailwind {
			ui.Info("  npm install")
			ui.Info("  npm run dev")
			ui.Info("  # In another terminal:")
			ui.Info("  scribe serve")
		} else {
			ui.Info("  scribe serve")
		}
		
		ui.Info("Then view your site at http://localhost:8080")
	} else {
		ui.Info("Next steps:")
		
		if useTailwind {
			ui.Info("  npm install")
			ui.Info("  npm run dev")
			ui.Info("  # In another terminal:")
			ui.Info("  scribe serve")
		} else {
			ui.Info("  scribe serve")
		}
		
		ui.Info("Then view your site at http://localhost:8080")
	}
	
	// Show additional Tailwind-specific instructions
	if useTailwind {
		ui.Divider()
		ui.Info("Tailwind CSS Setup:")
		ui.Info("1. Node.js is required to use Tailwind CSS")
		ui.Info("2. Run 'npm install' to install Tailwind CSS dependencies")
		ui.Info("3. Run 'npm run dev' to start the Tailwind CSS compiler")
		ui.Info("4. In a separate terminal, run 'scribe serve' to start the development server")
		ui.Info("5. For more detailed instructions, see the README.md file")
	}
	
	return nil
}