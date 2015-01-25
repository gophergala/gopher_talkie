package common

import (
	"database/sql"
	"encoding/base64"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path"
	"strings"
	"time"
)

type StoreSqlite struct {
	db *sql.DB
}

type SqliteStoreOptions struct {
	DBPath string
}

var (
	defaultOptions = SqliteStoreOptions{
		DBPath: path.Join(os.Getenv("HOME"), ".talkie", "talkie.db"),
	}

	createUsersTableStmt = `CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		"key" TEXT NOT NULL,
		"name" TEXT NOT NULL,
		"email" TEXT NOT NULL,
		"created_at" TEXT
		); CREATE UNIQUE INDEX IF NOT EXISTS users_idx1 ON users(email, key);`
	createMessagesTableStmt = `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		"from" TEXT NOT NULL,
		"to" TEXT NOT NULL,
		"duration" INTEGER,
		"content" TEXT,
		"created_at" TEXT,
		"played" INTEGER
		);`

	insertUserStmt         = `INSERT OR REPLACE INTO users (name, email, key) VALUES (?, ?, ?)`
	updateUserStmt         = `INSERT OR REPLACE INTO users (name, email, key) VALUES (?, ?, ?)`
	selectUserStmt         = `SELECT id, name, email, key FROM users WHERE id = ?`
	selectUserByKeyStmt    = `SELECT id, name, email, key FROM users WHERE key = ?`
	selectUserMessagesStmt = `SELECT id, "from", "to", "duration", "created_at", "played" FROM messages WHERE "to" = ?`
	selectUserByNameStmt   = `SELECT id, name, email, key FROM users WHERE name = ?`
	deleteUserStmt         = `DELETE users WHERE id = ?`
	deleteUserByKeyStmt    = `DELETE users WHERE key = ?`

	insertMessageStmt       = `INSERT OR REPLACE INTO messages ("from", "to", "duration", "content", "created_at", "played") VALUES (?, ?, ?, ?, ?, ?)`
	selectMessageStmt       = `SELECT id, "from", "to", "duration", "content", "created_at", "played" FROM messages WHERE id = ?`
	deleteMessageStmt       = `DELETE messages WHERE id = ?`
	updateMessagePlayedStmt = `UPDATE messages SET played = ? WHERE id = ?`

	ErrDBNotOpen      = errors.New("db not open")
	ErrNoResult       = errors.New("no result")
	ErrInvalidUser    = errors.New("invalid user")
	ErrInvalidMessage = errors.New("invalid message")
)

func NewStoreSqlite(options *SqliteStoreOptions) *StoreSqlite {
	if options == nil {
		options = &defaultOptions
	}
	db, err := sql.Open("sqlite3", options.DBPath)
	if err != nil {
		panic(err)
	}

	// create tables
	if _, err = db.Exec(createUsersTableStmt); err != nil {
		panic(err)
	}
	if _, err = db.Exec(createMessagesTableStmt); err != nil {
		panic(err)
	}

	return &StoreSqlite{
		db: db,
	}
}

func (s *StoreSqlite) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *StoreSqlite) AddUser(user *User) error {
	if user == nil {
		return ErrInvalidUser
	}
	if s.db == nil {
		return ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(insertUserStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(user.Name, user.Email, user.Key)
	if err != nil {
		return err
	}

	// update UserID with the ROWID
	user.UserID, _ = result.LastInsertId()

	return nil
}

func (s *StoreSqlite) FindUser(userID int64) (*User, error) {
	if s.db == nil {
		return nil, ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(selectUserStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		user := User{}
		err := s.scanUserFromRows(rows, &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, ErrNoResult
}

func (s *StoreSqlite) FindUserByKey(key string) (*User, error) {
	if s.db == nil {
		return nil, ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(selectUserByKeyStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		user := User{}
		err := s.scanUserFromRows(rows, &user)
		if err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, ErrNoResult
}

func (s *StoreSqlite) GetUserMessages(key string) ([]*Message, error) {
	if s.db == nil {
		return nil, ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(selectUserMessagesStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(key)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := Message{}
		err := s.scanMessageFromRows(rows, &msg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (s *StoreSqlite) AddMessage(msg *Message) error {
	if msg == nil {
		return ErrInvalidMessage
	}
	if s.db == nil {
		return ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(insertMessageStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	content := base64.StdEncoding.EncodeToString(msg.Content)
	result, err := stmt.Exec(msg.From.Key, msg.To.Key, msg.Duration.Seconds(), content, msg.CreatedAt.Format(time.RFC3339), msg.Played)
	if err != nil {
		return err
	}

	// update MsgID with the ROWID
	msg.MessageID, _ = result.LastInsertId()

	return nil
}

func (s *StoreSqlite) DeleteMessage(msgID int64) error {
	if s.db == nil {
		return ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(deleteMessageStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(msgID); err != nil {
		return err
	}
	return nil
}

func (s *StoreSqlite) GetMessage(msgID int64) (*Message, error) {
	if s.db == nil {
		return nil, ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(selectMessageStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(msgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		msg := Message{}
		err := s.scanMessageFromRows(rows, &msg)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	}
	return nil, ErrNoResult
}

func (s *StoreSqlite) UpdateMessagePlayed(msgID int64, played bool) error {
	if s.db == nil {
		return ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(updateMessagePlayedStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(played, msgID); err != nil {
		return err
	}
	return nil
}

func (s *StoreSqlite) FindUserByName(name string) ([]*User, error) {
	if s.db == nil {
		return nil, ErrDBNotOpen
	}
	stmt, err := s.db.Prepare(selectUserByNameStmt)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	if rows.Next() {
		user := User{}
		err := s.scanUserFromRows(rows, &user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, ErrNoResult
}

func (s *StoreSqlite) scanUserFromRows(rows *sql.Rows, user *User) error {
	if rows.Err() != nil {
		return rows.Err()
	}
	if user == nil {
		return ErrInvalidUser
	}
	var params []interface{}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	for i := range columns {
		col := columns[i]
		switch strings.ToLower(col) {
		case "id":
			params = append(params, &user.UserID)
		case "name":
			params = append(params, &user.Name)
		case "email":
			params = append(params, &user.Email)
		case "key":
			params = append(params, &user.Key)
		}
	}
	return rows.Scan(params...)
}

func (s *StoreSqlite) scanMessageFromRows(rows *sql.Rows, msg *Message) error {
	if rows.Err() != nil {
		return rows.Err()
	}
	if msg == nil {
		return ErrInvalidMessage
	}

	var from, to, content, createdAt string
	var duration int64
	var params []interface{}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	for i := range columns {
		col := columns[i]
		switch strings.ToLower(col) {
		case "id":
			params = append(params, &msg.MessageID)
		case "from":
			params = append(params, &from)
		case "to":
			params = append(params, &to)
		case "duration":
			params = append(params, &duration)
		case "created_at":
			params = append(params, &createdAt)
		case "played":
			params = append(params, &msg.Played)
		case "content":
			params = append(params, &content)
		}
	}
	err = rows.Scan(params...)
	if err != nil {
		return err
	}

	msg.From, err = s.FindUserByKey(from)
	msg.To, err = s.FindUserByKey(to)
	msg.Content, _ = base64.StdEncoding.DecodeString(content)
	msg.Duration = time.Duration(duration) * time.Second
	msg.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		msg.CreatedAt = time.Unix(0, 0)
	}
	return nil
}
