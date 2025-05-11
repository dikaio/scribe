package content

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestURLGeneration(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "url-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test content structure
	contentDir := filepath.Join(tempDir, "content")
	postsDir := filepath.Join(contentDir, "posts")
	pagesDir := filepath.Join(contentDir, "pages")
	nestedDir := filepath.Join(contentDir, "topics", "tech")

	// Create directories
	for _, dir := range []string{contentDir, postsDir, pagesDir, nestedDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create test files with front matter
	testFiles := map[string]string{
		filepath.Join(contentDir, "about.md"): `---
title: About Page
date: 2023-01-01T12:00:00Z
---
# About

This is the about page.`,

		filepath.Join(postsDir, "first-post.md"): `---
title: First Post
date: 2023-01-02T12:00:00Z
tags:
  - test
  - example
---
# First Post

This is the first post.`,

		filepath.Join(nestedDir, "golang.md"): `---
title: Go Programming
date: 2023-01-03T12:00:00Z
slug: golang-intro
---
# Go Programming

This is an introduction to Go programming.`,
	}

	for path, content := range testFiles {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	// Test cases for different URL configurations
	testCases := []struct {
		name          string
		filePath      string
		baseURL       string
		trailingSlash bool
		expectedURL   string
		expectedPerma string
	}{
		{
			name:          "About page with trailing slash",
			filePath:      filepath.Join(contentDir, "about.md"),
			baseURL:       "https://example.com",
			trailingSlash: true,
			expectedURL:   "about/",
			expectedPerma: "https://example.com/about/",
		},
		{
			name:          "About page without trailing slash",
			filePath:      filepath.Join(contentDir, "about.md"),
			baseURL:       "https://example.com",
			trailingSlash: false,
			expectedURL:   "about",
			expectedPerma: "https://example.com/about",
		},
		{
			name:          "Post with trailing slash",
			filePath:      filepath.Join(postsDir, "first-post.md"),
			baseURL:       "https://example.com",
			trailingSlash: true,
			expectedURL:   "posts/first-post/",
			expectedPerma: "https://example.com/posts/first-post/",
		},
		{
			name:          "Post without trailing slash",
			filePath:      filepath.Join(postsDir, "first-post.md"),
			baseURL:       "https://example.com",
			trailingSlash: false,
			expectedURL:   "posts/first-post",
			expectedPerma: "https://example.com/posts/first-post",
		},
		{
			name:          "Nested page with custom slug and trailing slash",
			filePath:      filepath.Join(nestedDir, "golang.md"),
			baseURL:       "https://example.com",
			trailingSlash: true,
			expectedURL:   "topics/tech/golang-intro/",
			expectedPerma: "https://example.com/topics/tech/golang-intro/",
		},
		{
			name:          "Nested page with custom slug without trailing slash",
			filePath:      filepath.Join(nestedDir, "golang.md"),
			baseURL:       "https://example.com",
			trailingSlash: false,
			expectedURL:   "topics/tech/golang-intro",
			expectedPerma: "https://example.com/topics/tech/golang-intro",
		},
		{
			name:          "Base URL with trailing slash",
			filePath:      filepath.Join(contentDir, "about.md"),
			baseURL:       "https://example.com/",
			trailingSlash: true,
			expectedURL:   "about/",
			expectedPerma: "https://example.com/about/",
		},
		{
			name:          "Base URL with trailing slash and URL without trailing slash",
			filePath:      filepath.Join(contentDir, "about.md"),
			baseURL:       "https://example.com/",
			trailingSlash: false,
			expectedURL:   "about",
			expectedPerma: "https://example.com/about",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			page, err := LoadPage(tc.filePath, tc.baseURL, tc.trailingSlash)
			if err != nil {
				t.Fatalf("Failed to load page: %v", err)
			}

			// Check URL format
			if page.URL != tc.expectedURL {
				t.Errorf("Expected URL to be %q, got %q", tc.expectedURL, page.URL)
			}

			// Check permalink format
			if page.Permalink != tc.expectedPerma {
				t.Errorf("Expected permalink to be %q, got %q", tc.expectedPerma, page.Permalink)
			}

			// Verify other properties as needed
			// Content should have been loaded correctly
			if !strings.Contains(page.Content, "This is") {
				t.Errorf("Expected content to contain 'This is', got %q", page.Content)
			}

			// HTML should have been generated
			if !strings.Contains(page.HTML, "<h1") {
				t.Errorf("Expected HTML to contain '<h1', got %q", page.HTML)
			}

			// Check if post detection works
			expectedIsPost := strings.Contains(tc.filePath, "/posts/")
			if page.IsPost != expectedIsPost {
				t.Errorf("Expected IsPost to be %v, got %v", expectedIsPost, page.IsPost)
			}
		})
	}
}