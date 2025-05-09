package cli

import (
	"fmt"
	"strings"
)

// Simplified UI functions

// Header prints a header
func Header(text string) {
	fmt.Printf("\n=== %s ===\n\n", text)
}

// Info prints an info message
func Info(message string) {
	fmt.Println(message)
}

// Success prints a success message
func Success(message string) {
	fmt.Printf("Success: %s\n", message)
}

// Warning prints a warning message
func Warning(message string) {
	fmt.Printf("Warning: %s\n", message)
}

// Error prints an error message
func Error(message string) {
	fmt.Printf("Error: %s\n", message)
}

// Prompt displays a prompt and returns the user input
func Prompt(message string, defaultValue string) string {
	var promptText string
	if defaultValue != "" {
		promptText = fmt.Sprintf("%s [%s]: ", message, defaultValue)
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

// PromptWithValidation displays a prompt, validates the input, and returns the validated input
func PromptWithValidation(message string, defaultValue string, validator func(string) error) string {
	for {
		input := Prompt(message, defaultValue)
		
		// Use default value if input is empty
		if input == "" && defaultValue != "" {
			input = defaultValue
		}
		
		// Validate input
		if validator != nil {
			if err := validator(input); err != nil {
				Error(err.Error())
				continue
			}
		}
		
		return input
	}
}

// Required returns a validator that checks if input is not empty
func Required(fieldName string) func(string) error {
	return func(input string) error {
		if input == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
		return nil
	}
}

// PromptTags captures a comma-separated list of tags
func PromptTags(message string, defaultTags []string) []string {
	// Format default tags for display
	defaultTagsStr := strings.Join(defaultTags, ", ")
	
	// Prompt for input
	input := Prompt(message, defaultTagsStr)
	
	// Use default if empty
	if input == "" {
		return defaultTags
	}
	
	// Parse tags
	result := make([]string, 0)
	tags := strings.Split(input, ",")
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	
	// If no valid tags were entered, use default
	if len(result) == 0 {
		return defaultTags
	}
	
	return result
}

// ConfirmYesNo asks for confirmation and returns a boolean
func ConfirmYesNo(message string, defaultYes bool) bool {
	defaultStr := "n"
	if defaultYes {
		defaultStr = "y"
	}
	
	// Prompt for y/n
	for {
		input := strings.ToLower(Prompt(fmt.Sprintf("%s (y/n)", message), defaultStr))
		
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

// Divider prints a horizontal divider
func Divider() {
	fmt.Println(strings.Repeat("-", 40))
}