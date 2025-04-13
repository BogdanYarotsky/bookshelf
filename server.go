package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Server struct {
	db DB
}

func (s *Server) handleBook(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			var names []string
			rows, err := s.db.Query("SELECT first_name from first_table")
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
		}
	case "POST":
		{
			defer r.Body.Close()
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "todo", 500)
			}

			body := string(bodyBytes)
			_, err = s.db.Exec("INSERT INTO first_table (first_name) VALUES ($1)", body)
			if err != nil {
				http.Error(w, "Insert failed.", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Name added successfully")
		}
	default:
		{
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
