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
	
	// Print options with improved visual formatting
	fmt.Println("Choose an option:")
	for i, option := range options {
		var marker string
		if i == defaultIndex {
			marker = Bold + Cyan + " ● " + Reset
		} else {
			marker = " ○ "
		}
		
		suffix := ""
		if i == defaultIndex {
			suffix = Cyan + " (default)" + Reset
		}
		
		fmt.Printf("%s%d. %s%s\n", 
			marker, 
			i+1, 
			option.Label,
			suffix)
	}
	
	// Construct default prompt
	defaultPrompt := ""
	if defaultIndex >= 0 && defaultIndex < len(options) {
		defaultPrompt = strconv.Itoa(defaultIndex + 1)
	}
	
	// Get user selection
	for {
		fmt.Println()
		input := Prompt("Enter option number", defaultPrompt)
		
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
	defaultStr := "n"
	if defaultYes {
		defaultStr = "y"
	}
	
	// Create the prompt message with visual indicators
	fmt.Println()
	fmt.Println(Bold + message + Reset)
	
	// Show options with visual indicators
	if defaultYes {
		fmt.Println(Bold + Cyan + " ● " + Reset + "Yes" + Cyan + " (default)" + Reset)
		fmt.Println(" ○ No")
	} else {
		fmt.Println(" ○ Yes")
		fmt.Println(Bold + Cyan + " ● " + Reset + "No" + Cyan + " (default)" + Reset)
	}
	
	// Prompt for y/n
	fmt.Println()
	
	for {
		input := strings.ToLower(Prompt("Enter y/n", defaultStr))
		
		// Handle default
		if input == "" {
			return defaultYes
		}
		
		// Parse response
		if input == "y" || input == "yes" || input == "1" {
			return true
		} else if input == "n" || input == "no" || input == "2" {
			return false
		}
		
		Error("Please enter 'y' or 'n'")
	}
}