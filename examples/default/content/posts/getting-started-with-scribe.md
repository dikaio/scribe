---
title: Getting Started with Scribe
description: A guide to setting up your first Scribe site
date: 2023-05-20T14:00:00Z
tags:
  - tutorial
  - guide
  - scribe
---

# Getting Started with Scribe

In this post, I'll walk you through the process of setting up your first website with Scribe, a lightweight static site generator built in Go.

## What is Scribe?

Scribe is a minimalist static site generator that converts Markdown content with YAML front matter into beautiful HTML websites. It's built entirely with Go's standard library, which means it's fast and has no external dependencies.

## Installation

To install Scribe, you'll need Go installed on your machine. Then, you can run:

```bash
go install github.com/dikaio/scribe@latest
```

## Creating a New Site

Once you have Scribe installed, you can create a new site with:

```bash
scribe new site my-site
```

This will create a new directory with the basic structure you need to get started.

## Directory Structure

Your new site will have the following structure:

```
my-site/
├── config.yml         # Site configuration
├── content/           # Content files (Markdown)
│   ├── posts/         # Blog posts
│   └── *.md           # Regular pages
├── layouts/           # Custom templates (optional)
├── static/            # Static files
└── themes/            # Site themes
```

## Writing Content

Content in Scribe is written in Markdown with YAML front matter. Here's an example:

```markdown
---
title: My First Post
description: A short description
date: 2023-05-15T10:00:00Z
tags:
  - example
  - tutorial
---

# My First Post

This is the content of my first post. You can use **Markdown** formatting here.
```

## Building Your Site

To build your site, run:

```bash
scribe build
```

This will generate your site in the `public` directory.

## Development Server

During development, you can use the built-in server to preview your site:

```bash
scribe serve
```

This will start a development server with live reload at http://localhost:8080.

## Customizing Your Site

Scribe allows you to customize your site through themes and templates. You can modify the existing theme or create your own.

## Conclusion

Scribe makes it easy to create and maintain a static website. With its simple approach and fast performance, it's perfect for blogs, portfolios, and small business sites.

Give it a try and let me know what you think!