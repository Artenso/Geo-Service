package model

type User struct {
	Name string `json:"username" db:"name"`
	Pass []byte `json:"password" db:"password"`
}
