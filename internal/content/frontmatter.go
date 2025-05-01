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
	currentKey := ""
	inList := false
	listItems := []string{}

	for i := 0; i < len(lines); i++ {
		line := bytes.TrimSpace(lines[i])
		if len(line) == 0 || bytes.HasPrefix(line, []byte("#")) {
			continue
		}

		// If line starts with a dash (list item) and is indented, it's part of a list
		if bytes.HasPrefix(line, []byte("  -")) || bytes.HasPrefix(line, []byte("\t-")) {
			if !inList {
				inList = true
				listItems = []string{}
			}
			item := strings.TrimSpace(string(line[1:]))
			listItems = append(listItems, item)
			
			// If this is the last line or the next line isn't a list item, 
			// end the list and store it
			if i == len(lines)-1 || 
				!(bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("-")) || 
				  bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("  -")) || 
				  bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("\t-"))) {
				jsonMap[currentKey] = listItems
				inList = false
			}
			continue
		}

		// Regular key-value line
		parts := bytes.SplitN(line, []byte(":"), 2)
		if len(parts) != 2 {
			continue
		}

		key := string(bytes.TrimSpace(parts[0]))
		value := bytes.TrimSpace(parts[1])
		currentKey = key

		// If value is empty and the next line is a list item, this is a list declaration
		if len(value) == 0 && i+1 < len(lines) && 
			(bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("-")) || 
			 bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("  -")) || 
			 bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("\t-"))) {
			continue
		}

		// Handle single-line lists
		if bytes.HasPrefix(value, []byte("[")) && bytes.HasSuffix(value, []byte("]")) {
			// Extract items between brackets and split by comma
			items := bytes.TrimSpace(value[1 : len(value)-1])
			if len(items) > 0 {
				itemList := []string{}
				for _, item := range bytes.Split(items, []byte(",")) {
					itemList = append(itemList, string(bytes.TrimSpace(item)))
				}
				jsonMap[key] = itemList
			} else {
				jsonMap[key] = []string{}
			}
			continue
		}

		// Handle inline list item
		if bytes.HasPrefix(value, []byte("-")) {
			listItems = []string{strings.TrimSpace(string(value[1:]))}
			inList = true
			
			// Check if this is a single item list
			if i == len(lines)-1 || 
				!(bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("-")) || 
				  bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("  -")) || 
				  bytes.HasPrefix(bytes.TrimSpace(lines[i+1]), []byte("\t-"))) {
				jsonMap[key] = listItems
				inList = false
			}
			continue
		}

		// Convert special values
		strValue := string(value)
		switch strValue {
		case "true":
			jsonMap[key] = true
		case "false":
			jsonMap[key] = false
		case "null", "nil":
			jsonMap[key] = nil
		default:
			// Try to parse as number
			if len(strValue) > 0 {
				if strings.Contains(strValue, ".") {
					if f, err := json.Number(strValue).Float64(); err == nil {
						jsonMap[key] = f
						continue
					}
				} else {
					if i, err := json.Number(strValue).Int64(); err == nil {
						jsonMap[key] = i
						continue
					}
				}
			}
			// Default to string
			jsonMap[key] = strValue
		}
	}

	return json.Marshal(jsonMap)
}
