import * as ReactDOM from "react-dom";
import * as React from "react";
import TodoBox from "./app/TodoBox";
import TodoModel from "./app/TodoModel";

const model = new TodoModel();
function r() {
    ReactDOM.render(
        <TodoBox model={model} pollInterval={5000} />,
        document.getElementById("app")
    );
}
model.subscribe(r);
