const e = React.createElement;

class TodoApp extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      newTodo: "",
      todos: [],
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
        Accept: "application/json",
        "Content-Type": "application/json"
      }
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

  addTodo(title) {
    if (title === "") {
      return;
    }
    const todo = {title: title, completed: false};

    return fetch("/api/todos", {
      credentials: "same-origin",
      method: "POST",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
        "X-CSRF-Token": this.token
      },
      body: JSON.stringify(todo)
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

  renderTodos() {
    const {todos, editText} = this.state;

    return todos.map(t => {
      return e("li", {key: t.id},
        e("div", {className: "view"},
          e("input", {className: "toggle", type: "checkbox"}),
          e("label", {}, t.title),
          e("button", {className: "destroy"}, "Delete"),
        ),
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
      e("div", {}, e("ul", {id: "todo-list"}, this.renderTodos())),
    );
  }
}

ReactDOM.render(e(TodoApp), document.getElementById("app"));
