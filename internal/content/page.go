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

// LoadPage loads a page from a file
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

	// Determine URL and permalink
	url := slug
	if isPost {
		url = filepath.Join("posts", slug)
	}

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
