package usecases

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	entities "gobootcamp.com/entity"
)

// ReadCsv - Read from a CSV iterates line by line and return a pointer to an array of pokemons error if fails to open the file (path not found or whateva), or if its corrupt
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

// FindPoke - Based on ID, iterates thru slice of pokemons and find the correct one, if the ID isnt found, then error is returned
func FindPoke(pokemons *[]entities.Pokemon, id int) (entities.Pokemon, error) {
	for _, pokemon := range *pokemons {
		if pokemon.ID == id {
			return pokemon, nil
		}
	}
	return entities.Pokemon{}, errors.New("the pokemon you are looking for doesnt exists")
}

type pokeInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type pokeAPIResp struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []pokeInfo `json:"results"`
}

func createPokemonFile(pokes []pokeInfo, fileName string) error {
	pathFile := fmt.Sprint("./files/", fileName)
	csvFile, err := os.Create(pathFile)
	if err != nil {
		fmt.Println("Error creating the file:", pathFile, err)
		return err
	}

	csvwriter := csv.NewWriter(csvFile)

	err = csvwriter.Write([]string{"id", "name"})
	if err != nil {
		fmt.Println("Error inserting the headers of the file:", err)
		return err
	}
	for indx, poke := range pokes {
		err = csvwriter.Write([]string{fmt.Sprint(indx + 1), poke.Name})
		if err != nil {
			fmt.Printf("Error inserting the pokemon %s (ID: %d) \n", poke.Name, indx+1)
		}
	}
	csvwriter.Flush()
	csvFile.Close()
	return nil
}

// Consume the poke API and saves the result into a CSV file https://pokeapi.co/api/v2/pokemon/?offset=0&limit=151
func GetPokesFromAPI() (map[string]bool, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/?offset=0&limit=151")
	if err != nil {
		fmt.Println("Error consuming the API: ", err)
		return map[string]bool{"saved": false}, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error code: ", resp.StatusCode, "\nError: ", resp.Status)
		return map[string]bool{"saved": false}, nil
	}

	defer resp.Body.Close()
	var jsonResp pokeAPIResp
	// Decoding the response to the struct
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		fmt.Println("Error decoding the body ", err)
		return map[string]bool{"saved": false}, err
	}
	// Creating the csv file
	err = createPokemonFile(jsonResp.Results, "pokemons.csv")
	if err != nil {
		fmt.Println("Error creating the csv file:", err)
		return map[string]bool{"saved": false}, err
	}

	return map[string]bool{"saved": true}, nil
}
