package apihandler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kylods/kbot3/pkg/models"
)

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

func (s *Server) wsDiscordHandler(w http.ResponseWriter, r *http.Request) {
	guildID := r.PathValue("id")

	var gConfig models.Guild

	s.db.Where(&models.Guild{GuildID: guildID}).First(&gConfig)
}

func (s *Server) getDiscordGuildsHandler(w http.ResponseWriter, r *http.Request) {
	guildSlice := s.discordClient.GetGuilds()
	respondJSON(w, 200, guildSlice)
}
