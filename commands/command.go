package commands

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	AppCommand discordgo.ApplicationCommand
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func panicResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: "Something went wrong",
		},
	})
	time.Sleep(time.Second * 10)
	s.InteractionResponseDelete(s.State.User.ID, i.Interaction)
}

// passing the state and guild ID is sloppy but Member.roles only returns role ID's and not the structs

func hasAdmin(state *discordgo.State, roles []string, gID string) bool {
	for _, v := range roles {
		r, _ := state.Role(gID, v)
		if strings.EqualFold(r.Name, "admin") {
			return true
		}
	}
	return false
}
