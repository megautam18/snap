package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type acrequest struct {
	AccessCode string `json:"code"`
}

func AccessCodeHandler(w http.ResponseWriter, r *http.Request) {
	/*client:=&http.Client{
		Timeout:30*time.Second,
	}*/
	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect method", http.StatusMethodNotAllowed)
		return
	}
	var reqdata acrequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqdata)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	acode := os.Getenv("access_code")
	if reqdata.AccessCode == acode {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Invalid access code", http.StatusUnauthorized)
	}

}
