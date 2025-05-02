package content

import (
	"testing"
	"bytes"
)

func TestMarkdownToHTML(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		{
			name:     "Headers",
			markdown: "# Header 1\n## Header 2",
			want:     "<h1>Header 1</h1>\n<h2>Header 2</h2>",
		},
		{
			name:     "Bold",
			markdown: "This is **bold** text",
			want:     "<p>This is <strong>bold</strong> text</p>",
		},
		{
			name:     "Italic",
			markdown: "This is *italic* text",
			want:     "<p>This is <em>italic</em> text</p>",
		},
		{
			name:     "Links",
			markdown: "This is a [link](https://example.com)",
			want:     "<p>This is a <a href=\"https://example.com\">link</a></p>",
		},
		{
			name:     "Lists",
			markdown: "- Item 1\n- Item 2",
			want:     "<ul>\n<li>Item 1</li>\n<li>Item 2</li>\n</ul>",
		},
		{
			name:     "Code blocks",
			markdown: "```\ncode block\n```",
			want:     "<pre><code>\ncode block\n</code></pre>",
		},
		{
			name:     "Inline code",
			markdown: "This is `inline code`",
			want:     "<p>This is <code>inline code</code></p>",
		},
		{
			name:     "Complex example with link",
			markdown: "# Title\n\nThis is a paragraph with a [link](https://example.com) and **bold** text.\n\n- List item with *italic*\n- Another item",
			want:     "<h1>Title</h1>\n\n<p>This is a paragraph with a <a href=\"https://example.com\">link</a> and <strong>bold</strong> text.</p>\n\n<ul>\n<li>List item with <em>italic</em></li>\n<li>Another item</li>\n</ul>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MarkdownToHTML([]byte(tt.markdown))
			if !bytes.Equal(got, []byte(tt.want)) {
				t.Errorf("MarkdownToHTML() = %q, want %q", got, tt.want)
			}
		})
	}
}