import React from "react";
const TodoList = ({ fileteredTodos, onToggleTodo}) => {
    return (
        <div>
            {fileteredTodos.map((todo) => (
                <div key={todo.id} style={{ display: 'flex', alignItems: 'center', padding: '5px 0' }}> 
                    
                    <label>
                        <input
                            type="checkbox"
                            checked={todo.completed}
                            onChange={() => onToggleTodo(todo.id)}
                            style={{ marginRight: '10px' }}
                        />
                        <span style={{ textDecoration: todo.completed ? 'line-through' : 'none' }}></span>
                        {todo.text}
                    </label>
                </div>
            ))}
        </div>
    );
}
export default TodoList;