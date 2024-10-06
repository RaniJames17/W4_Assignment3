package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Shift struct {
	ID        int       `json:"id"`
	Employee  string    `json:"employee"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

var shifts []Shift
var nextID = 1

func main() {
	http.HandleFunc("/shifts", shiftsHandler)
	http.HandleFunc("/shifts/", shiftsHandler)

	fmt.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func shiftsHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/shifts/")
	if id == "" || id == "/" {
		switch r.Method {
		case http.MethodGet:
			getAllShifts(w, r)
		case http.MethodPost:
			createShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getAllShifts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shifts)
}

func createShift(w http.ResponseWriter, r *http.Request) {

	var shift Shift
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shift.ID = nextID
	nextID++
	shifts = append(shifts, shift)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shift)
}
