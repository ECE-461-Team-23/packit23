package main

/*
Provide GET functionality for the following endpoints:
/packages
/package/(id)
/package/(id)/rate
/package/byName/(name)
*/

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/mysql/mysql"
)

func connect() {
	cleanup, err := mysql.RegisterDriver("cloudsql-mysql", cloudsqlconn.WithCredentialsFile("key.json"))
	if err != nil {
		// ... handle error
	}
	// call cleanup when you're done with the database connection
	defer cleanup()

	db, err := sql.Open(
		"cloudsql-mysql",
		"myuser:mypass@cloudsql-mysql(project:region:instance)/mydb",
	)
	// ... etc
}

// Get the packages from the registry
func handle_packages(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Return this package (ID)
func handle_packages_id(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Return the rating. Only use this if each metric was computed successfully.
func handle_packages_rate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// Return the history of this package (all versions).
func handle_packages_byname(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/packages", handle_packages)
	http.HandleFunc("/packages/id", handle_packages_id)
	http.HandleFunc("/packages/id/rate", handle_packages_rate)
	http.HandleFunc("/packages/byName/name", handle_packages_byname)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
