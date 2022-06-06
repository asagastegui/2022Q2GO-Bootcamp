package entities

type Pokemon struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PokeInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokeAPIResp struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []PokeInfo `json:"results"`
}
