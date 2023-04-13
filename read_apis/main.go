package main

/*
Provide GET functionality for the following endpoints:
/packages
/package/(id)
/package/(id)/rate
/package/byName/(name)
*/

import (
	"log"
	"net/http"

	r "github.com/packit461/packit23/read_apis/src"
)

func main() {
	log.Printf("Server started")

	router := r.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
