package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dikaio/scribe/internal/build"
	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
	"github.com/dikaio/scribe/internal/server"
	"github.com/dikaio/scribe/internal/templates"
)

// Version information set by build flags
var (
	// Version is the semantic version of the application
	Version = "v0.8.0"
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
	// Serve command
	a.Commands["serve"] = Command{
		Name:        "serve",
		Description: "Start a development server with live reload",
		Action:      a.cmdServe,
	}

	// Build command
	a.Commands["build"] = Command{
		Name:        "build",
		Description: "Build a static site",
		Action:      a.cmdBuild,
	}

	// New command for site and page creation
	a.Commands["new"] = Command{
		Name:        "new",
		Description: "Create a new site or page",
		Action:      a.cmdNew,
	}
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
	fmt.Printf("  %s new site             Create a new site with interactive prompts\n", a.Name)
	fmt.Printf("  %s new page [path]      Create a new page at the specified path\n", a.Name)
	fmt.Printf("  %s serve                Start development server for the current directory\n", a.Name)
	fmt.Printf("  %s build                Build the static site in the current directory\n", a.Name)

	fmt.Println("\nUse 'scribe --help' to display this help information.")
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
	
	// Enable template caching for production builds
	// This improves performance by not re-parsing templates unnecessarily
	builder.SetDevMode(false)

	// Build the site
	start := time.Now()
	err = builder.Build(sitePath)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	buildTime := time.Since(start)
	fmt.Printf("Site built successfully in %v! Output directory: '%s'\n", buildTime, cfg.OutputDir)
	return nil
}

// cmdServe implements the serve command, which starts a development server with live reload.
func (a *App) cmdServe(args []string) error {
	// Get site path and config
	sitePath, cfg, err := a.getSitePathAndConfig(args, "")
	if err != nil {
		return err
	}

	Info("Starting development server...")

	// Initialize the server (default port: 8080)
	port := 8080
	server := server.NewServer(cfg, port, false) // false = not quiet mode

	// Start the server
	err = server.Start(sitePath)
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// cmdRun is now an alias for cmdServe
func (a *App) cmdRun(args []string) error {
	return a.cmdServe(args)
}


// cmdNew implements commands for creating new sites or pages
func (a *App) cmdNew(args []string) error {
	// If no arguments, show help for new command
	if len(args) < 1 {
		fmt.Println("Usage:")
		fmt.Println("  scribe new site           Create a new site with interactive prompts")
		fmt.Println("  scribe new page [path]    Create a new page at the specified path")
		return nil
	}

	// Handle resource types
	resType := args[0]
	
	switch resType {
	case "site":
		// For sites, we need a name parameter
		if len(args) < 2 {
			return fmt.Errorf("site command requires a name argument")
		}
		siteName := args[1]
		return a.createNewSite(siteName)
		
	case "page":
		// For pages, second arg is the path
		if len(args) < 2 {
			return fmt.Errorf("page command requires a path argument")
		}
		path := args[1]
		
		// Title is optional, will be prompted if not provided
		title := ""
		if len(args) > 2 {
			title = args[2]
		}
		
		return a.createNewPage(title, path)
	
	default:
		return fmt.Errorf("unknown resource type: %s. Use 'site' or 'page'", resType)
	}
}

// createNewSite scaffolds a new site with default structure and templates
func (a *App) createNewSite(name string) error {
	// Use the enhanced version with improved UI
	return a.createSiteEnhanced(name)
}

// createNewPost creates a new blog post
func (a *App) createNewPost(initialTitle string, customPath string) error {
	Header("Create New Post")

	// Prompt for post title
	title := initialTitle
	if title == "" {
		title = PromptWithValidation("Post Title", "", Required("Post title"))
	}

	// Prompt for description (optional)
	description := Prompt("Description (optional)", "")

	// Prompt for tags
	defaultTags := []string{"uncategorized"}
	tags := PromptTags("Tags (comma-separated)", defaultTags)

	// Prompt for draft status
	draft := ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create post
	filePath, err := creator.CreateContent(content.PostType, title, description, tags, draft, customPath)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	Success(fmt.Sprintf("Post created successfully: %s", filePath))
	return nil
}

// createNewPage creates a new static page
func (a *App) createNewPage(initialTitle string, customPath string) error {
	Header("Create New Page")

	// Prompt for page title
	title := initialTitle
	if title == "" {
		title = PromptWithValidation("Page Title", "", Required("Page title"))
	}

	// Prompt for description (optional)
	description := Prompt("Description (optional)", "")

	// Prompt for draft status
	draft := ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create page
	filePath, err := creator.CreateContent(content.PageType, title, description, nil, draft, customPath)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	Success(fmt.Sprintf("Page created successfully: %s", filePath))
	return nil
}

// createNewContent creates generic content at a specific path
func (a *App) createNewContent(customPath string, initialTitle string) error {
	Header("Create New Content")
	
	// Determine if we should use post or page format
	contentType := content.PageType
	if strings.Contains(customPath, "posts/") || strings.HasPrefix(customPath, "posts") {
		contentType = content.PostType
	}

	// Prompt for title
	title := initialTitle
	if title == "" {
		if contentType == content.PostType {
			title = PromptWithValidation("Post Title", "", Required("Post title"))
		} else {
			title = PromptWithValidation("Content Title", "", Required("Content title"))
		}
	}

	// Prompt for description (optional)
	description := Prompt("Description (optional)", "")

	// Prompt for tags if it's a post
	var tags []string
	if contentType == content.PostType {
		defaultTags := []string{"uncategorized"}
		tags = PromptTags("Tags (comma-separated)", defaultTags)
	}

	// Prompt for draft status
	draft := ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create content
	filePath, err := creator.CreateContent(contentType, title, description, tags, draft, customPath)
	if err != nil {
		return fmt.Errorf("failed to create content: %w", err)
	}

	Success(fmt.Sprintf("Content created successfully: %s", filePath))
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

// createDefaultTemplates moved to createDefaultTemplates.go

// Test command removed as it was unimplemented
