package ui

import (
	"fmt"
	"strconv"
	"strings"
)

// Option represents a selectable option in a menu
type Option struct {
	Label string
	Value string
}

// SelectOption displays a selection menu and returns the selected option value
func SelectOption(title string, description string, options []Option, defaultIndex int) string {
	// Ensure defaultIndex is valid
	if defaultIndex < 0 || defaultIndex >= len(options) {
		defaultIndex = 0
	}
	
	// Print title and description
	Title(title, description)
	fmt.Println("Use ↑ and ↓ arrow keys to navigate, Enter to select, or type a number")
	
	// Try to use arrow keys if possible
	if arrowValue, ok := tryArrowKeySelect(options, defaultIndex); ok {
		return arrowValue
	}
	
	// Fall back to number-based selection if arrow keys don't work
	return numberBasedSelect(options, defaultIndex)
}

// tryArrowKeySelect attempts to use arrow keys for selection
// Returns the selected value and true if successful, or empty string and false if failed
func tryArrowKeySelect(options []Option, defaultIndex int) (string, bool) {
	// Put terminal in raw mode to read arrow keys
	rawCmd, err := SetRawMode()
	if err != nil {
		// If we can't set raw mode, fall back to number-based selection
		return "", false
	}
	
	// Ensure terminal is restored when function exits
	defer RestoreTerminal(rawCmd)
	
	// Current selection
	selected := defaultIndex
	
	// Display initial selection
	displayArrowSelection(options, selected)
	
	// Read keys until Enter is pressed
	for {
		// Read a key
		key, err := ReadKey()
		if err != nil {
			return "", false
		}
		
		// Process key
		if len(key) == 1 {
			// Enter key
			if key[0] == 13 || key[0] == 10 {
				// Return selected option
				return options[selected].Value, true
			}
			
			// Numbers 1-9 for direct selection
			if key[0] >= '1' && key[0] <= '9' {
				num := int(key[0] - '0') - 1
				if num < len(options) {
					// Update selection
					selected = num
					displayArrowSelection(options, selected)
				}
			}
			
			// ESC key
			if key[0] == 27 {
				return "", false
			}
		} else if len(key) == 3 {
			// Arrow up
			if key[0] == 27 && key[1] == 91 && key[2] == 65 {
				if selected > 0 {
					selected--
					displayArrowSelection(options, selected)
				}
			}
			
			// Arrow down
			if key[0] == 27 && key[1] == 91 && key[2] == 66 {
				if selected < len(options)-1 {
					selected++
					displayArrowSelection(options, selected)
				}
			}
		}
	}
}

// displayArrowSelection displays the options with the currently selected one highlighted
func displayArrowSelection(options []Option, selected int) {
	// Clear previous output (move up and clear lines)
	for i := 0; i < len(options); i++ {
		fmt.Print("\033[1A\033[2K")
	}
	
	// Display options
	for i, option := range options {
		if i == selected {
			fmt.Println(Cyan + " › " + Bold + option.Label + Reset)
		} else {
			fmt.Println("   " + option.Label)
		}
	}
}

// numberBasedSelect is a fallback for systems where arrow keys don't work
func numberBasedSelect(options []Option, defaultIndex int) string {
	// Display options
	for i, option := range options {
		var prefix string
		if i == defaultIndex {
			prefix = Bold + Cyan + " " + strconv.Itoa(i+1) + ". " + Reset
		} else {
			prefix = " " + strconv.Itoa(i+1) + ". "
		}
		fmt.Println(prefix + option.Label)
	}
	
	// Construct default prompt
	defaultPrompt := ""
	if defaultIndex >= 0 && defaultIndex < len(options) {
		defaultPrompt = strconv.Itoa(defaultIndex + 1)
	}
	
	// Get user selection
	for {
		fmt.Println()
		input := Prompt("Enter selection", defaultPrompt)
		
		// Handle default
		if input == "" && defaultPrompt != "" {
			return options[defaultIndex].Value
		}
		
		// Parse selection
		if selection, err := strconv.Atoi(input); err == nil {
			if selection > 0 && selection <= len(options) {
				return options[selection-1].Value
			}
		}
		
		Error("Invalid selection. Please enter a number between 1 and " + strconv.Itoa(len(options)))
	}
}

// ConfirmYesNo asks for confirmation and returns a boolean
func ConfirmYesNo(message string, defaultYes bool) bool {
	// Try arrow key select for yes/no
	yesNoOptions := []Option{
		{Label: "Yes", Value: "yes"},
		{Label: "No", Value: "no"},
	}
	
	defaultIndex := 1 // Default to No
	if defaultYes {
		defaultIndex = 0 // Default to Yes
	}
	
	// Print title and message
	fmt.Println()
	fmt.Println(Bold + message + Reset)
	fmt.Println("Use ↑ and ↓ arrow keys to navigate, Enter to select, or y/n")
	
	// Try to use arrow keys if possible
	if value, ok := tryArrowKeySelect(yesNoOptions, defaultIndex); ok {
		return value == "yes"
	}
	
	// Fall back to traditional y/n input if arrow keys don't work
	defaultStr := "n"
	if defaultYes {
		defaultStr = "y"
	}
	
	// Create the prompt message
	promptMsg := message + " (y/n)"
	
	for {
		input := strings.ToLower(Prompt(promptMsg, defaultStr))
		
		// Handle default
		if input == "" {
			return defaultYes
		}
		
		// Parse response
		if input == "y" || input == "yes" {
			return true
		} else if input == "n" || input == "no" {
			return false
		}
		
		Error("Please enter 'y' or 'n'")
	}
}