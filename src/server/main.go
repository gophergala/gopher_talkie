package main

import (
	"bytes"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gophergala/gopher_talkie/src/common"
	"github.com/gophergala/gopher_talkie/src/crypto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

var (
	store     common.Store
	serverKey string
)

func init() {
	os.MkdirAll(path.Join(os.Getenv("HOME"), ".talkie"), 0755)
	store = common.NewStoreSqlite(&common.SqliteStoreOptions{
		DBPath: path.Join(os.Getenv("HOME"), ".talkie", "talkie.db"),
	})
}

func parseJSON(body io.Reader, v interface{}) error {
	d, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(d, v); err != nil {
		return err
	}
	return nil
}

func responseSuccess(w http.ResponseWriter, obj interface{}) {
	res := &Response{
		Success: true,
		Data:    obj,
	}
	responseJSON(w, res)
}

func responseJSON(w http.ResponseWriter, obj interface{}) {
	d, _ := json.Marshal(obj)
	w.Header().Set("Content-Type", "application/json")
	w.Write(d)
}

func responseError(w http.ResponseWriter, err error) {
	res := &Response{
		Success: false,
		Error:   err.Error(),
	}
	responseJSON(w, res)
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusNotFound)
		return
	}

	var user common.User
	if r.Header.Get("Content-Type") == "application/json" {
		err := parseJSON(r.Body, &user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		key := r.PostFormValue("key")
		user = common.User{
			Name:  name,
			Email: email,
			Key:   key,
		}
	}

	keys, _ := crypto.GPGListPublicKeys(user.Key)
	if keys == nil || len(keys) == 0 {
		err := crypto.GPGRecvKey(user.Key)
		if err != nil {
			responseError(w, err)
			return
		}
	}

	if err := store.AddUser(&user); err != nil {
		responseError(w, err)
		return
	}
	responseSuccess(w, &user)
}

func send(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "method not allowed", http.StatusNotFound)
		return
	}

	var msg common.Message
	if r.Header.Get("Content-Type") == "application/json" {
		err := parseJSON(r.Body, &msg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err := store.AddMessage(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	responseSuccess(w, nil)
}

func messages(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	messages, err := store.GetUserMessages(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := &Response{
		Success: true,
		Data:    &messages,
	}
	d, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.FormValue("encrypt") == "1" {
		// encrypt all response content so that only the owner of the key can see the messages!!
		encryptedData, err := crypto.GPGEncrypt(serverKey, key, bytes.NewReader(d))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = w.Write(encryptedData)
	} else {
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(d)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func download(w http.ResponseWriter, r *http.Request) {
	msgID, err := strconv.ParseInt(r.FormValue("id"), 10, 4)
	// since := r.FormValue("since")
	msg, err := store.GetMessage(msgID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if msg == nil {
		http.Error(w, "message not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(msg.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "talkie-server"
	app.Usage = "Secure voicing messaging for geeks"
	app.Version = "0.1.0"
	app.Author = "Tom Li"
	app.Email = "nklizhe@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "host",
			Value: "0.0.0.0",
		},
		cli.IntFlag{
			Name:  "port",
			Value: 3333,
		},
		cli.StringFlag{
			Name: "server-key",
		},
	}
	app.Action = func(c *cli.Context) {
		serverKey = c.String("server-key")

		http.HandleFunc("/register", register)
		http.HandleFunc("/send", send)
		http.HandleFunc("/messages", messages)
		http.HandleFunc("/m", download)

		addr := fmt.Sprintf("%s:%d", c.String("host"), c.Int("port"))
		fmt.Printf("Listening %s...", addr)
		log.Fatal(http.ListenAndServe(addr, nil))
	}
	app.Run(os.Args)
}
