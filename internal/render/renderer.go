package render

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
)

// Renderer handles rendering pages to HTML files
type Renderer struct {
	templateManager *TemplateManager
	config          config.Config
	devMode         bool
}

// NewRenderer creates a new renderer
func NewRenderer(cfg config.Config) *Renderer {
	return &Renderer{
		templateManager: NewTemplateManager(cfg),
		config:          cfg,
		devMode:         false,
	}
}

// SetDevMode enables or disables development mode (disables caching)
func (r *Renderer) SetDevMode(enabled bool) {
	r.devMode = enabled
	if enabled {
		r.templateManager.DisableCaching()
	} else {
		r.templateManager.EnableCaching()
	}
}

// Init initializes the renderer
func (r *Renderer) Init(sitePath string) error {
	return r.templateManager.LoadTemplates(sitePath)
}

// createOutputFile creates output file and ensures directory exists
func (r *Renderer) createOutputFile(outputPath string) (*os.File, error) {
	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}
	
	return f, nil
}

// RenderPage renders a page to an HTML file
func (r *Renderer) RenderPage(page content.Page, outputPath string) error {
	// Create layout name based on page's layout or default to "single"
	layoutName := page.Layout
	if layoutName == "" {
		if page.IsPost {
			layoutName = "single"
		} else {
			layoutName = "page"
		}
	}

	// Get template
	tmpl, err := r.templateManager.GetTemplate(layoutName)
	if err != nil {
		// Fallback to single template
		tmpl, err = r.templateManager.GetTemplate("single")
		if err != nil {
			return err
		}
	}

	// Create output file
	f, err := r.createOutputFile(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Prepare template data
	data := map[string]interface{}{
		"Site":    r.config,
		"Page":    page,
		"Content": template.HTML(page.HTML),
	}

	// Execute template
	return tmpl.Execute(f, data)
}

// RenderList renders a list page (e.g., index, tag list)
func (r *Renderer) RenderList(title string, pages []content.Page, outputPath string) error {
	// Get template
	tmpl, err := r.templateManager.GetTemplate("list")
	if err != nil {
		return err
	}

	// Create output file
	f, err := r.createOutputFile(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Prepare template data
	data := map[string]interface{}{
		"Site":  r.config,
		"Title": title,
		"Pages": pages,
	}

	// Execute template
	return tmpl.Execute(f, data)
}

// RenderHome renders the home page
func (r *Renderer) RenderHome(pages []content.Page, outputPath string) error {
	// Get template
	tmpl, err := r.templateManager.GetTemplate("home")
	if err != nil {
		// Fallback to list template
		tmpl, err = r.templateManager.GetTemplate("list")
		if err != nil {
			return err
		}
	}

	// Create output file
	f, err := r.createOutputFile(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Prepare template data
	data := map[string]interface{}{
		"Site":  r.config,
		"Title": r.config.Title,
		"Pages": pages,
	}

	// Execute template
	return tmpl.Execute(f, data)
}
