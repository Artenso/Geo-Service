package model

type User struct {
	ID   int64  `json:"id" db:"id"`
	Name string `json:"username" db:"name"`
	Pass []byte `json:"password" db:"password"`
}
