package controller

import (
	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router) {
	router.Path("/pokemons").HandlerFunc(GetPokemons).Methods("GET")
	router.Path("/pokemons/{id:[0-9]+}").HandlerFunc(GetPokemon).Methods("GET")
	router.Path("/api/pokemons").HandlerFunc(GetAPIPokemon).Methods("GET")
}
