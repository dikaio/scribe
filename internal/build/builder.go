package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/dikaio/scribes/internal/config"
	"github.com/dikaio/scribes/internal/content"
	"github.com/dikaio/scribes/internal/render"
)

// Builder handles site building
type Builder struct {
	config   config.Config
	renderer *render.Renderer
	pages    []content.Page
	tags     map[string][]content.Page
}

// NewBuilder creates a new site builder
func NewBuilder(cfg config.Config) *Builder {
	return &Builder{
		config:   cfg,
		renderer: render.NewRenderer(cfg),
		pages:    []content.Page{},
		tags:     make(map[string][]content.Page),
	}
}

// Build builds the site
func (b *Builder) Build(sitePath string) error {
	// Initialize renderer
	if err := b.renderer.Init(sitePath); err != nil {
		return err
	}

	// Load content
	if err := b.loadContent(sitePath); err != nil {
		return err
	}

	// Create output directory
	outputPath := filepath.Join(sitePath, b.config.OutputDir)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	// Copy static files
	if err := b.copyStaticFiles(sitePath, outputPath); err != nil {
		return err
	}

	// Generate pages
	if err := b.generatePages(outputPath); err != nil {
		return err
	}

	// Generate tag pages
	if err := b.generateTagPages(outputPath); err != nil {
		return err
	}

	// Generate home page
	if err := b.generateHomePage(outputPath); err != nil {
		return err
	}

	return nil
}

// loadContent loads all content files
func (b *Builder) loadContent(sitePath string) error {
	contentPath := filepath.Join(sitePath, b.config.ContentDir)

	return filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		// Load page
		page, err := content.LoadPage(path, b.config.BaseURL)
		if err != nil {
			return err
		}

		// Skip draft pages in production
		if page.Draft {
			fmt.Printf("Skipping draft: %s\n", page.Title)
			return nil
		}

		// Add page to collection
		b.pages = append(b.pages, page)

		// Add page to tags
		for _, tag := range page.Tags {
			b.tags[tag] = append(b.tags[tag], page)
		}

		return nil
	})
}

// copyStaticFiles copies static files to the output directory
func (b *Builder) copyStaticFiles(sitePath, outputPath string) error {
	// Copy theme static files
	themeStaticPath := filepath.Join(sitePath, "themes", b.config.Theme, "static")
	if _, err := os.Stat(themeStaticPath); err == nil {
		if err := copyDir(themeStaticPath, outputPath); err != nil {
			return err
		}
	}

	// Copy site static files (overrides theme files)
	siteStaticPath := filepath.Join(sitePath, b.config.StaticDir)
	if _, err := os.Stat(siteStaticPath); err == nil {
		if err := copyDir(siteStaticPath, outputPath); err != nil {
			return err
		}
	}

	return nil
}

// generatePages generates all content pages
func (b *Builder) generatePages(outputPath string) error {
	for _, page := range b.pages {
		// Determine output file path
		outputFile := filepath.Join(outputPath, page.URL, "index.html")

		// Render page
		err := b.renderer.RenderPage(page, outputFile)
		if err != nil {
			return err
		}

		fmt.Printf("Generated: %s\n", page.URL)
	}

	return nil
}

// generateTagPages generates tag listing pages
func (b *Builder) generateTagPages(outputPath string) error {
	// Create tags directory
	tagsDir := filepath.Join(outputPath, "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	// Generate main tags index
	allTags := make([]string, 0, len(b.tags))
	for tag := range b.tags {
		allTags = append(allTags, tag)
	}
	sort.Strings(allTags)

	// Generate individual tag pages
	for tag, pages := range b.tags {
		// Sort pages by date (newest first)
		sort.Slice(pages, func(i, j int) bool {
			return pages[i].Date.After(pages[j].Date)
		})

		// Create tag directory
		tagDir := filepath.Join(tagsDir, tag)
		if err := os.MkdirAll(tagDir, 0755); err != nil {
			return err
		}

		// Render tag page
		outputFile := filepath.Join(tagDir, "index.html")
		title := fmt.Sprintf("Tag: %s", tag)
		if err := b.renderer.RenderList(title, pages, outputFile); err != nil {
			return err
		}

		fmt.Printf("Generated tag page: %s\n", tag)
	}

	return nil
}

// generateHomePage generates the home page
func (b *Builder) generateHomePage(outputPath string) error {
	// Filter and sort posts (newest first)
	posts := []content.Page{}
	for _, page := range b.pages {
		if page.IsPost {
			posts = append(posts, page)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// Render home page
	outputFile := filepath.Join(outputPath, "index.html")
	return b.renderer.RenderHome(posts, outputFile)
}

// copyDir recursively copies a directory tree
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from source
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip source root
		if rel == "." {
			return nil
		}

		// Get destination path
		dstPath := filepath.Join(dst, rel)

		// Create directories
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	return err
}
