# References

This directory contains cloned repositories used for reference and inspiration. These external codebases help with ideation of elegant solutions for features to be implemented in the main project.

The contents of repositories cloned into this directory are intentionally git-ignored to avoid committing external code to this project.

## Planned Features

### HTML Sanitizer Module
Create a simplified HTML sanitizer inspired by bluemonday but:
- Using only the Go standard library
- More elegant and maintainable design
- Focused on core sanitization needs for Markdown-to-HTML conversion
- Support for configurable allowed tags and attributes
- Secure by default while remaining performant

### Enhanced Markdown Parser
Develop improved markdown parsing capabilities inspired by goldmark but with:
- Cleaner, more maintainable architecture
- Better extension points for custom syntax
- Simpler AST representation
- More efficient parsing algorithms
- Support for common extensions (tables, strikethrough, task lists)
- Focus on readability and maintainability over extreme optimization

### HTMX-inspired Server-Side Interactions
Implement server-side dynamic content capabilities inspired by HTMX but:
- Using only Go standard library for backend implementation
- Minimal JavaScript on the client side
- Simple, elegant API for declaring interactive elements
- Efficient partial page updates without full page reloads
- Clean separation between content and interactivity
- Focus on performance and reduced network overhead

## Usage

Clone reference repositories here:

```bash
cd references
git clone https://github.com/example/repo
```

Then use these codebases as reference when discussing implementations with Claude.