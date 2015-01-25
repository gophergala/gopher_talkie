package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/gopher_talkie/src/common"
	_ "os"
)

func NewPlayCommand(this *App) cli.Command {
	return cli.Command{
		Name:  "play",
		Usage: "play a message",
		Action: func(c *cli.Context) {
			this.play(c)
		},
	}
}

func playMessage(msg *common.Message) {
	fmt.Printf("Decrypting message...\n")
	fmt.Printf("Playing...\n")
}

func (this *App) play(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Printf("Playing the first message...\n")
		return
	}
	msgID := c.Args()[0]
	fmt.Printf("Playing message %s...\n", msgID)
}
