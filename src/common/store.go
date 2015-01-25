package common

type Store interface {
	AddUser(user *User) error
	FindUser(userID int64) (*User, error)
	FindUserByName(name string) ([]*User, error)
	FindUserByKey(key string) (*User, error)
	GetUserMessages(key string) ([]*Message, error)

	AddMessage(msg *Message) error
	UpdateMessagePlayed(msgID int64, played bool) error
	DeleteMessage(msgID int64) error
	GetMessage(msgID int64) (*Message, error)

	Close()
}
