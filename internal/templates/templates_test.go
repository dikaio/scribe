package templates

import (
	"html/template"
	"testing"
)

func TestDefaultTemplates(t *testing.T) {
	// Test that all default templates can be parsed
	tmpl := template.New("test")
	
	// Add helper functions
	tmpl.Funcs(template.FuncMap{
		"formatDate": func(date interface{}) string {
			return "2025-01-01" // Mock date for testing
		},
		"now": func() interface{} {
			return struct {
				Format func(string) string
			}{
				Format: func(string) string {
					return "2025"
				},
			}
		},
		"sub": func(a, b int) int {
			return a - b 
		},
	})
	
	// Add BaseTemplate
	var err error
	tmpl, err = tmpl.Parse(BaseTemplate)
	if err != nil {
		t.Fatalf("Failed to parse BaseTemplate: %v", err)
	}
	
	// Add other templates
	templateStrings := map[string]string{
		"SingleTemplate": SingleTemplate,
		"ListTemplate":   ListTemplate,
		"HomeTemplate":   HomeTemplate,
		"PageTemplate":   PageTemplate,
	}
	
	for name, templateString := range templateStrings {
		_, err = tmpl.New(name).Parse(templateString)
		if err != nil {
			t.Fatalf("Failed to parse %s: %v", name, err)
		}
	}
}


func TestStyleCSS(t *testing.T) {
	// Just check that StyleCSS is not empty
	if StyleCSS == "" {
		t.Error("StyleCSS is empty")
	}
}

func TestSampleContent(t *testing.T) {
	// Check that sample content is not empty
	if SamplePost == "" {
		t.Error("SamplePost is empty")
	}
	
	if SamplePage == "" {
		t.Error("SamplePage is empty")
	}
}