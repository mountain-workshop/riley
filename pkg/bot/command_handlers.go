package bot

import (
	"fmt"
	"strconv"
	"strings"

	"git.iratepublik.com/sudermans/discord-house-cup/pkg/model"
	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, app *App){
	"help":          helpHandler,
	"register-team": registerTeamHandler,
	"list-teams":    listTeamHandler,
}

func registerTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling register-team")

	teamID, err := strconv.ParseUint(i.Data.Options[0].RoleValue(nil, "").ID, 10, 64)
	if err != nil {
		klog.Error(err)
		return
	}

	guildID, err := strconv.ParseUint(i.Interaction.GuildID, 10, 64)
	if err != nil {
		klog.Error(err)
		return
	}

	margs := []interface{}{
		i.Data.Options[0].RoleValue(nil, "").ID,
	}

	var msgFormat string

	_, err = app.DB().NewInsert().Model(&model.Team{
		DiscordRoleID:  teamID,
		DiscordGuildID: guildID,
	}).Exec(app.ctx)
	if err != nil {
		klog.Error(err)
		if strings.Contains(err.Error(), "23505") {
			msgFormat = " Team <@&%s> already exists"
		} else {
			msgFormat = " Unknown error registering team <@&%s>"
		}
	} else {
		msgFormat = " Team <@&%s> successfully registered\n"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: fmt.Sprintf(
				msgFormat,
				margs...,
			),
		},
	})
}

func listTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling list-teams")
	var returnMessage string

	guildID, err := strconv.ParseUint(i.Interaction.GuildID, 10, 64)
	if err != nil {
		klog.Error(err)
		returnMessage = " An unexpected error occurred getting your guild id"
	}

	teams := make([]model.Team, 0)
	if err := app.DB().NewSelect().Model(&teams).Where("discord_guild_id = ?", guildID).Scan(app.ctx); err != nil {
		klog.Error(err)
		returnMessage = " An unexpected error occurred listing teams"
	} else {
		returnMessage = "Here are the currently registered teams:\n"
		for _, team := range teams {
			returnMessage = returnMessage + fmt.Sprintf("    <@&%d>\n", team.DiscordRoleID)
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: returnMessage,
		},
	})
}

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling help interaction")

	helpText := "\nHi! I'm here to help you track points!\n\n"
	helpText = helpText + "Here's the various commands I support: \n"

	for _, cmd := range commands {
		helpText = helpText + fmt.Sprintf("    /%s - %s\n", cmd.Name, cmd.Description)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: helpText,
		},
	})
}
