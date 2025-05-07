package render

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dikaio/scribe/internal/config"
)

// TemplateCache represents a cached template
type TemplateCache struct {
	Template *template.Template
	ModTime  time.Time
	Files    []string
}

// TemplateManager manages template loading and rendering
type TemplateManager struct {
	templates    map[string]*template.Template
	cache        map[string]TemplateCache
	config       config.Config
	funcMap      template.FuncMap
	cacheMutex   sync.RWMutex
	cachingEnabled bool
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(cfg config.Config) *TemplateManager {
	// Define template functions
	funcMap := template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("January 2, 2006")
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
	}

	return &TemplateManager{
		templates:    make(map[string]*template.Template),
		cache:        make(map[string]TemplateCache),
		config:       cfg,
		funcMap:      funcMap,
		cachingEnabled: true,
	}
}

// DisableCaching disables template caching (for development mode)
func (tm *TemplateManager) DisableCaching() {
	tm.cachingEnabled = false
}

// EnableCaching enables template caching (for production mode)
func (tm *TemplateManager) EnableCaching() {
	tm.cachingEnabled = true
}

// getFileModTime gets the latest modification time of a file or files
func getFileModTime(files ...string) (time.Time, error) {
	var latest time.Time

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			return latest, err
		}

		if info.ModTime().After(latest) {
			latest = info.ModTime()
		}
	}

	return latest, nil
}

// templateNeedsUpdate checks if any template files have been modified
func (tm *TemplateManager) templateNeedsUpdate(name string, files []string) (bool, time.Time, error) {
	// Always return true if caching is disabled
	if !tm.cachingEnabled {
		return true, time.Time{}, nil
	}

	tm.cacheMutex.RLock()
	cachedTemplate, exists := tm.cache[name]
	tm.cacheMutex.RUnlock()

	if !exists {
		return true, time.Time{}, nil
	}

	// Check if file list has changed
	if len(cachedTemplate.Files) != len(files) {
		return true, time.Time{}, nil
	}

	for i, file := range files {
		if cachedTemplate.Files[i] != file {
			return true, time.Time{}, nil
		}
	}

	// Get the most recent modification time
	latestMod, err := getFileModTime(files...)
	if err != nil {
		return true, time.Time{}, err
	}

	// Check if any file has been modified since the template was cached
	return latestMod.After(cachedTemplate.ModTime), latestMod, nil
}

// LoadTemplates loads all templates from the layouts directory
func (tm *TemplateManager) LoadTemplates(sitePath string) error {
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
	
	// Parse all template combinations, using cache where possible
	for name, files := range layoutTemplates {
		// Check if the template needs to be reloaded
		needsUpdate, modTime, err := tm.templateNeedsUpdate(name, files)
		if err != nil {
			return fmt.Errorf("error checking template modification time: %v", err)
		}

		if !needsUpdate {
			// Use cached template
			tm.cacheMutex.RLock()
			tm.templates[name] = tm.cache[name].Template
			tm.cacheMutex.RUnlock()
			continue
		}

		// Parse the template set
		tmpl, err := template.New(filepath.Base(files[0])).Funcs(tm.funcMap).ParseFiles(files...)
		if err != nil {
			return fmt.Errorf("error parsing template %s: %v", name, err)
		}
		
		// Update the template in the current instance
		tm.templates[name] = tmpl

		// Update the cache if caching is enabled
		if tm.cachingEnabled {
			// If modTime is zero, get it now
			if modTime.IsZero() {
				modTime, err = getFileModTime(files...)
				if err != nil {
					return fmt.Errorf("error getting template modification time: %v", err)
				}
			}

			// Cache the template
			tm.cacheMutex.Lock()
			tm.cache[name] = TemplateCache{
				Template: tmpl,
				ModTime:  modTime,
				Files:    append([]string{}, files...), // Copy the files slice
			}
			tm.cacheMutex.Unlock()
		}
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
