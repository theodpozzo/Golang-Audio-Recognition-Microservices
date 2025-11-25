package main

import (
	"log"
	"net/http"
	"search/search"
)

func main() {
	log.Fatal(http.ListenAndServe(":3001", search.Router()))
}
