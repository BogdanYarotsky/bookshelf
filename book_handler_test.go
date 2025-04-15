package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

func TestThatNameIsPreserved(t *testing.T) {
	ctx := context.Background()
	conn := setupTestDB(ctx, t)
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	t.Cleanup(func() {
		_ = tx.Rollback(ctx)
	})

	handler := BookHandler{db: tx}

	// Insert 2 books via POST
	books := []string{"Golang Mastery", "The Pragmatic Programmer"}

	for _, b := range books {
		req := httptest.NewRequest(http.MethodPost, "/book", strings.NewReader(b))
		w := httptest.NewRecorder()
		handler.Handle(w, req)
		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("POST failed, got %d", resp.StatusCode)
		}
	}

	// Now GET all books
	req := httptest.NewRequest(http.MethodGet, "/book", nil)
	w := httptest.NewRecorder()
	handler.Handle(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET failed, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	for _, b := range books {
		if !strings.Contains(string(body), b) {
			t.Errorf("expected response to contain %q, got %q", b, string(body))
		}
	}
}

func setupTestDB(ctx context.Context, t *testing.T) *pgx.Conn {
	t.Helper()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Fatal("DATABASE_URL not set")
	}

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		t.Fatalf("failed to connect to db: %v", err)
	}

	t.Cleanup(func() {
		conn.Close(ctx)
	})

	return conn
}
