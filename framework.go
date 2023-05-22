package mtn_go_discord_framework

import (
	"log"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	discordgo.ApplicationCommand
	Handler func(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string)
}

type ButtonHandler struct {
	CustomID string
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string)
}

var (
	commandsToRegister = make([]SlashCommand, 0)
	commandsMap        = make(map[string]SlashCommand)
	handlerMap         = make(map[string]ButtonHandler)
	initCommandsOnce   sync.Once

	debug          bool
	testingGuildID string
	token          string

	discordSession *discordgo.Session
	ready          = false
	initDone       = false

	systemBusy = false
	mutex      = &sync.Mutex{}
)

const (
	EphemeralFlag = 64
)

// Initializes the framework and returns a discord session, needed for the other functions
func InitFramework(debugMode bool, testingGuildId string, botToken string) *discordgo.Session {
	debug = debugMode
	testingGuildID = testingGuildId
	token = botToken

	discord, err := discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatal(err)
	}
	discordSession = discord
	discordSession.AddHandler(handleCommand)

	ready = true
	return discordSession
}

// Registers a slash command with the framework
func RegisterSlashCommandWithFramework(command SlashCommand) {
	if !ready {
		log.Println("Framework not ready yet, cannot register command")
		return
	}
	if initDone {
		log.Println("Framework already initialized, cannot register command")
		return
	}
	commandsToRegister = append(commandsToRegister, command)
}

// Registers multiple slash commands with the framework
func RegisterSlashCommandsWithFramework(commands []SlashCommand) {
	if !ready {
		log.Println("Framework not ready yet, cannot register commands")
		return
	}
	if initDone {
		log.Println("Framework already initialized, cannot register commands")
		return
	}
	commandsToRegister = append(commandsToRegister, commands...)
}

// Registers a button handler with the framework
func RegisterButtonHandlerWithFramework(handler ButtonHandler) {
	if !ready {
		log.Println("Framework not ready yet, cannot register command")
		return
	}
	if initDone {
		log.Println("Framework already initialized, cannot register command")
		return
	}
	handlerMap[handler.CustomID] = handler
}

// Registers multiple button handlers with the framework
func RegisterButtonHandlersWithFramework(handlers []ButtonHandler) {
	if !ready {
		log.Println("Framework not ready yet, cannot register commands")
		return
	}
	if initDone {
		log.Println("Framework already initialized, cannot register commands")
		return
	}
	for _, handler := range handlers {
		handlerMap[handler.CustomID] = handler
	}
}

// Launches the framework and registers all commands and handlers
func StartFramework() {
	if !ready {
		log.Println("Framework not ready yet, cannot start it. Call InitFramework first")
		return
	}
	// check if discord session is initialized if not initialize it
	if discordSession == nil {
		log.Println("Discord session not initialized, initializing it now")
		InitFramework(debug, testingGuildID, token)
	}

	initDone = true
	initCommandsOnce.Do(func() {
		for _, command := range commandsToRegister {
			commandsMap[command.Name] = command
		}
		log.Println("Initialized commands")
	})
	err := discordSession.Open()
	if err != nil {
		log.Fatal(err)
	}
	registerCommands(discordSession)
}

// Shuts down the framework and closes the discord session
func ShutdownFramework() {
	if !ready || !initDone {
		log.Println("Framework not started, cannot shut it down")
		return
	}
	deleteCommands(discordSession)
	discordSession.Close()
	initDone = false
	// reset initCommandsOnce
	initCommandsOnce = sync.Once{}
}

func handleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {

	// check if i type is ApplicationCommand (slash command)
	case discordgo.InteractionApplicationCommand:
		if command, ok := commandsMap[i.ApplicationCommandData().Name]; ok {
			command.Handler(s, i)
			return
		}
		log.Printf("Unknown command '%s'", i.ApplicationCommandData().Name)

	// check if i type is MessageComponent
	case discordgo.InteractionMessageComponent:
		// check if i message component type is Button
		if i.MessageComponentData().ComponentType == discordgo.ButtonComponent {
			// check if button handler exists
			if handler, ok := handlerMap[i.MessageComponentData().CustomID]; ok {
				handler.Handler(s, i)
				return
			}
			// if handler does not exist, try to find handler without customid suffix in syntax "handler-customid"
			splitstrings := strings.Split(i.MessageComponentData().CustomID, "-")
			if len(splitstrings) > 1 {
				if handler, ok := handlerMap[splitstrings[0]]; ok {
					handler.Handler(s, i, splitstrings[1])
					return
				}
			}
			log.Printf("Unknown button '%s'", i.MessageComponentData().CustomID)
		}
		log.Printf("Unknown message component type '%v'", i.MessageComponentData().ComponentType)
	default:
		log.Printf("Unknown interaction type '%s'", i.Type)
	}
}

func registerCommands(s *discordgo.Session) {
	guildid := ""
	if debug {
		log.Println("Registering commands in DEBUG mode")
		guildid = testingGuildID
	}
	for name, command := range commandsMap {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildid, &command.ApplicationCommand)
		if err != nil {
			log.Printf("Cannot create '%s' command: %v", command.Name, err)
			continue
		}
		// Update the command in the global commands map
		command.ApplicationCommand = *cmd
		commandsMap[name] = command
	}
	log.Println("Registered commands")
}

func deleteCommands(s *discordgo.Session) {
	guildid := ""
	if debug {
		log.Println("Deleting commands in DEBUG mode")
		guildid = testingGuildID
	}
	for _, command := range commandsMap {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildid, command.ApplicationCommand.ID)
		if err != nil {
			log.Printf("Cannot delete '%s' command: %v", command.Name, err)
		}
	}
	log.Println("Deleted commands")
}
