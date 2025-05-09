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
	
	// Handle subdirectories correctly
	// Keep directory structure for all content
	if dir != "." {
		return filepath.Join(dir, slug)
	}
	
	// For all cases, use the slug (which is either the custom slug provided 
	// in frontmatter or the filename without extension)
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
	// A file is a post if it's in any directory named "posts"
	relativePath := filePath
	contentIdx := strings.Index(filePath, "/content/")
	if contentIdx >= 0 {
		relativePath = filePath[contentIdx+9:] // +9 for "/content/"
	}
	
	isPost := strings.HasPrefix(relativePath, "posts/") || 
		strings.Contains(relativePath, "/posts/")

	// Generate slug from filename if not specified
	slug := frontMatter.Slug
	if slug == "" {
		// This is just for the page metadata - the URL is handled separately
		baseName := filepath.Base(filePath)
		extName := filepath.Ext(baseName)
		slug = strings.TrimSuffix(baseName, extName)
	}

	// Determine URL from the file path, preserving directory structure
	url := extractContentPath(filePath, slug)

	// Add trailing slash for cleaner URLs (remove .html extension)
	// but don't add double slashes
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}

	// For permalinks, join baseURL and url properly
	permalink := baseURL
	if !strings.HasSuffix(permalink, "/") {
		permalink += "/"
	}
	// Remove leading slash from url if it exists to avoid double slashes
	cleanURL := strings.TrimPrefix(url, "/")
	permalink = permalink + cleanURL

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
