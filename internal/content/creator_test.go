package content

import (
	"os"
	"strings"
	"testing"
)

func TestCreateContent(t *testing.T) {
	// Create temporary directory for tests
	tempDir, err := os.MkdirTemp("", "scribe-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up after test

	// Create a new Creator
	creator := NewCreator(tempDir)

	tests := []struct {
		name        string
		contentType ContentType
		title       string
		description string
		tags        []string
		draft       bool
		customPath  string
		expectedDir string // Relative to content directory
		wantErr     bool
	}{
		{
			name:        "Create post",
			contentType: PostType,
			title:       "Test Post",
			description: "A test post",
			tags:        []string{"test", "example"},
			draft:       false,
			customPath:  "",
			expectedDir: "posts",
			wantErr:     false,
		},
		{
			name:        "Create page",
			contentType: PageType,
			title:       "Test Page",
			description: "A test page",
			tags:        nil,
			draft:       true,
			customPath:  "",
			expectedDir: "", // Root content directory
			wantErr:     false,
		},
		{
			name:        "Create content in custom directory",
			contentType: PageType,
			title:       "Custom Dir",
			description: "In a custom directory",
			tags:        nil,
			draft:       false,
			customPath:  "articles/tech",
			expectedDir: "articles/tech",
			wantErr:     false,
		},
		{
			name:        "Create content with custom path and filename",
			contentType: PageType,
			title:       "Custom Path",
			description: "With custom filename",
			tags:        nil,
			draft:       false,
			customPath:  "docs/guides/getting-started.md",
			expectedDir: "docs/guides",
			wantErr:     false,
		},
		{
			name:        "Empty title",
			contentType: PostType,
			title:       "",
			description: "Should fail",
			tags:        nil,
			draft:       false,
			customPath:  "",
			expectedDir: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create content
			filePath, err := creator.CreateContent(tt.contentType, tt.title, tt.description, tt.tags, tt.draft, tt.customPath)

			// Check error expectations
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip further checks if expected error
			if tt.wantErr {
				return
			}

			// Verify file exists
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Errorf("CreateContent() file was not created at %s", filePath)
			}
			
			// Verify the file is in the expected directory
			expectedBasePath := "content"
			if tt.expectedDir != "" {
				expectedBasePath = "content/" + tt.expectedDir
			}
			
			if !strings.Contains(filePath, expectedBasePath) {
				t.Errorf("File created in wrong directory, expected path containing '%s', got '%s'", expectedBasePath, filePath)
			}
			
			// For custom paths with a specific filename
			if strings.HasSuffix(tt.customPath, ".md") {
				expectedFilename := tt.customPath
				if !strings.HasSuffix(filePath, expectedFilename) {
					t.Errorf("File created with wrong filename, expected path ending with '%s', got '%s'", expectedFilename, filePath)
				}
			} else if tt.customPath == "" || !strings.HasSuffix(tt.customPath, "/") {
				// For default paths or directory paths without a specific filename
				expectedFilename := generateSlug(tt.title) + ".md"
				if !strings.HasSuffix(filePath, expectedFilename) {
					t.Errorf("File created with wrong filename, expected path ending with '%s', got '%s'", expectedFilename, filePath)
				}
			}

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("Failed to read created file: %v", err)
				return
			}

			// Verify content contains expected elements
			contentStr := string(content)

			// Check title
			if !strings.Contains(contentStr, "title: "+tt.title) {
				t.Errorf("Content doesn't contain expected title. Content: %s", contentStr)
			}

			// Check description if provided
			if tt.description != "" && !strings.Contains(contentStr, "description: "+tt.description) {
				t.Errorf("Content doesn't contain expected description. Content: %s", contentStr)
			}

			// Check draft status
			draftStr := "draft: false"
			if tt.draft {
				draftStr = "draft: true"
			}
			if !strings.Contains(contentStr, draftStr) {
				t.Errorf("Content doesn't contain expected draft status. Content: %s", contentStr)
			}

			// Check tags if any
			if len(tt.tags) > 0 {
				if !strings.Contains(contentStr, "tags:") {
					t.Errorf("Content doesn't contain tags section. Content: %s", contentStr)
				}

				for _, tag := range tt.tags {
					if !strings.Contains(contentStr, "  - "+tag) {
						t.Errorf("Content doesn't contain expected tag '%s'. Content: %s", tag, contentStr)
					}
				}
			}

			// Check content type-specific text
			var expectedText string
			if tt.contentType == PostType {
				expectedText = "Write your post content here."
			} else {
				expectedText = "Write your page content here."
			}
			if !strings.Contains(contentStr, expectedText) {
				t.Errorf("Content doesn't contain expected text '%s'. Content: %s", expectedText, contentStr)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		title string
		want  string
	}{
		{"Test Title", "test-title"},
		{"Hello, World!", "hello-world"},
		{"Multiple   Spaces", "multiple-spaces"},
		{"special-characters!@#$%^&*()", "special-characters"},
		{"", ""},
		{"123 Numbers", "123-numbers"},
		{"trailing -", "trailing"},
		{"- leading", "leading"},
		{"mixed CASE", "mixed-case"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got := generateSlug(tt.title)
			if got != tt.want {
				t.Errorf("generateSlug() = %v, want %v", got, tt.want)
			}
		})
	}
}
