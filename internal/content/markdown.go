package content

import (
	"regexp"
	"strings"
)

// MarkdownToHTML converts Markdown content to HTML
func MarkdownToHTML(markdown []byte) []byte {
	// Simple Markdown to HTML conversion
	html := string(markdown)

	// Headers
	headerRe := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	html = headerRe.ReplaceAllStringFunc(html, func(match string) string {
		submatches := headerRe.FindStringSubmatch(match)
		level := len(submatches[1])
		text := strings.TrimSpace(submatches[2])
		return "<h" + string(rune('0'+level)) + ">" + text + "</h" + string(rune('0'+level)) + ">"
	})

	// Bold
	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
	html = boldRe.ReplaceAllString(html, "<strong>$1</strong>")

	// Italic
	italicRe := regexp.MustCompile(`\*(.+?)\*`)
	html = italicRe.ReplaceAllString(html, "<em>$1</em>")

	// Links
	linkRe := regexp.MustCompile(`\[(.+?)\]\((.+?)\)`)
	html = linkRe.ReplaceAllString(html, "<a href=\"$2\">$1</a>")

	// Lists
	listItemRe := regexp.MustCompile(`(?m)^-\s+(.+)$`)
	html = listItemRe.ReplaceAllString(html, "<li>$1</li>")

	// Group list items
	lines := strings.Split(html, "\n")
	result := []string{}
	inList := false

	for _, line := range lines {
		if strings.HasPrefix(line, "<li>") {
			if !inList {
				result = append(result, "<ul>")
				inList = true
			}
			result = append(result, line)
		} else {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			result = append(result, line)
		}
	}

	if inList {
		result = append(result, "</ul>")
	}

	html = strings.Join(result, "\n")

	// Code blocks
	codeBlockRe := regexp.MustCompile("```([\\s\\S]*?)```")
	html = codeBlockRe.ReplaceAllStringFunc(html, func(match string) string {
		code := codeBlockRe.FindStringSubmatch(match)[1]
		return "<pre><code>" + code + "</code></pre>"
	})

	// Inline code
	inlineCodeRe := regexp.MustCompile("`([^`]+)`")
	html = inlineCodeRe.ReplaceAllString(html, "<code>$1</code>")

	// Paragraphs
	paragraphRe := regexp.MustCompile(`(?m)^([^<].+)$`)
	html = paragraphRe.ReplaceAllString(html, "<p>$1</p>")

	// Clean up empty paragraphs
	html = strings.ReplaceAll(html, "<p></p>", "")

	return []byte(html)
}
