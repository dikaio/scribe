package console

import (
	"time"

	"github.com/dikaio/scribe/internal/content"
)

// MockPost creates a mock post for testing
func MockPost() content.Page {
	return content.Page{
		Title:       "Test Post",
		Date:        time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Path:        "content/post.md",
		URL:         "posts/test-post",
		Tags:        []string{"test", "post"},
		Draft:       false,
		IsPost:      true,
		Description: "This is a test post",
	}
}

// MockPage creates a mock page for testing
func MockPage() content.Page {
	return content.Page{
		Title:       "About Page",
		Date:        time.Date(2023, 1, 2, 12, 0, 0, 0, time.UTC),
		Path:        "content/about.md",
		URL:         "about",
		Tags:        []string{},
		Draft:       false,
		IsPost:      false,
		Description: "This is an about page",
	}
}