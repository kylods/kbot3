package discordclient

import (
	"log"
	"math/rand/v2"
	"runtime"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jogramming/dca"
	"github.com/kylods/kbot3/pkg/models"
	"gorm.io/gorm"
)

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, m *discordgo.MessageCreate, c *Client, gConfig *models.Guild)
}

var commands = []Command{
	{
		Name:        "help",
		Description: "Returns a list of KBot's commands",
		Handler:     commandHelpHandler,
	},
	{
		Name:        "setprefix",
		Description: "Set the command prefix for the server",
		Handler:     commandSetprefixHandler,
	},
	{
		Name:        "download",
		Description: "Provides a download for KBot Media Player",
		Handler:     commandDownloadHandler,
	},
	{
		Name:        "about",
		Description: "Provides some information about KBot",
		Handler:     commandAboutHandler,
	},
}

type Client struct {
	session  *discordgo.Session
	db       *gorm.DB
	version  string
	commands []Command
	ready    bool
}

func NewDiscordClient(token string, version string, db *gorm.DB) *Client {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Could not initialize Discord session")
	}

	return &Client{
		session:  session,
		version:  version,
		db:       db,
		commands: commands,
		ready:    false,
	}
}

func (c *Client) Run() {
	// Add event handlers
	c.session.AddHandler(c.messageCreate)
	c.session.AddHandler(c.readyHandler)
	c.session.AddHandler(c.createGuildHandler)

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

func (c *Client) createGuildHandler(s *discordgo.Session, g *discordgo.GuildCreate) {
	/*This event can be sent in three different scenarios:

	    When a user is initially connecting, to lazily load and backfill information for all unavailable guilds sent in the Ready event. Guilds that are unavailable due to an outage will send a Guild Delete event.
	    When a Guild becomes available again to the client.
	    When the current user joins a new Guild.

		During an outage, the guild object in scenarios 1 and 3 may be marked as unavailable.
	*/

	configTemplate := models.Guild{
		GuildID:        g.ID,
		Name:           g.Name,
		CommandPrefix:  '!',
		CommandChannel: "",
		DjRoles:        "",
		LoopEnabled:    false,
	}
	var gConfig models.Guild

	c.db.Where(&models.Guild{GuildID: g.ID}).Attrs(configTemplate).FirstOrCreate(&gConfig)

	log.Printf("Initialized Guild %v", gConfig.Name)
}

func (c *Client) readyHandler(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged into Discord as %s", r.User.String())
	c.ready = true
}

func (c *Client) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !c.ready {
		log.Println("Incoming message ignored, still initializing...")
		return
	}

	if m.Author.Bot {
		return
	}

	// This 'should' equate to a DM... IDK when else it'd happen
	if m.GuildID == "" {
		easterEggDM(s, m)
		return
	}

	// Check for command prefix
	if len(m.Content) < 2 { // Commands need at least a prefix & one letter.
		log.Print(len(m.Content))
		return
	}

	var gConfig models.Guild
	c.db.Where(&models.Guild{GuildID: m.GuildID}).First(&gConfig)

	if m.Content[0] == byte(gConfig.CommandPrefix) {
		msgSlice := strings.Split(m.Content, " ")

		for _, cmd := range commands {
			if msgSlice[0][1:] == cmd.Name {
				log.Printf("Processing command from %s in %s: %s", m.Author.Username, gConfig.Name, cmd.Name)
				cmd.Handler(s, m, c, &gConfig)
				return
			}
		}
		log.Printf("Couldn't find command from %s: %s", m.Author.Username, msgSlice[0][1:])
	}
}

func easterEggDM(s *discordgo.Session, m *discordgo.MessageCreate) {
	var msg string
	const maxRoll int = 10000
	roll := rand.IntN(maxRoll) + 1

	maxRollStr := strconv.Itoa(maxRoll)
	rollStr := strconv.Itoa(roll)

	if roll == maxRoll {
		msg = "<https://www.youtube.com/watch?v=dQw4w9WgXcQ>"
	} else {
		msg = "You rolled " + rollStr + " out of " + maxRollStr + "!\n\n*Try for " + maxRollStr + "!*"
	}

	s.ChannelMessageSend(m.ChannelID, msg)
}

func commandHelpHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *Client, gConfig *models.Guild) {
	reply := "**KBot Commands**"
	for _, cmd := range c.commands {
		reply += "\n`" + cmd.Name + "`: " + cmd.Description
	}
	s.ChannelMessageSend(m.ChannelID, reply)
}

func commandAboutHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *Client, gConfig *models.Guild) {
	reply := "# KBot " + c.version + `
	A Discord extension for KBot Media Player, written by Kuelos

	Github: <https://github.com/kylods/kbot>

	Runtime: ` + "`" + runtime.Version() + "`" + `
	Discordgo: ` + "`" + discordgo.VERSION + "`" + `
	dca: ` + "`" + dca.LibraryVersion + "`" + `
	`

	s.ChannelMessageSend(m.ChannelID, reply)
}

func commandSetprefixHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *Client, gConfig *models.Guild) {
	stringSlice := strings.Split(m.Content, " ")
	if len(stringSlice) < 2 {
		s.ChannelMessageSend(m.ChannelID, "No prefix given")
		return
	}
	if len(stringSlice[1]) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Prefix must be a single character")
		return
	}

	gConfig.CommandPrefix = rune(stringSlice[1][0])
	c.db.Save(gConfig)

	s.ChannelMessageSend(m.ChannelID, "Updated command prefix of "+gConfig.Name+" to `"+string(gConfig.CommandPrefix)+"`")
}

func commandDownloadHandler(s *discordgo.Session, m *discordgo.MessageCreate, c *Client, gConfig *models.Guild) {
	s.ChannelMessageSend(m.ChannelID, "need to build this still :)")
}
