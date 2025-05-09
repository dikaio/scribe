package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
)

// Options for the release script
type Options struct {
	ReleaseType string
	DryRun      bool
	Help        bool
}

func main() {
	// Parse command line options
	options := parseOptions()

	if options.Help {
		showHelp()
		return
	}

	// Ensure we're in the project root
	if err := changeToProjectRoot(); err != nil {
		exitWithError(err.Error())
	}

	// Check if working directory is clean
	if clean, err := isWorkingDirClean(); err != nil {
		exitWithError(fmt.Sprintf("Error checking working directory: %v", err))
	} else if !clean {
		exitWithError("Working directory is not clean. Commit or stash changes first.")
	}

	// Check if we're on the main branch
	if branch, err := getCurrentBranch(); err != nil {
		exitWithError(fmt.Sprintf("Error getting current branch: %v", err))
	} else if branch != "main" {
		exitWithError(fmt.Sprintf("Not on main branch (currently on '%s'). Switch to main branch first.", branch))
	}

	// Pull latest changes
	printYellow("Pulling latest changes from remote...")
	if err := pullLatestChanges(); err != nil {
		exitWithError(fmt.Sprintf("Error pulling latest changes: %v", err))
	}

	// Run tests
	printYellow("Running tests...")
	if err := runTests(); err != nil {
		exitWithError(fmt.Sprintf("Tests failed: %v", err))
	}

	// Ensure go.mod exists
	if !fileExists("go.mod") {
		printYellow("Creating go.mod file for zero-dependency project...")
		if err := initGoMod(); err != nil {
			exitWithError(fmt.Sprintf("Error initializing go.mod: %v", err))
		}
	}

	// Read current version
	printYellow("Reading current version...")
	currentVersion, err := getCurrentVersion()
	if err != nil {
		exitWithError(fmt.Sprintf("Error getting current version: %v", err))
	}
	printGreen(fmt.Sprintf("Current version: v%s", currentVersion))

	// Calculate next version
	newVersion := calculateNextVersion(currentVersion, options.ReleaseType)
	printGreen(fmt.Sprintf("New version: v%s", newVersion))

	// Generate changelog
	printYellow("Generating changelog...")
	changelogContent, err := generateChangelog(newVersion)
	if err != nil {
		exitWithError(fmt.Sprintf("Error generating changelog: %v", err))
	}

	printYellow("Generated changelog:")
	fmt.Println(changelogContent)

	// Update version in code
	printYellow("Updating version in code...")
	if err := updateVersionInCode(newVersion); err != nil {
		exitWithError(fmt.Sprintf("Error updating version in code: %v", err))
	}

	// Update or create CHANGELOG.md
	printYellow("Updating CHANGELOG.md...")
	if err := updateChangelog(changelogContent); err != nil {
		exitWithError(fmt.Sprintf("Error updating changelog: %v", err))
	}

	// If dry run, revert changes and exit
	if options.DryRun {
		printYellow("Dry run mode. Changes will not be committed or pushed.")
		if err := revertChanges(); err != nil {
			exitWithError(fmt.Sprintf("Error reverting changes: %v", err))
		}
		printGreen("Dry run completed successfully.")
		return
	}

	// Commit changes
	printYellow("Committing changes...")
	if err := commitChanges(newVersion); err != nil {
		exitWithError(fmt.Sprintf("Error committing changes: %v", err))
	}

	// Create and push tag
	printYellow(fmt.Sprintf("Creating and pushing tag v%s...", newVersion))
	if err := createAndPushTag(newVersion); err != nil {
		exitWithError(fmt.Sprintf("Error creating and pushing tag: %v", err))
	}

	printGreen(fmt.Sprintf("Release v%s created and pushed successfully!", newVersion))
	printYellow("GitHub Actions will now build and publish the release.")
	printYellow("You can monitor the progress at: https://github.com/dikaio/scribe/actions")
}

