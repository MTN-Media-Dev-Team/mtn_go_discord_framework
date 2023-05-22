package mtn_go_discord_framework

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

type SlashCommand struct {
	Name               string
	Description        string
	Handler            func(s *discordgo.Session, i *discordgo.InteractionCreate, options *OptionContainer)
	RequiredOptions    []OptionRequirement
	applicationCommand discordgo.ApplicationCommand
}

func (s SlashCommand) generateOptions() []*discordgo.ApplicationCommandOption {
	options := make([]*discordgo.ApplicationCommandOption, 0)
	for _, option := range s.RequiredOptions {
		options = append(options, &discordgo.ApplicationCommandOption{
			Name:        option.Name,
			Description: option.Description,
			Type:        discordgo.ApplicationCommandOptionType(option.Type),
			Required:    option.Required,
		})
	}
	return options
}

func (s SlashCommand) generateApplicationCommand() discordgo.ApplicationCommand {
	return discordgo.ApplicationCommand{
		Name:        s.Name,
		Description: s.Description,
		Options:     s.generateOptions(),
	}
}

func (s SlashCommand) validateOptions(session *discordgo.Session, i *discordgo.InteractionCreate) (OptionContainer, error) {
	returnContainer := OptionContainer{
		Options: make(map[string]CommandOption),
	}
	// Mapping interaction options by their name
	interactionOptionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, interactionOption := range i.ApplicationCommandData().Options {
		interactionOptionsMap[interactionOption.Name] = interactionOption
	}
	for _, option := range s.RequiredOptions {
		if option.Required {
			if interactionOption, ok := interactionOptionsMap[option.Name]; ok {
				// Using a type switch instead of multiple case statements
				switch option.Type {
				case discordgo.ApplicationCommandOptionString:
					returnContainer.Options[option.Name] = StringOption{
						Value: interactionOption.StringValue(),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionInteger:
					returnContainer.Options[option.Name] = IntegerOption{
						Value: interactionOption.IntValue(),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionBoolean:
					returnContainer.Options[option.Name] = BooleanOption{
						Value: interactionOption.BoolValue(),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionChannel:
					returnContainer.Options[option.Name] = ChannelOption{
						Value: interactionOption.ChannelValue(session),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionRole:
					returnContainer.Options[option.Name] = RoleOption{
						Value: interactionOption.RoleValue(session, i.GuildID),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionMentionable:
					returnContainer.Options[option.Name] = UserOption{
						Value: interactionOption.UserValue(session),
						Name:  option.Name,
					}
				default:
					return returnContainer, ErrInvalidOptionType
				}
			} else {
				return returnContainer, ErrMissingRequiredOption
			}
		} else {
			if interactionOption, ok := interactionOptionsMap[option.Name]; ok {
				// Using a type switch instead of multiple case statements
				switch option.Type {
				case discordgo.ApplicationCommandOptionString:
					returnContainer.Options[option.Name] = StringOption{
						Value: interactionOption.StringValue(),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionInteger:
					returnContainer.Options[option.Name] = IntegerOption{
						Value: interactionOption.IntValue(),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionBoolean:
					returnContainer.Options[option.Name] = BooleanOption{
						Value: interactionOption.BoolValue(),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionChannel:
					returnContainer.Options[option.Name] = ChannelOption{
						Value: interactionOption.ChannelValue(session),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionRole:
					returnContainer.Options[option.Name] = RoleOption{
						Value: interactionOption.RoleValue(session, i.GuildID),
						Name:  option.Name,
					}
				case discordgo.ApplicationCommandOptionMentionable:
					returnContainer.Options[option.Name] = UserOption{
						Value: interactionOption.UserValue(session),
						Name:  option.Name,
					}
				default:
					return returnContainer, ErrInvalidOptionType
				}
			} else {
				returnContainer.Options[option.Name] = option.Default
			}
		}
	}
	return returnContainer, nil
}

// declare errors
var (
	ErrInvalidOptionType     = errors.New("invalid option type")
	ErrMissingRequiredOption = errors.New("missing required option")
)

type OptionRequirement struct {
	Required    bool
	Name        string
	Description string
	Type        discordgo.ApplicationCommandOptionType
	Default     CommandOption
}

type OptionContainer struct {
	Options map[string]CommandOption
}
