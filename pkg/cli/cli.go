package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dikaio/scribes/internal/build"
	"github.com/dikaio/scribes/internal/config"
	"github.com/dikaio/scribes/internal/console"
	"github.com/dikaio/scribes/internal/content"
	"github.com/dikaio/scribes/internal/server"
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
		Version:  "0.1.0",
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
		Description: "Create a new site, post, page, theme, plugin, or partial",
		Action:      a.cmdNew,
	}

	// Test command
	a.Commands["test"] = Command{
		Name:        "test",
		Description: "Run the tests",
		Action:      a.cmdTest,
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

// cmdBuild implements the build command, which generates the static site.
// It takes an optional path argument (or uses the current directory if not provided).
func (a *App) cmdBuild(args []string) error {
	// Determine site path
	sitePath := "."
	if len(args) > 0 {
		sitePath = args[0]
	}

	fmt.Printf("Building site from '%s'...\n", sitePath)

	// Load the site configuration
	cfg, err := config.LoadConfig(sitePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config file not found in '%s', make sure this is a valid Scribes site directory", sitePath)
		}
		return fmt.Errorf("failed to load configuration: %w", err)
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
	// Determine site path
	sitePath := "."
	if len(args) > 0 {
		sitePath = args[0]
	}

	fmt.Printf("Starting development server for '%s'...\n", sitePath)

	// Load the site configuration
	cfg, err := config.LoadConfig(sitePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config file not found in '%s', make sure this is a valid Scribes site directory", sitePath)
		}
		return fmt.Errorf("failed to load configuration: %w", err)
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
	// Determine site path
	sitePath := "."
	if len(args) > 0 {
		sitePath = args[0]
	}

	fmt.Printf("Starting console for '%s'...\n", sitePath)

	// Load the site configuration
	cfg, err := config.LoadConfig(sitePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("config file not found in '%s', make sure this is a valid Scribes site directory", sitePath)
		}
		return fmt.Errorf("failed to load configuration: %w", err)
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
	case "theme":
		fmt.Printf("Creating new theme: %s\n", name)
		return errors.New("not implemented yet")
	case "plugin":
		fmt.Printf("Creating new plugin: %s\n", name)
		return errors.New("not implemented yet")
	case "partial":
		fmt.Printf("Creating new partial: %s\n", name)
		return errors.New("not implemented yet")
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

	fmt.Printf("Creating new Scribes site in '%s'...\n", sitePath)

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
	// Sample post
	samplePost := `---
title: Welcome to Scribes
description: A sample post to get you started
date: 2025-05-01T12:00:00Z
tags:
  - welcome
  - scribes
draft: false
---

# Welcome to Scribes!

This is a sample post to help you get started with Scribes, a lightweight static site generator.

## Features

- Markdown support
- Fast and lightweight
- No external dependencies
- Simple to use

Enjoy creating content with Scribes!
`
	postPath := filepath.Join(sitePath, "content", "posts", "welcome.md")
	if err := os.WriteFile(postPath, []byte(samplePost), 0644); err != nil {
		return err
	}

	// Sample page
	samplePage := `---
title: About
description: About this site
draft: false
---

# About

This is an about page for your Scribes site. You can add information about yourself or your project here.

## Contact

Feel free to reach out with any questions or feedback.
`
	pagePath := filepath.Join(sitePath, "content", "about.md")
	return os.WriteFile(pagePath, []byte(samplePage), 0644)
}

// createDefaultTemplates creates default templates for a new site
func (a *App) createDefaultTemplates(sitePath string) error {
	// Base template
	baseTemplate := `<!DOCTYPE html>
<html lang="{{.Site.Language}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .Title}}{{.Title}} | {{end}}{{.Site.Title}}</title>
    <meta name="description" content="{{if .Description}}{{.Description}}{{else}}{{.Site.Description}}{{end}}">
    <link rel="stylesheet" href="/css/style.css">
</head>
<body>
    <header>
        <div class="container">
            <h1><a href="/">{{.Site.Title}}</a></h1>
            <nav>
                <ul>
                    <li><a href="/">Home</a></li>
                    <li><a href="/about/">About</a></li>
                </ul>
            </nav>
        </div>
    </header>
    <main class="container">
        {{block "content" .}}{{end}}
    </main>
    <footer>
        <div class="container">
            <p>&copy; {{.Site.Title}}</p>
        </div>
    </footer>
</body>
</html>`

	// Single post template
	singleTemplate := `{{define "content"}}
<article>
    <header>
        <h1>{{.Page.Title}}</h1>
        <p class="meta">
            <time>{{formatDate .Page.Date}}</time>
            {{if .Page.Tags}}
            | Tags: 
            {{range .Page.Tags}}
            <a href="/tags/{{.}}/">{{.}}</a>
            {{end}}
            {{end}}
        </p>
    </header>
    <div class="content">
        {{.Content}}
    </div>
</article>
{{end}}`

	// List template
	listTemplate := `{{define "content"}}
<h1>{{.Title}}</h1>
<div class="post-list">
    {{range .Pages}}
    <article class="post-summary">
        <h2><a href="/{{.URL}}/">{{.Title}}</a></h2>
        <p class="meta">
            <time>{{formatDate .Date}}</time>
            {{if .Tags}}
            | Tags: 
            {{range .Tags}}
            <a href="/tags/{{.}}/">{{.}}</a>
            {{end}}
            {{end}}
        </p>
        <p>{{.Description}}</p>
    </article>
    {{end}}
</div>
{{end}}`

	// Home template
	homeTemplate := `{{define "content"}}
<h1>Recent Posts</h1>
<div class="post-list">
    {{range .Pages}}
    <article class="post-summary">
        <h2><a href="/{{.URL}}/">{{.Title}}</a></h2>
        <p class="meta">
            <time>{{formatDate .Date}}</time>
            {{if .Tags}}
            | Tags: 
            {{range .Tags}}
            <a href="/tags/{{.}}/">{{.}}</a>
            {{end}}
            {{end}}
        </p>
        <p>{{.Description}}</p>
    </article>
    {{end}}
</div>
{{end}}`

	// Page template
	pageTemplate := `{{define "content"}}
<article>
    <header>
        <h1>{{.Page.Title}}</h1>
    </header>
    <div class="content">
        {{.Content}}
    </div>
</article>
{{end}}`

	// CSS file
	cssContent := `/* Basic styles for Scribes default theme */
:root {
    --primary-color: #0077cc;
    --text-color: #333;
    --background-color: #fff;
    --light-gray: #f5f5f5;
    --border-color: #ddd;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
    line-height: 1.6;
    color: var(--text-color);
    background-color: var(--background-color);
}

.container {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 20px;
}

header {
    background-color: var(--light-gray);
    padding: 20px 0;
    margin-bottom: 40px;
    border-bottom: 1px solid var(--border-color);
}

header h1 {
    font-size: 2rem;
}

header h1 a {
    color: var(--text-color);
    text-decoration: none;
}

nav ul {
    list-style: none;
    display: flex;
    gap: 20px;
}

nav a {
    color: var(--primary-color);
    text-decoration: none;
}

main {
    min-height: 70vh;
    margin-bottom: 40px;
}

footer {
    background-color: var(--light-gray);
    padding: 20px 0;
    border-top: 1px solid var(--border-color);
    text-align: center;
}

h1, h2, h3, h4, h5, h6 {
    margin-bottom: 1rem;
    line-height: 1.25;
}

p, ul, ol {
    margin-bottom: 1.5rem;
}

a {
    color: var(--primary-color);
}

.post-list {
    display: flex;
    flex-direction: column;
    gap: 30px;
}

.post-summary {
    padding-bottom: 20px;
    border-bottom: 1px solid var(--border-color);
}

.meta {
    color: #666;
    font-size: 0.9rem;
    margin-bottom: 1rem;
}

article .content {
    margin-top: 20px;
}

/* Code blocks */
pre {
    background-color: var(--light-gray);
    padding: 1rem;
    overflow-x: auto;
    border-radius: 4px;
    margin-bottom: 1.5rem;
}

code {
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
    font-size: 0.9em;
    background-color: var(--light-gray);
    padding: 0.2em 0.4em;
    border-radius: 3px;
}

pre code {
    padding: 0;
    background-color: transparent;
}`

	// Save templates
	templates := map[string]string{
		filepath.Join(sitePath, "themes", "default", "layouts", "base.html"):   baseTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "single.html"): singleTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "list.html"):   listTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "home.html"):   homeTemplate,
		filepath.Join(sitePath, "themes", "default", "layouts", "page.html"):   pageTemplate,
	}

	for path, content := range templates {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create template '%s': %w", path, err)
		}
	}

	// Create static directory with CSS
	cssPath := filepath.Join(sitePath, "themes", "default", "static", "css")
	if err := os.MkdirAll(cssPath, 0755); err != nil {
		return fmt.Errorf("failed to create CSS directory: %w", err)
	}

	cssFilePath := filepath.Join(cssPath, "style.css")
	return os.WriteFile(cssFilePath, []byte(cssContent), 0644)
}

func (a *App) cmdTest(args []string) error {
	fmt.Println("Running tests...")
	return errors.New("not implemented yet")
}
