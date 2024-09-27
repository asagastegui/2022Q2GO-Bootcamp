package infrastructure

import (
	"gobootcamp.com/controller"

	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router) {
	router.Path("/pokemons").HandlerFunc(controller.GetPokemons).Methods("GET")
	router.Path("/pokemons/{id:[0-9]+}").HandlerFunc(controller.GetPokemon).Methods("GET")
	router.Path("/api/pokemons").HandlerFunc(controller.GetAPIPokemon).Methods("GET")
	router.Path("/pokemons/{type:odd|even}/{items:[0-9]+}/{items_per_worker:[0-9]+}").HandlerFunc(controller.GetConcurrentPokemon).Methods("GET")
}
