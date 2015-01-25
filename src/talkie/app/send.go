package app

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/gopher_talkie/src/audio"
	"github.com/gophergala/gopher_talkie/src/common"
	"github.com/gophergala/gopher_talkie/src/crypto"
	"github.com/nklizhe/gopass"
	"os"
	"path"
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
	if err := this.setup(c); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		return
	}

	if this.user == nil {
		panic(ErrNoUser)
	}

	var recipient *common.User
	if len(c.Args()) == 0 {
		fmt.Printf("Sending message to yourself...\n")
		recipient = this.user
	} else {
		to := c.Args()[0]
		recipient, _ = this.store.FindUserByKey(to)
		if recipient == nil {
			keys, err := crypto.GPGListPublicKeys(to)
			if err == nil && len(keys) == 1 {
				// add to store
				recipient = &common.User{
					Name:  keys[0].Name,
					Email: keys[0].Email,
					Key:   keys[0].PublicKey,
				}
				this.store.AddUser(recipient)

				fmt.Printf("Sending message to %s <%s>...\n", recipient.Name, recipient.Email)
			}
		}
	}

	if recipient == nil {
		fmt.Printf("No recipient found!\n")
		return
	}

	fmt.Printf("Press any key to start recording...\n")
	gopass.GetCh()

	// create a temp file
	fileName := path.Join(os.TempDir(), fmt.Sprintf("%s.aiff", uuid.NewUUID().String()))

	// create a signal
	sig := make(chan int)

	// create a callback func
	cb := func(samples int) {
		maxSamples := int(float64(this.maxDuration/time.Second) * audio.DefaultSampleRate)
		remain := time.Duration(float64(maxSamples-samples)/audio.DefaultSampleRate) * time.Second
		fmt.Printf("\rRecording...%.1f seconds left", remain.Seconds())
	}

	// record
	fmt.Printf("\rRecording...%.1f seconds left", this.maxDuration.Seconds())
	options := audio.RecordOptions{
		FilePath:    fileName,
		MaxDuration: this.maxDuration,
		StopSignal:  sig,
		Callback:    cb,
	}

	if err := audio.RecordAIFF(options); err != nil {
		fmt.Printf("Error recording message! %s", err.Error())
		return
	}

	// encrypt message
	rd, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}
	defer rd.Close()

	fmt.Printf("\rRecorded.\nEncrypting message...\n")
	content, err := crypto.GPGEncrypt(this.user.Key, recipient.Key, rd)
	if err != nil {
		fmt.Printf("Error encrypting message! %s", err.Error())
		return
	}

	msg := &common.Message{
		From:      this.user,
		To:        recipient,
		CreatedAt: time.Now(),
		Content:   content,
	}
	this.store.AddMessage(msg) // Store message before send

	fmt.Printf("Sending...\n")

	err = this.client.Send(msg)
	if err != nil {
		fmt.Printf("Error sending message! %s\n...will retry later.\n", err.Error())
	}
	fmt.Println("Done")
}
