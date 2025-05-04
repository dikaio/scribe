package cli

import (
	"fmt"

	"github.com/dikaio/scribe/internal/content"
	"github.com/dikaio/scribe/internal/ui"
)

// createPostEnhanced is an enhanced version of createNewPost that uses
// the improved UI components for a better interactive experience
func (a *App) createPostEnhanced(initialTitle string) error {
	ui.Header("Create New Post")

	// Prompt for post title
	title := initialTitle
	if title == "" {
		title = ui.PromptWithValidation("Post Title", "", ui.Required("Post title"))
	}

	// Prompt for description (optional)
	description := ui.Prompt("Description (optional)", "")

	// Prompt for tags
	defaultTags := []string{"uncategorized"}
	tags := ui.PromptTags("Tags (comma-separated)", defaultTags)

	// Prompt for draft status
	draft := ui.ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create post
	filePath, err := creator.CreateContent(content.PostType, title, description, tags, draft)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	ui.Success(fmt.Sprintf("Post created successfully: %s", filePath))
	return nil
}

// createPageEnhanced is an enhanced version of createNewPage that uses
// the improved UI components for a better interactive experience
func (a *App) createPageEnhanced(initialTitle string) error {
	ui.Header("Create New Page")

	// Prompt for page title
	title := initialTitle
	if title == "" {
		title = ui.PromptWithValidation("Page Title", "", ui.Required("Page title"))
	}

	// Prompt for description (optional)
	description := ui.Prompt("Description (optional)", "")

	// Prompt for draft status
	draft := ui.ConfirmYesNo("Save as draft?", false)

	// Create content creator
	creator := content.NewCreator(".")

	// Create page
	filePath, err := creator.CreateContent(content.PageType, title, description, nil, draft)
	if err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}

	ui.Success(fmt.Sprintf("Page created successfully: %s", filePath))
	return nil
}