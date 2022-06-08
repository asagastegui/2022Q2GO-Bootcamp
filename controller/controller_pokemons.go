package controller

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
	"sync"
	"time"

	entities "gobootcamp.com/entity"
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

func GetConcurrentPokemon(w http.ResponseWriter, r *http.Request) {

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	numberOfWorkers := 1
	fileName := "./files/pokemons.csv"
	src := make(chan entities.Pokemon)
	dst := make(chan entities.Pokemon)

	start := time.Now()
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
		ReadCsvConcurrent(fileName, src, cancel)
		wg.Done()
	}()
	go func() {
		wg.Wait()

		close(dst)
	}()
	// drain the output

	itemsProcessed := 0
	var pokes []entities.Pokemon
	w.Header().Set("Content-Type", "application/json")
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
	// log.Println("Reach EOF or workers hit the limit")
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

func ReadCsvConcurrent(filename string, src chan entities.Pokemon, cancel context.CancelFunc) {
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
