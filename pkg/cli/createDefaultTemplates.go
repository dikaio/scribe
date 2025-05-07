package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dikaio/scribe/internal/templates"
)

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
		
		// No config file needed for modern Tailwind CSS 4.1
		
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