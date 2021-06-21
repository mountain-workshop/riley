package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"k8s.io/klog"
)

var helpCommand = discordgo.ApplicationCommand{
	Name:        "help",
	Description: "Help Command",
}

var registerTeamCommand = discordgo.ApplicationCommand{
	Name:        "register-team",
	Description: "Command for associating a team with a Discord Role",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionRole,
			Name:        "discord-role",
			Description: "The role to associate with this team",
			Required:    true,
		},
	},
}

var commands = []*discordgo.ApplicationCommand{
	&helpCommand,
	&registerTeamCommand,
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
