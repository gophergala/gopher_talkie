package app

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/gopher_talkie/src/audio"
	_ "github.com/gophergala/gopher_talkie/src/common"
	"github.com/gophergala/gopher_talkie/src/crypto"
	"io"
	"os"
	"path"
	"strconv"
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
	if err := this.setup(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		return
	}

	messages, err := this.client.GetMessages(this.user)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	if len(messages) > 0 {
		fmt.Printf("You have new messages:\n\n")
		for i := range messages {
			m := messages[i]
			fmt.Printf("  (%d) %s <%s> - %s\n", i+1, m.From.Name, m.From.Email, m.CreatedAt.Format("Jan 02"))
		}

		for {
			fmt.Printf("Enter number (%d - %d) to listen, or Q)uit > ", 1, len(messages))
			var choice string
			fmt.Scanf("%v", &choice)
			if choice == "Q" {
				break
			}

			idx, _ := strconv.ParseInt(choice, 10, 4)
			if int(idx) > 0 && int(idx) <= len(messages) {
				// play message at idx-1
				m := messages[int(idx)-1]

				// download content
				if len(m.Content) == 0 {
					m.Content, err = this.client.DownloadMessage(m.MessageID)
					if err != nil {
						fmt.Printf("Error download message! %s\n", err.Error())
						continue
					}
				}

				// decrypt content
				content, err := crypto.GPGDecrypt(this.user.Key, bytes.NewReader(m.Content))
				if err != nil {
					fmt.Printf("Error decrypt message! %s\n", err.Error())
					continue
				}

				// save it to a temp file
				tempfile := path.Join(os.TempDir(), fmt.Sprintf("%s.aiff", uuid.NewUUID().String()))
				defer func() {
					os.RemoveAll(tempfile)
				}()
				f, err := os.Create(tempfile)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					continue
				}
				defer f.Close()

				_, err = io.Copy(f, bytes.NewBuffer(content))
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					continue
				}

				// play
				fmt.Printf("Playing...")
				err = audio.PlayAIFF(tempfile, nil)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					continue
				}
				fmt.Println()
			}
		}

	} else {
		fmt.Println("No messages.")
	}
}
