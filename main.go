package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
	http.HandleFunc("/shifts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllShifts(w, r)
		case http.MethodPost:
			createShift(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/shifts/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/shifts/")
		switch r.Method {
		case http.MethodGet:
			getShiftByID(w, r, id)
		case http.MethodPut:
			updateShift(w, r, id)
		case http.MethodDelete:
			deleteShift(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func getAllShifts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(shifts)
}

func createShift(w http.ResponseWriter, r *http.Request) {
	var newshift Shift
	if err := json.NewDecoder(r.Body).Decode(&newshift); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse start_time and end_time
	startTime, err := time.Parse(time.RFC3339, newshift.StartTime.Format(time.RFC3339))
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, newshift.EndTime.Format(time.RFC3339))
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	newshift.ID = nextID
	nextID++
	newshift.StartTime = startTime
	newshift.EndTime = endTime
	shifts = append(shifts, newshift)
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newshift)
}

func getShiftByID(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid shift ID", http.StatusBadRequest)
		return
	}

	for _, shift := range shifts {
		if shift.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shift)
			return
		}
	}
	http.NotFound(w, r)
}

func updateShift(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid shift ID", http.StatusBadRequest)
		return
	}

	var updatedShift Shift
	if err := json.NewDecoder(r.Body).Decode(&updatedShift); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, shift := range shifts {
		if shift.ID == id {
			shifts[i].Employee = updatedShift.Employee
			shifts[i].StartTime = updatedShift.StartTime
			shifts[i].EndTime = updatedShift.EndTime
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(shifts[i])
			return
		}
	}
	http.NotFound(w, r)
}

func deleteShift(w http.ResponseWriter, r *http.Request, idStr string) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid shift ID", http.StatusBadRequest)
		return
	}

	for i, shift := range shifts {
		if shift.ID == id {
			shifts = append(shifts[:i], shifts[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}
