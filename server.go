package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

type server struct {
	router *chi.Mux
	db     *sql.DB
}

type subscriber struct {
	Email string
	Name  string
}

func (server *server) restartSchema(w http.ResponseWriter, r *http.Request) {
	_, err := server.db.Exec(
		`
    DROP TABLE IF EXISTS newsletter;
    CREATE TABLE newsletter
    (
        email   text,
        name    text
    );
    DROP TABLE IF EXISTS festival;
    CREATE TABLE festival
    (
        email   text,
        name    text
    );
    `,
	) //TODO: unique on emails?
	fmt.Println(err) //TODO: replace with logging
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

	_, err = server.db.Exec(
		`INSERT INTO newsletter(email, name) VALUES ($1, $2)`,
		sub.Email, sub.Name,
	)
	fmt.Println(err) //TODO: replace with logging
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
	server.db, err = sql.Open("sqlite3", "./otq.db")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open database: %v\n", err)
		os.Exit(1)
	}
	defer server.db.Close()

	server.router.Get("/", root)
	server.router.Patch("/restartSchema", server.restartSchema) // TODO: remove before going live
	server.router.Post("/subscribe", server.subscribe)

	fmt.Println("Opening server on :8080")
	http.ListenAndServe(":8080", server.router)
}
