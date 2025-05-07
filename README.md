# Scribe

A lightweight static site generator built in pure Go with zero external dependencies.

> **IMPORTANT:** This codebase is not intended for production use. It has been developed as an educational project leveraging Claude Code. While functional, it lacks the testing, security reviews, and production hardening necessary for real-world deployment.

## Overview

Scribe is a minimalist static site generator that transforms Markdown content with YAML front matter into elegant HTML websites. Built entirely with Go's standard library, it delivers exceptional performance without external dependencies.

## Features

- **Markdown to HTML conversion** - Write your content in Markdown, Scribe handles the rest
- **YAML front matter** - Add metadata to your content like title, date, tags, etc.
- **Templating system** - Use Go's html/template for layouts and themes
- **Live reload development server** - See changes as you make them
- **Command-line interface** - Simple commands for common operations
- **Zero dependencies** - Built using only the Go standard library

## Installation

Scribe is designed to be installed globally and used as a command-line tool. We offer two simple installation methods:

### Option 1: Install with Go (All Platforms)

```bash
go install github.com/dikaio/scribe@latest
```

This installs the latest version globally. Ensure your Go bin directory is in your PATH.

### Option 2: Install with Homebrew (macOS & Linux)

```bash
brew tap dikaio/tap
brew install scribe
```

This handles installation and dependencies automatically.

### Verifying Installation

After installation, verify Scribe is properly installed:

```bash
scribe version
```

You should see the version information displayed.

## Quick Start

### Create a new site

```bash
# Create a new site interactively
scribe new mysite

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


## Site Structure

```
mysite/
├── config.jsonc       # Site configuration
├── content/           # Content files (Markdown)
│   ├── posts/         # Blog posts (displayed in blog index)
│   ├── articles/      # Example custom content section
│   │   └── tech/      # Nested directories supported
│   └── *.md           # Regular pages
├── layouts/           # Custom template layouts (optional)
├── static/            # Static files (copied as-is)
├── themes/            # Site themes
│   └── default/       # Default theme
│       ├── layouts/   # Theme layouts
│       └── static/    # Theme static files
├── src/               # Source files (only with Tailwind CSS)
│   └── input.css      # Tailwind CSS input file
├── tailwind.config.js # Tailwind CSS config (only with Tailwind CSS)
└── package.json       # Node.js package file (only with Tailwind CSS)
```

## Content Format

Scribe uses Markdown files with YAML front matter for content:

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
  "title": "My Scribe Site",
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
| `scribe run [path]`       | Run development server and file watchers    |
| `scribe new [name]`       | Create a new site interactively             |
| `scribe new post [title]` | Create a new blog post                      |
| `scribe new page [title]` | Create a new page                           |
| `scribe help`             | Show help information                       |

## Templating

Scribe uses Go's `html/template` package for templating. Base templates include:

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

### CSS Frameworks

When creating a new site, Scribe offers two CSS framework options:

#### Default CSS

A simple, lightweight CSS framework with no external dependencies.

#### Tailwind CSS

A utility-first CSS framework that allows for rapid UI development. When choosing Tailwind CSS:

1. Node.js is required to use Tailwind CSS
2. The site will include a basic Tailwind CSS configuration
3. Dependencies will be automatically installed during site creation
4. To start development:
   ```bash
   # Run both Tailwind watcher and development server
   scribe run
   
   # Or run each process separately
   npm run dev     # In terminal 1
   scribe serve    # In terminal 2
   ```
5. Tailwind CSS styles will be compiled from `src/input.css` to `static/css/style.css`
6. For production, run `npm run build` to create an optimized CSS file

### Custom Themes

To create a custom theme:

1. Create a new directory under `themes/`
2. Add your templates in the `layouts/` subdirectory
3. Add your static assets in the `static/` subdirectory
4. Update your site's `config.jsonc` to use your theme name

### Custom Layouts

To override a theme's templates:

1. Create a file with the same name in your site's `layouts/` directory
2. Scribe will use your custom template instead of the theme's

### File-Based Routing

Scribe uses a file-based routing system similar to Next.js or Astro. URLs are derived directly from the content directory structure:

- Files at the root of the `content` directory become pages at the site root
  - Example: `content/about.md` → `/about/`
- Files in subdirectories maintain their directory structure in URLs
  - Example: `content/articles/tech/golang.md` → `/articles/tech/golang/`
- The special `posts` directory is used for blog posts and will be included in the homepage listing
  - Example: `content/posts/welcome.md` → `/posts/welcome/`

This intuitive system makes it easy to organize your content in logical sections while maintaining clean URLs. You can create any directory structure you need, and Scribe will automatically generate the corresponding URLs.

## Performance

Based on our benchmarking:

| Content Size | Build Time (Parallelized) |
| ------------ | ------------------------- |
| 100 pages    | ~19ms                     |
| 500 pages    | ~71ms                     |
| 1000 pages   | ~115ms                    |

Performance compared to other static site generators:

1. **Hugo** (Go): Generally processes 1000 pages in ~50-100ms, so our implementation is competitive
2. **Jekyll** (Ruby): Typically takes several seconds (5-10s) for 1000 pages, making our implementation ~50-80x faster
3. **Gatsby** (JavaScript/React): Often takes 15-45 seconds for a full build of 1000 pages, making ours ~150-400x faster
4. **Elevator** (JavaScript): Usually takes 2-5 seconds for 1000 pages, making ours ~20-40x faster
5. **Zola** (Rust): Processes 1000 pages in ~150-250ms, making our implementation comparable or slightly faster

Our implementation is significantly faster than most JavaScript and Ruby-based generators, competitive with Rust-based ones, and in the same performance tier as Hugo.

  The key advantage remains that we achieved this with zero external dependencies, using only
  Go's standard library, which is quite impressive for the performance level we've reached.

### Roadmap

#### Proposed

The following optimizations could bring Scribe performance to match or exceed Hugo:

- [ ] **Memory Pooling**: Implement object pools for frequently created/destroyed objects to reduce GC pressure
- [ ] **Incremental Builds**: Add file modification time tracking to only rebuild changed content
- [ ] **Template Caching**: Pre-compile and cache templates rather than loading them on each build
- [ ] **Custom Markdown Parser**: Replace our regex-based parser with a more optimized implementation
- [ ] **Concurrent File I/O**: Use async I/O patterns to overlap CPU and I/O work
- [ ] **Optimized Front Matter Parsing**: Replace our YAML parser with a faster implementation
- [ ] **Output Caching**: Cache rendered HTML for pages that haven't changed
- [ ] **Lazy Loading**: Only load dependencies when needed rather than upfront
- [ ] **Binary Template Storage**: Store compiled templates in binary format for faster loading
- [ ] **Optimized String Handling**: Reduce string allocations and use byte slices where possible
- [ ] **Write comprehensive documentation**: Outline the design, architecture, and usage of Scribe
- [ ] RSS Feeds
- [ ] CI/CD Pipeline

#### Planned

Optimizations we're actively working on:

- [x] **Parallel Processing**: Loading and processing content in parallel with worker pools
- [x] **Template Caching**: Pre-compile and cache templates for improved performance
- [ ] Implement plugin architecture
- [ ] Responsive Images (similar to Next.js Image component)
- [ ] Sitemap generation
- [ ] SEO Optimizations
- [ ] Add support for advanced content features like series and custom taxonomies
- [ ] Add more test coverage for all components
- [ ] Build additional themes
- [ ] Performance Profiling and Monitoring

## Contributing

Contributions are welcome! Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file for the detailed contribution process, including making changes, commits, and release procedures.

## License

This project is licensed under the MIT License - see the LICENSE file for details.