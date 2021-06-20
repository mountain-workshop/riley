package bot

import (
	"embed"
	"io/fs"
	"path"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed embed
	embedFS      embed.FS
	unwrapFSOnce sync.Once
	unwrappedFS  fs.FS
)

func FS() fs.FS {
	unwrapFSOnce.Do(func() {
		fsys, err := fs.Sub(embedFS, "embed")
		if err != nil {
			panic(err)
		}
		unwrappedFS = fsys
	})
	return unwrappedFS
}

type AppConfig struct {
	Guild         string
	CleanUpOnExit bool
	BotToken      string
	Env           string
	Debug         bool `yaml:"debug"`

	Service string

	DB struct {
		DSN string `yaml:"dsn"`
	} `yaml:"db"`
}

func ReadConfig(fsys fs.FS, env, service, token, guild string) (*AppConfig, error) {
	b, err := fs.ReadFile(fsys, path.Join("config", env+".yaml"))
	if err != nil {
		return nil, err
	}

	cfg := new(AppConfig)
	if err := yaml.Unmarshal(b, cfg); err != nil {
		return nil, err
	}

	cfg.Service = service
	cfg.Env = env
	cfg.Guild = guild
	cfg.BotToken = token

	return cfg, nil
}
