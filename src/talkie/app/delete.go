package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
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
	if err := this.setup(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}
	fmt.Printf("Delete message...\n")

	// TODO: need implement
}
