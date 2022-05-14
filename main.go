package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/pokemons", getPokemons).Methods("GET")

	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
	fmt.Println("Ejemplo")
}

func getPokemons(w http.ResponseWriter, r *http.Request) {
	response := readCsv("..\\files\\pokemons.csv")
	fmt.Println(response)
	// w.Write(response)
	json.NewEncoder(w).Encode(response)
}

func readCsv(csvPath string) []pokemon {
	fmt.Println("CSV file path:", csvPath)
	bulbasor := pokemon{ID: 1, Name: "bulbasor"}
	charmander := pokemon{ID: 4, Name: "charmander"}
	return []pokemon{bulbasor, charmander}
}

// Structs
type pokemon struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
