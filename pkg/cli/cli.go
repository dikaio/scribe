package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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
	Version = "v0.4.13"
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

	// Run command (serve + Tailwind watch)
	a.Commands["run"] = Command{
		Name:        "run",
		Description: "Run development server and file watchers concurrently",
		Action:      a.cmdRun,
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
	fmt.Printf("  %s run                  Run development server and file watchers (for Tailwind CSS sites)\n", a.Name)
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

// cmdRun implements the run command, which starts both the Tailwind CSS watcher and the development server.
// It automatically detects if the site uses Tailwind CSS.
func (a *App) cmdRun(args []string) error {
	// Get site path and config
	sitePath, cfg, err := a.getSitePathAndConfig(args, "")
	if err != nil {
		return err
	}

	// Check if the site uses Tailwind CSS by looking for package.json and src/input.css
	packageJsonPath := filepath.Join(sitePath, "package.json")
	inputCssPath := filepath.Join(sitePath, "src", "input.css")
	useTailwind := false
	
	// If both package.json and input.css exist, assume it's a Tailwind site
	if _, err := os.Stat(packageJsonPath); err == nil {
		if _, err := os.Stat(inputCssPath); err == nil {
			useTailwind = true
		}
	}

	if !useTailwind {
		// If no Tailwind CSS is used, just start the server
		ui.Info("No Tailwind CSS configuration detected. Starting development server only.")
		return a.cmdServe(args)
	}

	// Verify npm is installed
	npmCmd := exec.Command("npm", "--version")
	npmErr := npmCmd.Run()
	if npmErr != nil {
		return fmt.Errorf("npm not found. Please install Node.js and npm to use the run command with Tailwind CSS: %w", npmErr)
	}

	// Start Tailwind CSS watcher in the background
	ui.Info("Starting Tailwind CSS watcher and development server...")
	
	// Create a channel to catch signals
	done := make(chan struct{})

	// Start Tailwind CSS watcher in a goroutine
	go func() {
		defer close(done)
		
		tailwindCmd := exec.Command("npm", "run", "dev")
		tailwindCmd.Stdout = os.Stdout
		tailwindCmd.Stderr = os.Stderr
		tailwindCmd.Dir = sitePath
		
		err := tailwindCmd.Start()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to start Tailwind CSS watcher: %v", err))
			return
		}
		
		// Wait for the process to finish (which won't happen unless there's an error)
		tailwindCmd.Wait()
	}()

	// Give Tailwind a moment to start
	time.Sleep(1 * time.Second)

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

// createDefaultTemplates moved to createDefaultTemplates.go

// Test command removed as it was unimplemented
