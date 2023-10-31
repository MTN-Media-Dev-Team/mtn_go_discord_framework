# MTN Go Discord Framework

Framework to simplifly the creation of discord bots with slash commands and buttons in go build ontop of discordgo.
All discordgo features are still available and can be used directly via discordgo, this just helps with the startup and creation of slash commands and buttons.

## Example Bot with SlashCommand and Button

### Full Example Bot

```go

func main() {

    log.Println("Starting bot")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    // Init the framework and returns a discordgo.Session if you need it
    discord := mtn_go_discord_framework.InitFramework(os.Getenv("DEBUG"), os.Getenv("TESTING_GUILD_ID"), os.Getenv("BOT_TOKEN"))

    // Register Button Handlers
    mtn_go_discord_framework.RegisterButtonHandlersWithFramework(getButtonHandlers())
    // Register Slash Commands
	mtn_go_discord_framework.RegisterSlashCommandsWithFramework(getCommands())

    // Launch Framework and Start Bot
	mtn_go_discord_framework.StartFramework()

    // Wait for interrupt
    <-ctx.Done()

    // Stop Bot
    mtn_go_discord_framework.ShutdownFramework()
}

func getButtonHandlers() []mtn_go_discord_framework.ButtonHandler {
	handlers := make([]mtn_go_discord_framework.ButtonHandler, 0)
    // example button handler
	handlers = append(handlers, mtn_go_discord_framework.ButtonHandler{
		Handler:  func(s *discordgo.Session, i *discordgo.InteractionCreate, args ...string) {
            // do something
        },
		CustomID: "example_button", // must be mapped, use - seperator to append args to custom id, can then be used in the handler with args[0] etc.
	})
	return handlers
}


func getCommands() []mtn_go_discord_framework.SlashCommand {
	commandList := make([]mtn_go_discord_framework.SlashCommand, 0)
	commandList = append(commandList, testCommand())
	return commandList
}

func testCommand() mtn_go_discord_framework.SlashCommand {
	// define handler function
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate, options *mtn_go_discord_framework.OptionContainer) {
		mtn_go_discord_framework.SendEphemeralResponse(s, i, "Test command")
	}
	command := mtn_go_discord_framework.SlashCommand{
		Name:        "test",
		Description: "Test command",
		Handler:     handler,
	}
	return command
}

```

### Example Slash Command with Options

```go

func exampleOptionsCommand() mtn_go_discord_framework.SlashCommand {
	handler := func(s *discordgo.Session, i *discordgo.InteractionCreate, options *mtn_go_discord_framework.OptionContainer) {
		// get options from container. If required options are not provided, the handler will not be called, so no need to check for nil. If required is false default values will be used if not provided by user
		stringvalue := options.Options["stringvalue"].GetValue().(string)
		numbervalue := options.Options["numbervalue"].GetValue().(float64)
		boolvalue := options.Options["boolvalue"].GetValue().(bool)

        // check if user had roles or admin if needed
		var authenticated bool
        authenticated = mtn_go_discord_framework.CheckForRoles(s, i, GROUP_ID1, GROUP_ID2) // unlimited roles can be provided
        authenticated = mtn_go_discord_framework.CheckForAdmin(s, i)
        authenticated = mtn_go_discord_framework.CheckForRolesOrAdmin(s, i, GROUP_ID1, GROUP_ID2) // unlimited roles can be provided

        // do something

	}
	command := mtn_go_discord_framework.SlashCommand{
		Name:        "exampleoptions",
		Description: "Example command with options",
		RequiredOptions: []mtn_go_discord_framework.OptionRequirement{
			{
				Name:        "stringvalue",
				Description: "StringValue",
				Type:        discordgo.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "numbervalue",
				Description: "NumberValue",
				Type:        discordgo.ApplicationCommandOptionInteger,
				Required:    true,
			},
            // If Required is false, you must provide a default value
			{
				Name:        "boolvalue",
				Description: "BoolValue",
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Required:    false,
				Default: mtn_go_discord_framework.BooleanOption{
					Value: true,
					Name:  "boolvalue",
				},
			},
		},
		Handler: handler,
	}
	return command
}
```

### Example Responses

```go
// Send a response to the user that is only visible to them
mtn_go_discord_framework.SendEphemeralResponse(s *discordgo.Session, i *discordgo.InteractionCreate, "Ephemeral Response")

// Send embed to user that is only visible to them
mtn_go_discord_framework.SendEphemeralEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed)

// Send defer response to discord, that the interaction is handled without a response
mtn_go_discord_framework.SendDeferResponse(s *discordgo.Session, i *discordgo.InteractionCreate)

// Other reponses to be done directly via discordgo
```
