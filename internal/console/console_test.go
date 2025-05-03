package console

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
)

func setupTestConsole(t *testing.T) (*Console, string, func()) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "console-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	err = os.MkdirAll(contentDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create content dir: %v", err)
	}

	// Create a test markdown file (post)
	postContent := `---
title: Test Post
date: 2023-01-01T12:00:00Z
tags:
  - test
  - post
draft: false
---

# Test Post

This is a test post.
`
	err = os.WriteFile(filepath.Join(contentDir, "post.md"), []byte(postContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write post.md: %v", err)
	}

	// Create a test markdown file (page)
	pageContent := `---
title: About Page
date: 2023-01-02T12:00:00Z
draft: false
---

# About

This is an about page.
`
	err = os.WriteFile(filepath.Join(contentDir, "about.md"), []byte(pageContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write about.md: %v", err)
	}

	// Create config
	cfg := config.Config{
		Title:      "Test Site",
		BaseURL:    "http://example.com/",
		ContentDir: "content",
		Theme:      "default",
	}

	// Create console
	console := NewConsole(cfg, tempDir, 8090)

	// Override loadSiteStats for testing
	console.OverrideLoadSiteStats()

	// Initialize templates
	err = console.initTemplates()
	if err != nil {
		t.Fatalf("Failed to initialize templates: %v", err)
	}

	// Return console and cleanup function
	return console, tempDir, func() {
		console.RestoreLoadSiteStats()
		os.RemoveAll(tempDir)
	}
}

func TestNewConsole(t *testing.T) {
	cfg := config.Config{
		Title:  "Test Site",
		Theme:  "default",
		BaseURL: "http://example.com/",
	}

	// Create console
	c := NewConsole(cfg, "/path/to/site", 8090)

	// Check that the console was created successfully
	if c == nil {
		t.Fatal("Expected console to be created, got nil")
	}

	// Check configuration
	if c.config.Title != "Test Site" {
		t.Errorf("Expected title to be 'Test Site', got '%s'", c.config.Title)
	}

	// Check site path
	if c.sitePath != "/path/to/site" {
		t.Errorf("Expected site path to be '/path/to/site', got '%s'", c.sitePath)
	}

	// Check port
	if c.port != 8090 {
		t.Errorf("Expected port to be 8090, got %d", c.port)
	}
}

func TestInitTemplates(t *testing.T) {
	cfg := config.Config{
		Title:  "Test Site",
		Theme:  "default",
		BaseURL: "http://example.com/",
	}

	// Create console
	c := NewConsole(cfg, "/path/to/site", 8090)

	// Initialize templates
	err := c.initTemplates()
	if err != nil {
		t.Fatalf("Failed to initialize templates: %v", err)
	}

	// Check that templates were initialized
	if c.tmpl == nil {
		t.Error("Expected templates to be initialized")
	}
}

func TestHandleDashboard(t *testing.T) {
	// Setup test console
	c, _, cleanup := setupTestConsole(t)
	defer cleanup()

	// Create a request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	c.handleDashboard(rr, req)

	// Check response status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	// Check content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("Expected Content-Type text/html, got %s", contentType)
	}

	// Check response body contains expected elements
	body := rr.Body.String()
	expectedElements := []string{
		"Dashboard", 
		"Test Site",
		"Recent Posts",
		"Recent Pages",
		"Test Post",
		"About Page",
	}

	for _, element := range expectedElements {
		if !strings.Contains(body, element) {
			t.Errorf("Expected response to contain %q", element)
		}
	}

	// Test with invalid path
	req, _ = http.NewRequest("GET", "/invalid", nil)
	rr = httptest.NewRecorder()
	c.handleDashboard(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("Expected status NotFound for invalid path, got %v", rr.Code)
	}
}

