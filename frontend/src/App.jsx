import React, { useEffect, useState } from "react";

const API = import.meta.env.VITE_API_BASE || "http://localhost:8080";

export default function App() {
  const [todos, setTodos] = useState([]);
  const [title, setTitle] = useState("");

  useEffect(() => {
    fetch(`${API}/api/todos`).then((r) => r.json()).then(setTodos);
  }, []);

  const add = async (e) => {
    e.preventDefault();
    if (!title.trim()) return;
    const res = await fetch(`${API}/api/todos`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ title })
    });
    const t = await res.json();
    setTodos([...todos, t]);
    setTitle("");
  };

  const toggle = async (id) => {
    await fetch(`${API}/api/todos/${id}`, { method: "PUT" });
    setTodos(todos.map(t => t.id === id ? { ...t, completed: !t.completed } : t));
  };

  return (
    <div style={{ padding: 20, fontFamily: "Arial, sans-serif" }}>
      <h1>Todo</h1>
      <form onSubmit={add}>
        <input value={title} onChange={e => setTitle(e.target.value)} placeholder="Add a todo" />
        <button type="submit">Add</button>
      </form>
      <ul>
        {todos.map(t => (
          <li key={t.id}>
            <label>
              <input type="checkbox" checked={t.completed} onChange={() => toggle(t.id)} />
              {" "}
              {t.title}
            </label>
          </li>
        ))}
      </ul>
    </div>
  );
}