// parseOptions parses command line options
func parseOptions() Options {
	options := Options{
		ReleaseType: "patch", // Default to patch
	}

	// Define command line flags
	flag.StringVar(&options.ReleaseType, "type", "patch", "Release type: patch, minor, major")
	flag.StringVar(&options.ReleaseType, "t", "patch", "Release type (shorthand)")
	flag.BoolVar(&options.DryRun, "dry-run", false, "Do everything except the actual release")
	flag.BoolVar(&options.DryRun, "d", false, "Dry run (shorthand)")
	flag.BoolVar(&options.Help, "help", false, "Show this help message")
	flag.BoolVar(&options.Help, "h", false, "Show help (shorthand)")

	// Parse flags
	flag.Parse()

	// Validate release type
	if options.ReleaseType != "patch" && options.ReleaseType != "minor" && options.ReleaseType != "major" {
		exitWithError(fmt.Sprintf("Invalid release type: %s. Use patch, minor, or major.", options.ReleaseType))
	}

	return options
}

// showHelp displays help information
func showHelp() {
	fmt.Printf("%sScribe Release Automation%s\n\n", colorYellow, colorReset)
	fmt.Println("Usage: go run scripts/release.go [OPTIONS]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -t, --type TYPE    Release type: patch, minor, major (default: patch)")
	fmt.Println("  -d, --dry-run      Do everything except the actual release")
	fmt.Println("  -h, --help         Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run scripts/release.go                # Creates a patch release")
	fmt.Println("  go run scripts/release.go --type minor   # Creates a minor release")
	fmt.Println("  go run scripts/release.go --type major   # Creates a major release")
	fmt.Println("  go run scripts/release.go --dry-run      # Simulates the release process")
}

// changeToProjectRoot ensures we're running from the project root
func changeToProjectRoot() error {
	// When running with go run, we're already in the project root
	// No need to change directory, just confirm we have a git repo
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not in a git repository: %v", err)
	}
	return nil
}

// isWorkingDirClean checks if the Git working directory is clean
func isWorkingDirClean() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return len(output) == 0, nil
}

// getCurrentBranch gets the current Git branch
func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// pullLatestChanges pulls the latest changes from the remote
func pullLatestChanges() error {
	cmd := exec.Command("git", "pull", "origin", "main")
	return cmd.Run()
}

// runTests runs all tests
func runTests() error {
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// initGoMod initializes a new go.mod file
func initGoMod() error {
	initCmd := exec.Command("go", "mod", "init", "github.com/dikaio/scribe")
	if err := initCmd.Run(); err != nil {
		return err
	}

	tidyCmd := exec.Command("go", "mod", "tidy")
	return tidyCmd.Run()
}

// getCurrentVersion gets the current version from the Go code
func getCurrentVersion() (string, error) {
	cliCodePath := "pkg/cli/cli.go"
	data, err := os.ReadFile(cliCodePath)
	if err != nil {
		return "", err
	}

	// Use regexp to find the version line
	re := regexp.MustCompile(`Version\s*=\s*"v([0-9]+\.[0-9]+\.[0-9]+)"`)
	matches := re.FindSubmatch(data)
	if len(matches) < 2 {
		return "0.1.0", nil // Default version if not found
	}

	return string(matches[1]), nil
}

// calculateNextVersion calculates the next version based on the current version and release type
func calculateNextVersion(currentVersion, releaseType string) string {
	parts := strings.Split(currentVersion, ".")
	if len(parts) != 3 {
		// If the current version is invalid, start from 0.1.0
		parts = []string{"0", "1", "0"}
	}

	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])

	switch releaseType {
	case "patch":
		patch++
	case "minor":
		minor++
		patch = 0
	case "major":
		major++
		minor = 0
		patch = 0
	}

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// getPreviousTag gets the most recent Git tag
func getPreviousTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		// No tags exist
		return "", nil
	}
	return strings.TrimSpace(string(output)), nil
}

