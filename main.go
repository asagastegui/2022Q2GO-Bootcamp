package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Path("/pokemons").HandlerFunc(getPokemons).Methods("GET")
	router.Path("/pokemons/{id:[0-9]+}").HandlerFunc(getPokemons).Methods("GET")

	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getPokemons(w http.ResponseWriter, r *http.Request) {
	pokemons, err := readCsv("./files/pokemons.csv")
	if err != nil {
		fmt.Println("Error: ", err)
		json.NewEncoder(w).Encode(map[string]string{"error": "There was an error reading the csv file, pls contact the administrator"})
	}

	id, IDexists := mux.Vars(r)["id"]
	if IDexists {
		// Look for specific pokemon
		idInt, _ := strconv.Atoi(id)
		poke, err := findPoke(pokemons, idInt)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(map[string]string{"error": "The pokemon you are looking for doesnt exists"})
			return
		}
		json.NewEncoder(w).Encode(poke)
		return
	}
	// return all pokemons
	json.NewEncoder(w).Encode(pokemons)
}

func findPoke(pokemons *[]pokemon, id int) (pokemon, error) {
	for _, pokemon := range *pokemons {
		if pokemon.ID == id {
			return pokemon, nil
		}
	}
	return pokemon{}, errors.New("the pokemon you are looking for doesnt exists")
}

func readCsv(csvPath string) (*[]pokemon, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	// closing the file at the end of readCsv
	defer file.Close()
	csvReader := csv.NewReader(file)

	var pokemons []pokemon
	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Getting ID and Name
		id, err := strconv.Atoi(line[0])
		if err != nil {
			fmt.Println("Unknown ID - ", line)
			continue
		}
		name := line[1]
		pokemons = append(pokemons, pokemon{ID: id, Name: name})
	}
	return &pokemons, nil
}

// Structs
type pokemon struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
