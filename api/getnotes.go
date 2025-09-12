package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type note struct {
	Content string `json:"content"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	cxt := context.Background()
	dburl := os.Getenv("db_url")
	if dburl == "" {
		http.Error(w, "Database link failed", http.StatusInternalServerError)
		log.Print("Database link failed")
		return
	}
	dbobj, err := sql.Open("postgres", dburl)

	if err != nil {
		http.Error(w, "Database disruption", http.StatusInternalServerError)
		log.Print("Database disruption")
		return
	}
	defer dbobj.Close()
	var content string

	sc_err := dbobj.QueryRowContext(cxt, "SELECT content FROM notepad WHERE id=1").Scan(&content)
	if sc_err != nil && sc_err != sql.ErrNoRows {
		http.Error(w, "Content not present in the database", http.StatusInternalServerError)
		log.Print("Content not present in the database")
		return
	}

	note := note{Content: content}
	if sc_err == sql.ErrNoRows {
		note.Content = ""
	}
	jsonResponse, json_err := json.Marshal(note)

	if json_err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		log.Print("Failed to marshal JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
