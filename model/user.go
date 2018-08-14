package model

import (
	// "database/sql"
	// "time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID   int64  `db:"user_id" json:"id"`
	Name string `json:"name"`
}

func UsersAll(dbx *sqlx.DB) (users []User, err error) {
	if err := dbx.Select(&users, "select * from users"); err != nil {
		return nil, err
	}
	return users, nil
}
