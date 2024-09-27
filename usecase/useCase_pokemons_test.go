package usecases

import (
	"testing"

	entities "gobootcamp.com/entity"
)

var pokes []entities.Pokemon = []entities.Pokemon{
	{
		ID:   1,
		Name: "bulbasaur",
	},
	{
		ID:   2,
		Name: "ivysaur",
	},
	{
		ID:   3,
		Name: "venusaur",
	},
	{
		ID:   4,
		Name: "charmander",
	},
	{
		ID:   5,
		Name: "charmeleon",
	},
	{
		ID:   6,
		Name: "charizard",
	},
	{
		ID:   7,
		Name: "squirtle",
	},
	{
		ID:   8,
		Name: "wartortle",
	},
	{
		ID:   9,
		Name: "blastoise",
	},
}

//TestFindPokeSuccess Array should contain the ID to find
func TestFindPokeSuccess(t *testing.T) {
	// Creating test data
	wants := 5
	poke, _ := FindPoke(&pokes, wants)
	if poke.ID != wants {
		t.Errorf("found: %d, want %d", poke.ID, wants)
	}
}

// TestFindPokeNotFound Array does not contain the ID to find
func TestFindPokeNotFound(t *testing.T) {
	idToFound := 151
	wants := "the pokemon you are looking for doesnt exists"

	_, err := FindPoke(&pokes, idToFound)
	errFound := err.Error()
	if errFound != wants {
		t.Errorf("found: %s, wants: %s", errFound, wants)
	}
}
