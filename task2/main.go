package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Item struct
type Item struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Amount int     `json:"amount"`
	Price  float64 `json:"price"`
}

// DB var
var DB []Item

// GetItems func
func GetItems(w http.ResponseWriter, r *http.Request) {
	jsonEncoder := json.NewEncoder(w)

	if len(DB) == 0 {
		w.WriteHeader(http.StatusNotFound)
		jsonEncoder.Encode(map[string]string{"Error": "No one items in stock!"})
	} else {
		jsonEncoder.Encode(DB)
	}
}

// GetItem func
func GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}
	jsonEncoder := json.NewEncoder(w)
	found := false

	for _, item := range DB {
		if id == item.ID {
			found = true

			jsonEncoder.Encode(item)
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		jsonEncoder.Encode(map[string]string{"Error": "Item with that id not found!"})
	}
}

// AddItem func
func AddItem(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var item Item

	json.Unmarshal(reqBody, &item)

	DB = append(DB, item)

	w.WriteHeader(http.StatusCreated)
}

// UpdateItem func
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}
	jsonEncoder := json.NewEncoder(w)
	found := false

	for index, item := range DB {
		if id == item.ID {
			found = true

			reqBody, _ := ioutil.ReadAll(r.Body)
			json.Unmarshal(reqBody, &DB[index])

			w.WriteHeader(http.StatusAccepted)
			jsonEncoder.Encode(DB[index])
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		jsonEncoder.Encode(map[string]string{"Error": "Item with that id not found!"})
	}
}

// DeleteItem func
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
	}
	jsonEncoder := json.NewEncoder(w)
	found := false

	for index, item := range DB {
		if id == item.ID {
			found = true
			DB = append(DB[:index], DB[index+1:]...)

			w.WriteHeader(http.StatusAccepted)
			jsonEncoder.Encode(item)

			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		jsonEncoder.Encode(map[string]string{"Error": "Item with that id not found!"})
	}
}

func main() {
	DB = []Item{}

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/items", GetItems).Methods("GET")
	r.HandleFunc("/item/{id}", GetItem).Methods("GET")
	r.HandleFunc("/item", AddItem).Methods("POST")
	r.HandleFunc("/item/{id}", UpdateItem).Methods("PUT")
	r.HandleFunc("/item/{id}", DeleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
