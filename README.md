# Scribes

A lightweight static site generator built in pure Go with zero external dependencies.

## Overview

Scribes converts Markdown content with YAML front matter into HTML websites using customizable templates. It's designed to be simple, fast, and dependency-free, using only the Go standard library.

## Features

- **Markdown to HTML conversion** - Write your content in Markdown, Scribes handles the rest
- **YAML front matter** - Add metadata to your content like title, date, tags, etc.
- **Templating system** - Use Go's html/template for layouts and themes
- **Live reload development server** - See changes as you make them
- **Command-line interface** - Simple commands for common operations
- **Web management console** - Manage your site through a browser interface
- **Zero dependencies** - Built using only the Go standard library

## Installation

### Prerequisites

- Go 1.16 or later

### Building from source

```bash
# Clone the repository
git clone https://github.com/yourusername/scribes.git
cd scribes

# Build the binary
go build -o scribe ./cmd/scribe

# Install globally (optional)
go install ./cmd/scribe
```

## Quick Start

### Create a new site

```bash
# Create a new site in the current directory
scribe new site mysite

# Navigate to the site directory
cd mysite

# Start the development server
scribe serve
```

Open your browser at [http://localhost:8080](http://localhost:8080) to see your new site.

### Create content

```bash
# Create a new post
scribe new post "My First Post"

# Create a new page
scribe new page "About Me"
```

### Build your site

```bash
# Build the site
scribe build
```

The built site will be in the `public` directory by default.

### Launch the web console

```bash
# Start the web-based management console
scribe console
```

Open your browser at [http://localhost:8090](http://localhost:8090) to access the console.

## Site Structure

```
mysite/
├── config.jsonc       # Site configuration
├── content/           # Content files (Markdown)
│   ├── posts/         # Blog posts
│   └── *.md           # Regular pages
├── layouts/           # Custom template layouts (optional)
├── static/            # Static files (copied as-is)
└── themes/            # Site themes
    └── default/       # Default theme
        ├── layouts/   # Theme layouts
        └── static/    # Theme static files
```

## Content Format

Scribes uses Markdown files with YAML front matter for content:

```markdown
---
title: My First Post
description: A short description
date: 2025-05-01T08:23:09-07:00
tags: 
  - tag1
  - tag2
draft: false
---

# Main Content Here

This is the body of the post written in Markdown.
```

## Configuration

Site configuration is stored in `config.jsonc`:

```json
{
  "title": "My Scribes Site",
  "baseURL": "http://example.com/",
  "theme": "default",
  "language": "en",
  "contentDir": "content",
  "layoutDir": "layouts",
  "staticDir": "static",
  "outputDir": "public",
  "author": "Your Name",
  "description": "Your site description",
  "summaryLength": 70
}
```

## Commands

| Command                   | Description                                 |
| ------------------------- | ------------------------------------------- |
| `scribe build [path]`     | Build the site                              |
| `scribe serve [path]`     | Start a development server with live reload |
| `scribe console [path]`   | Start the web management console            |
| `scribe new site [name]`  | Create a new site                           |
| `scribe new post [title]` | Create a new blog post                      |
| `scribe new page [title]` | Create a new page                           |
| `scribe help`             | Show help information                       |

## Templating

Scribes uses Go's `html/template` package for templating. Base templates include:

- `base.html` - The base template that defines the overall structure
- `single.html` - Template for individual posts/pages
- `list.html` - Template for content lists (tags, etc.)
- `home.html` - The homepage template

Template data includes:

- `.Site` - Site configuration
- `.Page` - Current page information
- `.Content` - Rendered page content
- `.Pages` - List of pages (for list/home templates)

## Customization

### Custom Themes

To create a custom theme:

1. Create a new directory under `themes/`
2. Add your templates in the `layouts/` subdirectory
3. Add your static assets in the `static/` subdirectory
4. Update your site's `config.jsonc` to use your theme name

### Custom Layouts

To override a theme's templates:

1. Create a file with the same name in your site's `layouts/` directory
2. Scribes will use your custom template instead of the theme's

### Performance

Based on our measured performance with 500 markdown files:

  1. Hugo (Go): Generally processes 500 pages in ~30-100ms, so our implementation is competitive
   but likely 1.5-3x slower
  2. Jekyll (Ruby): Typically takes several seconds (2-5s) for 500 pages, making our
  implementation ~30-70x faster
  1. Gatsby (JavaScript/React): Often takes 10-30 seconds for a full build of 500 pages, making
  ours ~150-400x faster
  1. Eleventy (JavaScript): Usually takes 1-3 seconds for 500 pages, making ours ~15-40x faster
  2. Zola (Rust): Processes 500 pages in ~100-200ms, so our implementation is comparable or
  slightly faster

  Our implementation is significantly faster than most JavaScript and Ruby-based generators,
  competitive with Rust-based ones, and in the same performance tier as Hugo.

  The key advantage remains that we achieved this with zero external dependencies, using only
  Go's standard library, which is quite impressive for the performance level we've reached.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.