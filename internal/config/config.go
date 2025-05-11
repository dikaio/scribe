package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the site configuration
type Config struct {
	Title         string   `json:"title" yaml:"title"`
	BaseURL       string   `json:"baseURL" yaml:"baseURL"`
	Theme         string   `json:"theme" yaml:"theme"`
	Language      string   `json:"language" yaml:"language"`
	ContentDir    string   `json:"contentDir" yaml:"contentDir"`
	LayoutDir     string   `json:"layoutDir" yaml:"layoutDir"`
	StaticDir     string   `json:"staticDir" yaml:"staticDir"`
	OutputDir     string   `json:"outputDir" yaml:"outputDir"`
	Author        string   `json:"author" yaml:"author"`
	Description   string   `json:"description" yaml:"description"`
	SummaryLength int      `json:"summaryLength" yaml:"summaryLength"`
	Tags          []string `json:"tags" yaml:"tags"`
	TrailingSlash bool     `json:"trailingSlash" yaml:"trailingSlash"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Title:         "Scribe",
		BaseURL:       "http://example.com/",
		Theme:         "default",
		Language:      "en",
		ContentDir:    "content",
		LayoutDir:     "layouts",
		StaticDir:     "static",
		OutputDir:     "public",
		Author:        "",
		Description:   "",
		SummaryLength: 70,
		Tags:          []string{},
		TrailingSlash: true, // Default to trailing slashes for backward compatibility
	}
}

// determineConfigType determines the file type (YAML or JSON) based on extension
func determineConfigType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yml", ".yaml":
		return "yaml"
	case ".json", ".jsonc":
		return "json"
	default:
		// Default to JSON for backward compatibility
		return "json"
	}
}

// findConfigFile tries to find a configuration file in the site path
func findConfigFile(sitePath string) (string, error) {
	// Try YAML files first (preferred)
	yamlPaths := []string{
		filepath.Join(sitePath, "config.yml"),
	}

	// Then try JSON files (for backward compatibility)
	jsonPaths := []string{
		filepath.Join(sitePath, "config.jsonc"),
		filepath.Join(sitePath, "config.json"),
	}

	// Check for files in preferred order
	for _, path := range append(yamlPaths, jsonPaths...) {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// If no config file found, return default path (will be created)
	return filepath.Join(sitePath, "config.yml"), nil
}

// LoadConfig loads the site configuration from a file
func LoadConfig(sitePath string) (Config, error) {
	config := DefaultConfig()

	configPath, err := findConfigFile(sitePath)
	if err != nil {
		return config, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	// Determine format based on file extension
	configType := determineConfigType(configPath)
	switch configType {
	case "yaml":
		err = yaml.Unmarshal(data, &config)
	case "json":
		err = json.Unmarshal(data, &config)
	default:
		return config, errors.New("unsupported configuration format")
	}

	return config, err
}

// Save writes the configuration to a file
func (c Config) Save(sitePath string) error {
	// Get existing config path or use default
	configPath, err := findConfigFile(sitePath)
	if err != nil {
		// Use default config path with YAML format
		configPath = filepath.Join(sitePath, "config.yml")
	}

	// Determine format based on file extension
	configType := determineConfigType(configPath)

	var data []byte
	switch configType {
	case "yaml":
		data, err = yaml.Marshal(c)
	case "json":
		data, err = json.MarshalIndent(c, "", "  ")
	default:
		return errors.New("unsupported configuration format")
	}

	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
