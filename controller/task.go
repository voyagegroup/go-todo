package controller

import (
	"encoding/json"
	"net/http"

	"github.com/voyagegroup/go-todo/model"

	"github.com/jmoiron/sqlx"
)

// Todo はTodoへのリクエストに関する制御をします
type Todo struct {
	DB *sqlx.DB
}

type ErrorMessage struct {
	M string `json:"message"`
}

// GetはDBからユーザを取得して結果を返します
func (t *Todo) Get(w http.ResponseWriter, r *http.Request) error {
	todos, err := model.TodosAll(t.DB)
	if err != nil {
		return err
	}
	return JSON(w, 200, todos)
}

func (t *Todo) Put(w http.ResponseWriter, r *http.Request) error {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		result, err := todo.Update(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		todo.ID, err = result.LastInsertId()
		return err
	}); err != nil {
		return err
	}

	return JSON(w, http.StatusOK, todo)
}

// PostはタスクをDBに追加します
// todoをJSONとして受け取ることを想定しています。
func (t *Todo) Post(w http.ResponseWriter, r *http.Request) error {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		result, err := todo.Insert(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		todo.ID, err = result.LastInsertId()
		return err
	}); err != nil {
		return err
	}

	return JSON(w, http.StatusCreated, todo)
}

func (t *Todo) Delete(w http.ResponseWriter, r *http.Request) error {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		_, err := todo.Delete(tx)
		if err != nil {
			return err
		}
		return tx.Commit()
	}); err != nil {
		return err
	}

	return JSON(w, http.StatusOK, todo)
}

func (t *Todo) Toggle(w http.ResponseWriter, r *http.Request) error {
	return JSON(w, http.StatusNotImplemented, nil)
}

func (t *Todo) Search(w http.ResponseWriter, r *http.Request) error {
	q := r.URL.Query().Get("q")
	if q == "" {
		return JSON(w, http.StatusBadRequest, &ErrorMessage{"Required parameter q is missing."})
	}

	qs, e := model.NewSearchQuery(q)
	if e != nil {
		return JSON(w, http.StatusBadRequest, &ErrorMessage{e.Error()})
	}

	todos, err := model.Search(t.DB, qs)
	if err != nil {
		return err
	}
	return JSON(w, 200, todos)
}
