package templates

import (
	"log"
	"sync"
)

var (
	// Tailwind template strings loaded from embedded files
	TailwindInputCSS     string
	TailwindBaseTemplate string
	TailwindSingleTemplate string
	TailwindListTemplate string
	TailwindHomeTemplate string
	TailwindPageTemplate string
	TailwindPackageJSON string
	TailwindGitignore string
	TailwindREADME string

	// Initialization once
	tailwindTemplatesOnce sync.Once
)

// loadTailwindTemplates loads all tailwind templates from embedded files
func loadTailwindTemplates() {
	var err error

	TailwindInputCSS, err = GetTailwindTemplate("input.css")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind input CSS: %v", err)
	}

	TailwindBaseTemplate, err = GetTailwindTemplate("base.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind base template: %v", err)
	}

	TailwindSingleTemplate, err = GetTailwindTemplate("single.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind single template: %v", err)
	}

	TailwindListTemplate, err = GetTailwindTemplate("list.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind list template: %v", err)
	}

	TailwindHomeTemplate, err = GetTailwindTemplate("home.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind home template: %v", err)
	}

	TailwindPageTemplate, err = GetTailwindTemplate("page.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind page template: %v", err)
	}

	TailwindPackageJSON, err = GetTailwindTemplate("package.json")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind package.json: %v", err)
	}

	TailwindGitignore, err = GetTailwindTemplate("gitignore")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind gitignore: %v", err)
	}

	TailwindREADME, err = GetTailwindTemplate("README.md")
	if err != nil {
		log.Printf("Warning: Failed to load embedded tailwind README: %v", err)
	}
}

// Default templates for site creation - will load from embedded files
func init() {
	tailwindTemplatesOnce.Do(loadTailwindTemplates)
}