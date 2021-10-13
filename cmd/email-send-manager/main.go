package main

import (
	"context"
	"os"

	"github.com/urfave/cli/v2"

	"email-send-manager/internal/app"
	"email-send-manager/pkg/logger"
)

// VERSION 版本号，可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "1.0.0"

func main() {
	ctx := logger.NewTagContext(context.Background(), "__main__")

	app := cli.NewApp()
	app.Name = "email-send-manager"
	app.Version = VERSION
	app.Usage = "RBAC scaffolding based on GIN + GORM + CASBIN + WIRE."
	app.Commands = []*cli.Command{
		newWebCmd(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.WithContext(ctx).Errorf(err.Error())
	}
}

func newWebCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "conf",
				Aliases:  []string{"c"},
				Usage:    "App configuration file(.json,.yaml,.toml)",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "www",
				Usage: "Static site directory",
			},
		},
		Action: func(c *cli.Context) error {
			return app.Run(ctx,
				app.SetConfigFile(c.String("conf")),
				app.SetWWWDir(c.String("www")),
				app.SetVersion(VERSION))
		},
	}
}
