package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: helloserver [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	greeting = flag.String("g", "Hello", "Greet with `greeting`")
	addr     = flag.String("addr", "localhost:8080", "address to serve")
)

func main() {
	//flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		usage()
	}

	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			Database: "bookshelf",
			User:     "bookshelf",
			Password: "bookshelf",
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer pool.Close()

	http.HandleFunc("/name", func(w http.ResponseWriter, r *http.Request) {
		conn, err := pool.Acquire()
		if err != nil {
			http.Error(w, "Could not connect.", 500)
			return
		}
		defer conn.Close()

		if r.Method == "GET" {
			var names []string
			rows, err := conn.Query("SELECT first_name from first_table")
			if err != nil {
				log.Printf("Query error: %v", err)
				http.Error(w, "Query failed.", http.StatusInternalServerError)
				return
			}
			defer rows.Close()
			for rows.Next() {
				var name string
				err = rows.Scan(&name)
				if err != nil {
					http.Error(w, "todo", 500)
					return
				}
				names = append(names, name)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(names)

		} else if r.Method == "POST" {
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "todo", 500)
			}

			body := string(bodyBytes)
			_, err = conn.Exec("INSERT INTO first_table (first_name) VALUES ($1)", body)
			if err != nil {
				http.Error(w, "Insert failed.", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Name added successfully")
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Printf("serving http://%s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
