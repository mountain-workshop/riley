package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

var commands = []*discordgo.ApplicationCommand{
	&helpCommand,
	&registerTeamCommand,
	&listTeamCommand,
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, app *App){
	"help":          helpHandler,
	"register-team": registerTeamHandler,
	"list-teams":    listTeamHandler,
}

func (app *App) registerCommands() error {
	klog.Info("registering commands")
	for _, v := range commands {
		klog.Infof("registering command: %s", v.Name)
		_, err := app.Discord.ApplicationCommandCreate(app.Discord.State.User.ID, app.cfg.Guild, v)
		if err != nil {
			return fmt.Errorf("cannot create '%s' command: %v", v.Name, err)
		}
	}
	return nil
}

func (app *App) removeAllCommands() error {
	klog.Info("deleting all commands")
	commands, err := app.Discord.ApplicationCommands(app.Discord.State.User.ID, app.cfg.Guild)
	if err != nil {
		return err
	}

	for _, c := range commands {
		klog.Infof("deleting command %s", c.Name)
		err := app.Discord.ApplicationCommandDelete(app.Discord.State.User.ID, app.cfg.Guild, c.ID)
		if err != nil {
			klog.Errorf("Cannot delete '%v' command: %v", c.Name, err)
		}
	}

	return nil
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
