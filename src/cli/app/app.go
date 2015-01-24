package app

import (
	_ "fmt"
	"github.com/codegangsta/cli"
	"os"
)

const (
	Version = "0.1.0"
	Author  = "nklizhe@gmail.com"
)

type App struct {
	app *cli.App
}

func NewApp() *App {

	this := &App{}
	app := cli.NewApp()
	app.Name = "talkie"
	app.Usage = "Secure voicing messaging for geeks"
	app.Action = func(c *cli.Context) {
		// default
		this.list(c)
	}
	app.Version = Version
	app.Author = Author
	app.Commands = []cli.Command{
		NewListCommand(this),
		NewSendCommand(this),
		NewListenCommand(this),
		NewDeleteCommand(this),
	}

	this.app = app
	return this
}

func (this *App) Run() {
	this.app.Run(os.Args)
}
