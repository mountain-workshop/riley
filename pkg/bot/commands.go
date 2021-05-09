package bot

import (
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Help Command",
		},
		{
			Name:        "create-team",
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
		},
		// 	{
		// 		Name:        "subcommands",
		// 		Description: "Subcommands and command groups example",
		// 		Options: []*discordgo.ApplicationCommandOption{
		// 			// When a command has subcommands/subcommand groups
		// 			// It must not have top-level options, they aren't accesible in the UI
		// 			// in this case (at least not yet), so if a command has
		// 			// subcommands/subcommand any groups registering top-level options
		// 			// will cause the registration of the command to fail

		// 			{
		// 				Name:        "scmd-grp",
		// 				Description: "Subcommands group",
		// 				Options: []*discordgo.ApplicationCommandOption{
		// 					// Also, subcommand groups aren't capable of
		// 					// containing options, by the name of them, you can see
		// 					// they can only contain subcommands
		// 					{
		// 						Name:        "nst-subcmd",
		// 						Description: "Nested subcommand",
		// 						Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 					},
		// 				},
		// 				Type: discordgo.ApplicationCommandOptionSubCommandGroup,
		// 			},
		// 			// Also, you can create both subcommand groups and subcommands
		// 			// in the command at the same time. But, there's some limits to
		// 			// nesting, count of subcommands (top level and nested) and options.
		// 			// Read the intro of slash-commands docs on Discord dev portal
		// 			// to get more information
		// 			{
		// 				Name:        "subcmd",
		// 				Description: "Top-level subcommand",
		// 				Type:        discordgo.ApplicationCommandOptionSubCommand,
		// 			},
		// 		},
		// 	},
		// 	{
		// 		Name:        "responses",
		// 		Description: "Interaction responses testing initiative",
		// 		Options: []*discordgo.ApplicationCommandOption{
		// 			{
		// 				Name:        "resp-type",
		// 				Description: "Response type",
		// 				Type:        discordgo.ApplicationCommandOptionInteger,
		// 				Choices: []*discordgo.ApplicationCommandOptionChoice{
		// 					{
		// 						Name:  "Acknowledge",
		// 						Value: 2,
		// 					},
		// 					{
		// 						Name:  "Channel message",
		// 						Value: 3,
		// 					},
		// 					{
		// 						Name:  "Channel message with source",
		// 						Value: 4,
		// 					},
		// 					{
		// 						Name:  "Acknowledge with source",
		// 						Value: 5,
		// 					},
		// 				},
		// 				Required: true,
		// 			},
		// 		},
		// 	},
		// 	{
		// 		Name:        "followups",
		// 		Description: "Followup messages",
		// 	},
	}
)
