import React from "react";
const TodoFileter = ({ filter, onFilterChange, onClearCompleted }) => {
    const filterMap = {
        All: '所有',
        Active: '进行中',
        Completed: '已完成',
    };
    return (
        <div style={{ display: 'flex', justifyContent: 'space-between', marginTop: '10px' }}>
            <div>
                {Object.keys(filterMap).map((type) => (
                    <button
                        key={type}
                        onClick={() => onFilterChange(type)}
                        style={{ marginRight: '10px',
                            padding: '5px 10px',
                            border: 'none',
                            borderRadius: '4px',
                            backgroundColor: filter === type ? '#007bff' : '#e0e0e0',
                            color: filter === type ? '#fff' : '#000',
                            cursor: 'pointer',
                        }}>
                        {filterMap[type]}
                        </button>
                ))
                         }
            </div>
            <button onClick={onClearCompleted} style={{ padding: '5px 10px',
                border: 'none',
                borderRadius: '4px',
                backgroundColor: '#dc3545',
                color: '#fff',
                cursor: 'pointer',
            }}>清空</button>
        </div>
    );
}
export default TodoFileter;