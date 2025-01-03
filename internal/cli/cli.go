package cli

import (
	"os"

	"github.com/urfave/cli/v2"
)

const AppName = "configweb"

type App struct {
	app *cli.App
}

func NewApp() *App {
	return &App{
		app: &cli.App{
			Name:  AppName,
			Usage: "A config joiner",
			Commands: []*cli.Command{
				newCommand,
			},
		},
	}
}

func (a *App) Run() error {
	return a.app.Run(os.Args)
}
