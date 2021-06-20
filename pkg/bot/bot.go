package bot

import (
	"context"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/urfave/cli/v2"
	"k8s.io/klog"
)

type appCtxKey struct{}

func AppFromContext(ctx context.Context) *App {
	return ctx.Value(appCtxKey{}).(*App)
}

func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)
	return ctx
}

type App struct {
	Discord *discordgo.Session

	ctx context.Context
	cfg *AppConfig

	stopping uint32
	stopCh   chan struct{}

	onStop      appHooks
	onAfterStop appHooks

	discordOnce sync.Once

	// lazy init
	dbOnce sync.Once
	db     *bun.DB
}

func New(ctx context.Context, cfg *AppConfig) *App {
	app := &App{
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
	app.ctx = ContextWithApp(ctx, app)
	app.initDiscord()
	return app
}

func StartCLI(c *cli.Context) (context.Context, *App, error) {
	return Start(c.Context, c.Command.Name, c.String("env"), c.String("token"), c.String("guild"))
}

func Start(ctx context.Context, service, envName, token, guild string) (context.Context, *App, error) {
	cfg, err := ReadConfig(FS(), envName, service, token, guild)
	if err != nil {
		return nil, nil, err
	}
	return StartConfig(ctx, cfg)
}

func StartConfig(ctx context.Context, cfg *AppConfig) (context.Context, *App, error) {
	rand.Seed(time.Now().UnixNano())

	app := New(ctx, cfg)
	if err := onStart.Run(ctx, app); err != nil {
		return nil, nil, err
	}
	return app.ctx, app, nil
}

func (app *App) Stop() {
	_ = app.onStop.Run(app.ctx, app)
	_ = app.onAfterStop.Run(app.ctx, app)
}

func (app *App) OnStop(name string, fn HookFunc) {
	app.onStop.Add(newHook(name, fn))
}

func (app *App) OnAfterStop(name string, fn HookFunc) {
	app.onAfterStop.Add(newHook(name, fn))
}

func (app *App) Context() context.Context {
	return app.ctx
}

func (app *App) Config() *AppConfig {
	return app.cfg
}

func (app *App) Running() bool {
	return !app.Stopping()
}

func (app *App) Stopping() bool {
	return atomic.LoadUint32(&app.stopping) == 1
}

func (app *App) IsDebug() bool {
	return app.cfg.Debug
}

func (app *App) DB() *bun.DB {
	app.dbOnce.Do(func() {
		config, err := pgx.ParseConfig(app.cfg.DB.DSN)
		if err != nil {
			panic(err)
		}

		config.PreferSimpleProtocol = true
		sqldb := stdlib.OpenDB(*config)

		db := bun.NewDB(sqldb, pgdialect.New())
		if app.IsDebug() {
			db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose()))
		}

		app.db = db

	})
	return app.db
}

func (app *App) initDiscord() {
	if app.cfg.Service != "run" {
		klog.Infof("skipping discord init for non-bot command")
		return
	}
	app.discordOnce.Do(func() {
		var err error
		app.Discord, err = discordgo.New("Bot " + app.cfg.BotToken)
		if err != nil {
			klog.Fatalf("invalid bot parameters: %v", err)
		}

		app.OnStop("app.Discord.Close", func(ctx context.Context, _ *App) error {
			return app.Discord.Close()
		})

		app.Discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlers[i.Data.Name]; ok {
				h(s, i, app)
			}
		})

		app.Discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
			klog.Infof("bot is up!")
		})
		if err := app.Discord.Open(); err != nil {
			klog.Fatalf("cannot open the session: %v", err)
		}

		if err := app.registerCommands(); err != nil {
			klog.Fatal(err)
		}

		app.OnStop("app.removeAllCommands", func(ctx context.Context, _ *App) error {
			return app.removeAllCommands()
		})
	})
}

func WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}
