package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/voyagegroup/go-todo/model"

	"github.com/jmoiron/sqlx"
)

// TodoComment はTodoCommentへのリクエストに関する制御をします
type TodoComment struct {
	DB *sqlx.DB
}

// Get はDBからTodoCommentを取得して結果を返します
func (t *TodoComment) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["todo_id"])
	if err != nil {
		return err
	}
	todos, err := model.TodoCommentsAll(t.DB, todoID)
	if err != nil {
		return err
	}
	return JSON(w, 200, todos)
}

// Post はDBにTodoCommentを挿入して結果を返します
func (t *TodoComment) Post(w http.ResponseWriter, r *http.Request) error {
	var todoComment model.TodoComment
	if err := json.NewDecoder(r.Body).Decode(&todoComment); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		result, err := todoComment.Insert(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		todoComment.ID, err = result.LastInsertId()
		return err
	}); err != nil {
		return err
	}

	return JSON(w, http.StatusCreated, todoComment)
}
