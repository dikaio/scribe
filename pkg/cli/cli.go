package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dikaio/scribe/internal/build"
	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
	"github.com/dikaio/scribe/internal/server"
	"github.com/dikaio/scribe/internal/templates"
	"github.com/dikaio/scribe/internal/ui"
)

// Version information set by build flags
var (
	// Version is the semantic version of the application
	Version = "v0.4.10"
	// Commit is the git commit SHA at build time
	Commit = "none" 
	// Date is the build date
	Date = "unknown"
)

// App represents the CLI application
type App struct {
	Name     string
	Version  string
	Commands map[string]Command
}

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Action      func([]string) error
}

// NewApp creates a new CLI application
func NewApp() *App {
	app := &App{
		Name:     "scribe",
		Version:  Version,
		Commands: make(map[string]Command),
	}

	// Register commands
	app.registerCommands()

	return app
}

// registerCommands registers all available commands
func (a *App) registerCommands() {
	// Build command
	a.Commands["build"] = Command{
		Name:        "build",
		Description: "Build the site",
		Action:      a.cmdBuild,
	}

	// Serve command
	a.Commands["serve"] = Command{
		Name:        "serve",
		Description: "Start a development server with live reload",
		Action:      a.cmdServe,
	}


	// New site command
	a.Commands["new"] = Command{
		Name:        "new",
		Description: "Create a new site, post, or page interactively",
		Action:      a.cmdNew,
	}

	// NOTE: Test command removed as it was unimplemented
}

// Run executes the CLI application
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		a.showHelp()
		return nil
	}

	cmdName := args[1]

	if cmdName == "help" || cmdName == "-h" || cmdName == "--help" {
		a.showHelp()
		return nil
	}

	if cmdName == "version" || cmdName == "-v" || cmdName == "--version" {
		fmt.Printf("%s version %s\n", a.Name, a.Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", Date)
		return nil
	}

	cmd, exists := a.Commands[cmdName]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	// Execute the command
	return cmd.Action(args[2:])
}

// showHelp displays help information
func (a *App) showHelp() {
	fmt.Printf("%s - A lightweight static site generator\n\n", a.Name)
	fmt.Println("Usage:")
	fmt.Printf("  %s [command] [arguments]\n\n", a.Name)
	fmt.Println("Available commands:")

	for _, cmd := range a.Commands {
		fmt.Printf("  %-10s %s\n", cmd.Name, cmd.Description)
	}
	
	fmt.Println("\nExamples:")
	fmt.Printf("  %s new                  Create a new site in the current directory with interactive prompts\n", a.Name)
	fmt.Printf("  %s new my-site          Create a new site in 'my-site' directory with interactive prompts\n", a.Name)
	fmt.Printf("  %s serve                Start development server for the current directory\n", a.Name)
	fmt.Printf("  %s build                Build the site in the current directory\n", a.Name)

	fmt.Println("\nUse 'scribe help [command]' for more information about a command.")
}

// Command implementations

// getSitePathAndConfig is a helper function to determine the site path from args
// and load the site configuration. It standardizes this common operation across commands.
func (a *App) getSitePathAndConfig(args []string, action string) (string, config.Config, error) {
	// Determine site path
	sitePath := "."
	if len(args) > 0 {
		sitePath = args[0]
	}

	if action != "" {
		fmt.Printf("%s site from '%s'...\n", action, sitePath)
	}

	// Load the site configuration
	cfg, err := config.LoadConfig(sitePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", cfg, fmt.Errorf("config file not found in '%s', make sure this is a valid Scribe site directory", sitePath)
		}
		return "", cfg, fmt.Errorf("failed to load configuration: %w", err)
	}

	return sitePath, cfg, nil
}

// cmdBuild implements the build command, which generates the static site.
// It takes an optional path argument (or uses the current directory if not provided).
func (a *App) cmdBuild(args []string) error {
	sitePath, cfg, err := a.getSitePathAndConfig(args, "Building")
	if err != nil {
		return err
	}

	// Initialize the builder
	builder := build.NewBuilder(cfg)

	// Build the site
	err = builder.Build(sitePath)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	fmt.Printf("Site built successfully! Output directory: '%s'\n", cfg.OutputDir)
	return nil
}

// cmdServe implements the serve command, which starts a development server with live reload.
// It takes an optional path argument (or uses the current directory if not provided).
func (a *App) cmdServe(args []string) error {
	// Don't show the "Starting development server..." message
	sitePath, cfg, err := a.getSitePathAndConfig(args, "")
	if err != nil {
		return err
	}

	// Initialize the server (default port: 8080)
	port := 8080
	server := server.NewServer(cfg, port, true) // true = quiet mode

	// Start the server
	err = server.Start(sitePath)
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}


// cmdNew implements commands for creating new resources (site, post, page, etc.)
func (a *App) cmdNew(args []string) error {
	// If no arguments, assume we're creating a new site
	if len(args) < 1 {
		return a.createNewSite("")
	}

	// If first arg starts with a letter and not "site", "post", or "page", 
	// treat it as a site name
	firstArg := args[0]
	if firstArg != "site" && firstArg != "post" && firstArg != "page" {
		return a.createNewSite(firstArg)
	}

	// Otherwise, handle the classic way
	resType := firstArg
	name := ""
	if len(args) > 1 {
		name = args[1]
	}

	switch resType {
	case "site":
		return a.createNewSite(name)
	case "post":
		return a.createNewPost(name)
	case "page":
		return a.createNewPage(name)
	default:
		return fmt.Errorf("unknown resource type: %s", resType)
	}
}

