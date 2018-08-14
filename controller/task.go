package controller

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/voyagegroup/go-todo/model"

	"github.com/jmoiron/sqlx"
)

// Todo はTodoへのリクエストに関する制御をします
type Todo struct {
	DB *sqlx.DB
}

// ToggleAll はDBの全ての完了状態を変更し結果を返します
func (t *Todo) ToggleAll(w http.ResponseWriter, r *http.Request) error {
	c := struct {
		Completed bool `json:"completed"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		return err
	}
	var result sql.Result
	var err error
	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		result, err = model.TodosToggleAll(tx, c.Completed)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return JSON(w, 200, result)
}

// DeleteCompleted はDBから完了状態のタスクを全て削除し結果を返します
func (t *Todo) DeleteCompleted(w http.ResponseWriter, r *http.Request) error {
	var result sql.Result
	var err error
	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		result, err = model.TodosDeleteCompleted(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return JSON(w, 200, result)
}

// Get はDBからTodoを取得して結果を返します
func (t *Todo) Get(w http.ResponseWriter, r *http.Request) error {
	title := r.FormValue("title")
	if title == "" {
		todos, err := model.TodosAll(t.DB)
		if err != nil {
			return err
		}
		return JSON(w, 200, todos)
	}
	todos, err := model.TodoByTitle(t.DB, title)
	if err != nil {
		return err
	}
	return JSON(w, 200, todos)
}

// Put はDBのTodoを更新し結果を返します
func (t *Todo) Put(w http.ResponseWriter, r *http.Request) error {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		_, err := todo.Update(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return JSON(w, http.StatusOK, todo)
}

// Post はタスクをDBに追加し結果を返します
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

// Delete はDBのTodoを削除し結果を返します
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

// Toggle はDBのTodoの完了状態を変更します
func (t *Todo) Toggle(w http.ResponseWriter, r *http.Request) error {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		_, err := todo.Toggle(tx)
		if err != nil {
			return err
		}
		return tx.Commit()
	}); err != nil {
		return err
	}

	return JSON(w, http.StatusNotImplemented, nil)
}
