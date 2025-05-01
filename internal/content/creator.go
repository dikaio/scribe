package content

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ContentType represents the type of content
type ContentType string

const (
	// PostType represents blog post content
	PostType ContentType = "post"
	// PageType represents static page content
	PageType ContentType = "page"
)

// Creator handles content creation
type Creator struct {
	sitePath string
}

// NewCreator creates a new content creator
func NewCreator(sitePath string) *Creator {
	return &Creator{
		sitePath: sitePath,
	}
}

// CreateContent creates a new content file
func (c *Creator) CreateContent(contentType ContentType, title, description string, tags []string, draft bool) (string, error) {
	// Generate slug from title
	slug := generateSlug(title)
	if slug == "" {
		return "", fmt.Errorf("invalid title")
	}

	// Prepare front matter
	fm := FrontMatter{
		Title:       title,
		Description: description,
		Date:        time.Now(),
		Tags:        tags,
		Draft:       draft,
		Slug:        slug,
	}

	// Determine content directory and file path
	var contentDir, filePath string
	if contentType == PostType {
		contentDir = filepath.Join(c.sitePath, "content", "posts")
		filePath = filepath.Join(contentDir, slug+".md")
	} else {
		contentDir = filepath.Join(c.sitePath, "content")
		filePath = filepath.Join(contentDir, slug+".md")
	}

	// Create content directory if it doesn't exist
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Generate content
	content := formatContent(fm, contentType)

	// Write to file
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filePath, nil
}

// formatContent formats content with front matter
func formatContent(fm FrontMatter, contentType ContentType) string {
	var sb strings.Builder

	// Write front matter
	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: %s\n", fm.Title))

	if fm.Description != "" {
		sb.WriteString(fmt.Sprintf("description: %s\n", fm.Description))
	}

	sb.WriteString(fmt.Sprintf("date: %s\n", fm.Date.Format(time.RFC3339)))

	if len(fm.Tags) > 0 {
		sb.WriteString("tags:\n")
		for _, tag := range fm.Tags {
			sb.WriteString(fmt.Sprintf("  - %s\n", tag))
		}
	}

	sb.WriteString(fmt.Sprintf("draft: %t\n", fm.Draft))

	if fm.Slug != "" {
		sb.WriteString(fmt.Sprintf("slug: %s\n", fm.Slug))
	}

	sb.WriteString("---\n\n")

	// Write default content
	sb.WriteString(fmt.Sprintf("# %s\n\n", fm.Title))

	if contentType == PostType {
		sb.WriteString("Write your post content here.\n")
	} else {
		sb.WriteString("Write your page content here.\n")
	}

	return sb.String()
}

// generateSlug generates a URL-friendly slug from a title
func generateSlug(title string) string {
	if title == "" {
		return ""
	}

	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Replace consecutive spaces
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove non-alphanumeric characters
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	slug = result.String()

	// Remove consecutive hyphens (again, after filtering)
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Trim hyphens from ends
	slug = strings.Trim(slug, "-")

	return slug
}
