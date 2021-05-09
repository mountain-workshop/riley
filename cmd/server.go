package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"git.iratepublik.com/sudermans/discord-house-cup/pkg/bot"
)

var s *discordgo.Session
var token string
var guild string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&token, "token", "", "Bot access token")
	serverCmd.PersistentFlags().StringVar(&guild, "guild", "", "Test guild ID. If not passed - bot registers commands globally")

	envMap := map[string]string{
		"DISCORD_BOT_TOKEN": "token",
		"DISCORD_GUILD_ID":  "guild",
	}

	for env, flagName := range envMap {
		flag := serverCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			klog.Errorf("Could not find flag %s", flagName)
			continue
		}
		flag.Usage = fmt.Sprintf("%v [%v]", flag.Usage, env)
		if value := os.Getenv(env); value != "" {
			err := flag.Value.Set(value)
			if err != nil {
				klog.Errorf("Error setting flag %v to %s from environment variable %s", flag, value, env)
			}
		}
	}
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the discord bot server.",
	Long:  `Run the discord bot server`,
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		s, err = discordgo.New("Bot " + token)
		if err != nil {
			log.Fatalf("Invalid bot parameters: %v", err)
		}

		s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := bot.CommandHandlers[i.Data.Name]; ok {
				h(s, i)
			}
		})
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")
		s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			log.Println("Bot is up!")
		})
		err := s.Open()
		if err != nil {
			klog.Fatalf("Cannot open the session: %v", err)
		}

		for _, v := range bot.Commands {
			_, err := s.ApplicationCommandCreate(s.State.User.ID, guild, v)
			if err != nil {
				klog.Fatalf("Cannot create '%v' command: %v", v.Name, err)
			}
		}

		defer s.Close()

		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)
		<-stop
		log.Println("Gracefully shutting down...")
	},
}
