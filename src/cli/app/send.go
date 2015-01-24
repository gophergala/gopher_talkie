package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/nklizhe/gopass"
	"os"
	"time"
)

func NewSendCommand(this *App) cli.Command {
	return cli.Command{
		Name:  "send",
		Usage: "record and send a voice message",
		Action: func(c *cli.Context) {
			this.send(c)
		},
	}
}

func (this *App) send(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Printf("Usage: %s send <to>\n", os.Args[0])
		return
	}
	to := c.Args()[0]

	fmt.Printf("Press any key to start recording your message...\n")
	gopass.GetCh()

	for i := 15; i > 0; i-- {
		fmt.Printf("\rRecording...%d seconds left", i)
		time.Sleep(time.Duration(1) * time.Second)
	}
	fmt.Printf("\rRecorded\n")

	fmt.Printf("Sending to %s...\n", to)
	fmt.Printf("Done\n")
}
