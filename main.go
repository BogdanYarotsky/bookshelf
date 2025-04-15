package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: helloserver [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	addr = flag.String("addr", "localhost:8080", "address to serve")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) != 0 {
		usage()
	}

	connString := os.Getenv("DATABASE_URL")
	if len(connString) == 0 {
		log.Fatalln("Got empty db connection string")
	}

	log.Println(connString)

	pool, err := pgxpool.New(
		context.Background(),
		connString)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer pool.Close()

	book := BookHandler{db: pool}
	http.HandleFunc("/book", book.Handle)

	log.Printf("serving http://%s\n", *addr)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
