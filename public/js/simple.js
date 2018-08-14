const e = React.createElement;

class TodoApp extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      newTodo: "",
      todos: [],
      todoComments: new Immutable.Map(),
      editText: "",
    };

    this.token = this.fetchToken();
    this.load();
  }

  fetchToken() {
    fetch("/token", {credentials: "same-origin"})
      .then(x => x.json())
      .then(json => {
        if (json === null) {
          return;
        }
        this.token = json.token;
      })
      .catch(err => {
        console.error("fetch error", err);
      });
  }

  load() {
    return fetch("/api/todos", {
      credentials: "same-origin",
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json"
      },
    })
      .then(resp => {
        if (resp.status !== 200) {
          throw new Error(resp.statusText);
        }
        return resp;
      })
      .catch(err => {
        console.error("get todo error: ", err);
      })
      .then(x => x.json())
      .then(json => {
        if (json === null) {
          return;
        }
        this.setState({todos: json});
      });
  }

  loadComments(todo) {
    const {todoComments} = this.state;

    return fetch(`/api/todos/${todo.id}/comments`, {
      credentials: "same-origin",
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json"
      },
    })
      .then(resp => {
        if (resp.status !== 200) {
          throw new Error(resp.statusText);
        }
        return resp;
      })
      .catch(err => {
        console.error("get todo error: ", err);
      })
      .then(x => x.json())
      .then(json => {
        if (json === null) {
          return;
        }
        this.setState({
          todoComments:
            todoComments.set(todo.id, json)
        });
      });
  }

  addComment(todo, comment) {
    const {todoComments} = this.state;

    return fetch(`/api/todos/${todo.id}/comments`, {
      credentials: "same-origin",
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify({
        todo_id: todo.id,
        comment,
      }),
    })
      .then(resp => {
        if (resp.status === 201) {
          return resp;
        }

        throw new Error(resp.statusText);
      })
      .then(x => x.json())
      .then(data => {
        let comments=todoComments.get(todo.id)
        if (!comments) {
          todoComments = 
            todoComments.set(todo.id, comments = []);
        }
        comments.push(data)
        this.setState({todoComments});
      })
      .catch(err => {
        console.error("post comment error: ", err);
      });
  }

  addTodo(title) {
    if (title === "") {
      return;
    }
    const todo = {title: title, completed: false};

    return fetch("/api/todos", {
      credentials: "same-origin",
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify(todo),
    })
      .then(resp => {
        if (resp.status === 201) {
          return resp;
        }

        throw new Error(resp.statusText);
      })
      .then(x => x.json())
      .then(data => {
        this.setState({
          todos: [...this.state.todos, data],
          newTodo: ""
        });
      })
      .catch(err => {
        console.error("post todo error: ", err);
      });
  }

  update(todo) {
    const {todos} = this.state;

    return fetch("/api/todos", {
      credentials: "same-origin",
      method: "PUT",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify(todo),
    })
    .then(res => res.json())
    .then(newTodo => {
      this.setState({
        todos: todos.map(c=>c.id===todo.id?newTodo:c)
      });
    });
  }

  toggleAll(completed) {
    const {todos} = this.state;

    return fetch("/api/todos/toggle_all", {
      credentials: "same-origin",
      method: "POST",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify({completed}),
    })
    .then(() => {
      this.setState({
        todos: todos.map(c=>({...c, completed}))
      });
    });
  }

  destroy(todo) {
    const {todos} = this.state;

    return fetch("/api/todos", {
      credentials: "same-origin",
      method: "DELETE",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify(todo),
    }).then(() => {
      this.setState({
        todos: todos.filter(candidate => {
            return candidate !== todo;
        })
      });
    });
  }

  toggle(todoToToggle) {
    const {todos} = this.state;

    return fetch("/api/todos/toggle", {
      credentials: "same-origin",
      method: "PUT",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify(todoToToggle),
    })
      .then(() => {
        this.setState({
          todos: todos.map(todo => {
            return todo !== todoToToggle ? todo : Object.assign({}, todo, { completed: !todo.completed });
          })
        });
      });
  }

  renderTodos() {
    const {todos, editText, todoComments} = this.state;

    return todos.map(t => {
      const tc = todoComments.get(t.id);
      let writing_comment = "";
      return e("li", {key: t.id},
        e("div", {className: "view"},
          e("input", {className: "toggle", type: "checkbox", checked: t.completed, onChange: () => {
            this.toggle(t);
          }}),
          e("input", {className: "title", type: "text", defaultValue: t.title, onChange: (e) => {
            t.title = e.target.value;
          }}),
          e("button", {
            className: "update",
            onClick: () => {
              this.update(t);
            }
          }, "Update"),
          e("button", {
            className: "destroy",
            onClick: () => {
              this.destroy(t);
            }
          }, "Delete"),
        ),
        e("div", {className: "comment_box"},
          tc?(
            e("div", {},
              e("ul", {className: "comment_list"}, tc.map(
                ({id, comment}) => e("li", {className: "comment", key: t.id + "_c" + id}, comment)
              )),
              e("div", {className: "comment_input_box"}, 
                e("input", {className: "comment_input", type: "text", onChange: (e) => {
                  writing_comment = e.target.value;
                }}),
                e("button", {
                  className: "send",
                  onClick: () => {
                    this.addComment(t, writing_comment);
                  }
                }, "Send"),
              )
            )
          ):(
            e("button", {className: "comment_open", onClick: (e) => {
              this.loadComments(t);
            }}, `${t.comment_cnt}件のコメント...`)
          )
        ),
        // 何のためにあるのかよく分からなかった…グローバルのステート…const↓
        e("input", {className: "edit", value: editText}),
      );
    });
  }

  render() {
    const {newTodo} = this.state;

    return e("div", {},
      e("header", {id: "header"},
        e("input", {
          id: "new-todo",
          placeholder: "What needs to be done?",
          value: newTodo,
          autoFocus: true,
          onChange: event => {
            this.setState({newTodo: event.target.value});
          }
        }),
        e("button", {
            onClick: () => {
              this.addTodo(newTodo);
            }
          },
          "Add"
        ),
      ),
      e("div", {},
        e("label", {}, 
          "toggleAll",
          e("input", {
            type: "checkbox",
            onClick: (e) => {
              this.toggleAll(e.target.checked);
            }
          })
        )
      ),
      e("div", {}, e("ul", {id: "todo-list"}, this.renderTodos()),
      ),
    );
  }
}

ReactDOM.render(e(TodoApp), document.getElementById("app"));
