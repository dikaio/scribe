package templates

// Console templates for the web-based management console

// ConsoleBaseTemplate is the base HTML template for the console
const ConsoleBaseTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - Scribe Console</title>
    <style>
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
            display: flex;
            min-height: 100vh;
        }
        
        .sidebar {
            width: 200px;
            background-color: var(--light-gray);
            padding: 20px;
            border-right: 1px solid var(--border-color);
        }
        
        .sidebar h1 {
            font-size: 1.5rem;
            margin-bottom: 20px;
            color: var(--primary-color);
        }
        
        .sidebar nav ul {
            list-style: none;
        }
        
        .sidebar nav li {
            margin-bottom: 10px;
        }
        
        .sidebar nav a {
            color: var(--text-color);
            text-decoration: none;
            display: block;
            padding: 5px 0;
        }
        
        .sidebar nav a:hover {
            color: var(--primary-color);
        }
        
        .main-content {
            flex: 1;
            padding: 20px;
            max-width: 1000px;
        }
        
        .header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--border-color);
        }
        
        .header h2 {
            font-size: 1.8rem;
        }
        
        .btn {
            display: inline-block;
            background-color: var(--primary-color);
            color: white;
            padding: 8px 16px;
            border-radius: 4px;
            text-decoration: none;
            border: none;
            cursor: pointer;
            font-size: 14px;
        }
        
        .btn:hover {
            opacity: 0.9;
        }
        
        .card {
            background-color: white;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            padding: 20px;
            margin-bottom: 20px;
        }
        
        .card h3 {
            margin-bottom: 15px;
            border-bottom: 1px solid var(--border-color);
            padding-bottom: 10px;
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
        }
        
        table th,
        table td {
            text-align: left;
            padding: 10px;
            border-bottom: 1px solid var(--border-color);
        }
        
        table th {
            background-color: var(--light-gray);
        }
        
        .form-group {
            margin-bottom: 15px;
        }
        
        label {
            display: block;
            margin-bottom: 5px;
            font-weight: bold;
        }
        
        input[type="text"],
        input[type="url"],
        textarea,
        select {
            width: 100%;
            padding: 8px;
            border: 1px solid var(--border-color);
            border-radius: 4px;
            font-family: inherit;
            font-size: inherit;
        }
        
        textarea {
            min-height: 200px;
        }
        
        .dashboard-grid {
            display: grid;
            grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
            gap: 20px;
        }
    </style>
</head>
<body>
    <div class="sidebar">
        <h1>Scribe Console</h1>
        <nav>
            <ul>
                <li><a href="/">Dashboard</a></li>
                <li><a href="/content">Content</a></li>
                <li><a href="/settings">Settings</a></li>
                <li><a href="/build">Build Site</a></li>
                <li><a href="/" target="_blank">View Site</a></li>
            </ul>
        </nav>
    </div>
    <div class="main-content">
        {{template "content" .}}
    </div>
</body>
</html>`

// ConsoleTemplates contains all console view templates
var ConsoleTemplates = map[string]string{
	"dashboard": `{{define "content"}}
<div class="header">
    <h2>Dashboard</h2>
    <a href="/content/new" class="btn">New Content</a>
</div>

