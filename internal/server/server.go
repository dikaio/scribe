package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dikaio/scribes/internal/build"
	"github.com/dikaio/scribes/internal/config"
)

// Server represents the development server
type Server struct {
	config  config.Config
	builder *build.Builder
	watcher *build.Watcher
	port    int
}

// NewServer creates a new development server
func NewServer(cfg config.Config, port int) *Server {
	return &Server{
		config:  cfg,
		builder: build.NewBuilder(cfg),
		port:    port,
	}
}

// Start starts the development server
func (s *Server) Start(sitePath string) error {
	// Create output directory if it doesn't exist
	outputPath := filepath.Join(sitePath, s.config.OutputDir)
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return err
	}

	// Build site initially
	fmt.Println("Building site...")
	if err := s.builder.Build(sitePath); err != nil {
		return err
	}

	// Create watcher
	s.watcher = build.NewWatcher(s.builder, sitePath)

	// Start file watcher in background
	go func() {
		err := s.watcher.Watch(time.Second, func() error {
			return s.builder.Build(sitePath)
		})
		if err != nil {
			log.Printf("Watcher error: %s\n", err)
		}
	}()

	// Start HTTP server
	fmt.Printf("Server running at http://localhost:%d/\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), http.FileServer(http.Dir(outputPath)))
}
