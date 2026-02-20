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
	Body      string `json:"body"`
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
	rows, err := db.Query(ctx, "SELECT id, title, body, completed, created_at FROM todos ORDER BY id")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	out := make([]Todo, 0)
	for rows.Next() {
		var t Todo
		var ts time.Time
		if err := rows.Scan(&t.ID, &t.Title, &t.Body, &t.Completed, &ts); err != nil {
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
	var input struct {
		Title     string `json:"title"`
		Body      string `json:"body"`
		Timestamp string `json:"timestamp"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", 400)
		return
	}
	if input.Title == "" {
		http.Error(w, "title required", 400)
		return
	}
	var id int64
	var createdAt time.Time
	err := db.QueryRow(ctx,
		"INSERT INTO todos (title, body, created_at) VALUES ($1, $2, $3) RETURNING id, created_at",
		input.Title, input.Body, input.Timestamp,
	).Scan(&id, &createdAt)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	t := Todo{
		ID:        id,
		Title:     input.Title,
		Body:      input.Body,
		Completed: false,
		CreatedAt: createdAt.Format(time.RFC3339),
	}
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
