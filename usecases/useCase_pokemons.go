package usecases

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"gobootcamp.com/entities"
)

// Read from a CSV iterates line by line and return a pointer to an array of pokemons
// error if fails to open the file (path not found or whateva), or if its corrupt
func ReadCsv(csvPath string) (*[]entities.Pokemon, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	// closing the file at the end of readCsv
	defer file.Close()
	csvReader := csv.NewReader(file)

	var pokemons []entities.Pokemon
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
		pokemons = append(pokemons, entities.Pokemon{ID: id, Name: name})
	}
	return &pokemons, nil
}

func FindPoke(pokemons *[]entities.Pokemon, id int) (entities.Pokemon, error) {
	for _, pokemon := range *pokemons {
		if pokemon.ID == id {
			return pokemon, nil
		}
	}
	return entities.Pokemon{}, errors.New("the pokemon you are looking for doesnt exists")
}
