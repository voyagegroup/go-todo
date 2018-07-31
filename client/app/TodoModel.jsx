const get = path => {
    return fetch(path, {
        credentials: "same-origin",
        method: "GET",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json"
        }
    })
        .then(resp => checkStatus(resp, 200))
        .then(x => x.json())
        .catch(err => {
            console.error("get todo error: ", err);
        });
};

const internalFetch = (method, path, token, status, data) => {
    return fetch(path, {
        credentials: "same-origin",
        method,
        headers: {
            Accept: "application/json",
            "Content-Type": "application/json",
            "X-CSRF-Token": token
        },
        body: JSON.stringify(data)
    })
        .then(resp => checkStatus(resp, status))
        .then(x => x.json())
        .catch(err => {
            console.error("post todo error: ", err);
        });
};

const post = (path, token, status, data) => {
    return internalFetch("POST", path, token, status, data);
};

const put = (path, token, status, data) => {
    return internalFetch("PUT", path, token, status, data);
};

const deleteFetch = (path, token, status, data) => {
    return internalFetch("DELETE", path, token, status, data);
};

const checkStatus = (resp, code) => {
    if (resp.status === code) {
        return resp;
    } else {
        const error = new Error(resp.statusText);
        error.resp = resp;
        throw error;
    }
};

class TodoModel {
    constructor() {
        this.onChanges = [];
        this.todos = [];
        this.token = this.fetchToken();
    }

    subscribe(onChange) {
        this.onChanges.push(onChange);
        this.inform();
    }

    inform() {
        this.onChanges.forEach(cb => cb());
    }

    fetchToken() {
        fetch("/token", { credentials: "same-origin" })
            .then(x => x.json())
            .then(json => {
                if (json == null) {
                    return;
                }
                this.token = json.token;
            })
            .catch(err => {
                console.error("fetch error", err);
            });
    }

    load() {
        get("/api/todos").then(json => {
            if (json == null) {
                return;
            }
            this.todos = json;
            this.inform();
        });
    }

    addTodo(title) {
        const todo = {
            title: title,
            completed: false
        };

        post("/api/todos", this.token, 201, todo).then(data => {
            this.todos = this.todos.concat([data]);
            this.inform();
        });
    }

    toggleAll(checked) {
        post("/api/todos/toggleall", this.token, 200, { checked }).then(() => {
            this.todos = this.todos.map(todo => {
                return Object.assign({}, todo, { completed: checked });
            });
            this.inform();
        });
    }

    toggle(todoToToggle) {
        post("/api/todos/toggle", this.token, 200, todoToToggle).then(() => {
            this.todos = this.todos.map(todo => {
                return todo !== todoToToggle
                    ? todo
                    : Object.assign({}, todo, { completed: !todo.completed });
            });
            this.inform();
        });
    }

    destroy(todo) {
        deleteFetch("/api/todos", this.token, 200, todo).then(() => {
            this.todos = this.todos.filter(candidate => {
                return candidate !== todo;
            });
            this.inform();
        });
    }

    save(todoToSave, text) {
        const toSave = Object.assign({}, todoToSave, { title: text });

        put("/api/todos", this.token, 200, toSave).then(() => {
            this.todos = this.todos.map(todo => {
                return todo !== todoToSave ? todo : toSave;
            });
            this.inform();
        });
    }

    clearCompleted() {
        const todosToDelete = this.todos.filter(todo => {
            return todo.completed;
        });

        deleteFetch("/api/todos/multi", this.token, 200, todosToDelete).then(
            () => {
                this.todos = this.todos.filter(todo => {
                    return !todo.completed;
                });
                this.inform();
            }
        );
    }
}

export default TodoModel;
