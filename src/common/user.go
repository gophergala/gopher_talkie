package common

type User struct {
	UserID int64  `json:"id"`
	Key    string `json:"key"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
