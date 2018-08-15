package controller

import (
	"net/http"

	"github.com/voyagegroup/go-todo/model"

	"github.com/jmoiron/sqlx"
)

// User はUserへのリクエストに関する制御をします
type User struct {
	DB *sqlx.DB
}
type userName struct {
	Name string `json:"name"`
}

// Get はDBからユーザを取得して結果を返します
func (t *User) Get(w http.ResponseWriter, r *http.Request) error {
	users, err := model.UsersAll(t.DB)
	if err != nil {
		return err
	}
	var userNames []userName
	for _, user := range users {
		userNames = append(userNames, userName{user.Name})
	}

	return JSON(w, 200, userNames)
}
