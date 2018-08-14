package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Todoは管理するタスク
type Todo struct {
	ID         int64      `db:"todo_id" json:"id"`
	Title      string     `json:"title"`
	Completed  bool       `json:"completed"`
	Created    *time.Time `json:"created"`
	Updated    *time.Time `json:"updated"`
	CommentCnt int        `db:"comment_cnt" json:"comment_cnt"`
}

const allColStmt = "t.todo_id, t.title, t.completed, t.created, t.updated, SUM(CASE WHEN c.comment_id IS NULL THEN 0 ELSE 1 END) AS comment_cnt"
const fromStmt = "todos t LEFT OUTER JOIN todo_comments c ON t.todo_id = c.todo_id"

func TodosAll(dbx *sqlx.DB) (todos []Todo, err error) {
	err = dbx.Select(&todos, "SELECT "+allColStmt+" from "+fromStmt+" GROUP BY t.todo_id")
	return
}

func TodoOne(dbx *sqlx.DB, id int64) (*Todo, error) {
	var todo Todo
	if err := dbx.Get(&todo, `select `+allColStmt+` from `+fromStmt+` where t.todo_id = ?`, id); err != nil {
		return nil, err
	}
	return &todo, nil
}

func TodoByTitle(dbx *sqlx.DB, title string) (todos []Todo, err error) {
	err = dbx.Select(&todos, `select `+allColStmt+` from `+fromStmt+` WHERE t.title LIKE concat("%", ?, "%")`, title)
	return
}

// TodosToggleAllは全部のtoggleのステータスをトグルします
func TodosToggleAll(tx *sqlx.Tx, checked bool) (sql.Result, error) {
	stmt, err := tx.Prepare(`update todos set completed = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	fmt.Printf("[%#v]", checked)
	return stmt.Exec(checked)
}

func TodosDeleteCompleted(tx *sqlx.Tx) (sql.Result, error) {
	return tx.Exec("DELETE FROM todos WHERE completed=1")
}

func (t *Todo) Update(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update todos set title = ? where todo_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.Title, t.ID)
}

func (t *Todo) Insert(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	insert into todos (title, completed)
	values(?, ?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.Title, t.Completed)
}

// Toggle は指定されたタスクについて現在の状態と入れ替えます。
func (t *Todo) Toggle(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update todos set completed=?
	where todo_id=?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(!t.Completed, t.ID)
}

func (t *Todo) Delete(tx *sqlx.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`delete from todos where todo_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.ID)
}

// TodosDeleteAllはすべてのタスクを消去します。
// テストのために使用されます。
func TodosDeleteAll(tx *sqlx.Tx) (sql.Result, error) {
	return tx.Exec(`truncate table todos`)
}
