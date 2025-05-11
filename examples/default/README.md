# Scribe Default Starter

This is a starter template for creating a website with [Scribe](https://github.com/dikaio/scribe), a lightweight static site generator built in Go.

## Getting Started

### Prerequisites

- Go 1.19 or higher
- Scribe installed (`go install github.com/dikaio/scribe@latest`)

### Using This Starter

1. **Create a new site using this starter:**

   ```bash
   # Clone this starter template
   git clone https://github.com/dikaio/scribe.git
   cp -r scribe/examples/default my-site
   cd my-site
   ```

   Or with the Scribe CLI (if implemented):

   ```bash
   scribe new site my-site --starter default
   cd my-site
   ```

2. **Start the development server:**

   ```bash
   scribe serve
   ```

   Your site will be available at http://localhost:8080.

3. **Build the site for production:**

   ```bash
   scribe build
   ```

   This will generate your site in the `public` directory.

## Directory Structure

```
my-site/
├── config.yml         # Site configuration
├── content/           # Content files (Markdown)
│   ├── _index.md      # Homepage content
│   ├── about.md       # About page
│   ├── contact.md     # Contact page
│   └── posts/         # Blog posts
│       ├── welcome.md
│       ├── getting-started-with-scribe.md
│       └── markdown-syntax-guide.md
├── layouts/           # Custom templates (optional)
├── static/            # Static files (CSS, JS, images)
└── themes/            # Site themes
```

## Customization

### Configuration

Edit the `config.yml` file to customize your site settings:

```yaml
title: Your Site Title
baseURL: https://example.com/
theme: default
language: en
# ... more settings
```

### Content

All content is written in Markdown with YAML front matter:

```markdown
---
title: Page Title
description: Page description
date: 2023-05-01T12:00:00Z
tags:
  - tag1
  - tag2
---

# Content here

Your Markdown content goes here.
```

### Themes

The default theme is located in the `themes/default` directory. You can customize it or create your own theme.

## Features

- **Simple Content Management**: Write content in Markdown with YAML front matter
- **Fast Build Times**: Generate your site in milliseconds
- **Flexible Templating**: Customize your site's appearance with Go templates
- **No External Dependencies**: Built with Go's standard library

## Learn More

For more information about Scribe, check out the [Scribe documentation](https://github.com/dikaio/scribe).