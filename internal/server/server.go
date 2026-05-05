package server

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"devmemory/internal/config"
	"devmemory/internal/service"
)

//go:embed static
var staticFS embed.FS

type Server struct {
	svc    *service.MemoryService
	port   int
	router *http.ServeMux
}

func New(svc *service.MemoryService, port int) *Server {
	srv := &Server{
		svc:    svc,
		port:   port,
		router: http.NewServeMux(),
	}
	srv.setupRoutes()
	return srv
}

func (s *Server) setupRoutes() {
	// API routes
	s.router.HandleFunc("GET /api/health", s.handleHealth)
	s.router.HandleFunc("GET /api/entries", s.handleListEntries)
	s.router.HandleFunc("POST /api/entries", s.handleCreateEntry)
	s.router.HandleFunc("GET /api/entries/{id}", s.handleGetEntry)
	s.router.HandleFunc("PUT /api/entries/{id}", s.handleUpdateEntry)
	s.router.HandleFunc("DELETE /api/entries/{id}", s.handleDeleteEntry)
	s.router.HandleFunc("GET /api/search", s.handleSearch)
	s.router.HandleFunc("GET /api/today", s.handleToday)
	s.router.HandleFunc("POST /api/entries/{id}/copy", s.handleCopy)
	s.router.HandleFunc("GET /api/export/today", s.handleExportToday)
	s.router.HandleFunc("GET /api/export/json", s.handleExportJSON)
	s.router.HandleFunc("POST /api/import/json", s.handleImportJSON)

	// Static files
	staticContent, _ := fs.Sub(staticFS, "static")
	s.router.Handle("/", http.FileServer(http.FS(staticContent)))
}

func (s *Server) Start(openBrowser bool) error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	url := fmt.Sprintf("http://%s", addr)
	log.Printf("DevMemory server running at %s", url)

	if openBrowser {
		go openURL(url)
	}

	httpServer := &http.Server{
		Handler:      corsMiddleware(s.router),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	httpServer.Close()
	return nil
}

func openURL(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	time.Sleep(300 * time.Millisecond)
	if err := execCommand(cmd, args...); err != nil {
		log.Printf("Could not open browser: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ConfigPath returns the default config for the server.
func ConfigPath() string {
	cfg := config.DefaultConfig()
	return cfg.DataDir
}
