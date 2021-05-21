package bot

import (
	"github.com/bwmarrin/discordgo"
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
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "team-name",
			Description: "Team name if using a different name than the role",
			Required:    false,
		},
	},
}

var commands = []*discordgo.ApplicationCommand{
	&helpCommand,
	&registerTeamCommand,
}
