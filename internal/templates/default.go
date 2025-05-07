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
            <h1><a href="/">Scribe</a></h1>
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
            <p>&copy; Scribe - A lightweight static site generator</p>
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
            <time>{{formatDate .Page.Date}}</time> • 2 min read
            {{if .Page.Tags}}
            <br>
            <span class="tags">
                {{range .Page.Tags}}
                <span class="tag">{{.}}</span>
                {{end}}
            </span>
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
            <time>{{formatDate .Date}}</time> • 2 min read
        </p>
        <p>{{.Description}}</p>
        <p><a href="/{{.URL}}/" class="read-more">Read more →</a></p>
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
            <time>{{formatDate .Date}}</time> • 2 min read
        </p>
        <p>{{.Description}}</p>
        <p><a href="/{{.URL}}/" class="read-more">Read more →</a></p>
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
    --primary-color: #2a6ec9;
    --text-color: #333;
    --background-color: #fff;
    --light-gray: #f5f5f5;
    --border-color: #ddd;
    --meta-color: #666;
    --tag-bg: #f0f0f0;
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
    background-color: var(--background-color);
    padding: 20px 0;
    margin-bottom: 40px;
    border-bottom: 1px solid var(--border-color);
}

header h1 {
    font-size: 1.8rem;
    font-weight: 700;
}

header h1 a {
    color: var(--text-color);
    text-decoration: none;
}

nav ul {
    list-style: none;
    display: flex;
    gap: 20px;
    margin-top: 10px;
}

nav a {
    color: var(--text-color);
    text-decoration: none;
    font-weight: 500;
}

nav a:hover {
    color: var(--primary-color);
}

main {
    min-height: 70vh;
    margin-bottom: 60px;
}

footer {
    background-color: var(--background-color);
    padding: 20px 0;
    border-top: 1px solid var(--border-color);
    text-align: center;
    color: var(--meta-color);
    font-size: 0.9rem;
}

h1, h2, h3, h4, h5, h6 {
    margin-bottom: 1rem;
    line-height: 1.25;
    font-weight: 600;
}

main > h1 {
    font-size: 2.5rem;
    margin-bottom: 2rem;
    color: #222;
}

p, ul, ol {
    margin-bottom: 1.5rem;
}

a {
    color: var(--primary-color);
    text-decoration: none;
}

a:hover {
    text-decoration: underline;
}

.post-list {
    display: flex;
    flex-direction: column;
    gap: 40px;
}

.post-summary {
    padding-bottom: 30px;
    border-bottom: 1px solid var(--border-color);
}

.post-summary h2 {
    font-size: 1.8rem;
    margin-bottom: 0.5rem;
}

.post-summary h2 a {
    color: var(--text-color);
    text-decoration: none;
}

.post-summary h2 a:hover {
    color: var(--primary-color);
}

.meta {
    color: var(--meta-color);
    font-size: 0.9rem;
    margin-bottom: 1rem;
}

.read-more {
    display: inline-block;
    font-weight: 500;
    margin-top: 0.5rem;
}

article .content {
    margin-top: 30px;
}

.tags {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
    margin-top: 5px;
}

.tag {
    display: inline-block;
    background-color: var(--tag-bg);
    padding: 3px 8px;
    border-radius: 4px;
    font-size: 0.8rem;
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
description: Lorem ipsum dolor sit amet, consectetur adipiscing elit
date: 2025-05-01T12:00:00Z
tags:
  - welcome
  - scribe
draft: false
---

# Welcome to Scribe!

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam in dui mauris. Vivamus hendrerit arcu sed erat molestie vehicula. Sed auctor neque eu tellus rhoncus ut eleifend nibh porttitor.

## Lorem Ipsum

Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur.

- Consectetur adipiscing elit
- Sed auctor neque eu tellus
- Vivamus hendrerit arcu
- Nam tincidunt congue enim

Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante. Donec eu libero sit amet quam egestas semper. Aenean ultricies mi vitae est. Mauris placerat eleifend leo.

## Action Button

[Action](#) 

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam in dui mauris.
`

// SamplePage is the default about page
const SamplePage = `---
title: About
description: Information about Scribe
draft: false
---

# About Scribe

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut elit tellus, luctus nec ullamcorper mattis, pulvinar dapibus leo. Sed non mauris vitae erat consequat auctor eu in elit.

## Our Mission

Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante. Donec eu libero sit amet quam egestas semper.

## Contact

[Action](#)

Mauris placerat eleifend leo. Quisque sit amet est et sapien ullamcorper pharetra. Vestibulum erat wisi, condimentum sed, commodo vitae, ornare sit amet, wisi.
`