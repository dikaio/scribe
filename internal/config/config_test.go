package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "scribe-config-default-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test configuration
	testCfg := DefaultConfig()

	// Save without an existing config file should create config.yml by default
	err = testCfg.Save(tempDir)
	if err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}

	// Verify YAML config file was created
	yamlPath := filepath.Join(tempDir, "config.yml")
	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		t.Fatalf("Default YAML config file was not created at %s", yamlPath)
	}

	// Verify JSON config file was NOT created
	jsonPath := filepath.Join(tempDir, "config.jsonc")
	if _, err := os.Stat(jsonPath); !os.IsNotExist(err) {
		t.Fatalf("JSON config file was created when it shouldn't have been: %s", jsonPath)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Check default values
	if cfg.Title != "Scribe" {
		t.Errorf("Expected default title to be 'Scribe', got %s", cfg.Title)
	}

	if cfg.Theme != "default" {
		t.Errorf("Expected default theme to be 'default', got %s", cfg.Theme)
	}

	if cfg.OutputDir != "public" {
		t.Errorf("Expected default output directory to be 'public', got %s", cfg.OutputDir)
	}

	if !cfg.TrailingSlash {
		t.Errorf("Expected default TrailingSlash to be true, got %t", cfg.TrailingSlash)
	}
}

func TestSaveAndLoadYAMLConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "scribe-config-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test configuration
	testCfg := Config{
		Title:         "Test Scribe",
		BaseURL:       "http://test.example.com/",
		Theme:         "test-theme",
		Language:      "fr",
		ContentDir:    "test-content",
		LayoutDir:     "test-layouts",
		StaticDir:     "test-static",
		OutputDir:     "test-public",
		Author:        "Test Author",
		Description:   "Test Description",
		SummaryLength: 100,
		Tags:          []string{"test", "example"},
		TrailingSlash: false, // Set to false to test saving/loading this value
	}

	// Create a YAML config file manually
	configPath := filepath.Join(tempDir, "config.yml")
	yamlData, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create YAML config file: %v", err)
	}
	yamlData.Close()

	// Save configuration (will use YAML since config.yml exists)
	err = testCfg.Save(tempDir)
	if err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}

	// Verify config file was created with correct extension
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("YAML config file was not created at %s", configPath)
	}

	// Load configuration
	loadedCfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Verify loaded configuration matches original
	if loadedCfg.Title != testCfg.Title {
		t.Errorf("Loaded title %s doesn't match original %s", loadedCfg.Title, testCfg.Title)
	}
	if loadedCfg.BaseURL != testCfg.BaseURL {
		t.Errorf("Loaded BaseURL %s doesn't match original %s", loadedCfg.BaseURL, testCfg.BaseURL)
	}
	if loadedCfg.Theme != testCfg.Theme {
		t.Errorf("Loaded theme %s doesn't match original %s", loadedCfg.Theme, testCfg.Theme)
	}
	if loadedCfg.Language != testCfg.Language {
		t.Errorf("Loaded language %s doesn't match original %s", loadedCfg.Language, testCfg.Language)
	}
	if loadedCfg.ContentDir != testCfg.ContentDir {
		t.Errorf("Loaded ContentDir %s doesn't match original %s", loadedCfg.ContentDir, testCfg.ContentDir)
	}
	if loadedCfg.LayoutDir != testCfg.LayoutDir {
		t.Errorf("Loaded LayoutDir %s doesn't match original %s", loadedCfg.LayoutDir, testCfg.LayoutDir)
	}
	if loadedCfg.StaticDir != testCfg.StaticDir {
		t.Errorf("Loaded StaticDir %s doesn't match original %s", loadedCfg.StaticDir, testCfg.StaticDir)
	}
	if loadedCfg.OutputDir != testCfg.OutputDir {
		t.Errorf("Loaded OutputDir %s doesn't match original %s", loadedCfg.OutputDir, testCfg.OutputDir)
	}
	if loadedCfg.Author != testCfg.Author {
		t.Errorf("Loaded Author %s doesn't match original %s", loadedCfg.Author, testCfg.Author)
	}
	if loadedCfg.Description != testCfg.Description {
		t.Errorf("Loaded Description %s doesn't match original %s", loadedCfg.Description, testCfg.Description)
	}
	if loadedCfg.SummaryLength != testCfg.SummaryLength {
		t.Errorf("Loaded SummaryLength %d doesn't match original %d", loadedCfg.SummaryLength, testCfg.SummaryLength)
	}

	// Check tags length
	if len(loadedCfg.Tags) != len(testCfg.Tags) {
		t.Errorf("Loaded Tags count %d doesn't match original %d", len(loadedCfg.Tags), len(testCfg.Tags))
	} else {
		// Check individual tags
		for i, tag := range testCfg.Tags {
			if loadedCfg.Tags[i] != tag {
				t.Errorf("Loaded tag %s doesn't match original %s", loadedCfg.Tags[i], tag)
			}
		}
	}

	// Check trailing slash setting
	if loadedCfg.TrailingSlash != testCfg.TrailingSlash {
		t.Errorf("Loaded TrailingSlash %t doesn't match original %t", loadedCfg.TrailingSlash, testCfg.TrailingSlash)
	}
}

