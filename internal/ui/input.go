package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ValidatorFunc is a function that validates input
type ValidatorFunc func(string) error

// PromptWithValidation displays a prompt, validates the input, and returns the validated input
func PromptWithValidation(message string, defaultValue string, validator ValidatorFunc) string {
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

// MultiLinePrompt captures multi-line input from the user
// The input is terminated when the user enters a line containing only a dot (.)
func MultiLinePrompt(message string) string {
	fmt.Println(Bold + message + Reset + " (end with a single '.' on a line)")
	fmt.Println("─────────────────")
	
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for {
		scanner.Scan()
		text := scanner.Text()
		
		// Check for termination
		if text == "." {
			break
		}
		
		lines = append(lines, text)
	}
	
	fmt.Println("─────────────────")
	return strings.Join(lines, "\n")
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

// Required returns a validator that checks if input is not empty
func Required(fieldName string) ValidatorFunc {
	return func(input string) error {
		if input == "" {
			return fmt.Errorf("%s is required", fieldName)
		}
		return nil
	}
}