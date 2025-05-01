package content

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// FrontMatter represents the metadata at the beginning of content files
type FrontMatter struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Tags        []string  `json:"tags"`
	Draft       bool      `json:"draft"`
	Layout      string    `json:"layout"`
	Slug        string    `json:"slug"`
}

// ParseFrontMatter extracts and parses front matter from content
func ParseFrontMatter(content []byte) (FrontMatter, []byte, error) {
	var frontMatter FrontMatter

	// Check if the content has front matter (starts with ---)
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return frontMatter, content, errors.New("no front matter found")
	}

	// Find the end of front matter
	parts := bytes.SplitN(content[4:], []byte("---\n"), 2)
	if len(parts) != 2 {
		return frontMatter, content, errors.New("invalid front matter format")
	}

	rawYAML := parts[0]
	bodyContent := parts[1]

	// Convert YAML to JSON-compatible format
	jsonData, err := yamlToJSON(rawYAML)
	if err != nil {
		return frontMatter, content, err
	}

	// Parse front matter
	err = json.Unmarshal(jsonData, &frontMatter)
	return frontMatter, bodyContent, err
}

// Simple YAML to JSON converter (limited to basic front matter needs)
func yamlToJSON(yamlData []byte) ([]byte, error) {
	lines := bytes.Split(yamlData, []byte("\n"))
	jsonMap := make(map[string]interface{})

	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
			continue
		}

		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}

		key := string(bytes.TrimSpace(parts[0]))
		value := bytes.TrimSpace(parts[1])

		// Handle lists
		if bytes.HasPrefix(value, []byte("-")) {
			var listItems []string
			listItems = append(listItems, strings.TrimSpace(string(value[1:])))
			jsonMap[key] = listItems
		} else {
			jsonMap[key] = string(value)
		}
	}

	return json.Marshal(jsonMap)
}
