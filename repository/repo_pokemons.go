package repository

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	entities "gobootcamp.com/entity"
)

type repoPokemon struct {
}

type RepositoryPokemons interface {
	ReadCsvConcurrent(string, chan entities.Pokemon, context.CancelFunc)
	ReadCsv(string) (*[]entities.Pokemon, error)
	ReadAPIPokemon() (entities.PokeAPIResp, error)
	CreatePokemonFile([]entities.PokeInfo, string) error
}

func NewRepoPokemon() RepositoryPokemons {
	return &repoPokemon{}
}

func (*repoPokemon) ReadCsvConcurrent(filename string, src chan entities.Pokemon, cancel context.CancelFunc) {
	csvfile, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(csvfile)
	go func() {
		defer csvfile.Close()
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
				fmt.Println("Error reading the excel file")
				break
			}
			name := record[1]
			id, err := strconv.Atoi(record[0])
			if err != nil {
				fmt.Println("Unknown ID - ", record)
				continue
			}
			src <- entities.Pokemon{ID: id, Name: name} // you might select on ctx.Done().
		}
		close(src) // close src to signal workers that no more job are incoming.
	}()
}

// ReadCsv - Read from a CSV iterates line by line and return a pointer to an array of pokemons error if fails to open the file (path not found or whateva), or if its corrupt
func (*repoPokemon) ReadCsv(csvPath string) (*[]entities.Pokemon, error) {
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

func (*repoPokemon) ReadAPIPokemon() (entities.PokeAPIResp, error) {
	resp, err := http.Get("https://pokeapi.co/api/v2/pokemon/?offset=0&limit=151")
	if err != nil {
		fmt.Println("Error consuming the API: ", err)
		return entities.PokeAPIResp{}, err
	}

	if resp.StatusCode != 200 {
		fmt.Println("Error code: ", resp.StatusCode, "\nError: ", resp.Status)
		return entities.PokeAPIResp{}, fmt.Errorf("error code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	var jsonResp entities.PokeAPIResp
	// Decoding the response to the struct
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		fmt.Println("Error decoding the body ", err)
		return entities.PokeAPIResp{}, err
	}
	return jsonResp, nil
}

func (*repoPokemon) CreatePokemonFile(pokes []entities.PokeInfo, fileName string) error {
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
