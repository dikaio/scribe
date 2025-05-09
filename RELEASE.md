# Release Process

This document describes the process for creating new releases of Scribe.

The release script is implemented as a Go package in `internal/release` with a command-line interface in `cmd/release`.

## Release Steps

1. Make sure all your changes are committed and pushed to GitHub:
   ```bash
   git add .
   git commit -m "your commit message"
   git push origin main
   ```

2. Run the appropriate release task depending on the type of changes:
   
   ```bash
   # For patch releases (bug fixes)
   task release:patch
   
   # For minor releases (new features that are backwards compatible)
   task release:minor
   
   # For major releases (breaking changes)
   task release:major
   ```

3. The release script will automatically:
   - Pull the latest changes
   - Run tests to ensure everything works
   - Update the version number in the code
   - Generate a changelog from commit messages
   - Commit these changes
   - Create and push a new version tag
   - Push all changes to GitHub

## Testing the Release Process

If you want to test the release process without actually making changes, use:

```bash
task release:dry-run
```

This will show you what would happen without actually committing or pushing anything.

## Release Types

- **Patch**: For bug fixes and minor improvements that don't add features or break existing functionality (e.g., 1.0.0 → 1.0.1)
- **Minor**: For new features that are backwards compatible (e.g., 1.0.0 → 1.1.0)
- **Major**: For breaking changes or major rewrites (e.g., 1.0.0 → 2.0.0)

## Commit Message Conventions

The changelog is generated automatically from commit messages. To ensure your changes are properly categorized, use these prefixes:

- `feat:` or `feature:` for new features
- `fix:` for bug fixes
- `refactor:` for code refactoring
- `docs:` for documentation changes
- `test:` for test-related changes
- `chore:` for maintenance tasks
- `ci:` for CI-related changes
- `build:` for build-related changes

Example: `feat: add dark mode support`