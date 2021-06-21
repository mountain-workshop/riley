package bot

import (
	"fmt"
	"strconv"

	"git.iratepublik.com/sudermans/discord-house-cup/pkg/model"
	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, app *App){
	"help":          helpHandler,
	"register-team": registerTeamHandler,
}

func registerTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
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

	if len(i.Data.Options) >= 2 {
		margs = append(margs, i.Data.Options[1].ChannelValue(nil).ID)
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

	app.DB().NewInsert().Model(&model.Team{
		DiscordRoleID: teamID,
	}).Exec(app.ctx)
}

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling help interaction")
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: `I'm here to help you track points. Here's the various commands I support:
			/help - This command`,
		},
	})
}
