package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/dikaio/scribe/internal/config"
	"github.com/dikaio/scribe/internal/content"
	"github.com/dikaio/scribe/internal/render"
	"github.com/dikaio/scribe/internal/sitemap"
)

// Builder handles site building
type Builder struct {
	config   config.Config
	renderer *render.Renderer
	pages    []content.Page
	tags     map[string][]content.Page
	quiet    bool
	devMode  bool
}

// NewBuilder creates a new site builder
func NewBuilder(cfg config.Config) *Builder {
	return &Builder{
		config:   cfg,
		renderer: render.NewRenderer(cfg),
		pages:    []content.Page{},
		tags:     make(map[string][]content.Page),
		quiet:    false,
		devMode:  false,
	}
}

// SetQuiet sets the quiet mode for the builder
func (b *Builder) SetQuiet(quiet bool) {
	b.quiet = quiet
}

// SetDevMode sets the development mode for the builder and renderer
// Development mode disables template caching for live reloading
func (b *Builder) SetDevMode(enabled bool) {
	b.devMode = enabled
	b.renderer.SetDevMode(enabled)
}

// Build builds the site
func (b *Builder) Build(sitePath string) error {
	// Initialize renderer
	if err := b.renderer.Init(sitePath); err != nil {
		return err
	}

	// Load content
	if err := b.loadContent(sitePath); err != nil {
		return err
	}

	// Create output directory
	outputPath := filepath.Join(sitePath, b.config.OutputDir)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	// Copy static files
	if err := b.copyStaticFiles(sitePath, outputPath); err != nil {
		return err
	}

	// Generate pages in parallel
	if err := b.generatePages(outputPath); err != nil {
		return err
	}

	// Generate tag pages
	if err := b.generateTagPages(outputPath); err != nil {
		return err
	}

	// Generate home page
	if err := b.generateHomePage(outputPath); err != nil {
		return err
	}

	// Generate sitemap
	if err := b.generateSitemap(outputPath); err != nil {
		return err
	}

	return nil
}

// loadContent loads all content files
func (b *Builder) loadContent(sitePath string) error {
	contentPath := filepath.Join(sitePath, b.config.ContentDir)
	
	// First, collect all markdown files
	var markdownFiles []string
	err := filepath.Walk(contentPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-markdown files
		if filepath.Ext(path) != ".md" {
			return nil
		}

		markdownFiles = append(markdownFiles, path)
		return nil
	})

	if err != nil {
		return err
	}

	// Create a worker function to load pages in parallel
	worker := func(workerID int, jobs <-chan interface{}, results chan<- interface{}, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		
		for job := range jobs {
			filePath := job.(string)
			
			// Load page
			page, err := content.LoadPage(filePath, b.config.BaseURL, b.config.TrailingSlash)
			if err != nil {
				errChan <- fmt.Errorf("error loading %s: %v", filePath, err)
				continue
			}

			// Skip draft pages in production
			if page.Draft {
				// Skip silently
				continue
			}

			// Add page to results
			results <- page
		}
	}

	// Process all markdown files in parallel
	jobsInterface := make([]interface{}, len(markdownFiles))
	for i, file := range markdownFiles {
		jobsInterface[i] = file
	}

	resultsInterface, errors := parallelExecutor(jobsInterface, worker)
	
	// Check for errors
	if len(errors) > 0 {
		return errors[0]
	}

	// Process results
	for _, result := range resultsInterface {
		page := result.(content.Page)
		
		// Add page to collection
		b.pages = append(b.pages, page)

		// Add page to tags
		for _, tag := range page.Tags {
			b.tags[tag] = append(b.tags[tag], page)
		}
	}
	
	return nil
}

// fileCopyJob represents a file copy operation
type fileCopyJob struct {
	SrcPath string
	DstPath string
}

// copyStaticFiles copies static files to the output directory
func (b *Builder) copyStaticFiles(sitePath, outputPath string) error {
	var copyJobs []interface{}

	// Collect theme static files
	themeStaticPath := filepath.Join(sitePath, "themes", b.config.Theme, "static")
	if _, err := os.Stat(themeStaticPath); err == nil {
		// Get theme files
		themeFiles, err := collectFilesToCopy(themeStaticPath, outputPath)
		if err != nil {
			return err
		}
		copyJobs = append(copyJobs, themeFiles...)
	}

	// Collect site static files (overrides theme files)
	siteStaticPath := filepath.Join(sitePath, b.config.StaticDir)
	if _, err := os.Stat(siteStaticPath); err == nil {
		// Get site files
		siteFiles, err := collectFilesToCopy(siteStaticPath, outputPath)
		if err != nil {
			return err
		}
		copyJobs = append(copyJobs, siteFiles...)
	}

	// Create a worker function to copy files in parallel
	worker := func(workerID int, jobs <-chan interface{}, results chan<- interface{}, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		
		for job := range jobs {
			copyJob := job.(fileCopyJob)
			
			// Create directory if needed
			dir := filepath.Dir(copyJob.DstPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				errChan <- fmt.Errorf("error creating directory for %s: %v", copyJob.DstPath, err)
				continue
			}
			
			// Copy file
			err := copyFile(copyJob.SrcPath, copyJob.DstPath)
			if err != nil {
				errChan <- fmt.Errorf("error copying %s: %v", copyJob.SrcPath, err)
				continue
			}
			
			// Report success
			rel, _ := filepath.Rel(outputPath, copyJob.DstPath)
			results <- rel
		}
	}

	// Execute jobs in parallel
	_, errors := parallelExecutor(copyJobs, worker)
	
	// Check for errors
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

// collectFilesToCopy collects files to copy from source directory to destination
func collectFilesToCopy(srcDir, dstDir string) ([]interface{}, error) {
	var jobs []interface{}
	
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories (they'll be created when copying files)
		if info.IsDir() {
			return nil
		}

		// Get relative path
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		// Create job
		jobs = append(jobs, fileCopyJob{
			SrcPath: path,
			DstPath: filepath.Join(dstDir, rel),
		})

		return nil
	})

	return jobs, err
}

