package content

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Page represents a content page
type Page struct {
	Title       string
	Description string
	Date        time.Time
	Tags        []string
	Draft       bool
	Layout      string
	Slug        string
	Content     string
	HTML        string
	Path        string
	URL         string
	Permalink   string
	IsPost      bool
}

// extractContentPath extracts the URL path from the file path
// It preserves directory structure within the content directory
func extractContentPath(filePath string, slug string) string {
	// Find the "content" directory in the path
	contentIdx := strings.Index(filePath, "/content/")
	if contentIdx < 0 {
		// Fallback: if "content" not found, just use the slug
		return slug
	}

	// Get the path after "content/"
	relativePath := filePath[contentIdx+9:] // +9 for "/content/"
	
	// Replace the file name with the slug
	dir := filepath.Dir(relativePath)
	if dir == "." {
		// File is directly in content directory
		return slug
	}
	
	// Special case for posts directory to maintain backward compatibility
	if strings.HasPrefix(dir, "posts") && filepath.Dir(dir) == "." {
		return filepath.Join("posts", slug)
	}
	
	// Join directory with slug for the final URL
	return filepath.Join(dir, slug)
}

func LoadPage(filePath string, baseURL string) (Page, error) {
	var page Page

	// Read file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return page, err
	}

	// Parse front matter
	frontMatter, content, err := ParseFrontMatter(data)
	if err != nil {
		return page, err
	}

	// Convert markdown to HTML
	html := MarkdownToHTML(content)

	// Determine if it's a post based on the path
	isPost := strings.Contains(filePath, "/posts/")

	// Generate slug from filename if not specified
	slug := frontMatter.Slug
	if slug == "" {
		baseName := filepath.Base(filePath)
		extName := filepath.Ext(baseName)
		slug = strings.TrimSuffix(baseName, extName)
	}

	// Determine URL from the file path, preserving directory structure
	url := extractContentPath(filePath, slug)

	permalink := filepath.Join(baseURL, url)

	// Create page
	page = Page{
		Title:       frontMatter.Title,
		Description: frontMatter.Description,
		Date:        frontMatter.Date,
		Tags:        frontMatter.Tags,
		Draft:       frontMatter.Draft,
		Layout:      frontMatter.Layout,
		Slug:        slug,
		Content:     string(content),
		HTML:        string(html),
		Path:        filePath,
		URL:         url,
		Permalink:   permalink,
		IsPost:      isPost,
	}

	return page, nil
}
