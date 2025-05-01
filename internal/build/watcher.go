package build

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Watcher watches for file changes
type Watcher struct {
	builder     *Builder
	sitePath    string
	lastBuild   time.Time
	changedDirs map[string]time.Time
}

// NewWatcher creates a new file watcher
func NewWatcher(builder *Builder, sitePath string) *Watcher {
	return &Watcher{
		builder:     builder,
		sitePath:    sitePath,
		lastBuild:   time.Now(),
		changedDirs: make(map[string]time.Time),
	}
}

// Watch starts watching for file changes
func (w *Watcher) Watch(interval time.Duration, rebuild func() error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("Watching for changes. Press Ctrl+C to stop.")

	for range ticker.C {
		changed, err := w.checkForChanges()
		if err != nil {
			return err
		}

		if changed {
			fmt.Println("Changes detected, rebuilding...")
			if err := rebuild(); err != nil {
				fmt.Printf("Error rebuilding: %s\n", err)
			} else {
				fmt.Println("Rebuild complete.")
				w.lastBuild = time.Now()
			}
		}
	}

	return nil
}

// checkForChanges checks if any files have changed since the last build
func (w *Watcher) checkForChanges() (bool, error) {
	// Directories to watch
	dirsToWatch := []string{
		filepath.Join(w.sitePath, "content"),
		filepath.Join(w.sitePath, "layouts"),
		filepath.Join(w.sitePath, "static"),
		filepath.Join(w.sitePath, "themes"),
		filepath.Join(w.sitePath, "config.jsonc"),
	}

	changed := false

	// Check each directory
	for _, dir := range dirsToWatch {
		fileChanged, err := w.hasChanged(dir)
		if err != nil {
			// Skip if directory doesn't exist
			if os.IsNotExist(err) {
				continue
			}
			return false, err
		}

		if fileChanged {
			changed = true
		}
	}

	return changed, nil
}

// hasChanged checks if a file or directory has changed
func (w *Watcher) hasChanged(path string) (bool, error) {
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// Check if it's a directory
	if info.IsDir() {
		return w.dirHasChanged(path)
	}

	// Check if file has changed
	return info.ModTime().After(w.lastBuild), nil
}

// dirHasChanged checks if any files in a directory have changed
func (w *Watcher) dirHasChanged(dir string) (bool, error) {
	changed := false

	// Walk through directory
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden files and directories
		if filepath.Base(path)[0] == '.' {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file has changed
		if info.ModTime().After(w.lastBuild) {
			changed = true
		}

		return nil
	})

	return changed, err
}