// createNewSite scaffolds a new site with default structure and templates
func (a *App) createNewSite(name string) error {
	// Use the enhanced version with improved UI
	return a.createSiteEnhanced(name)
}

// createNewPost creates a new blog post
func (a *App) createNewPost(initialTitle string) error {
	ui.Header("Create New Post")

	// Prompt for post title
	title := initialTitle
	if title == "" {
		title = ui.PromptWithValidation("Post Title", "", ui.Required("Post title"))
	}

	// Prompt for description (optional)
	description := ui.Prompt("Description (optional)", "")

	// Prompt for tags
	defaultTags := []string{"uncategorized"}
	tags := ui.PromptTags("Tags (comma-separated)", defaultTags)

	// Prompt for draft status
	draft := ui.ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create post
	filePath, err := creator.CreateContent(content.PostType, title, description, tags, draft)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	ui.Success(fmt.Sprintf("Post created successfully: %s", filePath))
	return nil
}

// createNewPage creates a new static page
func (a *App) createNewPage(initialTitle string) error {
	ui.Header("Create New Page")

	// Prompt for page title
	title := initialTitle
	if title == "" {
		title = ui.PromptWithValidation("Page Title", "", ui.Required("Page title"))
	}

	// Prompt for description (optional)
	description := ui.Prompt("Description (optional)", "")

	// Prompt for draft status
	draft := ui.ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create page
	filePath, err := creator.CreateContent(content.PageType, title, description, nil, draft)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	ui.Success(fmt.Sprintf("Page created successfully: %s", filePath))
	return nil
}

// createSampleContent creates sample content files for a new site
func (a *App) createSampleContent(sitePath string) error {
	// Write sample post
	postPath := filepath.Join(sitePath, "content", "posts", "welcome.md")
	if err := os.WriteFile(postPath, []byte(templates.SamplePost), 0644); err != nil {
		return err
	}

	// Write sample page
	pagePath := filepath.Join(sitePath, "content", "about.md")
	return os.WriteFile(pagePath, []byte(templates.SamplePage), 0644)
}

// createDefaultTemplates creates default templates for a new site
func (a *App) createDefaultTemplates(sitePath string, useTailwind bool) error {
	var templatePaths map[string]string
	
	// Define template paths based on the selected theme
	if useTailwind {
		// Use Tailwind CSS templates
		templatePaths = map[string]string{
			filepath.Join(sitePath, "themes", "default", "layouts", "base.html"):   templates.TailwindBaseTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "single.html"): templates.TailwindSingleTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "list.html"):   templates.TailwindListTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "home.html"):   templates.TailwindHomeTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "page.html"):   templates.TailwindPageTemplate,
		}
	} else {
		// Use default templates
		templatePaths = map[string]string{
			filepath.Join(sitePath, "themes", "default", "layouts", "base.html"):   templates.BaseTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "single.html"): templates.SingleTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "list.html"):   templates.ListTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "home.html"):   templates.HomeTemplate,
			filepath.Join(sitePath, "themes", "default", "layouts", "page.html"):   templates.PageTemplate,
		}
	}

	// Write template files
	for path, content := range templatePaths {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create template '%s': %w", path, err)
		}
	}

	// Create static directory with CSS
	cssPath := filepath.Join(sitePath, "themes", "default", "static", "css")
	if err := os.MkdirAll(cssPath, 0755); err != nil {
		return fmt.Errorf("failed to create CSS directory: %w", err)
	}

	// Write CSS file
	cssFilePath := filepath.Join(cssPath, "style.css")
	
	if useTailwind {
		// For Tailwind, create the necessary files
		
		// Create src directory for input CSS
		srcPath := filepath.Join(sitePath, "src")
		if err := os.MkdirAll(srcPath, 0755); err != nil {
			return fmt.Errorf("failed to create src directory: %w", err)
		}
		
		// Write Tailwind input.css
		inputFilePath := filepath.Join(srcPath, "input.css")
		if err := os.WriteFile(inputFilePath, []byte(templates.TailwindInputCSS), 0644); err != nil {
			return fmt.Errorf("failed to create Tailwind input.css: %w", err)
		}
		
		// Create empty output CSS file (will be populated by Tailwind CLI)
		if err := os.WriteFile(cssFilePath, []byte("/* Tailwind CSS styles will be generated here */"), 0644); err != nil {
			return fmt.Errorf("failed to create CSS file: %w", err)
		}
		
		// Write tailwind.config.js
		configFilePath := filepath.Join(sitePath, "tailwind.config.js")
		if err := os.WriteFile(configFilePath, []byte(templates.TailwindCSSConfig), 0644); err != nil {
			return fmt.Errorf("failed to create tailwind.config.js: %w", err)
		}
		
		// Write package.json
		packageFilePath := filepath.Join(sitePath, "package.json")
		if err := os.WriteFile(packageFilePath, []byte(templates.TailwindPackageJSON), 0644); err != nil {
			return fmt.Errorf("failed to create package.json: %w", err)
		}
		
		// Create README.md with Tailwind instructions
		readmeFilePath := filepath.Join(sitePath, "README.md")
		if err := os.WriteFile(readmeFilePath, []byte(templates.TailwindREADME), 0644); err != nil {
			return fmt.Errorf("failed to create README.md: %w", err)
		}
		
		return nil
	}
	
	// Write default CSS file
	return os.WriteFile(cssFilePath, []byte(templates.StyleCSS), 0644)
}

// Test command removed as it was unimplemented
