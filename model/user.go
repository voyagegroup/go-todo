package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// Todoは管理するタスク
type User struct {
	ID      int64      `db:"user_id" json:"id"`
	Name    string     `json:"name"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

func UsersAll(dbx *sqlx.DB) (users []User, err error) {
	err = dbx.Select(&users, "select * from users")
	return
}
