package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gobootcamp.com/repository"
	usecases "gobootcamp.com/usecase"

	"github.com/gorilla/mux"
)

func handleError(w http.ResponseWriter, errorStatus int, errorMsg string) {
	fmt.Println("Error: ", errorMsg)
	w.WriteHeader(errorStatus)
	w.Write([]byte(errorMsg))
}

// GetPokemons - Get all the pokemons (no params required)
func GetPokemons(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewRepoPokemon()
	pokes, err := repo.ReadCsv("./files/pokemons.csv")
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
	repo := repository.NewRepoPokemon()
	id, IDexists := mux.Vars(r)["id"]
	w.Header().Set("Content-Type", "application/json")
	if IDexists {
		pokemons, err := repo.ReadCsv("./files/pokemons.csv")
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
	repo := repository.NewRepoPokemon()
	saved, err := usecases.GetPokesFromAPI(repo)
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

func GetConcurrentPokemon(w http.ResponseWriter, r *http.Request) {
	repo := repository.NewRepoPokemon()
	params := mux.Vars(r)
	regType := params["type"]
	numItems, err := strconv.Atoi(params["items"])
	if err != nil {
		handleError(w, 400, "Param items invalid value")
		return
	}

	itemsPerWorker, err := strconv.Atoi(params["items_per_worker"])
	if err != nil {
		handleError(w, 400, "Param items_per_worker invalid value")
		return
	}
	// var pokes []entities.Pokemon
	start := time.Now()
	pokes, _ := usecases.GetPokesFromCsvConcurrent(itemsPerWorker, numItems, regType, repo)
	w.Header().Set("Content-Type", "application/json")
	jsonPoke, err := json.Marshal(pokes)
	if err != nil {
		fmt.Println(err)
		handleError(w, 500, "Error converting data to json")
		return
	}
	w.Write(jsonPoke)
	elapsed := time.Since(start)
	log.Printf("Time passed: %v\n", elapsed)
}
