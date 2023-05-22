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
		if interactionOption, ok := interactionOptionsMap[option.Name]; ok {
			value, err := assignOptionValue(interactionOption, option, session, i.GuildID)
			if err != nil {
				return returnContainer, err
			}
			returnContainer.Options[option.Name] = value
		} else if option.Required {
			return returnContainer, ErrMissingRequiredOption
		} else {
			returnContainer.Options[option.Name] = option.Default
		}
	}
	return returnContainer, nil
}

func assignOptionValue(interactionOption *discordgo.ApplicationCommandInteractionDataOption, option OptionRequirement, session *discordgo.Session, guildID string) (CommandOption, error) {
	switch option.Type {
	case discordgo.ApplicationCommandOptionString:
		return StringOption{
			Value: interactionOption.StringValue(),
			Name:  option.Name,
		}, nil
	case discordgo.ApplicationCommandOptionInteger:
		return IntegerOption{
			Value: interactionOption.IntValue(),
			Name:  option.Name,
		}, nil
	case discordgo.ApplicationCommandOptionNumber:
		return FloatOption{
			Value: interactionOption.FloatValue(),
			Name:  option.Name,
		}, nil
	case discordgo.ApplicationCommandOptionBoolean:
		return BooleanOption{
			Value: interactionOption.BoolValue(),
			Name:  option.Name,
		}, nil
	case discordgo.ApplicationCommandOptionChannel:
		return ChannelOption{
			Value: interactionOption.ChannelValue(session),
			Name:  option.Name,
		}, nil
	case discordgo.ApplicationCommandOptionRole:
		return RoleOption{
			Value: interactionOption.RoleValue(session, guildID),
			Name:  option.Name,
		}, nil
	case discordgo.ApplicationCommandOptionMentionable:
		return UserOption{
			Value: interactionOption.UserValue(session),
			Name:  option.Name,
		}, nil
	default:
		return nil, ErrInvalidOptionType
	}
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
