package usecases

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	entities "gobootcamp.com/entity"
)

// FindPoke - Based on ID, iterates thru slice of pokemons and find the correct one, if the ID isnt found, then error is returned
func FindPoke(pokemons *[]entities.Pokemon, id int) (entities.Pokemon, error) {
	for _, pokemon := range *pokemons {
		if pokemon.ID == id {
			return pokemon, nil
		}
	}
	return entities.Pokemon{}, errors.New("the pokemon you are looking for doesnt exists")
}

type RepositoryPokemons interface {
	ReadCsvConcurrent(string, chan entities.Pokemon, context.CancelFunc)
	ReadCsv(string) (*[]entities.Pokemon, error)
	ReadAPIPokemon() (entities.PokeAPIResp, error)
	CreatePokemonFile([]entities.PokeInfo, string) error
}

// Consume the poke API and saves the result into a CSV file https://pokeapi.co/api/v2/pokemon/?offset=0&limit=151
func GetPokesFromAPI(repo RepositoryPokemons) (map[string]bool, error) {
	jsonResp, err := repo.ReadAPIPokemon()
	if err != nil {
		fmt.Println("Error getting pokemosn from api:", err)
		return map[string]bool{"saved": false}, err
	}
	// Creating the csv file
	err = repo.CreatePokemonFile(jsonResp.Results, "pokemons.csv")
	if err != nil {
		fmt.Println("Error creating the csv file:", err)
		return map[string]bool{"saved": false}, err
	}

	return map[string]bool{"saved": true}, nil
}

func createWorkers(ctx context.Context, wg *sync.WaitGroup, numberOfWorkers, itemsPerWorker int, dst chan<- entities.Pokemon, src <-chan entities.Pokemon, regType string) {
	// declare the workers
	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(ctx, dst, src, regType, id, itemsPerWorker)
		}(i + 1)
	}
}

func worker(ctx context.Context, dst chan<- entities.Pokemon, src <-chan entities.Pokemon, regType string, workerID, itemsPerWorker int) {
	counter := 0
	for {
		select {
		case poke, ok := <-src: // you must check for readable state of the channel.
			if !ok {
				log.Println("worker ", workerID, "no more tasks in src (chan closed)")
				return
			}
			// If looking for evens
			// log.Println("worker", workerID, "counter", counter, "pokemon", poke.Name)
			if regType == "even" {
				if poke.ID%2 == 0 {
					dst <- poke // do somethingg useful.
					counter++
				}
			} else if regType == "odd" {
				// If looking for odds
				if poke.ID%2 != 0 {
					dst <- poke // do somethingg useful.
					counter++
				}
			}
			if counter == itemsPerWorker {
				// log.Println("Worker", workerID, " - Max task limit reached - shuting down..")
				return
			}

		case <-ctx.Done(): // if the context is cancelled, quit.
			log.Println("worker ", workerID, "closed by ctx done")
			return
		}
	}
}

func GetPokesFromCsvConcurrent(itemsPerWorker, numItems int, regType string, repo RepositoryPokemons) ([]entities.Pokemon, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	numberOfWorkers := 10
	fileName := "./files/pokemons.csv"
	src := make(chan entities.Pokemon)
	dst := make(chan entities.Pokemon)

	// Creating the workers
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		createWorkers(ctx, &wg, numberOfWorkers, itemsPerWorker, dst, src, regType)
		wg.Done()
	}()
	wg.Add(1)
	// Reading the CSV concurrently (sending src channel so he can send the pokemons read)
	go func() {
		repo.ReadCsvConcurrent(fileName, src, cancel)
		wg.Done()
	}()
	go func() {
		wg.Wait()

		close(dst)
	}()
	// drain the output

	itemsProcessed := 0
	var pokes []entities.Pokemon
	for poke := range dst {
		fmt.Println(poke.ID, poke.Name)
		pokes = append(pokes, poke)
		itemsProcessed++
		if itemsProcessed == numItems {
			// log.Println("All items processed")
			cancel()
			break
			// return all pokemons
		}
	}
	return pokes, nil
}
