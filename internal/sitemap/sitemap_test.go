package sitemap

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
)

func TestSitemapGeneration(t *testing.T) {
	// Create temp directory for test output
	tempDir, err := os.MkdirTemp("", "sitemap-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test config
	cfg := config.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	// Create test pages
	now := time.Now()
	testPages := []content.Page{
		{
			Title:     "Home",
			Date:      now,
			Draft:     false,
			Permalink: "https://example.com/",
			IsPost:    false,
		},
		{
			Title:     "About",
			Date:      now.Add(-24 * time.Hour),
			Draft:     false,
			Permalink: "https://example.com/about/",
			IsPost:    false,
		},
		{
			Title:     "First Post",
			Date:      now.Add(-48 * time.Hour),
			Draft:     false,
			Permalink: "https://example.com/posts/first-post/",
			IsPost:    true,
		},
		{
			Title:     "Draft Post",
			Date:      now,
			Draft:     true, // This should be excluded
			Permalink: "https://example.com/posts/draft-post/",
			IsPost:    true,
		},
	}

	// Create generator
	generator := NewGenerator(cfg)

	// Generate sitemap
	outputPath := filepath.Join(tempDir, "sitemap.xml")
	err = generator.Generate(testPages, outputPath)
	if err != nil {
		t.Fatalf("Failed to generate sitemap: %v", err)
	}

	// Read generated file
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated sitemap: %v", err)
	}

	// Parse XML
	var urlset URLSet
	err = xml.Unmarshal(data, &urlset)
	if err != nil {
		t.Fatalf("Failed to parse generated sitemap XML: %v", err)
	}

	// Verify XMLNS
	if urlset.XMLNS != "http://www.sitemaps.org/schemas/sitemap/0.9" {
		t.Errorf("Expected xmlns to be 'http://www.sitemaps.org/schemas/sitemap/0.9', got %s", urlset.XMLNS)
	}

	// Verify URLs
	expectedURLCount := 4 // 3 non-draft pages + homepage added by generator
	if len(urlset.URLs) != expectedURLCount {
		t.Errorf("Expected %d URLs, got %d", expectedURLCount, len(urlset.URLs))
	}

	// Verify homepage is present
	homepageFound := false
	for _, url := range urlset.URLs {
		if url.Loc == "https://example.com/" {
			homepageFound = true
		}
	}
	if !homepageFound {
		t.Errorf("Homepage URL not found in sitemap")
	}

	// Verify draft pages are excluded
	for _, url := range urlset.URLs {
		if url.Loc == "https://example.com/posts/draft-post/" {
			t.Errorf("Draft post should not be in sitemap")
		}
	}
}

func TestBasicURLGeneration(t *testing.T) {
	// Create temp directory for test output
	tempDir, err := os.MkdirTemp("", "sitemap-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with empty page list to check basic functionality
	cfg := config.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}

	generator := NewGenerator(cfg)
	outputPath := filepath.Join(tempDir, "sitemap.xml")
	err = generator.Generate([]content.Page{}, outputPath)
	if err != nil {
		t.Fatalf("Failed to generate sitemap: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Sitemap file was not created")
	}

	// Read and parse the file
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read generated sitemap: %v", err)
	}

	var urlset URLSet
	err = xml.Unmarshal(data, &urlset)
	if err != nil {
		t.Fatalf("Failed to parse generated sitemap XML: %v", err)
	}

	// Should only contain the homepage
	if len(urlset.URLs) != 1 {
		t.Errorf("Expected 1 URL (homepage), got %d", len(urlset.URLs))
	}
}

func TestBaseURLHandling(t *testing.T) {
	// Test trailing slash handling in base URL
	tempDir, err := os.MkdirTemp("", "sitemap-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test config without trailing slash
	cfg := config.Config{
		Title:   "Test Site",
		BaseURL: "https://example.com", // No trailing slash
	}

	generator := NewGenerator(cfg)
	
	// The base URL in the generator should have a trailing slash
	if generator.baseURL != "https://example.com/" {
		t.Errorf("Expected baseURL to have trailing slash, got %s", generator.baseURL)
	}
}