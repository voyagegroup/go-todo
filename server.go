package base

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/voyagegroup/go-todo/controller"
	"github.com/voyagegroup/go-todo/db"

	"github.com/gorilla/context"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// Server はベースアプリケーションのserverを示します
type Server struct {
	dbx    *sqlx.DB
	router *mux.Router
}

// TODO: dbxをstructから分離したほうが複数人数開発だと見通しがよいかもしれない

func (s *Server) Close() error {
	return s.dbx.Close()
}

// Init はServerを初期化する
func (s *Server) Init(dbconf, env string) {
	cs, err := db.NewConfigsFromFile(dbconf)
	if err != nil {
		log.Fatalf("cannot open database configuration. exit. %s", err)
	}
	dbx, err := cs.Open(env)
	if err != nil {
		log.Fatalf("db initialization failed: %s", err)
	}
	s.dbx = dbx
	s.router = s.Route()
}

// New はベースアプリケーションを初期化します
func New() *Server {
	return &Server{}
}

// csrfProtectKey should have 32 byte length.
var csrfProtectKey = []byte("32-byte-long-auth-key")

func (s *Server) Run(addr string) {
	log.Printf("start listening on %s", addr)
	// NOTE: when you serve on TLS, make csrf.Secure(true)
	CSRF := csrf.Protect(
		csrfProtectKey, csrf.Secure(false))
	http.ListenAndServe(addr, context.ClearHandler(CSRF(s.router)))
}

// Route はベースアプリケーションのroutingを設定します
func (s *Server) Route() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "pong")
	}).Methods("GET")
	router.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"token": csrf.Token(r),
		})
	}).Methods("GET")

	todo := &controller.Todo{DB: s.dbx}
	todoComment := &controller.TodoComment{DB: s.dbx}
	user := &controller.User{DB: s.dbx}

	// TODO ng?
	router.Handle("/api/users", handler(user.Get)).Methods("GET")
	router.Handle("/api/todos/completed", handler(todo.DeleteCompleted)).Methods("DELETE")
	router.Handle("/api/todos/toggle_all", handler(todo.ToggleAll)).Methods("POST")
	router.Handle("/api/todos", handler(todo.Get)).Methods("GET")
	router.Handle("/api/todos", handler(todo.Put)).Methods("PUT")                             //更新
	router.Handle("/api/todos", handler(todo.Post)).Methods("POST")                           //挿入
	router.Handle("/api/todos/{todo_id}/comments", handler(todoComment.Get)).Methods("GET")   //取得
	router.Handle("/api/todos/{todo_id}/comments", handler(todoComment.Post)).Methods("POST") //挿入
	router.Handle("/api/todos", handler(todo.Delete)).Methods("DELETE")
	router.Handle("/api/todos/toggle", handler(todo.Toggle)).Methods("PUT")

	// TODO return index.html
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/index.html")
	})
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))

	return router
}
