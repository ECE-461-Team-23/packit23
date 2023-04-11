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
	_ "github.com/go-sql-driver/mysql"
)

type RegistryPackage struct {
	ID          int
	NAME        string
	RATING_PK   int
	AUTHOR_PK   int
	URL         string
	BINARY_PK   int
	VERSION     string
	UPLOADED    int // datetime?
	IS_EXTERNAL bool
}

func connect() {
	cleanup, err := mysql.RegisterDriver("cloudsql-mysql", cloudsqlconn.WithCredentialsFile("key.json"))
	if err != nil {
		log.Fatal(err)
	}
	// call cleanup when you're done with the database connection
	defer cleanup()

	db, err := sql.Open(
		"cloudsql-mysql",
		"myuser:mypass@cloudsql-mysql(project:region:instance)/mydb",
	)

	if db != nil {
		fmt.Print("Db not nil!")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// Get the packages from the registry
func handle_packages(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user7:s$cret@tcp(127.0.0.1:3306)/testdb") // EMILE FIX PLS
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	res, err := db.Query(`BEGIN
		SELECT * FROM Registry AS A 
		INNER JOIN Binaries AS B 
			ON A.BINARY_PK = B.ID 
		INNER JOIN Users AS C 
			ON A.USER_PK = C.ID 
		INNER JOIN Ratings AS D 
			ON A.RATING_PK = D.ID 
		END;`)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
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
