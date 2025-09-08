package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
)

type ac_check struct {
	Accesscode string `json:"code"`
	Content    string `json:"content"`
}

func ACUHandler(w http.ResponseWriter, r *http.Request) {
	cxt := context.Background()
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
		return
	}
	var acc ac_check
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&acc)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	acode := os.Getenv("access_code")
	if acc.Accesscode != acode {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	dburl := os.Getenv("db_url")
	dbobj, err := sql.Open("postgres", dburl)
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer dbobj.Close()
	_, dberr := dbobj.ExecContext(cxt, "INSERT INTO notepad (content) VALUES ($1) ON CONFLICT (id) DO UPDATE SET content=EXCLUDED.content;", acc.Content)
	if dberr != nil {
		http.Error(w, "Failed to update the database", http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
