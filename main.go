package main

import (
	"flag"
	"fmt"
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
	flag.Usage = usage
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
	s := Server{db: pool}

	http.HandleFunc("/book", s.handleBook)

	log.Printf("serving http://%s\n", *addr)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
