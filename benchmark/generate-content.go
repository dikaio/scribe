package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// How many files to generate
	numFiles := 1000 // Generate 1000 files for a comprehensive benchmark

	// Base directory for posts
	postsDir := filepath.Join("content", "posts")

	// Template for markdown content
	template := `---
title: %s
description: This is a benchmark post %d
date: %s
tags:
  - tag%d
  - benchmark
  - test
draft: false
---

# %s

This is a benchmark post to test the performance of Scribes.

## Section 1

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam
euismod, nisl eget ultricies tincidunt, nisl nisl aliquam nisl, eget
ultricies nisl nisl eget nisl. Nullam euismod, nisl eget ultricies
tincidunt, nisl nisl aliquam nisl, eget ultricies nisl nisl eget nisl.

## Section 2

- List item 1
- List item 2
- List item 3

## Section 3

[Link to example.com](https://example.com)

## Code Example

` + "```" + `
func main() {
    fmt.Println("Hello, World!")
}
` + "```" + `
`

	// Generate files
	for i := 1; i <= numFiles; i++ {
		title := fmt.Sprintf("Benchmark Post %d", i)
		date := time.Now().AddDate(0, 0, -i).Format(time.RFC3339)
		content := fmt.Sprintf(template, title, i, date, i%10, title)
		
		filename := fmt.Sprintf("post-%04d.md", i)
		filepath := filepath.Join(postsDir, filename)
		
		err := os.WriteFile(filepath, []byte(content), 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", filepath, err)
		} else {
			fmt.Printf("Generated %s\n", filepath)
		}
	}
	
	fmt.Printf("Generated %d benchmark posts\n", numFiles)
}