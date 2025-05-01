# Scribes: Product Requirements Document

## 1. Introduction

### 1.1 Purpose
This document outlines the requirements for Scribes, a lightweight static site generator built using only the Go standard library. Scribes aims to provide a simple alternative to existing Go static site generators while maintaining core functionality needed for small to medium websites.

### 1.2 Product Scope
Scribes converts Markdown content with front matter into HTML websites using customizable templates. It targets developers and content creators who want a minimalist, dependency-free tool for generating static websites.

### 1.3 Definitions and Acronyms
- **SSG**: Static Site Generator
- **Markdown**: Lightweight markup language with plain text formatting syntax
- **Front Matter**: Metadata at the beginning of content files
- **Template**: HTML files with placeholders for dynamic content
- **Permalink**: Permanent URL for content

## 2. Product Overview

### 2.1 Product Perspective
Scribes is a standalone tool that takes content written in Markdown and converts it to a static website using templates. Unlike more complex SSGs, Scribes aims for simplicity and as close to zero external dependencies as possible utilizing the Go standard library.

### 2.2 Product Features
- Markdown to HTML conversion
- YAML front matter parsing
- Templating system using Go's html/template
- Content organization (posts, pages, tags)
- Basic theme support
- Development server with auto-rebuild
- Command-line interface for site operations (developer mode)
- Web interface for site operations (writer mode)

### 2.3 User Classes and Characteristics
- **Developers**: Familiar with command-line tools and interested in simple site generation
- **Content Creators**: Writers who prefer Markdown for content creation
- **Small Business Owners**: Looking for simple website solutions
- **Open Source Maintainers**: Need project documentation sites

### 2.4 Operating Environment
- Any system capable of running Go applications
- No external dependencies required
- Cross-platform compatibility (Windows, macOS, Linux)

## 3. System Features and Requirements

### 3.1 Content Management

#### 3.1.1 Markdown Processing
- **Requirement**: Convert Markdown syntax to HTML
- **Description**: Support common Markdown elements including headers, lists, links, images, code blocks, emphasis, and blockquotes
- **Priority**: High

#### 3.1.2 Front Matter Support
- **Requirement**: Parse YAML front matter from content files
- **Description**: Extract metadata like title, date, tags, and draft status from the beginning of content files
- **Priority**: High

#### 3.1.3 Content Types
- **Requirement**: Support both posts and regular pages
- **Description**: Differentiate between chronological content (posts) and static content (pages)
- **Priority**: Medium

#### 3.1.4 Draft Mode
- **Requirement**: Allow content to be marked as drafts
- **Description**: Exclude draft content from production builds while including it in development mode
- **Priority**: Medium

### 3.2 Site Structure

#### 3.2.1 Directory Organization
- **Requirement**: Provide a clear, intuitive directory structure
- **Description**: Organize content, layouts, and static files in separate directories
- **Priority**: High

#### 3.2.2 Content Organization
- **Requirement**: Support organization of content into logical sections
- **Description**: Allow posts to be organized by date and category/tags
- **Priority**: Medium

#### 3.2.3 URL Structure
- **Requirement**: Generate clean, customizable URLs
- **Description**: Create permalinks based on content hierarchy and front matter
- **Priority**: Medium

### 3.3 Templating System

#### 3.3.1 Layout Templates
- **Requirement**: Provide a base template system
- **Description**: Support template inheritance with base templates and content blocks
- **Priority**: High

#### 3.3.2 Page Templates
- **Requirement**: Support different templates for different content types
- **Description**: Allow specific templates for single posts, list pages, home page, etc.
- **Priority**: High

#### 3.3.3 Theme Support
- **Requirement**: Allow for swappable themes
- **Description**: Support theme directories with layouts and static assets
- **Priority**: Medium

### 3.4 Build System

#### 3.4.1 Static Site Generation
- **Requirement**: Convert all content to a static site
- **Description**: Process all content and templates to generate a complete static site
- **Priority**: High

#### 3.4.2 Asset Handling
- **Requirement**: Copy static assets to the output directory
- **Description**: Include CSS, JavaScript, images, and other static files in the final build
- **Priority**: High

#### 3.4.3 Incremental Builds
- **Requirement**: Support rebuilding only changed files
- **Description**: Detect file changes and rebuild only affected pages for faster development
- **Priority**: Low

### 3.5 Development Tools

#### 3.5.1 Live Preview Server
- **Requirement**: Provide a development server
- **Description**: Serve the site locally with automatic rebuilding on file changes
- **Priority**: High

#### 3.5.2 Site Initialization
- **Requirement**: Easy site creation command
- **Description**: Provide a command to scaffold a new site with default structure and example content
- **Priority**: Medium

#### 3.5.3 Content Creation
- **Requirement**: Simplify new content creation
- **Description**: Provide a command to create new posts/pages with proper front matter
- **Priority**: Medium

