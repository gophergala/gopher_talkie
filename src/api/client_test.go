package api

import (
	"github.com/gophergala/gopher_talkie/src/common"
	"github.com/stretchr/testify/assert"

	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true}`))
	})
	ts := httptest.NewServer(mux)
	u, err := url.Parse(ts.URL)
	assert.Nil(t, err)

	c := NewClient(u.Host)
	url := c.GetURL("hello", nil)
	assert.Equal(t, fmt.Sprintf("%s/hello", ts.URL), url)

	user := &common.User{
		Name:  "Tester1",
		Email: "tester1@example.com",
		Key:   "123456",
	}
	err = c.Register(user)
	assert.Nil(t, err)
}
