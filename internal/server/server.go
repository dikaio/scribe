package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dikaio/scribe/internal/build"
	"github.com/dikaio/scribe/internal/config"
)

// Server represents the development server
type Server struct {
	config  config.Config
	builder *build.Builder
	watcher *build.Watcher
	port    int
	quiet   bool
}

// NewServer creates a new development server
func NewServer(cfg config.Config, port int, quiet bool) *Server {
	builder := build.NewBuilder(cfg)
	builder.SetQuiet(quiet)
	
	return &Server{
		config:  cfg,
		builder: builder,
		port:    port,
		quiet:   quiet,
	}
}

// Start starts the development server
func (s *Server) Start(sitePath string) error {
	// Create output directory if it doesn't exist
	outputPath := filepath.Join(sitePath, s.config.OutputDir)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	// Show initial message
	if !s.quiet {
		fmt.Printf("Starting development server for site from '%s'...\n", sitePath)
	}
	
	// Build site initially
	if err := s.builder.Build(sitePath); err != nil {
		return err
	}

	// Create watcher with quiet mode
	s.watcher = build.NewWatcher(s.builder, sitePath)
	s.watcher.SetQuiet(s.quiet)

	// Start file watcher in background
	go func() {
		err := s.watcher.Watch(time.Second, func() error {
			return s.builder.Build(sitePath)
		})
		if err != nil {
			log.Printf("Watcher error: %s\n", err)
		}
	}()

	// Start HTTP server - always show the URL regardless of quiet mode
	fmt.Printf("Server running at http://localhost:%d/\n", s.port)
	if !s.quiet {
		fmt.Println("Watching for changes. Press Ctrl+C to stop.")
	}
	
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), http.FileServer(http.Dir(outputPath)))
}
