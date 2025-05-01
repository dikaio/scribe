package render

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/dikaio/scribes/internal/config"
	"github.com/dikaio/scribes/internal/content"
)

// Renderer handles rendering pages to HTML files
type Renderer struct {
	templateManager *TemplateManager
	config          config.Config
}

// NewRenderer creates a new renderer
func NewRenderer(cfg config.Config) *Renderer {
	return &Renderer{
		templateManager: NewTemplateManager(cfg),
		config:          cfg,
	}
}

// Init initializes the renderer
func (r *Renderer) Init(sitePath string) error {
	return r.templateManager.LoadTemplates(sitePath)
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

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Create output file
	f, err := os.Create(outputPath)
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
	return tmpl.ExecuteTemplate(f, "base", data)
}

// RenderList renders a list page (e.g., index, tag list)
func (r *Renderer) RenderList(title string, pages []content.Page, outputPath string) error {
	// Get template
	tmpl, err := r.templateManager.GetTemplate("list")
	if err != nil {
		return err
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	// Create output file
	f, err := os.Create(outputPath)
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
	return tmpl.ExecuteTemplate(f, "base", data)
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
	f, err := os.Create(outputPath)
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
	return tmpl.ExecuteTemplate(f, "base", data)
}
