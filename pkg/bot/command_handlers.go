package bot

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"help":          helpHandler,
	"register-team": registerTeamHandler,
}

func registerTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	klog.V(4).Info("handling register-team")

	teamID, err := strconv.ParseUint(i.Data.Options[0].RoleValue(nil, "").ID, 10, 64)
	if err != nil {
		klog.Error(err)
		return
	}

	margs := []interface{}{
		i.Data.Options[0].RoleValue(nil, "").ID,
	}

	msgformat := " We are going to register this team:\n"
	msgformat += "> role-id: <@&%s>\n"

	var teamName string
	if len(i.Data.Options) >= 2 {
		margs = append(margs, i.Data.Options[1].ChannelValue(nil).ID)
		teamName = i.Data.Options[1].ChannelValue(nil).ID
		msgformat += "> team-name: %s\n"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf(
				msgformat,
				margs...,
			),
		},
	})
	registerTeam(teamName, teamID)
}

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	klog.V(4).Info("handling help interaction")
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: `I'm here to help you track points. Here's the various commands I support:
			/help - This command`,
		},
	})
}
