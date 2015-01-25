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
	if this.user == nil {
		panic(ErrNoUser)
		return
	}

	if len(c.Args()) == 0 {
		fmt.Printf("Usage: %s send <to>\n", os.Args[0])
		return
	}
	to := c.Args()[0]

	// TODO: find key of the recipient

	fmt.Printf("Press any key to start recording your message...\n")
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

	content, err := crypto.GPGEncrypt(this.user.Key, to, rd)
	if err != nil {
		fmt.Printf("Error encrypt message! %s", err.Error())
		return
	}

	// ask user to select if there are multiple users
	fmt.Printf("Encrypting message...\n")
	msg := &common.Message{
		From: this.user,
		To: &common.User{
			Key: to,
		},
		CreatedAt: time.Now(),
		Content:   content,
	}
	this.store.AddMessage(msg) // Store message before send

	fmt.Printf("Sending...\n", msg.To.Name)
	fmt.Printf("Done\n")
}
