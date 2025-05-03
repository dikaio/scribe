package templates

// Default templates for site creation

// BaseTemplate is the base HTML template
const BaseTemplate = `<!DOCTYPE html>
<html lang="{{.Site.Language}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .Title}}{{.Title}} | {{end}}{{.Site.Title}}</title>
    <meta name="description" content="{{if .Description}}{{.Description}}{{else}}{{.Site.Description}}{{end}}">
    <link rel="stylesheet" href="/css/style.css">
</head>
<body>
    <header>
        <div class="container">
            <h1><a href="/">{{.Site.Title}}</a></h1>
            <nav>
                <ul>
                    <li><a href="/">Home</a></li>
                    <li><a href="/about/">About</a></li>
                </ul>
            </nav>
        </div>
    </header>
    <main class="container">
        {{block "content" .}}{{end}}
    </main>
    <footer>
        <div class="container">
            <p>&copy; {{.Site.Title}}</p>
        </div>
    </footer>
</body>
</html>`

// SingleTemplate is the template for individual posts
const SingleTemplate = `{{define "content"}}
<article>
    <header>
        <h1>{{.Page.Title}}</h1>
        <p class="meta">
            <time>{{formatDate .Page.Date}}</time>
            {{if .Page.Tags}}
            | Tags: 
            {{range .Page.Tags}}
            <a href="/tags/{{.}}/">{{.}}</a>
            {{end}}
            {{end}}
        </p>
    </header>
    <div class="content">
        {{.Content}}
    </div>
</article>
{{end}}`

// ListTemplate is the template for content lists
const ListTemplate = `{{define "content"}}
<h1>{{.Title}}</h1>
<div class="post-list">
    {{range .Pages}}
    <article class="post-summary">
        <h2><a href="/{{.URL}}/">{{.Title}}</a></h2>
        <p class="meta">
            <time>{{formatDate .Date}}</time>
            {{if .Tags}}
            | Tags: 
            {{range .Tags}}
            <a href="/tags/{{.}}/">{{.}}</a>
            {{end}}
            {{end}}
        </p>
        <p>{{.Description}}</p>
    </article>
    {{end}}
</div>
{{end}}`

// HomeTemplate is the template for the homepage
const HomeTemplate = `{{define "content"}}
<h1>Recent Posts</h1>
<div class="post-list">
    {{range .Pages}}
    <article class="post-summary">
        <h2><a href="/{{.URL}}/">{{.Title}}</a></h2>
        <p class="meta">
            <time>{{formatDate .Date}}</time>
            {{if .Tags}}
            | Tags: 
            {{range .Tags}}
            <a href="/tags/{{.}}/">{{.}}</a>
            {{end}}
            {{end}}
        </p>
        <p>{{.Description}}</p>
    </article>
    {{end}}
</div>
{{end}}`

// PageTemplate is the template for static pages
const PageTemplate = `{{define "content"}}
<article>
    <header>
        <h1>{{.Page.Title}}</h1>
    </header>
    <div class="content">
        {{.Content}}
    </div>
</article>
{{end}}`

// StyleCSS is the default stylesheet
const StyleCSS = `/* Basic styles for Scribe default theme */
:root {
    --primary-color: #0077cc;
    --text-color: #333;
    --background-color: #fff;
    --light-gray: #f5f5f5;
    --border-color: #ddd;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen, Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
    line-height: 1.6;
    color: var(--text-color);
    background-color: var(--background-color);
}

.container {
    max-width: 800px;
    margin: 0 auto;
    padding: 0 20px;
}

header {
    background-color: var(--light-gray);
    padding: 20px 0;
    margin-bottom: 40px;
    border-bottom: 1px solid var(--border-color);
}

header h1 {
    font-size: 2rem;
}

header h1 a {
    color: var(--text-color);
    text-decoration: none;
}

nav ul {
    list-style: none;
    display: flex;
    gap: 20px;
}

nav a {
    color: var(--primary-color);
    text-decoration: none;
}

main {
    min-height: 70vh;
    margin-bottom: 40px;
}

footer {
    background-color: var(--light-gray);
    padding: 20px 0;
    border-top: 1px solid var(--border-color);
    text-align: center;
}

h1, h2, h3, h4, h5, h6 {
    margin-bottom: 1rem;
    line-height: 1.25;
}

p, ul, ol {
    margin-bottom: 1.5rem;
}

a {
    color: var(--primary-color);
}

.post-list {
    display: flex;
    flex-direction: column;
    gap: 30px;
}

.post-summary {
    padding-bottom: 20px;
    border-bottom: 1px solid var(--border-color);
}

.meta {
    color: #666;
    font-size: 0.9rem;
    margin-bottom: 1rem;
}

article .content {
    margin-top: 20px;
}

/* Code blocks */
pre {
    background-color: var(--light-gray);
    padding: 1rem;
    overflow-x: auto;
    border-radius: 4px;
    margin-bottom: 1.5rem;
}

code {
    font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
    font-size: 0.9em;
    background-color: var(--light-gray);
    padding: 0.2em 0.4em;
    border-radius: 3px;
}

pre code {
    padding: 0;
    background-color: transparent;
}`

// SampleContent contains default content for new sites

// SamplePost is the default welcome post
const SamplePost = `---
title: Welcome to Scribe
description: A sample post to get you started
date: 2025-05-01T12:00:00Z
tags:
  - welcome
  - scribe
draft: false
---

# Welcome to Scribe!

This is a sample post to help you get started with Scribe, a lightweight static site generator.

## Features

- Markdown support
- Fast and lightweight
- No external dependencies
- Simple to use

Enjoy creating content with Scribe!
`

// SamplePage is the default about page
const SamplePage = `---
title: About
description: About this site
draft: false
---

# About

This is an about page for your Scribe site. You can add information about yourself or your project here.

## Contact

Feel free to reach out with any questions or feedback.
`