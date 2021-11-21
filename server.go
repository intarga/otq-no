package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
)

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello root"))
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r.Get("/", root)

	http.ListenAndServe(":8080", r)
}