func TestHandleContent(t *testing.T) {
	// Setup test console
	c, _, cleanup := setupTestConsole(t)
	defer cleanup()

	// Test with different content types
	testCases := []struct {
		name        string
		path        string
		expectedElements []string
	}{
		{
			name: "All content",
			path: "/content",
			expectedElements: []string{"Content", "Test Post", "About Page"},
		},
		{
			name: "Posts only",
			path: "/content?type=posts",
			expectedElements: []string{"Content", "Test Post"},
		},
		{
			name: "Pages only",
			path: "/content?type=pages",
			expectedElements: []string{"Content", "About Page"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.path, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			c.handleContent(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("Expected status OK, got %v", rr.Code)
			}

			body := rr.Body.String()
			for _, element := range tc.expectedElements {
				if !strings.Contains(body, element) {
					t.Errorf("Expected response to contain %q", element)
				}
			}
		})
	}
}

func TestHandleNewContent(t *testing.T) {
	// Setup test console
	c, _, cleanup := setupTestConsole(t)
	defer cleanup()

	// Test for post
	req, err := http.NewRequest("GET", "/content/new?type=post", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	c.handleNewContent(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "New Post") {
		t.Errorf("Expected response to contain 'New Post'")
	}

	// Test for page
	req, err = http.NewRequest("GET", "/content/new?type=page", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr = httptest.NewRecorder()
	c.handleNewContent(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	body = rr.Body.String()
	if !strings.Contains(body, "New Page") {
		t.Errorf("Expected response to contain 'New Page'")
	}

	// Test POST request (should redirect)
	req, err = http.NewRequest("POST", "/content/new", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr = httptest.NewRecorder()
	c.handleNewContent(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status SeeOther, got %v", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "/content" {
		t.Errorf("Expected redirect to /content, got %s", location)
	}
}

func TestHandleSettings(t *testing.T) {
	// Setup test console
	c, _, cleanup := setupTestConsole(t)
	defer cleanup()

	// Test GET request
	req, err := http.NewRequest("GET", "/settings", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	c.handleSettings(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %v", rr.Code)
	}

	body := rr.Body.String()
	expectedElements := []string{
		"Settings",
		"Site Title",
		"Base URL",
		"Test Site",
		"http://example.com/",
	}

	for _, element := range expectedElements {
		if !strings.Contains(body, element) {
			t.Errorf("Expected response to contain %q", element)
		}
	}

	// Test POST request (should redirect)
	req, err = http.NewRequest("POST", "/settings", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr = httptest.NewRecorder()
	c.handleSettings(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status SeeOther, got %v", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "/settings" {
		t.Errorf("Expected redirect to /settings, got %s", location)
	}
}

func TestHandleBuild(t *testing.T) {
	// Setup test console
	c, _, cleanup := setupTestConsole(t)
	defer cleanup()

	// Test request
	req, err := http.NewRequest("GET", "/build", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	c.handleBuild(rr, req)

	// Should redirect to dashboard
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected status SeeOther, got %v", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "/" {
		t.Errorf("Expected redirect to /, got %s", location)
	}
}

func TestGetRecentItems(t *testing.T) {
	// Create test items
	items := []content.Page{
		{Title: "Item 1", Date: parseDate(t, "2023-01-01T12:00:00Z")},
		{Title: "Item 2", Date: parseDate(t, "2023-01-02T12:00:00Z")},
		{Title: "Item 3", Date: parseDate(t, "2023-01-03T12:00:00Z")},
		{Title: "Item 4", Date: parseDate(t, "2023-01-04T12:00:00Z")},
		{Title: "Item 5", Date: parseDate(t, "2023-01-05T12:00:00Z")},
	}

	// Get recent items
	recent := getRecentItems(items, 3)

	// Check that we got the expected number of items
	if len(recent) != 3 {
		t.Errorf("Expected 3 items, got %d", len(recent))
	}

	// Check that items are sorted by date (newest first)
	if recent[0].Title != "Item 5" || recent[1].Title != "Item 4" || recent[2].Title != "Item 3" {
		t.Errorf("Items not sorted correctly: %v", recent)
	}

	// Test with fewer items than requested
	recent = getRecentItems(items[:2], 3)
	if len(recent) != 2 {
		t.Errorf("Expected 2 items, got %d", len(recent))
	}
}

// Helper function to parse time
func parseDate(t *testing.T, dateStr string) time.Time {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		t.Fatalf("Failed to parse date %q: %v", dateStr, err)
	}
	return date
}