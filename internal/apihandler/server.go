package apihandler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/kylods/kbot-backend/internal/discordclient"

	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	discordClient *discordclient.Client
	db            *gorm.DB
	httpServer    *http.Server
}

// Initializes a new APi server
func NewServer(port string, discordClient *discordclient.Client, db *gorm.DB) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	output := Server{
		discordClient: discordClient,
		db:            db,
		httpServer:    &server,
	}

	return &output
}
func (s *Server) Start() {
	log.Printf("API server is running on port %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("API server failed: %v", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down API server...")
	return s.httpServer.Shutdown(ctx)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"message": "Pong!"})
}

func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Failed to write JSON response: %v", err)
	}
}
