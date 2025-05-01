package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dikaio/scribes/internal/config"
)

// TemplateManager manages template loading and rendering
type TemplateManager struct {
	templates map[string]*template.Template
	config    config.Config
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(cfg config.Config) *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*template.Template),
		config:    cfg,
	}
}

// LoadTemplates loads all templates from the layouts directory
func (tm *TemplateManager) LoadTemplates(sitePath string) error {
	// Define template functions
	funcMap := template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("January 2, 2006")
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
	}

	// Load templates from site and theme
	themePath := filepath.Join(sitePath, "themes", tm.config.Theme, "layouts")
	siteLayoutPath := filepath.Join(sitePath, tm.config.LayoutDir)

	// First try to load base template from theme
	baseTemplatePath := filepath.Join(themePath, "base.html")
	if _, err := os.Stat(baseTemplatePath); os.IsNotExist(err) {
		// Then try site layouts
		baseTemplatePath = filepath.Join(siteLayoutPath, "base.html")
		if _, err := os.Stat(baseTemplatePath); os.IsNotExist(err) {
			return fmt.Errorf("base template not found")
		}
	}

	// Create a map to hold layout templates keyed by name
	layoutTemplates := make(map[string][]string)
	
	// Add base template to each layout
	layoutTemplates["base"] = []string{baseTemplatePath}
	
	// Collect theme templates
	themeLayoutFiles, err := filepath.Glob(filepath.Join(themePath, "*.html"))
	if err == nil {
		for _, file := range themeLayoutFiles {
			if file != baseTemplatePath {
				name := filepath.Base(file)
				name = strings.TrimSuffix(name, filepath.Ext(name))
				if _, exists := layoutTemplates[name]; !exists {
					layoutTemplates[name] = []string{baseTemplatePath}
				}
				layoutTemplates[name] = append(layoutTemplates[name], file)
			}
		}
	}
	
	// Collect site templates (overrides)
	siteLayoutFiles, err := filepath.Glob(filepath.Join(siteLayoutPath, "*.html"))
	if err == nil {
		for _, file := range siteLayoutFiles {
			if file != baseTemplatePath {
				name := filepath.Base(file)
				name = strings.TrimSuffix(name, filepath.Ext(name))
				if _, exists := layoutTemplates[name]; !exists {
					layoutTemplates[name] = []string{baseTemplatePath}
				}
				// Site templates override theme templates, so we replace instead of append
				// First keep the base template
				baseTemplate := layoutTemplates[name][0]
				layoutTemplates[name] = []string{baseTemplate, file}
			}
		}
	}
	
	// Parse all template combinations
	for name, files := range layoutTemplates {
		// Parse the template set
		tmpl, err := template.New(filepath.Base(files[0])).Funcs(funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", name, err)
		}
		
		tm.templates[name] = tmpl
	}

	return nil
}

// GetTemplate returns a template by name
func (tm *TemplateManager) GetTemplate(name string) (*template.Template, error) {
	tmpl, exists := tm.templates[name]
	if !exists {
		return nil, fmt.Errorf("template %s not found", name)
	}

	return tmpl, nil
}
