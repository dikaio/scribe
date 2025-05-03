# References Directory

## Purpose

This directory serves as a collection point for external repositories that contain implementations of features we're interested in adapting for our project. These repositories are kept separate from our main codebase and are used solely for reference and inspiration.

All contents in this directory (except this README) are git-ignored to prevent external code from being committed to our repository.

## Usage

Clone reference repositories into this directory:

```bash
cd references
git clone https://github.com/example/repo
```

These repositories will serve as the basis for discussions with Claude about creating more elegant implementations.

## Feature Roadmap

Below are features we plan to implement by studying existing libraries but reimagining them with improved architecture, better performance, and greater simplicity.

### 1. HTML Sanitizer Module

Inspired by: bluemonday

Our goals:
- Use only the Go standard library
- Create a more elegant and maintainable design
- Focus on core sanitization needs for Markdown-to-HTML conversion
- Support configurable allowed tags and attributes
- Ensure security by default while maintaining performance

### 2. Enhanced Markdown Parser

Inspired by: goldmark

Our goals:
- Design a cleaner, more maintainable architecture
- Create better extension points for custom syntax
- Implement a simpler AST representation
- Develop more efficient parsing algorithms
- Support common extensions (tables, strikethrough, task lists)
- Prioritize readability and maintainability over extreme optimization

### 3. HTMX-inspired Server-Side Interactions

Inspired by: HTMX

Our goals:
- Implement using only Go standard library for backend
- Require minimal JavaScript on the client side
- Design a simple, elegant API for declaring interactive elements
- Enable efficient partial page updates without full page reloads
- Maintain clean separation between content and interactivity
- Focus on performance and reduced network overhead

## Contributing

When adding a new reference repository:
1. Update this README with notes about which features you're interested in
2. Document key insights about the implementation in the reference repo
3. Keep notes about architectural decisions you want to improve upon