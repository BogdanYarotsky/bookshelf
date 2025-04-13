package main

import "github.com/jackc/pgx"

type DB interface {
	Query(query string, args ...interface{}) (*pgx.Rows, error)
	Exec(query string, arguments ...interface{}) (pgx.CommandTag, error)
}
