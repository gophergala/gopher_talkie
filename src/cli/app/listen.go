package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	_ "os"
)

func NewListenCommand(this *App) cli.Command {
	return cli.Command{
		Name:  "listen",
		Usage: "listen messages",
		Action: func(c *cli.Context) {
			this.listen(c)
		},
	}
}

func (this *App) listen(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Printf("Listening first message...\n")
		return
	}
	msgID := c.Args()[0]
	fmt.Printf("Listening message %s...\n", msgID)
}
