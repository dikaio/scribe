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

// createSiteEnhanced creates a new site with interactive prompts
// for site name and theme selection
func (a *App) createSiteEnhanced(initialName string) error {
	ui.Header("Create New Site")

	// Prompt for site name
	siteName := initialName
	
	if siteName == "" {
		siteName = ui.PromptWithValidation("What is your site named?", "", ui.Required("Site name is required"))
	}
	
	// Create site directory
	sitePath := siteName
	if err := os.MkdirAll(sitePath, 0755); err != nil {
		return fmt.Errorf("failed to create site directory: %w", err)
	}

	// Prompt for Tailwind usage
	useTailwind := ui.ConfirmYesNo("Would you like to use Tailwind?", false)

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
	cfg.Title = siteName
	cfg.Description = fmt.Sprintf("A %s site created with Scribe", siteName)
	
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
	
	// Create .gitignore and initialize git repository
	a.createGitignore(sitePath, useTailwind)
	a.initGitRepo(sitePath)

	// Automatically run npm install if using Tailwind
	if useTailwind {
		a.installNpmDependencies()
	}
	
	// Success message with specific instruction
	ui.Divider()
	fmt.Printf("Your new site \"%s\" is ready, to start: cd \"%s\" && scribe serve\n", siteName, siteName)
	
	return nil
}

// createGitignore creates a .gitignore file appropriate for the site type
func (a *App) createGitignore(sitePath string, useTailwind bool) {
	gitignorePath := filepath.Join(sitePath, ".gitignore")
	
	// Select the appropriate content based on site type
	var gitignoreContent string
	if useTailwind {
		gitignoreContent = templates.TailwindGitignore
	} else {
		gitignoreContent = "# Output directory\npublic/\n\n# IDE files\n.idea/\n.vscode/\n\n# System files\n.DS_Store\nThumbs.db\n"
	}
	
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		ui.Warning(fmt.Sprintf("Failed to create .gitignore file: %v", err))
	}
}

// initGitRepo initializes a git repository
func (a *App) initGitRepo(sitePath string) {
	ui.Info("Initializing git repository...")
	gitCmd := exec.Command("git", "init", sitePath)
	if err := gitCmd.Run(); err != nil {
		ui.Warning(fmt.Sprintf("Failed to initialize git repository: %v", err))
	}
}

// installNpmDependencies installs NPM dependencies for Tailwind
func (a *App) installNpmDependencies() {
	ui.Info("Installing Tailwind CSS dependencies...")
	
	// Run npm install
	npmCmd := exec.Command("npm", "install")
	npmCmd.Stdout = os.Stdout
	npmCmd.Stderr = os.Stderr
	
	if err := npmCmd.Run(); err != nil {
		ui.Warning(fmt.Sprintf("Failed to install dependencies: %v", err))
		ui.Info("Please run 'npm install' manually to complete setup.")
	} else {
		ui.Success("Dependencies installed successfully!")
	}
}

// showNextSteps is removed since we now display a simple success message directly
// in the createSiteEnhanced function