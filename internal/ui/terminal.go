package ui

import (
	"fmt"
	"os"
	"os/exec"
)

// SetRawMode puts the terminal into raw mode where we can read individual keystrokes
func SetRawMode() (*exec.Cmd, error) {
	// Only works on Unix-like systems
	cmd := exec.Command("stty", "-F", "/dev/tty", "raw", "-echo")
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

// RestoreTerminal restores the terminal to its previous state
func RestoreTerminal(cmd *exec.Cmd) error {
	if cmd == nil {
		return nil
	}
	// Restore terminal
	restoreCmd := exec.Command("stty", "-F", "/dev/tty", "sane")
	return restoreCmd.Run()
}

// ReadKey reads a single key press from stdin
func ReadKey() ([]byte, error) {
	// Read a single byte
	buf := make([]byte, 3)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// Clear clears the terminal screen
func Clear() {
	fmt.Print("\033[H\033[2J") // ANSI escape code to clear screen
}

// MoveCursor moves the cursor to a specific position
func MoveCursor(row, col int) {
	fmt.Printf("\033[%d;%dH", row, col)
}

// ClearLine clears the current line
func ClearLine() {
	fmt.Print("\033[2K\r")
}