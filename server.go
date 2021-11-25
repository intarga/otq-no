package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
)

type server struct {
	router *chi.Mux
	db     *pgx.Conn
}

type subscriber struct {
	Email string
	Name  string
}

func (server *server) restartSchema(w http.ResponseWriter, r *http.Request) {
	tag, err := server.db.Exec(context.Background(),
		`
        DROP SCHEMA public CASCADE;
        CREATE SCHEMA public;
        CREATE TABLE newsletter
        (
            email   text,
            name    text
        );
        CREATE TABLE festival
        (
            email   text,
            name    text
        );
        `,
	) //TODO: unique on emails?
	fmt.Printf("reset: %v\n", tag.String()) //TODO: replace with logging
	fmt.Println(err)
}

func (server *server) subscribe(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var sub subscriber

	err := dec.Decode(&sub) //TODO: sanitise input?
	if err != nil {
		fmt.Printf("Error decoding json: %v\n", err)
		//TODO: log
		return //TODO: return error message?
	}

	tag, err := server.db.Exec(context.Background(),
		`INSERT INTO newsletter(email, name) VALUES ($1, $2)`,
		sub.Email, sub.Name,
	)
	fmt.Printf("subscribe: %v\n", tag.String()) //TODO: replace with logging
	fmt.Println(err)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello root"))
}

func main() {
	var server server

	server.router = chi.NewRouter()
	server.router.Use(middleware.Logger)

	var err error
	//FIXME: db connection not concurrency-safe
	server.db, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer server.db.Close(context.Background())

	server.router.Get("/", root)
	server.router.Patch("/restartSchema", server.restartSchema)
	server.router.Post("/subscribe", server.subscribe)

	fmt.Println("Opening server on :8080")
	http.ListenAndServe(":8080", server.router)
}
