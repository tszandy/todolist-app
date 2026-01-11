import React, { useEffect, useState } from "react";
import { DragDropContext, Droppable, Draggable } from "react-beautiful-dnd";

const API = import.meta.env.VITE_API_BASE || "http://localhost:8080";

export default function App() {
  const [todos, setTodos] = useState([]);
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");

  useEffect(() => {
    fetch(`${API}/api/todos`).then((r) => r.json()).then(setTodos);
  }, []);

  const add = async (e) => {
    e.preventDefault();
    if (!title.trim()) return;
    const timestamp = new Date().toISOString();
    const res = await fetch(`${API}/api/todos`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ title, body, timestamp })
    });
    const t = await res.json();
    setTodos([...todos, t]);
    setTitle("");
    setBody("");
  };

  const toggle = async (id) => {
    await fetch(`${API}/api/todos/${id}`, { method: "PUT" });
    setTodos(todos.map(t => t.id === id ? { ...t, completed: !t.completed } : t));
  };

  const onDragEnd = (result) => {
    if (!result.destination) return;
    const items = Array.from(todos);
    const [reorderedItem] = items.splice(result.source.index, 1);
    items.splice(result.destination.index, 0, reorderedItem);
    setTodos(items);
  };

  return (
    <div style={{ padding: 20, fontFamily: "Arial, sans-serif" }}>
      <h1>Todo</h1>
      <form onSubmit={add}>
        <input
          value={title}
          onChange={e => setTitle(e.target.value)}
          placeholder="Add a todo title"
        />
        <input
          value={body}
          onChange={e => setBody(e.target.value)}
          placeholder="Add a todo body"
        />
        <button type="submit">Add</button>
      </form>
      <DragDropContext onDragEnd={onDragEnd}>
        <Droppable droppableId="todos">
          {(provided) => (
            <ul {...provided.droppableProps} ref={provided.innerRef}>
              {todos.map((t, idx) => (
                <Draggable key={t.id} draggableId={String(t.id)} index={idx}>
                  {(provided) => (
                    <li ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                      <label>
                        <input type="checkbox" checked={t.completed} onChange={() => toggle(t.id)} />
                        {" "}
                        <strong>{t.title}</strong>
                        <div style={{ color: "#555", fontSize: 14 }}>{t.body}</div>
                        <span style={{ marginLeft: 10, color: "#888", fontSize: 12 }}>
                          {t.timestamp ? new Date(t.timestamp).toLocaleString() : ""}
                        </span>
                      </label>
                    </li>
                  )}
                </Draggable>
              ))}
              {provided.placeholder}
            </ul>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  );
}
