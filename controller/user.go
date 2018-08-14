package controller

import (
	// "encoding/json"
	"net/http"

	"github.com/voyagegroup/go-todo/model"

	"github.com/jmoiron/sqlx"
)

type User struct {
	DB *sqlx.DB
}

func (u *User) Get(w http.ResponseWriter, r *http.Request) error {
	users, err := model.UsersAll(u.DB)
	if err != nil {
		return err
	}
	return JSON(w, 200, users)
}
