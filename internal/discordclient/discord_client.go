package discordclient

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	session *discordgo.Session
}

func NewDiscordClient(token string) *Client {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Could not initialize Discord session")
	}

	return &Client{session: session}
}

func (c *Client) Run() {
	// Add event handlers
	c.session.AddHandler(messageCreate)

	// Open a connection to Discord
	err := c.session.Open()
	if err != nil {
		log.Fatalf("Failed to open Discord session: %v", err)
	}
	log.Println("Discord client is running")
}

func (c *Client) Close() {
	c.session.Close()
}

func (c *Client) SendChannelMessage(channelId string, message string) {
	c.session.ChannelMessageSend(channelId, message)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// Handle messages here
	log.Printf("Message from %s: %s", m.Author.Username, m.Content)
}
