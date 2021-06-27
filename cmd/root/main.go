package main

import (
	"log"
	"os"

	"github.com/mountain-workshop/riley/cmd/root/migrations"
	"github.com/mountain-workshop/riley/pkg/bot"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"k8s.io/klog"
)

func main() {
	app := &cli.App{
		Name: "bot",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			newDBCommand(migrations.Migrations),
			newRunCommand(),
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newRunCommand() *cli.Command {
	return &cli.Command{
		Name:  "run",
		Usage: "Run the discord bot.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "token",
				Value: "",
				Usage: "discord bot token",
			},
			&cli.StringFlag{
				Name:  "guild",
				Value: "",
				Usage: "discord guild id for testing",
			},
		},
		Action: func(c *cli.Context) error {
			_, app, err := bot.StartCLI(c)
			if err != nil {
				return err
			}
			defer app.Stop()
			klog.Info(bot.WaitExitSignal())
			return nil
		},
	}
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Init(ctx, app.DB())
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()
					return migrations.Migrate(ctx, app.DB())
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Rollback(ctx, app.DB())
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Lock(ctx, app.DB())
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Unlock(ctx, app.DB())
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.CreateGo(ctx, app.DB(), c.Args().Get(0))
				},
			},
			{
				Name:  "create_sql",
				Usage: "create SQL migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.CreateSQL(ctx, app.DB(), c.Args().Get(0))
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.Status(ctx, app.DB())
				},
			},
			{
				Name:  "mark_completed",
				Usage: "mark migrations as completed without actually running them",
				Action: func(c *cli.Context) error {
					ctx, app, err := bot.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					return migrations.MarkCompleted(ctx, app.DB())
				},
			},
		},
	}
}
