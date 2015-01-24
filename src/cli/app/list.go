package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	_ "os"
)

func NewListCommand(app *App) cli.Command {
	return cli.Command{
		Name:  "list",
		Usage: "list all messages",
		Action: func(c *cli.Context) {
			app.list(c)
		},
	}
}

func (this *App) list(c *cli.Context) {
	fmt.Printf("List all message...\n")
}
