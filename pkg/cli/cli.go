package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dikaio/scribe/internal/build"
	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/console"
	"github.com/dikaio/scribe/internal/content"
	"github.com/dikaio/scribe/internal/server"
	"github.com/dikaio/scribe/internal/templates"
)

// Version information set by build flags
var (
	// Version is the semantic version of the application
	Version = "v0.2.0"
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

	// Console command
	a.Commands["console"] = Command{
		Name:        "console",
		Description: "Start the console",
		Action:      a.cmdConsole,
	}

	// New site command
	a.Commands["new"] = Command{
		Name:        "new",
		Description: "Create a new site, post, or page",
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
	sitePath, cfg, err := a.getSitePathAndConfig(args, "Starting development server for")
	if err != nil {
		return err
	}

	// Initialize the server (default port: 8080)
	port := 8080
	server := server.NewServer(cfg, port)

	// Start the server
	err = server.Start(sitePath)
	if err != nil {
		return fmt.Errorf("server failed: %w", err)
	}

	return nil
}

// cmdConsole implements the console command, which starts the web management interface.
// It takes an optional path argument (or uses the current directory if not provided).
func (a *App) cmdConsole(args []string) error {
	sitePath, cfg, err := a.getSitePathAndConfig(args, "Starting console for")
	if err != nil {
		return err
	}

	// Initialize console (default port: 8090)
	port := 8090
	console := console.NewConsole(cfg, sitePath, port)

	// Start the console
	return console.Start()
}

// cmdNew implements commands for creating new resources (site, post, page, etc.)
func (a *App) cmdNew(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing resource type")
	}

	resType := args[0]
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
	// Use current directory if no name provided
	sitePath := "."
	if name != "" {
		sitePath = name
		// Create site directory
		if err := os.MkdirAll(sitePath, 0755); err != nil {
			return fmt.Errorf("failed to create site directory: %w", err)
		}
	}

	fmt.Printf("Creating new Scribe site in '%s'...\n", sitePath)

	// Create directories
	dirs := []string{
		"content",
		"content/posts",
		"layouts",
		"static",
		"themes/default",
		"themes/default/layouts",
		"themes/default/static",
	}

	for _, dir := range dirs {
		path := filepath.Join(sitePath, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}
	}

	// Create default config file
	cfg := config.DefaultConfig()
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

	fmt.Println("Site created successfully!")
	fmt.Println("Run 'scribe serve' to start the development server.")
	return nil
}

// createNewPost creates a new blog post
func (a *App) createNewPost(title string) error {
	// Prompt for title if not provided
	if title == "" {
		fmt.Print("Enter post title: ")
		fmt.Scanln(&title)
		if title == "" {
			return fmt.Errorf("post title is required")
		}
	}

	// Create content creator
	creator := content.NewCreator(".")

	// Create post
	tags := []string{"uncategorized"}
	filePath, err := creator.CreateContent(content.PostType, title, "", tags, false)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	fmt.Printf("Post created successfully: %s\n", filePath)
	return nil
}

// createNewPage creates a new static page
func (a *App) createNewPage(title string) error {
	// Prompt for title if not provided
	if title == "" {
		fmt.Print("Enter page title: ")
		fmt.Scanln(&title)
		if title == "" {
			return fmt.Errorf("page title is required")
		}
	}

	// Create content creator
	creator := content.NewCreator(".")

	// Create page
	filePath, err := creator.CreateContent(content.PageType, title, "", nil, false)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	fmt.Printf("Page created successfully: %s\n", filePath)
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
func (a *App) createDefaultTemplates(sitePath string) error {
	// Define template paths
	templatePaths := map[string]string{
		filepath.Join(sitePath, "themes", "default", "layouts", "base.html"):   templates.BaseTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "single.html"): templates.SingleTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "list.html"):   templates.ListTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "home.html"):   templates.HomeTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "page.html"):   templates.PageTemplate,
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
	return os.WriteFile(cssFilePath, []byte(templates.StyleCSS), 0644)
}

// Test command removed as it was unimplemented
