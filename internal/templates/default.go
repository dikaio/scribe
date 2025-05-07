package templates

import (
	"log"
	"sync"
)

var (
	// Default template strings
	BaseTemplate   string
	SingleTemplate string
	ListTemplate   string
	HomeTemplate   string
	PageTemplate   string
	StyleCSS       string

	// Initialization once
	defaultTemplatesOnce sync.Once
)

// loadDefaultTemplates loads all default templates from embedded files
func loadDefaultTemplates() {
	var err error

	BaseTemplate, err = GetDefaultTemplate("base.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded base template: %v", err)
	}

	SingleTemplate, err = GetDefaultTemplate("single.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded single template: %v", err)
	}

	ListTemplate, err = GetDefaultTemplate("list.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded list template: %v", err)
	}

	HomeTemplate, err = GetDefaultTemplate("home.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded home template: %v", err)
	}

	PageTemplate, err = GetDefaultTemplate("page.html")
	if err != nil {
		log.Printf("Warning: Failed to load embedded page template: %v", err)
	}

	StyleCSS, err = GetDefaultTemplate("style.css")
	if err != nil {
		log.Printf("Warning: Failed to load embedded CSS: %v", err)
	}
}

// Default templates for site creation - will load from embedded files
func init() {
	defaultTemplatesOnce.Do(loadDefaultTemplates)
}

// SampleContent contains default content for new sites

// SamplePost is the default welcome post
const SamplePost = `---
title: Welcome to Scribe
description: Lorem ipsum dolor sit amet, consectetur adipiscing elit
date: 2025-05-01T12:00:00Z
tags:
  - welcome
  - scribe
draft: false
---

# Welcome to Scribe!

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam in dui mauris. Vivamus hendrerit arcu sed erat molestie vehicula. Sed auctor neque eu tellus rhoncus ut eleifend nibh porttitor.

## Lorem Ipsum

Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur.

- Consectetur adipiscing elit
- Sed auctor neque eu tellus
- Vivamus hendrerit arcu
- Nam tincidunt congue enim

Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante. Donec eu libero sit amet quam egestas semper. Aenean ultricies mi vitae est. Mauris placerat eleifend leo.

## Action Button

[Action](#) 

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam in dui mauris.
`

// SamplePage is the default about page
const SamplePage = `---
title: About
description: Information about Scribe
draft: false
---

# About Scribe

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Ut elit tellus, luctus nec ullamcorper mattis, pulvinar dapibus leo. Sed non mauris vitae erat consequat auctor eu in elit.

## Our Mission

Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Vestibulum tortor quam, feugiat vitae, ultricies eget, tempor sit amet, ante. Donec eu libero sit amet quam egestas semper.

## Contact

[Action](#)

Mauris placerat eleifend leo. Quisque sit amet est et sapien ullamcorper pharetra. Vestibulum erat wisi, condimentum sed, commodo vitae, ornare sit amet, wisi.
`