package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	usecases "gobootcamp.com/usecase"
)

func handleError(w http.ResponseWriter, errorStatus int, errorMsg string) {
	fmt.Println("Error: ", errorMsg)
	w.WriteHeader(errorStatus)
	w.Write([]byte(errorMsg))
}

// GetPokemons - Get all the pokemons (no params required)
func GetPokemons(w http.ResponseWriter, r *http.Request) {
	pokes, err := usecases.ReadCsv("./files/pokemons.csv")
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		handleError(w, 500, "There was an error reading the csv file, pls contact the administrator")
		return
	}
	// return all pokemons
	jsonPoke, err := json.Marshal(pokes)
	if err != nil {
		fmt.Println(err)
		handleError(w, 500, "Error converting data to json")
		return
	}
	w.Write(jsonPoke)
}

// GetPokemon - Get all pokemons and filter based on the ID of the required pokemon, if id is missing, error 400 will be throwned
func GetPokemon(w http.ResponseWriter, r *http.Request) {
	id, IDexists := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	if IDexists {
		pokemons, err := usecases.ReadCsv("./files/pokemons.csv")
		if err != nil {
			fmt.Println(err)
			handleError(w, 500, "There was an error reading the csv file, pls contact the administrator")
			return
		}
		// Look for specific pokemon
		idInt, _ := strconv.Atoi(id)
		poke, err := usecases.FindPoke(pokemons, idInt)
		if err != nil {
			fmt.Println(err)
			handleError(w, 500, "The pokemon you are looking for doesnt exists")
			return
		}
		jsonPoke, err := json.Marshal(poke)
		if err != nil {
			fmt.Println(err)
			handleError(w, 500, "Error converting data to json")
			return
		}
		w.Write(jsonPoke)
		return

	} else {
		handleError(w, 400, "ID parameter required")
		return
	}
}

// GetAPIPokemon - Get all the pokemons directly from the API
func GetAPIPokemon(w http.ResponseWriter, r *http.Request) {
	saved, err := usecases.GetPokesFromAPI()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		handleError(w, 500, "There was an error getting the pokes, pls contact the administrator")
		return
	}
	// return all pokemons
	jsonPoke, err := json.Marshal(saved)
	if err != nil {
		handleError(w, 500, "Error converting data to json")
		return
	}
	w.Write(jsonPoke)
}
