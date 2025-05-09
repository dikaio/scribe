package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dikaio/scribe/internal/templates"
)

// createDefaultTemplates creates default templates for a new site
func (a *App) createDefaultTemplates(sitePath string, _ bool) error {
	// Use default templates (ignoring tailwind parameter)
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
	
	// Write default CSS file
	return os.WriteFile(cssFilePath, []byte(templates.StyleCSS), 0644)
}