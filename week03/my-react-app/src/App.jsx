import React from "react";
import { useState} from "react";
import TodoInput from "./components/TodoInput";
import TodoList from "./components/TodoList";

import TodoFileter from "./components/TodoFilter";
const App = () => { 
    const [todos, setTodos] = useState([]);
    const [inputValue, setInputValue] = useState('');
    const [filter, setFilter] = useState('All');
    const addTodo = (e) => {
        if (e.key === 'Enter' && inputValue.trim() !== '') {
            const newTodo = {
                id: Date.now(),
                text: inputValue.trim(),
                completed: false,
            };
            setTodos([...todos, newTodo]);
            setInputValue('');
        }

};
    const toggleTodo = (id) => {
        setTodos(
            todos.map((todo) =>
                todo.id === id ? { ...todo, completed: !todo.completed } : todo
            )
        );
    };
    const filteredTodos = todos.filter((todo) => {
        if (filter === 'Active') return !todo.completed;
        if (filter === 'Completed') return todo.completed;
        return true;
    });
    return (
        <div style = {{ width: '400px', margin: '20px auto', padding: '20px', border: '1px solid #ccc', borderRadius: '8px', boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)' }}>
            <TodoInput inputValue={inputValue} onChange={(e) => setInputValue(e.target.value)} onKeyDown={addTodo} /> 
                <TodoFileter filter={filter} onFilterChange={setFilter} onClearCompleted={() => setTodos(todos.filter((todo) => !todo.completed))} />
                    <TodoList fileteredTodos={filteredTodos} onToggleTodo={toggleTodo} />
        </div>
    );
}
export default App;