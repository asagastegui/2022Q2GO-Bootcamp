package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gobootcamp.com/usecases"
)

// GetPokemons - Get all the pokemons (no params required)
func GetPokemons(w http.ResponseWriter, r *http.Request) {
	pokemons, err := usecases.ReadCsv("./files/pokemons.csv")
	if err != nil {
		fmt.Println("Error: ", err)
		// TODO: Fix to change to use this approach (as its better to check if encoding/marshalling actually worked)
		// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		json.NewEncoder(w).Encode(map[string]string{"error": "There was an error reading the csv file, pls contact the administrator"})
	}
	// TODO: Fix to change to use this approach (as its better to check if encoding/marshalling actually worked)
	// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
	// return all pokemons
	json.NewEncoder(w).Encode(pokemons)
}

// GetPokemon - Get all pokemons and filter based on the ID of the required pokemon, if id is missing, error 400 will be throwned
func GetPokemon(w http.ResponseWriter, r *http.Request) {
	id, IDexists := mux.Vars(r)["id"]
	if IDexists {
		pokemons, err := usecases.ReadCsv("./files/pokemons.csv")
		if err != nil {
			fmt.Println("Error: ", err)
			json.NewEncoder(w).Encode(map[string]string{"error": "There was an error reading the csv file, pls contact the administrator"})
		}
		// Look for specific pokemon
		idInt, _ := strconv.Atoi(id)
		poke, err := usecases.FindPoke(pokemons, idInt)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "The pokemon you are looking for doesnt exists"})
			return
		}
		// TODO: Fix to change to use this approach (as its better to check if encoding/marshalling actually worked)
		// https://stackoverflow.com/questions/31622052/how-to-serve-up-a-json-response-using-go
		json.NewEncoder(w).Encode(poke)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("400 - missing ID"))
	}
}
