package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/rs/cors"
)

type Todo struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
	CreatedAt string `json:"created_at"`
}

var db *pgx.Conn

func main() {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}
	var err error
	db, err = pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}
	defer db.Close(ctx)

	// ensure table exists (migrations recommended)
	_, err = db.Exec(ctx, `
    CREATE TABLE IF NOT EXISTS todos (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        completed BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT now()
    );
    `)
	if err != nil {
		log.Fatalf("failed to create table: %v", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/todos", listTodos).Methods("GET")
	r.HandleFunc("/api/todos", createTodo).Methods("POST")
	r.HandleFunc("/api/todos/{id}", toggleTodo).Methods("PUT")
	r.HandleFunc("/api/health", health).Methods("GET")

	// enable CORS for local dev; adjust allowed origins for production
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Handler:      handler,
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("backend listening on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func listTodos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := db.Query(ctx, "SELECT id, title, completed, created_at FROM todos ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	out := make([]Todo, 0)
	for rows.Next() {
		var t Todo
		var ts time.Time
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed, &ts); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		t.CreatedAt = ts.Format(time.RFC3339)
		out = append(out, t)
	}
	writeJSON(w, out)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var t Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if t.Title == "" {
		http.Error(w, "title required", 400)
		return
	}
	var id int64
	err := db.QueryRow(ctx, "INSERT INTO todos (title) VALUES ($1) RETURNING id", t.Title).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	t.ID = id
	writeJSON(w, t)
}

func toggleTodo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := db.Exec(ctx, "UPDATE todos SET completed = NOT completed WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]string{"status": "ok"})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}