## 4. Technical Requirements

### 4.1 Standard Library Only
- **Requirement**: Use only Go standard library packages
- **Description**: Avoid external dependencies to maintain simplicity and stability
- **Priority**: High

### 4.2 Cross-Platform Support
- **Requirement**: Support major operating systems
- **Description**: Ensure functionality on Windows, macOS, and Linux
- **Priority**: Medium

### 4.3 Performance
- **Requirement**: Fast build times
- **Description**: Generate sites quickly, even for medium-sized content collections
- **Priority**: Medium

### 4.4 Command-Line Interface
- **Requirement**: Provide a simple, intuitive CLI
- **Description**: Support commands for common operations with clear help documentation
- **Priority**: High

### 4.5 Web Interface
- **Requirement**: Provide a web interface for site operations
- **Description**: Support commands for common operations with clear help documentation
- **Priority**: Medium

## 5. User Interface

### 5.1 Command-Line Commands

#### 5.1.1 Build Command
- **Command**: `scribe build [path]`
- **Description**: Build the site
- **Options**: Optional path to build the site in (defaults to current directory)

#### 5.1.2 Serve Command
- **Command**: `scribe serve [path]`
- **Description**: Start a development server with live reload
- **Options**: Optional path to serve the site from (defaults to current directory)

#### 5.1.3 Console Command
- **Command**: `scribe console`
- **Description**: Start the console

#### 5.1.4 Test Command
- **Command**: `scribe test`
- **Description**: Run the tests

#### 5.1.5 New Site Command
- **Command**: `scribe new site [name]`
- **Description**: Create a new site with default structure and sample content
- **Options**: Optional name to create site in directory [name]

#### 5.1.6 New Post Command
- **Command**: `scribe new post [title]`
- **Description**: Create a new post
- **Options**: Optional title to create the post with

#### 5.1.7 New Page Command
- **Command**: `scribe new page [title]`
- **Description**: Create a new page
- **Options**: Optional title to create the page with

#### 5.1.8 New Theme Command
- **Command**: `scribe new theme [name]`
- **Description**: Create a new theme
- **Options**: Optional name to create the theme in

#### 5.1.9 New Plugin Command
- **Command**: `scribe new plugin [name]`
- **Description**: Create a new plugin
- **Options**: Optional name to create the plugin in

#### 5.1.10 New Partial Command
- **Command**: `scribe new partial [name]`
- **Description**: Create a new partial
- **Options**: Optional name to create the partial in

## 6. Directory Structure

```
[name]/
├── config.jsonc       # Site configuration
├── content/           # Content files (Markdown)
│   ├── posts/         # Blog posts
│   └── *.md           # Regular pages
├── layouts/           # Template layouts
│   ├── base.html      # Base template
│   ├── single.html    # Post/page template
│   ├── list.html      # List template
│   └── home.html      # Homepage template
├── static/            # Static files (copied as-is)
└── themes/            # Site themes
    ├── layouts/   # Theme layouts
    └── static/    # Theme static files
```

## 7. Content Format

### 7.1 Front Matter Example

```markdown
---
title: My First Post
description: A short description
date: 2025-05-01T08:23:09-07:00
tags: 
  - tag1
  - tag2
  - tag3
draft: false
---

# Main Content Here

This is the body of the post written in Markdown.
```

## 8. Implementation Plan

### 8.1 Phase 1: Core Functionality
- Markdown to HTML conversion
- Front matter parsing
- Basic templating
- Build command

### 8.2 Phase 2: Enhanced Features
- Theme support
- Tags and categories
- Template inheritance
- Asset handling

### 8.3 Phase 3: Development Tools
- Development server
- Live reload
- Content creation helpers
- Documentation

## 9. Success Metrics

### 9.1 Performance Benchmarks
- Site with 100 pages builds in under 5 seconds
- Memory usage remains under 100MB during build

### 9.2 User Adoption Metrics
- Number of GitHub stars
- Number of downloads
- Community contributions

## 10. Appendix

### 10.1 Comparison with Existing Solutions

| Feature         | Scribes             | Hugo         | Jekyll         |
| --------------- | ------------------- | ------------ | -------------- |
| Dependencies    | None (std lib only) | Few          | Ruby ecosystem |
| Build Speed     | Very Fast           | Very Fast    | Moderate       |
| Learning Curve  | Low                 | Medium       | Medium         |
| Templating      | Go templates        | Go templates | Liquid         |
| Extensibility   | Limited             | High         | High           |
| Theme Ecosystem | Small               | Large        | Large          |

### 10.2 Future Considerations
- Plugin system (while maintaining minimal dependencies)
- Multilingual support
- Taxonomy pages beyond tags
- Image processing capabilities
- RSS/Atom feed generation