package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/suzuken/go-todo/model"

	"github.com/jmoiron/sqlx"
)

// Todo はTodoへのリクエストに関する制御をします
type Todo struct {
	DB *sqlx.DB
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

	w.WriteHeader(http.StatusOK)
	return nil
}

func (t *Todo) DeleteMulti(w http.ResponseWriter, r *http.Request) error {
	var todos []model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todos); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		for _, todo := range todos {
			if _, err := todo.Delete(tx); err != nil {
				return err
			}
		}
		return tx.Commit()
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (t *Todo) Toggle(w http.ResponseWriter, r *http.Request) error {
	var todo model.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		result, err := todo.Toggle(tx)
		if err != nil {
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
		if n, err := result.RowsAffected(); err != nil {
			return err
		} else if n != 1 {
			return errors.New("no rows updated")
		}
		return nil
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (t *Todo) ToggleAll(w http.ResponseWriter, r *http.Request) error {
	var req = struct {
		Checked bool `json:"checked"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	if err := TXHandler(t.DB, func(tx *sqlx.Tx) error {
		if _, err := model.TodosToggleAll(tx, req.Checked); err != nil {
			return err
		}
		return tx.Commit()
	}); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
