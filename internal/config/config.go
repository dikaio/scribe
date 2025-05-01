package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config represents the site configuration
type Config struct {
	Title         string   `json:"title"`
	BaseURL       string   `json:"baseURL"`
	Theme         string   `json:"theme"`
	Language      string   `json:"language"`
	ContentDir    string   `json:"contentDir"`
	LayoutDir     string   `json:"layoutDir"`
	StaticDir     string   `json:"staticDir"`
	OutputDir     string   `json:"outputDir"`
	Author        string   `json:"author"`
	Description   string   `json:"description"`
	SummaryLength int      `json:"summaryLength"`
	Tags          []string `json:"tags"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Title:         "My Scribes Site",
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
	}
}

// LoadConfig loads the site configuration from a file
func LoadConfig(sitePath string) (Config, error) {
	config := DefaultConfig()

	configPath := filepath.Join(sitePath, "config.jsonc")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(data, &config)
	return config, err
}

// Save writes the configuration to a file
func (c Config) Save(sitePath string) error {
	configPath := filepath.Join(sitePath, "config.jsonc")
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
