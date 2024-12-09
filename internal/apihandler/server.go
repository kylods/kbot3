package apihandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kylods/kbot3/internal/discordclient"
	"github.com/kylods/kbot3/internal/middleware"

	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	discordClient *discordclient.Client
	db            *gorm.DB
	httpServer    *http.Server
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Upload hit!")

	// Parse our multipart form, 10 << 24 specifies a max
	// upload of 160MB files
	r.ParseMultipartForm(10 << 24)
	// FormFile returns the first file for the given key 'audioFile'
	// It also returns the FileHeader so we can grab the Filename,
	// the header & the size of the file
	file, handler, err := r.FormFile("audioFile")
	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MINE Header: %+v\n", handler.Header)

	// Create a temp file that follows a particular naming pattern
	tempFile, err := os.CreateTemp("temp-audio", "upload-")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// Read the contents of our file into a byte array,
	// then write it to the tempfile
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	respondJSON(w, 202, map[string]string{"message": "Upload Successful"})
}

// Initializes a new APi server
func NewServer(port string, discordClient *discordclient.Client, db *gorm.DB) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("GET /auth", authenticateHandler)
	mux.HandleFunc("POST /upload", uploadFile)

	server := http.Server{
		Addr:         ":" + port,
		Handler:      middleware.Logging(mux),
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

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, "To Be Implemented")
}
