import * as React from "react";
import PropTypes from "prop-types";
import classNames from "classnames";

class TodoFooter extends React.Component {
    handleClickAll(e) {
        e.preventDefault();
        this.props.onFooterClick("all");
    }

    handleClickActive(e) {
        e.preventDefault();
        this.props.onFooterClick("active");
    }

    handleClickCompleted(e) {
        e.preventDefault();
        this.props.onFooterClick("completed");
    }

    pluralize(count, word) {
        return count === 1 ? word : word + "s";
    }

    render() {
        const activeTodoWord = this.pluralize(this.props.count, "item");
        let clearButton = null;

        if (this.props.completedCount > 0) {
            clearButton = (
                <button
                    className="clear-completed"
                    onClick={this.props.onClearCompleted}
                >
                    Clear completed
                </button>
            );
        }

        const nowShowing = this.props.nowShowing;
        return (
            <footer id="footer">
                <span id="todo-count">
                    <strong>{this.props.count}</strong> {activeTodoWord} left
                </span>
                <ul id="filters">
                    <li>
                        <a
                            href="#/"
                            onClick={this.handleClickAll.bind(this)}
                            className={classNames({
                                selected: nowShowing === "all"
                            })}
                        >
                            All
                        </a>
                    </li>{" "}
                    <li>
                        <a
                            href="#/active"
                            onClick={this.handleClickActive.bind(this)}
                            className={classNames({
                                selected: nowShowing === "active"
                            })}
                        >
                            Active
                        </a>
                    </li>{" "}
                    <li>
                        <a
                            href="#/completed"
                            onClick={this.handleClickCompleted.bind(this)}
                            className={classNames({
                                selected: nowShowing === "completed"
                            })}
                        >
                            Completed
                        </a>
                    </li>
                </ul>
                {clearButton}
            </footer>
        );
    }
}

TodoFooter.propTypes = {
    count: PropTypes.number,
    completedCount: PropTypes.number,
    nowShowing: PropTypes.string,
    onClearCompleted: PropTypes.func,
    onFooterClick: PropTypes.func
};

export default TodoFooter;
