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

func helpHandler(s *discordgo.Session, i *discordgo.InteractionCreate, app *App) {
	klog.V(4).Info("handling help interaction")

	helpText := "\nHi! I'm here to help you track points!\n\n"
	helpText = helpText + "Here's the various commands I support: \n"

	for _, cmd := range commands {
		helpText = helpText + fmt.Sprintf("    /%s - %s\n", cmd.Name, cmd.Description)
	}

	respond(s, i.Interaction, helpText)
}
