package content

import (
	"testing"
)

func TestExtractContentPath(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		slug     string
		expected string
	}{
		{
			name:     "Root content file",
			filePath: "/path/to/site/content/page.md",
			slug:     "page",
			expected: "page",
		},
		{
			name:     "Post in posts directory",
			filePath: "/path/to/site/content/posts/article.md",
			slug:     "article",
			expected: "posts/article",
		},
		{
			name:     "File in custom directory",
			filePath: "/path/to/site/content/articles/tech/golang.md",
			slug:     "golang",
			expected: "articles/tech/golang",
		},
		{
			name:     "Deeply nested file",
			filePath: "/path/to/site/content/topics/programming/languages/go/basics.md",
			slug:     "basics",
			expected: "topics/programming/languages/go/basics",
		},
		{
			name:     "File with custom slug",
			filePath: "/path/to/site/content/articles/tech/javascript.md",
			slug:     "js-tutorial",
			expected: "articles/tech/js-tutorial",
		},
		{
			name:     "Fallback when content not in path",
			filePath: "/some/other/path/javascript.md",
			slug:     "js-tutorial",
			expected: "js-tutorial",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractContentPath(tt.filePath, tt.slug)
			if result != tt.expected {
				t.Errorf("extractContentPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}