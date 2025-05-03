# Release Process

This document outlines the fully automated process for creating a new release of Scribe.

## Automated Release Workflow

The release process is fully automated and includes:
- Version bumping (semantic versioning)
- Changelog generation
- Commit and tag creation
- Binary building for all platforms
- GitHub release creation
- Homebrew formula update

## Creating a Release

### Option 1: Using the Release Script (Recommended)

1. Ensure your working directory is clean and you're on the main branch
2. Run the release script with desired version bump type:

```bash
# For a patch release (0.1.0 -> 0.1.1)
./scripts/release.sh

# For a minor release (0.1.0 -> 0.2.0)
./scripts/release.sh --type minor

# For a major release (0.1.0 -> 1.0.0)
./scripts/release.sh --type major

# To perform a dry run (no actual changes)
./scripts/release.sh --dry-run
```

The script will:
1. Check your working directory is clean
2. Run all tests
3. Calculate the new version
4. Generate a changelog from commit messages
5. Update version in code
6. Commit changes
7. Create and push a Git tag
8. GitHub Actions will then build and publish the release

### Option 2: Using GitHub Actions

1. Go to the GitHub repository
2. Navigate to Actions > Pre-Release Testing
3. Click "Run workflow"
4. Select the version type (patch, minor, major)
5. Review the workflow results
6. If successful, manually run the release script or create a tag

## Versioning Convention

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Incompatible API changes
- **MINOR**: Add functionality in a backward-compatible manner
- **PATCH**: Backward-compatible bug fixes

## Commit Message Format

The changelog is automatically generated from commit messages using these prefixes:
- `feat:` or `feature:` - New features
- `fix:` - Bug fixes
- `refactor:` - Code changes that neither fix bugs nor add features
- `docs:` - Documentation changes
- `test:` - Adding or fixing tests
- `chore:` - Maintenance tasks

Example: `feat: add support for custom templates`

## Installation Methods

After the release process completes, users can install Scribe via:

### Option 1: Go Install

```bash
go install github.com/dikaio/scribe@latest
```

### Option 2: Homebrew

```bash
brew install scribe
```

## Troubleshooting

If the automated process fails:

1. Check GitHub Actions logs for errors
2. Run the release script with `--dry-run` to debug
3. Ensure all tests are passing
4. Verify you have clean working directory and are on main branch