# Template Development Guide

This guide explains how to develop templates for Scribe using the embedded HTML files.

## Template Structure

Templates are now stored as real HTML files in:
- `internal/templates/embedded/default/` - Default theme templates
- `internal/templates/embedded/tailwind/` - Tailwind CSS theme templates

## Editing Templates

You can edit the HTML templates directly with full syntax highlighting. The templates use Go's HTML template syntax.

## Building After Changes

After making changes to the template files, you need to rebuild the application:

```bash
task build
```

This compiles the templates into the binary using Go's embed package.

## Available Template Files

### Default Theme
- `base.html` - Base layout template
- `home.html` - Homepage template
- `list.html` - Content list template
- `single.html` - Single post template
- `page.html` - Static page template
- `style.css` - Default CSS

### Tailwind Theme
- `base.html` - Base layout template
- `home.html` - Homepage template
- `list.html` - Content list template
- `single.html` - Single post template
- `page.html` - Static page template
- `input.css` - Tailwind CSS input file
- `package.json` - NPM package configuration
- `gitignore` - Git ignore file
- `README.md` - Tailwind site README

## Template Functions

These functions are available in templates:

- `{{formatDate .Date}}` - Format a date as "January 2, 2006"
- `{{lower "TEXT"}}` - Convert to lowercase
- `{{upper "text"}}` - Convert to uppercase
- `{{title "text"}}` - Convert to title case
- `{{now.Format "2006"}}` - Get the current year (or other date format)
- `{{sub a b}}` - Subtract b from a

## Testing Your Changes

After building, create a new site and check the results:

```bash
bin/scribe new test-site
cd test-site
bin/scribe serve
```

Then visit http://localhost:8080 to see your changes.