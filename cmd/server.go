package cmd

import (
	"flag"
	"fmt"
	"os"

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
	dbUser        string
	dbPassword    string
	dbHost        string
	dbName        string
	dbPort        int
	dbSSLMode     string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.PersistentFlags().StringVar(&token, "token", "", "Bot access token")
	serverCmd.PersistentFlags().StringVar(&guild, "guild", "", "Test guild ID. If not passed - bot registers commands globally")
	serverCmd.PersistentFlags().BoolVarP(&cleanupOnExit, "cleanup-on-exit", "c", true, "If true, when the server shuts down, it will delete all slash commands assocated with the bot.")

	serverCmd.PersistentFlags().StringVar(&dbUser, "db-user", "discordhousecup", "The postgres database user.")
	serverCmd.PersistentFlags().StringVar(&dbPassword, "db-password", "", "The password for the postgres database")
	serverCmd.PersistentFlags().StringVar(&dbHost, "db-host", "", "The database host")
	serverCmd.PersistentFlags().StringVar(&dbName, "db-name", "", "The database name")
	serverCmd.PersistentFlags().IntVar(&dbPort, "db-port", 5432, "The port to connect to the database on")
	serverCmd.PersistentFlags().StringVar(&dbSSLMode, "db-ssl-mode", "require", "The database ssl mode")

	envMap := map[string]string{
		"DISCORD_BOT_TOKEN":       "token",
		"DISCORD_GUILD_ID":        "guild",
		"DISCORD_BOT_DB_USER":     "db-user",
		"DISCORD_BOT_DB_PASSWORD": "db-password",
		"DISCORD_BOT_DB_HOST":     "db-host",
		"DISCORD_BOT_DB_NAME":     "db-name",
		"DISCORD_BOT_DB_PORT":     "db-port",
		"DISCORD_BOT_SSL_MODE":    "db-ssl-mode",
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		err = bot.InitDatabase(dbHost, dbUser, dbPassword, dbName, dbSSLMode, dbPort)
		if err != nil {
			return fmt.Errorf("could not connect to database: %v", err)
		}

		s, err = discordgo.New("Bot " + token)
		if err != nil {
			return fmt.Errorf("invalid bot parameters: %v", err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		server := bot.Server{
			Session:       s,
			CleanUpOnExit: cleanupOnExit,
		}
		if err := server.Run(); err != nil {
			klog.Fatal(err)
		}
	},
}
