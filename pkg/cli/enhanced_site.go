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
		{Label: "Default CSS", Value: "default"},
		{Label: "Tailwind CSS", Value: "tailwind"},
	}
	cssChoice := ui.SelectOption("Select Theme", "Choose a theme for your site", cssOptions, 0)
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
	
	// Automatically change directory if site name is not empty
	// This will affect the current process but not the parent shell
	// (which is why we need to show instructions to the user)
	if siteName != "" {
		os.Chdir(sitePath)
	}
	
	// Automatically run npm install if using Tailwind
	if useTailwind {
		ui.Info("Installing Tailwind CSS dependencies...")
		
		// Run npm install
		npmCmd := exec.Command("npm", "install")
		npmCmd.Stdout = os.Stdout
		npmCmd.Stderr = os.Stderr
		
		err := npmCmd.Run()
		if err != nil {
			ui.Warning(fmt.Sprintf("Failed to install dependencies: %v", err))
			ui.Info("Please run 'npm install' manually to complete setup.")
		} else {
			ui.Success("Dependencies installed successfully!")
		}
	}
	
	// Navigation instructions with framework-specific guidance
	ui.Info("Next steps:")
	
	if siteName != "" {
		ui.Info(fmt.Sprintf("  cd %s  (if not already in that directory)", siteName))
	}
	
	if useTailwind {
		ui.Info("  scribe run     (run both Tailwind compiler and development server)")
		ui.Info("  - OR -")
		ui.Info("  npm run dev    (in one terminal)")
		ui.Info("  scribe serve   (in another terminal)")
	} else {
		ui.Info("  scribe serve   (start development server)")
	}
	
	ui.Info("Then view your site at http://localhost:8080")
	
	// Show additional Tailwind-specific instructions
	if useTailwind {
		ui.Divider()
		ui.Info("Tailwind CSS Setup:")
		ui.Info("1. Node.js is required to use Tailwind CSS")
		ui.Info("2. Run 'scribe run' to start both the Tailwind compiler and development server")
		ui.Info("3. For more detailed instructions, see the README.md file")
	}
	
	return nil
}