// generateChangelog generates changelog content from Git commit history
func generateChangelog(newVersion string) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## v%s (%s)\n\n", newVersion, time.Now().Format("2006-01-02")))

	// Get the previous tag
	prevTag, err := getPreviousTag()
	if err != nil {
		return "", err
	}

	gitLogRange := ""
	if prevTag != "" {
		gitLogRange = prevTag + "..HEAD"
	}

	// Add features
	features, err := getCommitsByType(gitLogRange, "feat", "feature")
	if err != nil {
		return "", err
	}
	if len(features) > 0 {
		sb.WriteString("### Features\n\n")
		for _, feature := range features {
			sb.WriteString(fmt.Sprintf("- %s\n", feature))
		}
		sb.WriteString("\n")
	}

	// Add bug fixes
	fixes, err := getCommitsByType(gitLogRange, "fix")
	if err != nil {
		return "", err
	}
	if len(fixes) > 0 {
		sb.WriteString("### Bug Fixes\n\n")
		for _, fix := range fixes {
			sb.WriteString(fmt.Sprintf("- %s\n", fix))
		}
		sb.WriteString("\n")
	}

	// Add other changes
	others, err := getCommitsByType(gitLogRange, "refactor", "docs", "chore", "test", "ci", "build")
	if err != nil {
		return "", err
	}
	if len(others) > 0 {
		sb.WriteString("### Other Changes\n\n")
		for _, other := range others {
			sb.WriteString(fmt.Sprintf("- %s\n", other))
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

// getCommitsByType gets commit messages by type (feat, fix, etc.)
func getCommitsByType(gitLogRange string, types ...string) ([]string, error) {
	var commits []string

	for _, t := range types {
		args := []string{"log", "--pretty=format:%s"}
		if gitLogRange != "" {
			args = append(args, gitLogRange)
		}
		args = append(args, fmt.Sprintf("--grep=^%s", t))

		cmd := exec.Command("git", args...)
		output, err := cmd.Output()
		if err != nil {
			// Skip errors from grep command when no commits match
			continue
		}

		if len(output) > 0 {
			lines := strings.Split(strings.TrimSpace(string(output)), "\n")
			for _, line := range lines {
				if line != "" {
					commits = append(commits, line)
				}
			}
		}
	}

	return commits, nil
}

// updateVersionInCode updates the version in cli.go
func updateVersionInCode(newVersion string) error {
	filePath := "pkg/cli/cli.go"
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Use regexp to find and replace the version line
	re := regexp.MustCompile(`(Version\s*=\s*)"v[0-9]+\.[0-9]+\.[0-9]+"`)
	updated := re.ReplaceAll(data, []byte(fmt.Sprintf(`$1"v%s"`, newVersion)))

	return os.WriteFile(filePath, updated, 0644)
}

// updateChangelog updates or creates the CHANGELOG.md file
func updateChangelog(content string) error {
	changelogPath := "CHANGELOG.md"
	if fileExists(changelogPath) {
		// Read existing changelog
		data, err := os.ReadFile(changelogPath)
		if err != nil {
			return err
		}

		// Prepend new content to existing content
		updatedContent := content + string(data)
		return os.WriteFile(changelogPath, []byte(updatedContent), 0644)
	}

	// Create new changelog
	return os.WriteFile(changelogPath, []byte(content), 0644)
}

// revertChanges reverts changes for dry run
func revertChanges() error {
	cmd := exec.Command("git", "checkout", "--", "pkg/cli/cli.go", "CHANGELOG.md")
	return cmd.Run()
}

// commitChanges commits the version and changelog changes
func commitChanges(version string) error {
	// Add files
	addCmd := exec.Command("git", "add", "pkg/cli/cli.go", "CHANGELOG.md")
	if err := addCmd.Run(); err != nil {
		return err
	}

	// Commit
	commitCmd := exec.Command("git", "commit", "-m", fmt.Sprintf("chore: release v%s", version))
	return commitCmd.Run()
}

// createAndPushTag creates and pushes a Git tag
func createAndPushTag(version string) error {
	// Create annotated tag
	tagCmd := exec.Command("git", "tag", "-a", fmt.Sprintf("v%s", version), "-m", fmt.Sprintf("Release v%s", version))
	if err := tagCmd.Run(); err != nil {
		return err
	}

	// Push changes
	pushCmd := exec.Command("git", "push", "origin", "main")
	if err := pushCmd.Run(); err != nil {
		return err
	}

	// Push tag
	pushTagCmd := exec.Command("git", "push", "origin", fmt.Sprintf("v%s", version))
	return pushTagCmd.Run()
}

// print helpers
func printRed(message string) {
	fmt.Printf("%s%s%s\n", colorRed, message, colorReset)
}

func printGreen(message string) {
	fmt.Printf("%s%s%s\n", colorGreen, message, colorReset)
}

func printYellow(message string) {
	fmt.Printf("%s%s%s\n", colorYellow, message, colorReset)
}

func exitWithError(message string) {
	printRed("Error: " + message)
	os.Exit(1)
}