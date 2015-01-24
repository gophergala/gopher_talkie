package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/gopher_talkie/common"
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

	go func() {
		for i := 15; i > 0; i-- {
			fmt.Printf("\rRecording...%d seconds left", i)
			time.Sleep(time.Duration(1) * time.Second)
		}
	}()

	users, err := this.store.FindUserByName(to)
	if err != nil || len(users) == 0 {
		fmt.Printf("\nUser '%s' not found!", to)
		return
	}

	// TODO: ask user to select if there are multiple users

	msg := common.NewMessage(this.user, users[0])
	f, err := this.record(time.Duration(15) * time.Second)
	if err != nil {
		fmt.Printf("\nError recording!%s", err.Error())
		return
	}
	fmt.Println(f)

	fmt.Printf("\rRecorded\n")
	msg.Content = []byte(f)    // TODO: encrypt message
	this.store.AddMessage(msg) // Save message

	fmt.Printf("Sending to %s...\n", msg.To.Name)
	fmt.Printf("Done\n")
}
