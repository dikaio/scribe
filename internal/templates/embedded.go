package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
)

//go:embed embedded/*
var EmbeddedFS embed.FS

// GetEmbeddedTemplate returns the content of an embedded template file
func GetEmbeddedTemplate(theme, name string) (string, error) {
	path := filepath.Join("embedded", theme, name)
	content, err := EmbeddedFS.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded template %s: %w", path, err)
	}
	return string(content), nil
}

// ListEmbeddedTemplates returns a list of all embedded templates for a theme
func ListEmbeddedTemplates(theme string) ([]string, error) {
	var templates []string
	
	dir := filepath.Join("embedded", theme)
	entries, err := fs.ReadDir(EmbeddedFS, dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded directory %s: %w", dir, err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}
	
	return templates, nil
}

// GetDefaultTemplate returns the default template content using the embedded files
func GetDefaultTemplate(name string) (string, error) {
	return GetEmbeddedTemplate("default", name)
}

// GetAllDefaultTemplates returns all embedded default templates
func GetAllDefaultTemplates() (map[string]string, error) {
	return getAllTemplates("default")
}

// getAllTemplates is a helper to get all templates for a theme
func getAllTemplates(theme string) (map[string]string, error) {
	templates := make(map[string]string)
	
	fileNames, err := ListEmbeddedTemplates(theme)
	if err != nil {
		return nil, err
	}
	
	for _, fileName := range fileNames {
		content, err := GetEmbeddedTemplate(theme, fileName)
		if err != nil {
			return nil, err
		}
		templates[fileName] = content
	}
	
	return templates, nil
}