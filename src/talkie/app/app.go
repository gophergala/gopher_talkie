package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/gopher_talkie/src/api"
	"github.com/gophergala/gopher_talkie/src/common"
	"github.com/gophergala/gopher_talkie/src/crypto"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const (
	Version            = "0.1.0"
	Author             = "Tom Li"
	Email              = "nklizhe@gmail.com"
	DefaultMaxDuration = time.Duration(15) * time.Second
)

var (
	ErrNoKeyFound = errors.New("no key found")
	ErrNoUser     = errors.New("no user")
)

type AppConfig struct {
	CurrentUser string `json:"current_user,omitempty"`
}

type App struct {
	app         *cli.App
	config      *AppConfig
	user        *common.User
	store       common.Store
	maxDuration time.Duration
	client      *api.Client
}

func NewApp() *App {

	this := &App{}
	app := cli.NewApp()
	app.Name = "talkie"
	app.Usage = "Secure voicing messaging for geeks"
	app.Action = func(c *cli.Context) {
		// default
		this.list(c)
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "server",
			Value: "130.211.156.226:3333",
		},
	}
	app.Version = Version
	app.Author = Author
	app.Email = Email
	app.Commands = []cli.Command{
		NewListCommand(this),
		NewSendCommand(this),
		// NewPlayCommand(this),
		// NewDeleteCommand(this),
	}

	os.MkdirAll(path.Join(os.Getenv("HOME"), ".talkie"), 0750)

	this.app = app
	this.store = common.NewStoreSqlite(&common.SqliteStoreOptions{
		DBPath: path.Join(os.Getenv("HOME"), ".talkie", "talkie.db"),
	})
	this.maxDuration = DefaultMaxDuration
	this.config, _ = this.loadConfig()

	return this
}

func (this *App) Run() {
	this.loadConfig()

	this.app.Run(os.Args)
}

func (this *App) loadConfig() (*AppConfig, error) {
	confdir := path.Join(os.Getenv("HOME"), ".talkie")
	if err := os.MkdirAll(confdir, 0750); err != nil {
		return nil, err
	}

	conffile := path.Join(confdir, "config")
	rd, err := os.Open(conffile)
	if err != nil {
		return nil, err
	}
	defer rd.Close()

	d, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	var cfg AppConfig
	err = json.Unmarshal(d, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (this *App) saveConfig(cfg *AppConfig) error {
	confdir := path.Join(os.Getenv("HOME"), ".talkie")
	if err := os.MkdirAll(confdir, 0750); err != nil {
		return err
	}
	conffile := path.Join(confdir, "config")
	wd, err := os.Create(conffile)
	if err != nil {
		return err
	}
	defer wd.Close()

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	_, err = wd.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (this *App) selectCurrentUser() *common.User {
	// List all users of gpg
	keys, err := crypto.GPGListSecretKeys("")
	if err != nil {
		return nil
	}
	if len(keys) == 0 {
		return nil
	}

	var idx int
	if len(keys) > 0 {
		fmt.Fprintf(os.Stderr, "There more than one keypairs, please select which one do you want to use:\n")
		for i := range keys {
			k := keys[i]
			fmt.Fprintf(os.Stderr, "  (%d) %s %s <%s>\n", i+1, k.PublicKey, k.Name, k.Email)
		}
		fmt.Fprintf(os.Stderr, "Enter number (%d - %d) > ", 1, len(keys))

		var choice int
		for {
			fmt.Scanf("%v", &choice)
			if choice > 0 && choice <= len(keys) {
				break
			}
		}
		idx = choice - 1
	}

	k := keys[idx]
	user := &common.User{
		Name:  k.Name,
		Email: k.Email,
		Key:   k.PublicKey,
	}
	err = this.store.AddUser(user)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
	}
	return user
}

func (this *App) setup(c *cli.Context) error {
	if this.client == nil {
		this.client = api.NewClient(c.GlobalString("server"))
	}

	if this.user == nil {
		// try load user from config
		if this.config != nil && this.config.CurrentUser != "" {
			this.user, _ = this.store.FindUserByKey(this.config.CurrentUser)
		}

		// still not found
		if this.user == nil {
			this.user = this.selectCurrentUser()
		}
	}

	if this.user == nil {
		return ErrNoUser
	}

	if err := this.client.Register(this.user); err != nil {
		return err
	}

	// save config
	if this.config == nil {
		this.config = &AppConfig{
			CurrentUser: this.user.Key,
		}
		this.saveConfig(this.config)
	}

	return nil
}
