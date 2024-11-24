package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/kylods/kbot-backend/internal/database"
	"github.com/kylods/kbot-backend/internal/discordclient"
)

const version string = "INDEV"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env, quitting...")
	}

	log.Printf("Starting KBot Server %s\n", version)

	// Initialize Discord client
	discordClient := discordclient.NewDiscordClient(os.Getenv("DISCORD_BOT_TOKEN"))
	go discordClient.Run() // Run Discord client in a separate goroutine

	// Initialize database
	db, err := database.Connect(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get SQL DB conneection: %v", err)
	}
	defer sqlDB.Close()

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)

	select {
	case <-interruptCh:
		log.Println("Received interrupt signal, closing...")
		os.Exit(0)
	}

}
