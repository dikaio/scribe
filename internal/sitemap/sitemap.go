package sitemap

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
)

// URLSet represents the root element of a sitemap
type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL    `xml:"url"`
}

// URL represents a URL entry in the sitemap
type URL struct {
	Loc        string  `xml:"loc"`
	LastMod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float64 `xml:"priority,omitempty"`
}

// Generator handles sitemap generation
type Generator struct {
	config  config.Config
	baseURL string
}

// NewGenerator creates a new sitemap generator
func NewGenerator(cfg config.Config) *Generator {
	// Ensure the base URL has a trailing slash
	baseURL := cfg.BaseURL
	if baseURL != "" && !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return &Generator{
		config:  cfg,
		baseURL: baseURL,
	}
}

// Generate creates a sitemap.xml file from a list of pages
func (g *Generator) Generate(pages []content.Page, outputPath string) error {
	// Create sitemap structure
	urlset := URLSet{
		XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]URL, 0, len(pages)+1), // +1 for the homepage
	}

	// Add homepage
	urlset.URLs = append(urlset.URLs, URL{
		Loc:        g.baseURL,
		LastMod:    time.Now().Format("2006-01-02"),
		ChangeFreq: "daily",
		Priority:   1.0,
	})

	// Add pages
	for _, page := range pages {
		// Skip draft pages
		if page.Draft {
			continue
		}

		// Create URL entry
		url := URL{
			Loc:     page.Permalink,
			LastMod: page.Date.Format("2006-01-02"),
		}

		// Set different priorities and change frequencies based on content type
		if page.IsPost {
			url.ChangeFreq = "weekly"
			url.Priority = 0.8
		} else {
			url.ChangeFreq = "monthly"
			url.Priority = 0.5
		}

		urlset.URLs = append(urlset.URLs, url)
	}

	// Ensure the directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory for sitemap: %w", err)
	}

	// Create output file
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create sitemap file: %w", err)
	}
	defer f.Close()

	// Write XML header
	f.WriteString(xml.Header)

	// Encode and write the sitemap
	encoder := xml.NewEncoder(f)
	encoder.Indent("", "  ")
	if err := encoder.Encode(urlset); err != nil {
		return fmt.Errorf("failed to encode sitemap: %w", err)
	}

	return nil
}