package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Solutions struct
type Solution struct {
	A     int `json:"a"`
	B     int `json:"b"`
	C     int `json:"c"`
	Roots int `json:"n_roots"`
}

// DB var
var DB []Solution

func solve(a int, b int, c int) int {
	if a == 0 {
		if b != 0 {
			return 1
		}

		return 0
	}

	d := (b * b) - (4 * a * c)

	if d > 0 {
		return 2
	} else if d == 0 {
		return 1
	}

	return 0
}

// Solve func
func Solve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := strconv.Atoi(vars["a"])
	if err != nil {
		log.Fatal(err)
	}
	b, err := strconv.Atoi(vars["b"])
	if err != nil {
		log.Fatal(err)
	}
	c, err := strconv.Atoi(vars["c"])
	if err != nil {
		log.Fatal(err)
	}

	roots := solve(a, b, c)

	DB = append(DB, Solution{A: a, B: b, C: c, Roots: roots})

	w.WriteHeader(http.StatusAccepted)
}

// GetSolution func
func GetSolution(w http.ResponseWriter, r *http.Request) {
	if len(DB) == 0 {
		w.WriteHeader(http.StatusNotFound)
	} else {
		solution := DB[len(DB)-1]

		json.NewEncoder(w).Encode(solution)
	}
}

func main() {
	DB = []Solution{}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/solve/{a}/{b}/{c}", Solve).Methods("POST")
	r.HandleFunc("/solution", GetSolution).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