// generatePages generates all content pages
func (b *Builder) generatePages(outputPath string) error {
	// Define a page rendering job
	type pageRenderJob struct {
		Page       content.Page
		OutputFile string
	}

	// Create a worker function to render pages in parallel
	worker := func(workerID int, jobs <-chan interface{}, results chan<- interface{}, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		
		for job := range jobs {
			renderJob := job.(pageRenderJob)
			
			// Create the directory for this page
			dir := filepath.Dir(renderJob.OutputFile)
			if err := os.MkdirAll(dir, 0755); err != nil {
				errChan <- fmt.Errorf("error creating directory for %s: %v", renderJob.Page.URL, err)
				continue
			}
			
			// Render the page
			err := b.renderer.RenderPage(renderJob.Page, renderJob.OutputFile)
			if err != nil {
				errChan <- fmt.Errorf("error rendering %s: %v", renderJob.Page.URL, err)
				continue
			}
			
			// Report success
			results <- renderJob.Page.URL
		}
	}

	// Create jobs for all pages
	jobs := make([]interface{}, len(b.pages))
	for i, page := range b.pages {
		// Clean URL for file path construction
		cleanURL := strings.TrimSuffix(page.URL, "/")
		outputFile := filepath.Join(outputPath, cleanURL, "index.html")
		jobs[i] = pageRenderJob{
			Page:       page,
			OutputFile: outputFile,
		}
	}

	// Execute jobs in parallel
	_, errors := parallelExecutor(jobs, worker)
	
	// Check for errors
	if len(errors) > 0 {
		return errors[0]
	}

	return nil
}

// generateTagPages generates tag listing pages
func (b *Builder) generateTagPages(outputPath string) error {
	// Create tags directory
	tagsDir := filepath.Join(outputPath, "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	// Generate main tags index
	allTags := make([]string, 0, len(b.tags))
	for tag := range b.tags {
		allTags = append(allTags, tag)
	}
	sort.Strings(allTags)

	// Define a tag page rendering job
	type tagRenderJob struct {
		Tag        string
		Pages      []content.Page
		OutputFile string
		Title      string
	}

	// Create jobs for each tag
	jobs := make([]interface{}, 0, len(b.tags))
	for tag, pages := range b.tags {
		// Sort pages by date (newest first)
		sort.Slice(pages, func(i, j int) bool {
			return pages[i].Date.After(pages[j].Date)
		})

		// Create tag directory
		tagDir := filepath.Join(tagsDir, tag)
		if err := os.MkdirAll(tagDir, 0755); err != nil {
			return err
		}

		// Add job
		outputFile := filepath.Join(tagDir, "index.html")
		title := fmt.Sprintf("Tag: %s", tag)
		jobs = append(jobs, tagRenderJob{
			Tag:        tag,
			Pages:      pages,
			OutputFile: outputFile,
			Title:      title,
		})
	}

	// Create a worker function to render tag pages in parallel
	worker := func(workerID int, jobs <-chan interface{}, results chan<- interface{}, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()
		
		for job := range jobs {
			renderJob := job.(tagRenderJob)
			
			// Render tag page
			err := b.renderer.RenderList(renderJob.Title, renderJob.Pages, renderJob.OutputFile)
			if err != nil {
				errChan <- fmt.Errorf("error rendering tag page %s: %v", renderJob.Tag, err)
				continue
			}
			
			// Report success
			results <- renderJob.Tag
		}
	}

	// Execute jobs in parallel
	_, errors := parallelExecutor(jobs, worker)
	
	// Check for errors
	if len(errors) > 0 {
		return errors[0]
	}

	// Removed generation reporting

	return nil
}

// generateHomePage generates the home page
func (b *Builder) generateHomePage(outputPath string) error {
	// Filter and sort posts (newest first)
	posts := []content.Page{}
	for _, page := range b.pages {
		if page.IsPost {
			posts = append(posts, page)
		}
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	// Render home page
	outputFile := filepath.Join(outputPath, "index.html")
	return b.renderer.RenderHome(posts, outputFile)
}

// copyDir recursively copies a directory tree
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from source
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Skip source root
		if rel == "." {
			return nil
		}

		// Get destination path
		dstPath := filepath.Join(dst, rel)

		// Create directories
		if info.IsDir() {
			return os.MkdirAll(dstPath, 0755)
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	return err
}

// generateSitemap generates a sitemap.xml file for the site
func (b *Builder) generateSitemap(outputPath string) error {
	// Create sitemap generator
	generator := sitemap.NewGenerator(b.config)

	// Generate sitemap.xml
	sitemapPath := filepath.Join(outputPath, "sitemap.xml")

	// Only include non-draft pages in the sitemap
	allPages := b.pages

	// Log sitemap generation if not in quiet mode
	if !b.quiet {
		fmt.Println("Generating sitemap.xml...")
	}

	// Generate the sitemap
	err := generator.Generate(allPages, sitemapPath)
	if err != nil {
		return fmt.Errorf("failed to generate sitemap: %w", err)
	}

	if !b.quiet {
		fmt.Println("Sitemap generated successfully at sitemap.xml")
	}

	return nil
}
