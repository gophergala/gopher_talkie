package common

type User struct {
	UserID    string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
}
