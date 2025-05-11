package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRouting(t *testing.T) {
	// Create a temporary directory structure that mimics the output of a Scribe build
	tempDir, err := os.MkdirTemp("", "routing-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the directory structure for different types of pages
	createTestSiteStructure(t, tempDir)

	// Create the file server using the temp directory
	fileServer := http.FileServer(http.Dir(tempDir))

	// Define test cases for different routes
	tests := []struct {
		name     string
		path     string
		expected int
	}{
		{"Root path", "/", http.StatusOK},
		{"About page", "/about/", http.StatusOK},
		{"Privacy page", "/privacy/", http.StatusOK},
		{"Posts listing", "/posts/", http.StatusOK},
		{"Single post", "/posts/article-one/", http.StatusOK},
		{"Another post", "/posts/article-two/", http.StatusOK},
		{"Custom section", "/news/", http.StatusOK},
		{"Custom section article", "/news/latest/", http.StatusOK},
		
		// Test redirects for directory paths without trailing slash
		{"About redirect", "/about", http.StatusMovedPermanently},
		{"Posts redirect", "/posts", http.StatusMovedPermanently},
		
		// Test non-existent paths
		{"Non-existent page", "/nonexistent/", http.StatusNotFound},
		{"Non-existent post", "/posts/nonexistent/", http.StatusNotFound},
	}

	// Run the tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			
			fileServer.ServeHTTP(w, req)
			
			if w.Code != tc.expected {
				t.Errorf("Expected status code %d for path %s, got %d", 
					tc.expected, tc.path, w.Code)
			}
		})
	}
}

func createTestSiteStructure(t *testing.T, root string) {
	// Define directory structure and files to create
	structure := map[string]string{
		// Root/home page
		"index.html": "<html><body>Home Page</body></html>",
		
		// Regular pages
		"about/index.html":   "<html><body>About Page</body></html>",
		"privacy/index.html": "<html><body>Privacy Page</body></html>",
		
		// Post listing and individual posts
		"posts/index.html":         "<html><body>Posts Listing</body></html>",
		"posts/article-one/index.html": "<html><body>Article One</body></html>",
		"posts/article-two/index.html": "<html><body>Article Two</body></html>",
		
		// Custom section
		"news/index.html":        "<html><body>News Section</body></html>",
		"news/latest/index.html": "<html><body>Latest News</body></html>",
	}
	
	// Create all directories and files
	for path, content := range structure {
		fullPath := filepath.Join(root, path)
		
		// Create directory
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		
		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}
}

// TestDirectAccessRouting tests that files are properly served when accessed directly
func TestDirectAccessRouting(t *testing.T) {
	// Create a temporary directory for direct file access tests
	tempDir, err := os.MkdirTemp("", "direct-access-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create static files
	staticFiles := map[string]string{
		"styles.css":  "body { color: black; }",
		"script.js":   "console.log('Hello');",
		"image.png":   "fake-image-content",
		"robots.txt":  "User-agent: *\nDisallow: /admin/",
		"sitemap.xml": "<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\"><url><loc>http://example.com/</loc></url></urlset>",
	}

	for filename, content := range staticFiles {
		path := filepath.Join(tempDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", path, err)
		}
	}

	// Create the file server
	fileServer := http.FileServer(http.Dir(tempDir))

	// Test cases for direct file access
	tests := []struct {
		name         string
		path         string
		expectedCode int
		contentType  string
	}{
		{"CSS file", "/styles.css", http.StatusOK, "text/css"},
		{"JavaScript file", "/script.js", http.StatusOK, "text/javascript"},
		{"Image file", "/image.png", http.StatusOK, "image/png"},
		{"Robots.txt", "/robots.txt", http.StatusOK, "text/plain"},
		{"Sitemap.xml", "/sitemap.xml", http.StatusOK, "application/xml"},
		{"Non-existent file", "/nonexistent.txt", http.StatusNotFound, ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			
			fileServer.ServeHTTP(w, req)
			
			if w.Code != tc.expectedCode {
				t.Errorf("Expected status code %d for path %s, got %d", 
					tc.expectedCode, tc.path, w.Code)
			}
			
			// Check content type if we expect a successful response
			if tc.expectedCode == http.StatusOK && tc.contentType != "" {
				contentType := w.Header().Get("Content-Type")
				if contentType == "" || contentType[:len(tc.contentType)] != tc.contentType {
					t.Errorf("Expected Content-Type to start with %q for path %s, got %q", 
						tc.contentType, tc.path, contentType)
				}
			}
		})
	}
}

// TestNestedDirectoryRouting verifies that nested directory structures work correctly
func TestNestedDirectoryRouting(t *testing.T) {
	// Create a temporary directory for nested directory tests
	tempDir, err := os.MkdirTemp("", "nested-directory-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Define a deeply nested structure
	nestedStructure := map[string]string{
		"categories/tech/programming/golang/index.html": "<html><body>Go Programming</body></html>",
		"categories/tech/programming/python/index.html": "<html><body>Python Programming</body></html>",
		"categories/tech/hardware/index.html":          "<html><body>Hardware</body></html>",
		"categories/lifestyle/cooking/recipes/pasta/index.html": "<html><body>Pasta Recipes</body></html>",
		"categories/index.html": "<html><body>All Categories</body></html>",
	}

	// Create all directories and files
	for path, content := range nestedStructure {
		fullPath := filepath.Join(tempDir, path)
		
		// Create directory
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		
		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file %s: %v", fullPath, err)
		}
	}

	// Create the file server
	fileServer := http.FileServer(http.Dir(tempDir))

	// Test cases for nested directories
	tests := []struct {
		name         string
		path         string
		expectedCode int
	}{
		{"Categories root", "/categories/", http.StatusOK},
		{"Tech programming golang", "/categories/tech/programming/golang/", http.StatusOK},
		{"Tech programming python", "/categories/tech/programming/python/", http.StatusOK},
		{"Tech hardware", "/categories/tech/hardware/", http.StatusOK},
		{"Lifestyle cooking recipes pasta", "/categories/lifestyle/cooking/recipes/pasta/", http.StatusOK},
		
		// Redirects for paths without trailing slash
		{"Categories redirect", "/categories", http.StatusMovedPermanently},
		{"Nested redirect", "/categories/tech/programming/golang", http.StatusMovedPermanently},
		
		// Non-existent paths
		{"Non-existent category", "/categories/nonexistent/", http.StatusNotFound},
		{"Partial path", "/categories/tech/nonexistent/", http.StatusNotFound},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()
			
			fileServer.ServeHTTP(w, req)
			
			if w.Code != tc.expectedCode {
				t.Errorf("Expected status code %d for path %s, got %d", 
					tc.expectedCode, tc.path, w.Code)
			}
		})
	}
}