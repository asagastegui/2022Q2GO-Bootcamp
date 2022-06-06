package main

import (
	"fmt"
	"log"
	"net/http"

	"gobootcamp.com/controller"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	controller.InitRouter(router)
	fmt.Println("Server at 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