func TestSaveAndLoadJSONConfig(t *testing.T) {
	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "scribe-config-json-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test configuration
	testCfg := Config{
		Title:         "Test Scribe",
		BaseURL:       "http://test.example.com/",
		Theme:         "test-theme",
		Language:      "fr",
		ContentDir:    "test-content",
		LayoutDir:     "test-layouts",
		StaticDir:     "test-static",
		OutputDir:     "test-public",
		Author:        "Test Author",
		Description:   "Test Description",
		SummaryLength: 100,
		Tags:          []string{"test", "example"},
		TrailingSlash: false, // Set to false to test saving/loading this value
	}

	// Create a JSON config file manually
	configPath := filepath.Join(tempDir, "config.jsonc")
	jsonData, err := os.Create(configPath)
	if err != nil {
		t.Fatalf("Failed to create JSON config file: %v", err)
	}
	jsonData.Close()

	// Save configuration (will use JSON since config.jsonc exists)
	err = testCfg.Save(tempDir)
	if err != nil {
		t.Fatalf("Failed to save configuration: %v", err)
	}

	// Verify config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf("JSON config file was not created at %s", configPath)
	}

	// Load configuration
	loadedCfg, err := LoadConfig(tempDir)
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	// Verify loaded configuration matches original
	if loadedCfg.Title != testCfg.Title {
		t.Errorf("Loaded title %s doesn't match original %s", loadedCfg.Title, testCfg.Title)
	}
	if loadedCfg.BaseURL != testCfg.BaseURL {
		t.Errorf("Loaded BaseURL %s doesn't match original %s", loadedCfg.BaseURL, testCfg.BaseURL)
	}
	if loadedCfg.Theme != testCfg.Theme {
		t.Errorf("Loaded theme %s doesn't match original %s", loadedCfg.Theme, testCfg.Theme)
	}
	if loadedCfg.Language != testCfg.Language {
		t.Errorf("Loaded language %s doesn't match original %s", loadedCfg.Language, testCfg.Language)
	}
	if loadedCfg.ContentDir != testCfg.ContentDir {
		t.Errorf("Loaded ContentDir %s doesn't match original %s", loadedCfg.ContentDir, testCfg.ContentDir)
	}
	if loadedCfg.LayoutDir != testCfg.LayoutDir {
		t.Errorf("Loaded LayoutDir %s doesn't match original %s", loadedCfg.LayoutDir, testCfg.LayoutDir)
	}
	if loadedCfg.StaticDir != testCfg.StaticDir {
		t.Errorf("Loaded StaticDir %s doesn't match original %s", loadedCfg.StaticDir, testCfg.StaticDir)
	}
	if loadedCfg.OutputDir != testCfg.OutputDir {
		t.Errorf("Loaded OutputDir %s doesn't match original %s", loadedCfg.OutputDir, testCfg.OutputDir)
	}
	if loadedCfg.Author != testCfg.Author {
		t.Errorf("Loaded Author %s doesn't match original %s", loadedCfg.Author, testCfg.Author)
	}
	if loadedCfg.Description != testCfg.Description {
		t.Errorf("Loaded Description %s doesn't match original %s", loadedCfg.Description, testCfg.Description)
	}
	if loadedCfg.SummaryLength != testCfg.SummaryLength {
		t.Errorf("Loaded SummaryLength %d doesn't match original %d", loadedCfg.SummaryLength, testCfg.SummaryLength)
	}

	// Check tags length
	if len(loadedCfg.Tags) != len(testCfg.Tags) {
		t.Errorf("Loaded Tags count %d doesn't match original %d", len(loadedCfg.Tags), len(testCfg.Tags))
	} else {
		// Check individual tags
		for i, tag := range testCfg.Tags {
			if loadedCfg.Tags[i] != tag {
				t.Errorf("Loaded tag %s doesn't match original %s", loadedCfg.Tags[i], tag)
			}
		}
	}

	// Check trailing slash setting
	if loadedCfg.TrailingSlash != testCfg.TrailingSlash {
		t.Errorf("Loaded TrailingSlash %t doesn't match original %t", loadedCfg.TrailingSlash, testCfg.TrailingSlash)
	}
}
