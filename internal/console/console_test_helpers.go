package console

import (
	"github.com/dikaio/scribe/internal/content"
)

// Mock loadSiteStats for testing
func (c *Console) mockLoadSiteStats(contentPath string) ([]content.Page, []content.Page, error) {
	// Create mock posts and pages
	posts := []content.Page{MockPost()}
	pages := []content.Page{MockPage()}
	
	return posts, pages, nil
}

// OverrideLoadSiteStats replaces the loadSiteStats method with a mock implementation for testing
func (c *Console) OverrideLoadSiteStats() {
	// Store a reference to the original method for later tests if needed
	c.originalLoadSiteStats = c.loadSiteStats
	
	// Replace with mock method
	c.loadSiteStats = c.mockLoadSiteStats
}

// RestoreLoadSiteStats restores the original loadSiteStats method
func (c *Console) RestoreLoadSiteStats() {
	if c.originalLoadSiteStats != nil {
		c.loadSiteStats = c.originalLoadSiteStats
		c.originalLoadSiteStats = nil
	}
}