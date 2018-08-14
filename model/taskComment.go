package model

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Todoは管理するタスク
type TodoComment struct {
	ID      int64      `db:"comment_id" json:"id"`
	TodoId  int64      `db:"todo_id" json:"todo_id"`
	Comment string     `json:"comment"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

func TodoCommentsAll(dbx *sqlx.DB, todo_id int) (todoComments []TodoComment, err error) {
	if err := dbx.Select(&todoComments, "SELECT * FROM todo_comments WHERE todo_id = ?", todo_id); err != nil {
		return nil, err
	}
	return todoComments, nil
}

func (t *TodoComment) Insert(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	insert into todo_comments (todo_id, comment)
	values(?, ?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.TodoId, t.Comment)
}

func (t *TodoComment) Delete(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`delete from todos where todo_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.ID)
}
