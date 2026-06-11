package main

import (
	"context"
	"curs/cmd/api-server/internal/app"
	"curs/cmd/api-server/internal/store"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

	store := store.New(conn)
	app := app.New(store)

	http.HandleFunc("GET /debug/info", getConfig)
	http.HandleFunc("POST /v1/lists", handleCreateList(app))
	http.HandleFunc("GET /v1/lists/{id}", handleGetListById(app))
	http.HandleFunc("GET /v1/lists", handleGetLists(app))
	http.HandleFunc("OPTIONS /{path...}", handleOptions)

	fmt.Println("Сервер запущен")
	http.ListenAndServe(":8090", corsMiddleware(http.DefaultServeMux))
}

func handleGetLists(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := a.GetLists(r.Context())

		if err != nil {

		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
		}
	}
}

func handleGetListById(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println(id)

		result, err := a.GetListById(r.Context(), id);
		fmt.Println(result)
		if err != nil {
		}

		var resp struct {
			List struct {
				Id string `json:"id"`
				Name string `json:"name"`
			} `json:"list"`
		}

		resp.List.Id = result.Id
		resp.List.Name = result.Name

		if err := json.NewEncoder(w).Encode(resp); err != nil {
		}
	}
}

func handleCreateList(a *app.App) http.HandlerFunc {
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

		list, err := a.CreateList(r.Context(), req.Name)

		if err != nil {
			// TODO
		}

		var resp struct {
			List struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"list"`
		}

		resp.List.ID = list.Id
		resp.List.Name = req.Name

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			// how to handle this error?
		}
	}
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusNoContent)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		next.ServeHTTP(w, r)
	})
}

func getConfig(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Last commit: ", gitCommit);
}
