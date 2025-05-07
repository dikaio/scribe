package templates

// TailwindInputCSS is the input CSS file for Tailwind
const TailwindInputCSS = `@import "tailwindcss";

/* This file only imports Tailwind CSS */
/* All styling should be done with Tailwind utility classes in HTML */
/* You can add custom Tailwind theme configuration below if needed */

@layer base {
  /* You can add custom base styles or theme customization here if needed */
}
`

// TailwindBaseTemplate is the base HTML template with Tailwind CSS
const TailwindBaseTemplate = `<!DOCTYPE html>
<html lang="{{.Site.Language}}">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{if .Title}}{{.Title}} | {{end}}{{.Site.Title}}</title>
    <meta name="description" content="{{if .Description}}{{.Description}}{{else}}{{.Site.Description}}{{end}}">
    <link rel="stylesheet" href="/css/style.css">
</head>
<body class="bg-white text-gray-800 font-sans">
    <header class="bg-gray-100 border-b border-gray-200 py-6 mb-10">
        <div class="container mx-auto px-4 max-w-4xl">
            <h1 class="text-3xl font-bold"><a href="/" class="text-gray-800 no-underline">Scribe</a></h1>
            <nav class="mt-2">
                <ul class="flex space-x-6">
                    <li><a href="/" class="text-blue-600 hover:text-blue-800 no-underline">Home</a></li>
                    <li><a href="/about/" class="text-blue-600 hover:text-blue-800 no-underline">About</a></li>
                </ul>
            </nav>
        </div>
    </header>
    <main class="container mx-auto px-4 max-w-4xl min-h-screen mb-10">
        {{block "content" .}}{{end}}
    </main>
    <footer class="bg-gray-100 border-t border-gray-200 py-6 text-center">
        <div class="container mx-auto px-4 max-w-4xl">
            <p>&copy; Scribe - A lightweight static site generator</p>
        </div>
    </footer>
</body>
</html>`

// TailwindSingleTemplate is the template for individual posts with Tailwind CSS
const TailwindSingleTemplate = `{{define "content"}}
<article>
    <header>
        <h1 class="text-4xl font-bold mb-2">{{.Page.Title}}</h1>
        <p class="text-gray-600 text-sm mb-6">
            <time>{{formatDate .Page.Date}}</time> • 2 min read
            {{if .Page.Tags}}
            <div class="flex flex-wrap gap-2 mt-2">
                {{range .Page.Tags}}
                <span class="bg-gray-100 text-gray-700 px-2 py-1 rounded text-xs">{{.}}</span>
                {{end}}
            </div>
            {{end}}
        </p>
    </header>
    <div class="prose lg:prose-xl max-w-none">
        {{.Content}}
    </div>
</article>
{{end}}`

// TailwindListTemplate is the template for content lists with Tailwind CSS
const TailwindListTemplate = `{{define "content"}}
<h1 class="text-4xl font-bold mb-8">{{.Title}}</h1>
<div class="space-y-10">
    {{range .Pages}}
    <article class="pb-8 border-b border-gray-200">
        <h2 class="text-2xl font-bold mb-1"><a href="/{{.URL}}/" class="text-gray-800 hover:text-blue-700 no-underline">{{.Title}}</a></h2>
        <p class="text-gray-600 text-sm mb-3">
            <time>{{formatDate .Date}}</time> • 2 min read
        </p>
        <p class="text-gray-700 mb-2">{{.Description}}</p>
        <p><a href="/{{.URL}}/" class="text-blue-600 hover:text-blue-800 font-medium">Read more →</a></p>
    </article>
    {{end}}
</div>
{{end}}`

// TailwindHomeTemplate is the template for the homepage with Tailwind CSS
const TailwindHomeTemplate = `{{define "content"}}
<h1 class="text-4xl font-bold mb-8">Recent Posts</h1>
<div class="space-y-10">
    {{range .Pages}}
    <article class="pb-8 border-b border-gray-200">
        <h2 class="text-2xl font-bold mb-1"><a href="/{{.URL}}/" class="text-gray-800 hover:text-blue-700 no-underline">{{.Title}}</a></h2>
        <p class="text-gray-600 text-sm mb-3">
            <time>{{formatDate .Date}}</time> • 2 min read
        </p>
        <p class="text-gray-700 mb-2">{{.Description}}</p>
        <p><a href="/{{.URL}}/" class="text-blue-600 hover:text-blue-800 font-medium">Read more →</a></p>
    </article>
    {{end}}
</div>
{{end}}`

// TailwindPageTemplate is the template for static pages with Tailwind CSS
const TailwindPageTemplate = `{{define "content"}}
<article>
    <header>
        <h1 class="text-4xl font-bold mb-6">{{.Page.Title}}</h1>
    </header>
    <div class="prose lg:prose-xl max-w-none">
        {{.Content}}
    </div>
</article>
{{end}}`

// TailwindPackageJSON is the package.json file for Tailwind CSS
const TailwindPackageJSON = `{
  "name": "scribe-tailwind-site",
  "version": "1.0.0",
  "description": "A site created with Scribe and Tailwind CSS",
  "scripts": {
    "dev": "npx tailwindcss -i ./src/input.css -o ./static/css/style.css --watch",
    "build": "npx tailwindcss -i ./src/input.css -o ./static/css/style.css --minify"
  },
  "devDependencies": {
    "tailwindcss": "^4.1.0",
    "@tailwindcss/cli": "^0.1.0"
  }
}`

// TailwindGitignore adds Node.js-specific entries to the .gitignore file
const TailwindGitignore = `# Output directory
public/

# IDE files
.idea/
.vscode/

# System files
.DS_Store
Thumbbs.db

# Node.js
node_modules/
package-lock.json
`

// TailwindREADME is the README.md file with Tailwind CSS instructions
const TailwindREADME = `# Scribe Site with Tailwind CSS 4.1

This site was created using Scribe static site generator with modern Tailwind CSS 4.1.

## Tailwind CSS Approach

This site uses a utility-first CSS approach with Tailwind CSS:

- All styling is done with Tailwind utility classes directly in HTML templates
- No custom CSS needed in most cases
- Rapid development with composable utility classes
- Optimized build process produces minimal CSS

## Development

To start developing your site:

### Initial Setup

1. Install dependencies:
   ` + "```" + `
   npm install
   ` + "```" + `

2. Start the Tailwind CSS compiler:
   ` + "```" + `
   npm run dev
   ` + "```" + `

3. In another terminal, start the Scribe development server:
   ` + "```" + `
   scribe serve
   ` + "```" + `

4. View your site at http://localhost:8080

Or for convenience, use the all-in-one command:
   ` + "```" + `
   scribe run
   ` + "```" + `

### Building for Production

To build your site for production:

1. Compile and minify the CSS:
   ` + "```" + `
   npm run build
   ` + "```" + `

2. Build the site:
   ` + "```" + `
   scribe build
   ` + "```" + `

The built site will be in the ` + "`public`" + ` directory.

### Modern Tailwind CSS Setup

This site uses the modern Tailwind CSS 4.1 approach:
- Simplified setup with just ` + "`@import \"tailwindcss\"`" + ` in the CSS
- Comes with both tailwindcss and @tailwindcss/cli for processing
- No configuration files needed (no tailwind.config.js)
- Customization available through ` + "`@layer`" + ` directives if needed
- Focus on using utility classes in HTML rather than writing custom CSS
`