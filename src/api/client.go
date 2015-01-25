package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gophergala/gopher_talkie/src/common"
	"github.com/gophergala/gopher_talkie/src/crypto"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrInvalidRequest        = errors.New("invalid request")
	ErrUnexpectedContentType = errors.New("unexpected content type")
)

type Client struct {
	serverAddr string
}

func NewClient(addr string) *Client {
	return &Client{
		serverAddr: addr,
	}
}

type RegisterResponse struct {
	Success bool         `json:"success"`
	Data    *common.User `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

func (c *Client) GetURL(p string, query *url.Values) string {
	if c.serverAddr == "" {
		c.serverAddr = "127.0.0.1:3333"
	}
	if query != nil && len(query.Encode()) > 0 {
		return fmt.Sprintf("http://%s/%s?%s", c.serverAddr, p, query.Encode())
	} else {
		return fmt.Sprintf("http://%s/%s", c.serverAddr, p)
	}
}

func (c *Client) Register(user *common.User) error {
	if user == nil {
		return ErrInvalidRequest
	}
	url := c.GetURL("register", nil)
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// fmt.Printf("register: %v\n", user)
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if res.Header.Get("Content-Type") != "application/json" {
		return ErrUnexpectedContentType
	}
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var reg RegisterResponse
	if err := json.Unmarshal(d, &reg); err != nil {
		return err
	}

	if !reg.Success {
		return errors.New(reg.Error)
	}

	// success, update userID
	if reg.Data != nil {
		user.UserID = reg.Data.UserID
	}
	return nil
}

type SendResponse struct {
	Success bool   `json:"success"`
	Data    int64  `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (c *Client) Send(msg *common.Message) error {
	if msg == nil {
		return ErrInvalidRequest
	}
	url := c.GetURL("send", nil)
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if res.Header.Get("Content-Type") != "application/json" {
		return ErrUnexpectedContentType
	}
	d, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var s SendResponse
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}

	if !s.Success {
		return errors.New(s.Error)
	}

	// success, update messageID
	msg.MessageID = s.Data
	return nil
}

type MessagesResponse struct {
	Success  bool              `json:"success"`
	Messages []*common.Message `json:"data,omitempty"`
	Error    string            `json:"error,omitempty"`
}

func (c *Client) GetMessages(user *common.User) ([]*common.Message, error) {
	if user == nil {
		return nil, ErrInvalidRequest
	}
	query := &url.Values{}
	query.Set("key", user.Key)
	url := c.GetURL("messages", query)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	var s MessagesResponse
	if res.Header.Get("Content-Type") == "application/octet-stream" {
		decryptedData, err := crypto.GPGDecrypt(user.Key, res.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(decryptedData, &s); err != nil {
			return nil, err
		}
	} else if res.Header.Get("Content-Type") == "application/json" {
		d, _ := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(d, &s); err != nil {
			return nil, err
		}
	}

	if !s.Success {
		return nil, errors.New(s.Error)
	}

	// success
	return s.Messages, nil
}

// download message content
func (c *Client) DownloadMessage(msgID int64) ([]byte, error) {
	query := &url.Values{}
	query.Set("id", fmt.Sprintf("%d", msgID))
	url := c.GetURL("m", query)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if res.Header.Get("Content-Type") != "application/octet-stream" {
		return nil, ErrUnexpectedContentType
	}

	// success
	return ioutil.ReadAll(res.Body)
}
