package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/urfave/cli/v2"

	"panda.com/api/service"
	"panda.com/api/server"
	version "panda.com/cmd"
	"panda.com/config"
	"panda.com/database/mysql"
)

var configPathFlag = &cli.StringFlag{
	Name:     "config-file",
	Usage:    "The filepath to a json file, flag is required",
	Required: true,
}

type Config struct {
	Port int `yaml:"port"`
	MySQL mysql.Config `yaml:"mysql"`
}

func main() {
	app := &cli.App{
		Name:  "panda-api",
		Usage: "Panda API Server",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8080,
				Usage:   "Port to run the server on",
			},
		},
		Action: exec,
		Version: version.String(),
	}

	if err := app.Run(os.Args); err != nil {
		slog.Error("running api application failed", "error", err)
	}

}

func exec(ctx *cli.Context) error {
	cfg := &Config{}
	if err := config.Load(ctx.String(configPathFlag.Name),cfg); err != nil {
		return fmt.Errorf("load config file error: %w", err)
	}

	db, err := mysql.NewMySQLDB(cfg.MySQL)
	if err != nil {
		return fmt.Errorf("init mysql db error: %w", err)
	}

	slog.Info("starting api server...", "port", cfg.Port)

	server.New(
		cfg.Port,
		service.NewService(db),
	).Run()
	return nil
}