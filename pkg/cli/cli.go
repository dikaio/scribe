package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	Version = "v0.3.0"
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

	// New project command
	a.Commands["new"] = Command{
		Name:        "new",
		Description: "Create a new project, post, or page interactively",
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
	fmt.Printf("%s - A lightweight static project generator\n\n", a.Name)
	fmt.Println("Usage:")
	fmt.Printf("  %s [command] [arguments]\n\n", a.Name)
	fmt.Println("Available commands:")

	for _, cmd := range a.Commands {
		fmt.Printf("  %-10s %s\n", cmd.Name, cmd.Description)
	}
	
	fmt.Println("\nExamples:")
	fmt.Printf("  %s new                  Create a new project in the current directory with interactive prompts\n", a.Name)
	fmt.Printf("  %s new my-project       Create a new project in 'my-project' directory with interactive prompts\n", a.Name)
	fmt.Printf("  %s serve                Start development server for the current directory\n", a.Name)
	fmt.Printf("  %s build                Build the project in the current directory\n", a.Name)

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
	// If no arguments, assume we're creating a new site
	if len(args) < 1 {
		return a.createNewProject("")
	}

	// If first arg starts with a letter and not "project", "post", or "page", 
	// treat it as a project name
	firstArg := args[0]
	if firstArg != "project" && firstArg != "post" && firstArg != "page" {
		return a.createNewProject(firstArg)
	}

	// Otherwise, handle the classic way
	resType := firstArg
	name := ""
	if len(args) > 1 {
		name = args[1]
	}

	switch resType {
	case "project":
		return a.createNewProject(name)
	case "post":
		return a.createNewPost(name)
	case "page":
		return a.createNewPage(name)
	default:
		return fmt.Errorf("unknown resource type: %s", resType)
	}
}

// createNewProject scaffolds a new project with default structure and templates
func (a *App) createNewProject(name string) error {
	// Prompt for name if not provided
	projectPath := "."
	if name == "" {
		fmt.Print("Enter project name (leave empty to use current directory): ")
		fmt.Scanln(&name)
	}
	
	if name != "" {
		projectPath = name
		// Create project directory
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return fmt.Errorf("failed to create project directory: %w", err)
		}
	}

	// Prompt for template selection
	template := "none"
	fmt.Println("Select a template:")
	fmt.Println("1) None (minimal)")
	fmt.Println("2) Blog")
	fmt.Println("3) Docs")
	fmt.Println("4) Kitchen sink (all features)")
	fmt.Print("Enter choice (1-4) [1]: ")
	var choice string
	fmt.Scanln(&choice)
	
	switch choice {
	case "2":
		template = "blog"
	case "3":
		template = "docs"
	case "4":
		template = "kitchen-sink"
	default:
		template = "none" // default to minimal if input is empty or invalid
	}
	
	fmt.Printf("Creating new Scribe project in '%s' with '%s' template...\n", projectPath, template)

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
		path := filepath.Join(projectPath, dir)
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
	
	if err := cfg.Save(projectPath); err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}

	// Create sample content
	if err := a.createSampleContent(projectPath); err != nil {
		return fmt.Errorf("failed to create sample content: %w", err)
	}

	// Create default templates
	if err := a.createDefaultTemplates(projectPath); err != nil {
		return fmt.Errorf("failed to create default templates: %w", err)
	}
	
	// Initialize git repository if requested
	initGit := false
	fmt.Print("Initialize git repository? [y/N]: ")
	var gitChoice string
	fmt.Scanln(&gitChoice)
	gitChoice = strings.ToLower(gitChoice)
	if gitChoice == "y" || gitChoice == "yes" {
		initGit = true
	}
	
	if initGit {
		fmt.Println("Initializing git repository...")
		gitCmd := exec.Command("git", "init", projectPath)
		err := gitCmd.Run()
		if err != nil {
			fmt.Printf("Warning: Failed to initialize git repository: %v\n", err)
		} else {
			// Create .gitignore
			gitignorePath := filepath.Join(projectPath, ".gitignore")
			gitignoreContent := "# Output directory\npublic/\n\n# IDE files\n.idea/\n.vscode/\n\n# System files\n.DS_Store\nThumbs.db\n"
			if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
				fmt.Printf("Warning: Failed to create .gitignore file: %v\n", err)
			}
		}
	}

	fmt.Println("Project created successfully!")
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
