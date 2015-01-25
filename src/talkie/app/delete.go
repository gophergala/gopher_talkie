package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	_ "os"
)

func NewDeleteCommand(this *App) cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "delete a message",
		Action: func(c *cli.Context) {
			this.delete(c)
		},
	}
}

func (this *App) delete(c *cli.Context) {
	fmt.Printf("Delete message...\n")
}
