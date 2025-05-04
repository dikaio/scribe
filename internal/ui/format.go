package ui

import (
	"fmt"
	"strings"
)

// ANSI color codes
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Underline  = "\033[4m"
	
	Black      = "\033[30m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	
	BgBlack    = "\033[40m"
	BgRed      = "\033[41m"
	BgGreen    = "\033[42m"
	BgYellow   = "\033[43m"
	BgBlue     = "\033[44m"
	BgMagenta  = "\033[45m"
	BgCyan     = "\033[46m"
	BgWhite    = "\033[47m"
)

// Header prints a formatted header text
func Header(text string) {
	width := len(text) + 4
	border := strings.Repeat("─", width)
	
	fmt.Println()
	fmt.Println(Cyan + border + Reset)
	fmt.Println(Cyan + "│ " + Bold + text + Reset + Cyan + " │" + Reset)
	fmt.Println(Cyan + border + Reset)
	fmt.Println()
}

// Title prints a title with optional description
func Title(title string, description string) {
	fmt.Println(Bold + title + Reset)
	if description != "" {
		fmt.Println(description)
	}
	fmt.Println()
}

// Success prints a success message
func Success(message string) {
	fmt.Println(Green + "✓ " + message + Reset)
}

// Info prints an info message
func Info(message string) {
	fmt.Println(Blue + "ℹ " + Reset + message)
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Println(Yellow + "⚠ " + message + Reset)
}

// Error prints an error message
func Error(message string) {
	fmt.Println(Red + "✗ " + message + Reset)
}

// Prompt displays a prompt and returns the user input
func Prompt(message string, defaultValue string) string {
	var promptText string
	if defaultValue != "" {
		promptText = fmt.Sprintf("%s %s[%s]%s: ", message, Cyan, defaultValue, Reset)
	} else {
		promptText = fmt.Sprintf("%s: ", message)
	}
	
	fmt.Print(promptText)
	
	var input string
	fmt.Scanln(&input)
	
	if input == "" {
		return defaultValue
	}
	return input
}

// PromptRequired displays a prompt that requires input and returns the user input
func PromptRequired(message string) string {
	for {
		fmt.Print(Bold + message + Reset + ": ")
		
		var input string
		fmt.Scanln(&input)
		
		if input != "" {
			return input
		}
		
		fmt.Println(Red + "This field is required. Please enter a value." + Reset)
	}
}

// Divider prints a horizontal divider
func Divider() {
	fmt.Println(strings.Repeat("─", 40))
}