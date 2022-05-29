package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	controllers "gobootcamp.com/controller"
)

func main() {
	router := mux.NewRouter()

	router.Path("/pokemons").HandlerFunc(controllers.GetPokemons).Methods("GET")
	router.Path("/pokemons/{id:[0-9]+}").HandlerFunc(controllers.GetPokemon).Methods("GET")
	router.Path("/api/pokemons").HandlerFunc(controllers.GetAPIPokemon).Methods("GET")

	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
