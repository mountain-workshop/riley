package cmd

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog"

	"git.iratepublik.com/sudermans/discord-house-cup/pkg/bot"
)

var (
	s             *discordgo.Session
	token         string
	guild         string
	cleanupOnExit bool
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&token, "token", "", "Bot access token")
	serverCmd.PersistentFlags().StringVar(&guild, "guild", "", "Test guild ID. If not passed - bot registers commands globally")
	serverCmd.PersistentFlags().BoolVarP(&cleanupOnExit, "cleanup-on-exit", "c", false, "If true, when the server shuts down, it will delete all slash commands assocated with the bot.")

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
	klog.InitFlags(nil)
	flag.Parse()
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
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
		s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			klog.Infof("Bot is up!")
		})
		err := s.Open()
		if err != nil {
			klog.Fatalf("Cannot open the session: %v", err)
		}

		defer s.Close()

		if err := registerCommands(); err != nil {
			klog.Error(err)
		}
		stop := make(chan os.Signal)
		signal.Notify(stop, os.Interrupt)
		<-stop
		klog.Info("Gracefully shutting down...")
		if cleanupOnExit {
			if err := removeAllCommands(); err != nil {
				klog.Error(err)
			}
		}
	},
}

func registerCommands() error {
	klog.Info("registering commands")
	for _, v := range bot.Commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guild, v)
		if err != nil {
			return fmt.Errorf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
	return nil
}

func removeAllCommands() error {
	klog.Info("deleting all commands")
	commands, err := s.ApplicationCommands(s.State.User.ID, guild)
	if err != nil {
		return err
	}

	for _, c := range commands {
		klog.Infof("deleting command %s", c.Name)
		err := s.ApplicationCommandDelete(s.State.User.ID, guild, c.ID)
		if err != nil {
			klog.Errorf("Cannot delete '%v' command: %v", c.Name, err)
		}
	}

	return nil
}
