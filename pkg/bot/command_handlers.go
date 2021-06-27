package bot

import (
	"fmt"
	"strconv"

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

	roleID, err := strconv.ParseUint(i.Data.Options[0].RoleValue(nil, "").ID, 10, 64)
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

	_, exists, err := app.createTeam(guildID, roleID)
	if exists {
		msgFormat = " Team <@&%s> is already registered"
	} else if err != nil {
		msgFormat = " Unknown error registering team <@&%s>"
	} else {
		msgFormat = " Team <@&%s> successfully registered\n"
	}

	respond(s, i.Interaction, fmt.Sprintf(msgFormat, margs...))
}

func listTeamHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling list-teams")
	var returnMessage string

	guildID, err := strconv.ParseUint(i.Interaction.GuildID, 10, 64)
	if err != nil {
		klog.Error(err)
		returnMessage = " An unexpected error occurred getting your guild id"
	}

	teams, err := app.listTeams(guildID)
	if err != nil {
		klog.Error(err)
		returnMessage = " An unexpected error occurred listing teams"
	} else {
		if len(teams) < 1 {
			returnMessage = "There are no teams registered"
		} else {
			returnMessage = "Here are the currently registered teams:\n"
			for _, team := range teams {
				returnMessage = returnMessage + fmt.Sprintf("    <@&%d>\n", team.DiscordRoleID)
			}
		}
	}

	respond(s, i.Interaction, returnMessage)
}

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling help interaction")

	helpText := "\nHi! I'm here to help you track points!\n\n"
	helpText = helpText + "Here's the various commands I support: \n"

	for _, cmd := range commands {
		helpText = helpText + fmt.Sprintf("    /%s - %s\n", cmd.Name, cmd.Description)
	}

	respond(s, i.Interaction, helpText)
}

func respond(s *discordgo.Session, i *discordgo.Interaction, response string) {
	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: response,
		},
	})

	if err != nil {
		klog.Error(err)
	}
}
