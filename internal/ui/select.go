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
	// Print title and description
	Title(title, description)
	
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