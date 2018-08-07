class Todo extends React.Component {
    constructor(props) {
        super(props);

        this.ESCAPE_KEY = 27;
        this.ENTER_KEY = 13;
        this.state = {
            editText: this.props.todo.title
        };

        this.handleSubmit = this.handleSubmit.bind(this);
        this.handleChange = this.handleChange.bind(this);
        this.handleEdit = this.handleEdit.bind(this);
        this.handleKeyDown = this.handleKeyDown.bind(this);
    }

    handleSubmit() {
        const val = this.state.editText.trim();
        if (val) {
            this.props.onSave(val);
            this.setState({ editText: val });
        } else {
            this.props.onDestroy();
        }
    }

    handleEdit() {
        this.props.onEdit();
        this.setState({ editText: this.props.todo.title });
    }

    handleKeyDown(event) {
        if (event.which === this.ESCAPE_KEY) {
            this.setState({ editText: this.props.todo.title });
            this.props.onCancel(event);
        } else if (event.which === this.ENTER_KEY) {
            this.handleSubmit();
        }
    }

    handleChange(event) {
        if (this.props.editing) {
            this.setState({ editText: event.target.value });
        }
    }

    shouldComponentUpdate(nextProps, nextState) {
        return (
            nextProps.todo !== this.props.todo ||
            nextProps.editing !== this.props.editing ||
            nextState.editText !== this.state.editText
        );
    }

    componentDidUpdate(prevProps) {
        if (!prevProps.editing && this.props.editing) {
            const node = ReactDOM.findDOMNode(this.refs.editField);
            node.focus();
            node.setSelectionRange(node.value.length, node.value.length);
        }
    }

    render() {
        return (
            <li
                className={classNames({
                    completed: this.props.todo.completed,
                    editing: this.props.editing
                })}
            >
                <div className="view">
                    <input
                        className="toggle"
                        type="checkbox"
                        checked={this.props.todo.completed}
                        onChange={this.props.onToggle}
                    />
                    <label onDoubleClick={this.handleEdit}>
                        {this.props.todo.title}
                    </label>
                    <button
                        className="destroy"
                        onClick={this.props.onDestroy}
                    />
                </div>
                <input
                    ref="editField"
                    className="edit"
                    value={this.state.editText}
                    onBlur={this.handleSubmit}
                    onChange={this.handleChange}
                    onKeyDown={this.handleKeyDown}
                />
            </li>
        );
    }
}

Todo.propTypes = {
    editing: PropTypes.bool,
    editText: PropTypes.string,
    onToggle: PropTypes.func,
    onSave: PropTypes.func,
    onDestroy: PropTypes.func,
    onEdit: PropTypes.func,
    onCancel: PropTypes.func,
    todo: PropTypes.object
};
