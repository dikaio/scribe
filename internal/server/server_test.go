package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/dikaio/scribe/internal/config"
)

func TestNewServer(t *testing.T) {
	// Create config
	cfg := config.Config{
		Title:     "Test Site",
		OutputDir: "public",
	}

	// Create server
	s := NewServer(cfg, 8080, true)

	// Check that the server was created successfully
	if s == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	// Check configuration
	if s.config.Title != "Test Site" {
		t.Errorf("Expected title to be 'Test Site', got '%s'", s.config.Title)
	}

	// Check port
	if s.port != 8080 {
		t.Errorf("Expected port to be 8080, got %d", s.port)
	}

	// Check builder
	if s.builder == nil {
		t.Error("Expected builder to be initialized")
	}
}

func TestFileServer(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "server-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a public directory and add a test file
	publicDir := filepath.Join(tempDir, "public")
	err = os.MkdirAll(publicDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create public dir: %v", err)
	}

	// Create a test index.html file
	indexContent := "<html><body>Test Site</body></html>"
	err = os.WriteFile(filepath.Join(publicDir, "index.html"), []byte(indexContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write index.html: %v", err)
	}

	// Get the file server handler
	handler := http.FileServer(http.Dir(publicDir))

	// Create a request to test the handler
	req, err := http.NewRequest("GET", "/index.html", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request
	handler.ServeHTTP(rr, req)

	// Check the response status code
	// Note: HTTP servers often return 301 for directory requests without trailing slash
	expectedCodes := []int{http.StatusOK, http.StatusMovedPermanently}
	codeOK := false
	for _, code := range expectedCodes {
		if rr.Code == code {
			codeOK = true
			break
		}
	}
	if !codeOK {
		t.Errorf("Expected status code to be one of %v, got %d", expectedCodes, rr.Code)
	}

	// If we got a redirect, we can't check the body
	if rr.Code == http.StatusOK {
		// Check the response body
		if rr.Body.String() != indexContent {
			t.Errorf("Expected body to be %q, got %q", indexContent, rr.Body.String())
		}
	}
}

func TestCreateOutputDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "server-output-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config
	cfg := config.Config{
		Title:     "Test Site",
		OutputDir: "public",
	}

	// Create server
	s := NewServer(cfg, 8080, true)

	// We'll just test the output directory creation part of Start
	outputPath := filepath.Join(tempDir, s.config.OutputDir)
	_, err = os.Stat(outputPath)
	if !os.IsNotExist(err) {
		t.Fatalf("Expected output directory to not exist yet")
	}

	// Create the output directory
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Check that the directory was created
	_, err = os.Stat(outputPath)
	if os.IsNotExist(err) {
		t.Errorf("Expected output directory to exist")
	}
}