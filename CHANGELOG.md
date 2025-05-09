## v0.7.0 (2025-05-08)

### Features

- feat: add release tasks to Taskfile
- feat: add Go implementation of release script

### Bug Fixes

- fix: remove unused path/filepath import from release script
- fix: improve project directory detection in release script
- fix: update release script to fix compilation errors

### Other Changes

- docs: add release documentation and update CLAUDE.md

## v0.5.2 (2025-05-07)

### Bug Fixes

- fix: correct imports and slug handling in tests
- fix: ensure test suite passes for URL path handling

## v0.5.1 (2025-05-07)

### Bug Fixes

- fix: resolve 404 errors for custom page paths
- fix: improved URL generation for pages created with custom directories

## v0.5.0 (2025-05-07)

### Features

- feat: simplified CLI interface with focused command set
- feat: improved interactive site creation experience
- feat: automatic Tailwind detection and integration in serve command

### Other Changes

- refactor: remove build and run commands for cleaner interface
- refactor: streamline new command to focus on site and page creation
- docs: update documentation to reflect simplified command structure
- test: add comprehensive tests for new CLI functionality

## v0.4.17 (2025-05-07)

### Bug Fixes

- fix: add missing strings import in CLI file
- fix: resolve build errors in goreleaser

## v0.4.16 (2025-05-07)

### Features

- feat: add flexible content creation with custom paths
- feat: add new generic content creation command
- feat: improve CLI with path-first syntax for content creation

### Other Changes

- refactor: update content creator to support file-based routing
- docs: enhance documentation for file-based routing and content creation

## v0.4.15 (2025-05-07)

### Bug Fixes

- fix: improve template caching test reliability
- fix: resolve intermittent test failures in CI

## v0.4.14 (2025-05-07)

### Features

- feat: add full file-based routing support for content directories
- feat: improve URL generation from file paths

### Bug Fixes

- fix: resolve issue with custom content directories not being rendered correctly

## v0.4.13 (2025-05-07)

### Bug Fixes

- fix: update template caching implementation marks in roadmap
- fix: refine CI workflow testing sequence

## v0.4.12 (2025-05-07)

### Features

- feat: implement template editing with proper syntax highlighting
- feat: optimize codebase with improved modularity and reduced size
- feat: enhance template management with embedded file structure
- feat: add developer documentation for template customization

### Bug Fixes

- fix: improve error handling with early returns
- fix: make task commands handle multiple directory scenarios
- fix: add missing template functions to tests

### Other Changes

- refactor: extract helper functions for improved code organization
- refactor: optimize template loading with better caching
- refactor: streamline site creation process

## v0.4.11 (2025-05-05)

### Features

- feat: enhance user experience
- feat: add Tailwind CSS support


## v0.4.10 (2025-05-05)

### Other Changes

- refactor: reduce codebase size and consolidate functionality


## v0.4.9 (2025-05-04)

### Features

- feat: completely remove all build output messages
- feat: further simplify server output to only show URL and watching message
- feat: add quiet mode to reduce verbose output from server
- feat: simplify site name prompt text
- feat: improve site creation messaging and instructions
- feat: simplify site creation to only ask for name and auto-initialize git
- feat: add arrow key navigation for interactive CLI selections
- feat: enhance CLI with improved UI components using standard library

### Other Changes

- chore: add test directory with sample site files
- refactor: move example site files from app/ to test/ directory and update server tests to support quiet mode
- refactor: improve selection UI with colored indicators instead of arrow keys


## v0.4.5 (2025-05-04)

### Bug Fixes

- fix: improve goreleaser configuration to fix failing CI

### Other Changes

- chore: bump version to v0.4.4


## v0.4.3 (2025-05-04)

### Bug Fixes

- fix: resolve GoReleaser dependency check issues

### Other Changes

- chore: update config.jsonc
- chore: bump version to v0.4.2 and add empty go.sum


## v0.3.2 (2025-05-04)

### Features

- feat: add Taskfile for development workflow

### Bug Fixes

- fix: correct GoReleaser configuration for zero-dependency project

### Other Changes

- refactor: remove console functionality


## v0.3.1 (2025-05-04)


## v0.3.0 (2025-05-03)

### Features

- feat: add manual homebrew formula update script

### Bug Fixes

- fix: ensure goreleaser config exists in release workflow


## v0.2.1 (2025-05-03)

### Features

- feat: automate Homebrew formula updates on release

### Bug Fixes

- fix: improve release script with robust version handling
- fix: correct version parsing in release script

### Other Changes

- chore: release v	Version = "v0.3.0
- chore: release v	Version = "v0.3.0
- chore: prepare for v0.2.0 release
- chore: release v.1.0
- chore: update Homebrew formula to v.1.0


# Changelog

All notable changes to this project will be documented in this file.

## v0.2.0 (2025-05-03)

### Features
- Automate Homebrew formula updates on release
- Implement fully automated release process

### Bug Fixes
- Fix version parsing in release script

### Other Changes
- Fix release tooling and version naming
- Standardize version format

## v0.1.0 (Initial Release)

### Features
- Basic static site generation
- Markdown to HTML conversion
- YAML front matter support
- Go template-based themes
- Development server with live reload
- Command-line interface
- Web management console
- Installation via Go install and Homebrew

### Coming Soon
- RSS feeds
- Improved markdown parsing
- HTML sanitization
- Better performance optimizations