<div class="dashboard-grid">
    <div class="card">
        <h3>Site Overview</h3>
        <ul>
            <li><strong>Title:</strong> {{.Site.Title}}</li>
            <li><strong>URL:</strong> {{.Site.BaseURL}}</li>
            <li><strong>Posts:</strong> {{.PostCount}}</li>
            <li><strong>Pages:</strong> {{.PageCount}}</li>
            <li><strong>Theme:</strong> {{.Site.Theme}}</li>
        </ul>
    </div>
    
    <div class="card">
        <h3>Recent Posts</h3>
        {{if .RecentPosts}}
        <table>
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Date</th>
                </tr>
            </thead>
            <tbody>
                {{range .RecentPosts}}
                <tr>
                    <td><a href="/content/edit?path={{.Path}}">{{.Title}}</a></td>
                    <td>{{formatDate .Date}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <p>No posts found.</p>
        {{end}}
    </div>
    
    <div class="card">
        <h3>Recent Pages</h3>
        {{if .RecentPages}}
        <table>
            <thead>
                <tr>
                    <th>Title</th>
                    <th>Date</th>
                </tr>
            </thead>
            <tbody>
                {{range .RecentPages}}
                <tr>
                    <td><a href="/content/edit?path={{.Path}}">{{.Title}}</a></td>
                    <td>{{formatDate .Date}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
        {{else}}
        <p>No pages found.</p>
        {{end}}
    </div>
</div>
{{end}}`,

	"content": `{{define "content"}}
<div class="header">
    <h2>Content</h2>
    <a href="/content/new" class="btn">New Content</a>
</div>

<div class="card">
    <div style="margin-bottom: 20px;">
        <a href="/content?type=all" {{if eq .ContentType "all"}}style="font-weight: bold;"{{end}}>All</a> |
        <a href="/content?type=posts" {{if eq .ContentType "posts"}}style="font-weight: bold;"{{end}}>Posts</a> |
        <a href="/content?type=pages" {{if eq .ContentType "pages"}}style="font-weight: bold;"{{end}}>Pages</a>
    </div>
    
    {{if .Items}}
    <table>
        <thead>
            <tr>
                <th>Title</th>
                <th>Type</th>
                <th>Date</th>
                <th>Status</th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td><a href="/content/edit?path={{.Path}}">{{.Title}}</a></td>
                <td>{{if .IsPost}}Post{{else}}Page{{end}}</td>
                <td>{{formatDate .Date}}</td>
                <td>{{if .Draft}}Draft{{else}}Published{{end}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
    {{else}}
    <p>No content found.</p>
    {{end}}
</div>
{{end}}`,

	"new_content": `{{define "content"}}
<div class="header">
    <h2>New {{if eq .ContentType "post"}}Post{{else}}Page{{end}}</h2>
</div>

<div class="card">
    <form action="/content/new?type={{.ContentType}}" method="post">
        <div class="form-group">
            <label for="title">Title</label>
            <input type="text" id="title" name="title" required>
        </div>
        
        <div class="form-group">
            <label for="description">Description</label>
            <input type="text" id="description" name="description">
        </div>
        
        {{if eq .ContentType "post"}}
        <div class="form-group">
            <label for="tags">Tags (comma separated)</label>
            <input type="text" id="tags" name="tags">
        </div>
        {{end}}
        
        <div class="form-group">
            <label for="content">Content</label>
            <textarea id="content" name="content" required></textarea>
        </div>
        
        <div class="form-group">
            <label>
                <input type="checkbox" name="draft"> Draft
            </label>
        </div>
        
        <button type="submit" class="btn">Create {{if eq .ContentType "post"}}Post{{else}}Page{{end}}</button>
    </form>
</div>
{{end}}`,

	"settings": `{{define "content"}}
<div class="header">
    <h2>Settings</h2>
</div>

<div class="card">
    <form action="/settings" method="post">
        <div class="form-group">
            <label for="title">Site Title</label>
            <input type="text" id="title" name="title" value="{{.Site.Title}}" required>
        </div>
        
        <div class="form-group">
            <label for="baseURL">Base URL</label>
            <input type="url" id="baseURL" name="baseURL" value="{{.Site.BaseURL}}" required>
        </div>
        
        <div class="form-group">
            <label for="description">Site Description</label>
            <input type="text" id="description" name="description" value="{{.Site.Description}}">
        </div>
        
        <div class="form-group">
            <label for="author">Author</label>
            <input type="text" id="author" name="author" value="{{.Site.Author}}">
        </div>
        
        <div class="form-group">
            <label for="theme">Theme</label>
            <select id="theme" name="theme">
                <option value="default" {{if eq .Site.Theme "default"}}selected{{end}}>Default</option>
                <!-- Add other themes here as they become available -->
            </select>
        </div>
        
        <button type="submit" class="btn">Save Settings</button>
    </form>
</div>
{{end}}`,
}