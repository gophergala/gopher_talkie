package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"
)

func randomDBPath() string {
	return path.Join(os.TempDir(), fmt.Sprintf("%06x.db", rand.Uint32()&0xFFFFFF))
}

func createStore(dbPath string) Store {
	if dbPath == "" {
		dbPath = randomDBPath()
	}
	err := os.RemoveAll(dbPath)
	if err != nil {
		panic(err)
	}
	return NewStoreSqlite(&SqliteStoreOptions{
		DBPath: dbPath,
	})
}

func createRandomUsers(store Store, n int) ([]*User, error) {
	var users []*User
	for i := 0; i < n; i++ {
		user := &User{
			Name:  fmt.Sprintf("Tester%3d", i),
			Email: fmt.Sprintf("tester%3d@example.com", i),
			Key:   fmt.Sprintf("%08x", rand.Uint32()),
		}
		err := store.AddUser(user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func createRandomMessages(store Store, from, to *User, n int) ([]*Message, error) {
	var messages []*Message

	for i := 0; i < n; i++ {
		msg := NewMessage(from, to)
		err := store.AddMessage(msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func TestNewStoreSqlite(t *testing.T) {
	dbPath := randomDBPath()

	store := createStore(dbPath)
	assert.NotNil(t, store)

	_, err := os.Stat(dbPath)
	assert.Nil(t, err, "database file should be created")
}

func TestAddUser(t *testing.T) {
	store := createStore("")
	assert.NotNil(t, store)
	defer store.Close()

	u := &User{
		Name:  "Tester",
		Email: "tester@example.com",
		Key:   "5678ABCD",
	}
	err := store.AddUser(u)
	assert.Nil(t, err)
	assert.True(t, u.UserID > 0)

	u1, err := store.FindUserByKey(u.Key)
	assert.Nil(t, err)
	assert.NotNil(t, u1)
	assert.Equal(t, u.UserID, u1.UserID)
	assert.Equal(t, u.Name, u1.Name)
	assert.Equal(t, u.Email, u1.Email)

	// u.Name = "Tester2"
	// err = store.UpdateUser(u)
	// assert.Nil(t, err)

	// u2, err := store.FindUserByKey(u.Key)
	// assert.Nil(t, err)
	// assert.Equal(t, u.Name, u2.Name)
}

func TestAddMessage(t *testing.T) {
	store := createStore("")
	assert.NotNil(t, store)
	defer store.Close()

	users, err := createRandomUsers(store, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))

	m := &Message{
		From:      users[0],
		To:        users[1],
		Duration:  time.Duration(1) * time.Second,
		Content:   []byte("test"),
		CreatedAt: time.Now(),
		Played:    false,
	}
	err = store.AddMessage(m)
	assert.Nil(t, err)
	assert.True(t, m.MessageID > 0)

	m1, err := store.GetMessage(m.MessageID)
	assert.Nil(t, err)
	assert.NotNil(t, m1)
	assert.Equal(t, m.From.Key, m1.From.Key)
	assert.Equal(t, m.To.Key, m1.To.Key)
	assert.Equal(t, m.Duration, m1.Duration)
	assert.Equal(t, m.Content, m1.Content)
	assert.Equal(t, m.CreatedAt.Unix()/1000, m1.CreatedAt.Unix()/1000)
	assert.Equal(t, m.Played, m1.Played)

	m.Played = true
	err = store.UpdateMessagePlayed(m.MessageID, true)
	assert.Nil(t, err)

	m2, err := store.GetMessage(m.MessageID)
	assert.Nil(t, err)
	assert.True(t, m2.Played)
}

func TestUserMessages(t *testing.T) {

	store := createStore("")
	assert.NotNil(t, store)
	defer store.Close()

	// create 10 random users
	users, err := createRandomUsers(store, 10)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(users))
	length := len(users)

	// random pick two users
	idx1 := rand.Intn(length)
	u1 := users[idx1]
	assert.NotNil(t, u1)

	idx2 := (idx1 + rand.Intn(length) + 1) % length
	u2 := users[idx2]
	assert.NotNil(t, u2)
	assert.NotEqual(t, u1.UserID, u2.UserID)

	// create 10 random messages
	messages, err := createRandomMessages(store, u1, u2, 10)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(messages))

	// get messages
	m2, err := store.GetUserMessages(u2.Key)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(m2))
}
