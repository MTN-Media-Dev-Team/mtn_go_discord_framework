package mtn_go_discord_framework

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// sendEphemeralResponse sends a response to an interaction that only the user can see
func SendEphemeralResponse(s *discordgo.Session, i *discordgo.InteractionCreate, msg string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
			Flags:   EphemeralFlag,
		},
	})
}

func SendEphemeralEmbed(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  EphemeralFlag,
		},
	})
}

// sendDeferResponse sends a deferred response to an interaction
func SendDeferResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
}

// busySleep pauses execution for up to 10 seconds if the system is busy
func BusySleep() {
	if systemBusy {
		counter := 0
		for systemBusy {
			time.Sleep(1 * time.Second)
			counter++
			if counter > 10 {
				systemBusy = false
				break
			}
		}
	}
}

// Check if the system is busy. If not, set it to busy and return true.
// If it is busy, return false.
func TryAcquireSystem() bool {
	mutex.Lock()
	defer mutex.Unlock()

	if systemBusy {
		return false
	}

	systemBusy = true
	return true
}

// Set systemBusy to false. Should be deferred after a successful tryAcquireSystem.
func ReleaseSystem() {
	mutex.Lock()
	systemBusy = false
	mutex.Unlock()
}

// checkForRoles checks if a member has one of the provided roles
func CheckForRoles(s *discordgo.Session, i *discordgo.InteractionCreate, roles ...string) bool {
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Printf("MTN Discord Framework - CheckForRoles: Error getting guild member: %v", err)
		return false
	}

	roleMap := make(map[string]bool)
	for _, role := range roles {
		roleMap[role] = true
	}

	for _, roleID := range member.Roles {
		if _, ok := roleMap[roleID]; ok {
			return true
		}
	}

	return false
}

func CheckForRolesOrAdmin(s *discordgo.Session, i *discordgo.InteractionCreate, roles ...string) bool {
	if CheckForAdmin(s, i) {
		return true
	}
	return CheckForRoles(s, i, roles...)
}

func CheckForAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		log.Printf("MTN Discord Framework - CheckForRoles: Error getting guild member: %v", err)
		return false
	}
	// check if user has administator on the guild
	if member.Permissions&discordgo.PermissionAdministrator == discordgo.PermissionAdministrator {
		return true
	}
	return false
}
