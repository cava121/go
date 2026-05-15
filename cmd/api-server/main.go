package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var gitCommit string

func main() {
	urlExample := os.Getenv("API_SERVER_DB_URL")

	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		fmt.Println(err);
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

    http.HandleFunc("GET /debug/info", getConfig)
	http.HandleFunc("POST /v1/lists", handleCreateList(conn))
	http.HandleFunc("GET /v1/lists/{id}", handleGetList(conn))

	fmt.Println("Сервер запущен")
    http.ListenAndServe(":8090", nil)
}

func handleGetList(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var name string;

		id := r.PathValue("id")

		if err := conn.QueryRow(r.Context(), `SELECT name FROM lists WHERE id = $1`, id).Scan(&name); err != nil {
			fmt.Println(err);
		}

		var resp struct {
			List struct {
				Id string `json:"id"`
				Name string `json:"name"`
			} `json:"list"`
		}

		resp.List.Id = id
		resp.List.Name = name

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			// how to handle this error?
		}
	}
}

func handleCreateList(conn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			slog.Error("failed decode create list request", slog.String("err", err.Error()))
			http.Error(w, "invalid body json", http.StatusBadRequest)
			return
		}

		if req.Name == "" {
			http.Error(w, "missed required field: name", http.StatusBadRequest)
			return
		}

		listID := uuid.NewString();
		if _, err := conn.Exec(r.Context(), `INSERT INTO lists(id, name) VALUES($1, $2)`, listID, req.Name); err != nil {
			fmt.Println(err);
			slog.Error("failed create list in db", slog.String("err", err.Error()))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		var resp struct {
			List struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"list"`
		}

		resp.List.ID = listID
		resp.List.Name = req.Name

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			// how to handle this error?
		}
	}
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Last commit: ", gitCommit);
}
