import React from "react";
const  TodoInput = ({ inputValue, onChange, onKeyDown }) => {
    return(
        <input
        type="txet"
        value={inputValue}
        onChange={onChange}
        onKeyDown={onKeyDown}
        placeholder="添加任务"
        style={{
            width: '100%',
            padding: '10px',
            fontSize: '16px',
            boxSizing: 'border-box',
        }}
        />
    )


}
export default TodoInput;