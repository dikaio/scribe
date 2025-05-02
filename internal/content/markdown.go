package content

import (
	"fmt"
	"regexp"
	"strings"
)

// MarkdownToHTML converts Markdown content to HTML
func MarkdownToHTML(markdown []byte) []byte {
	// Simple Markdown to HTML conversion
	html := string(markdown)

	// We'll first process the special blocks that we don't want to wrap in paragraphs
	// Process Code blocks
	codeBlockRe := regexp.MustCompile("```([\\s\\S]*?)```")
	codeBlocks := make(map[string]string)
	codeBlockCount := 0
	html = codeBlockRe.ReplaceAllStringFunc(html, func(match string) string {
		code := codeBlockRe.FindStringSubmatch(match)[1]
		placeholder := fmt.Sprintf("___CODE_BLOCK_%d___", codeBlockCount)
		codeBlocks[placeholder] = "<pre><code>" + code + "</code></pre>"
		codeBlockCount++
		return placeholder
	})

	// Process Headers
	headerRe := regexp.MustCompile(`(?m)^(#{1,6})\s+(.+)$`)
	headers := make(map[string]string)
	headerCount := 0
	html = headerRe.ReplaceAllStringFunc(html, func(match string) string {
		submatches := headerRe.FindStringSubmatch(match)
		level := len(submatches[1])
		text := strings.TrimSpace(submatches[2])
		placeholder := fmt.Sprintf("___HEADER_%d___", headerCount)
		headers[placeholder] = "<h" + string(rune('0'+level)) + ">" + text + "</h" + string(rune('0'+level)) + ">"
		headerCount++
		return placeholder
	})

	// Process regular inline elements
	// Bold
	boldRe := regexp.MustCompile(`\*\*(.+?)\*\*`)
	html = boldRe.ReplaceAllString(html, "<strong>$1</strong>")

	// Italic
	italicRe := regexp.MustCompile(`\*(.+?)\*`)
	html = italicRe.ReplaceAllString(html, "<em>$1</em>")

	// Links
	linkRe := regexp.MustCompile(`\[(.+?)\]\((.+?)\)`)
	html = linkRe.ReplaceAllString(html, "<a href=\"$2\">$1</a>")

	// Inline code
	inlineCodeRe := regexp.MustCompile("`([^`]+)`")
	html = inlineCodeRe.ReplaceAllString(html, "<code>$1</code>")

	// Process lists
	listItemRe := regexp.MustCompile(`(?m)^-\s+(.+)$`)
	listItems := make(map[string]string)
	listItemCount := 0
	html = listItemRe.ReplaceAllStringFunc(html, func(match string) string {
		submatches := listItemRe.FindStringSubmatch(match)
		placeholder := fmt.Sprintf("___LIST_ITEM_%d___", listItemCount)
		listItems[placeholder] = "<li>" + submatches[1] + "</li>"
		listItemCount++
		return placeholder
	})

	// Now wrap remaining text content in paragraphs
	// Split by newlines
	lines := strings.Split(html, "\n")
	for i, line := range lines {
		// Skip if line is empty or is a placeholder for headers, code blocks, or list items
		if line == "" || strings.HasPrefix(line, "___CODE_BLOCK_") || 
		   strings.HasPrefix(line, "___HEADER_") || 
		   strings.HasPrefix(line, "___LIST_ITEM_") {
			continue
		}
		
		// If not already wrapped with HTML tags, wrap in paragraph tags
		if !strings.HasPrefix(line, "<") {
			lines[i] = "<p>" + line + "</p>"
		}
	}
	html = strings.Join(lines, "\n")

	// Process list items and wrap in <ul>...</ul>
	listRe := regexp.MustCompile("(?s)((?:___LIST_ITEM_\\d+___\n?)+)")
	html = listRe.ReplaceAllStringFunc(html, func(match string) string {
		listHTML := "<ul>\n"
		listPlaceholders := strings.Split(strings.TrimSpace(match), "\n")
		for _, placeholder := range listPlaceholders {
			if item, ok := listItems[strings.TrimSpace(placeholder)]; ok {
				listHTML += item + "\n"
			}
		}
		listHTML += "</ul>"
		return listHTML
	})

	// Restore code blocks
	for placeholder, codeBlock := range codeBlocks {
		html = strings.Replace(html, placeholder, codeBlock, 1)
	}

	// Restore headers
	for placeholder, header := range headers {
		html = strings.Replace(html, placeholder, header, 1)
	}

	// Clean up empty paragraphs
	html = strings.ReplaceAll(html, "<p></p>", "")

	return []byte(html)
}