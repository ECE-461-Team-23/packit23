package src

/*
Provide GET functionality for the following endpoints:
/packages
/package/(id)
/package/(id)/rate
/package/byName/(name)
*/

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/packit461/packit23/package_rater/internal/logger"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/mysql/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func connect() {
	cleanup, err := mysql.RegisterDriver("cloudsql-mysql", cloudsqlconn.WithCredentialsFile("key.json"))
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}
	mysql.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instanceConnectionName, opts...)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		dbUser, dbPwd, dbName)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
}

func return_error_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal error"))
}

// Get the packages from the registry
func handle_packages(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user7:s$cret@tcp(127.0.0.1:3306)/testdb") // EMILE FIX PLS
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Received %s request", r.Method))
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}
	logger.Info(headers)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Info("\nError reading body of request\n")
		return_error_packet(w, r)
		return
	}
	logger.Info(fmt.Sprintf("Body:\n%s\n", body))
	// TODO: ADD PAGINATION STUFF
	res, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry;`)
	defer res.Close()

	if err != nil {
		log.Fatal(err)
	}

	// --------- DEBUGGING/EXPERIMENTAL CODE TO VIEW RETURN ---------
	for res.Next() {
		var pack PackageMetadata // I have no idea what to put here
		err := res.Scan(&pack.ID, &pack.VERSION, &pack.NAME)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%v\n", pack)
	}
	// --------------------------------------------------------------
}

// `BEGIN
//	SELECT * FROM Registry AS A
//	INNER JOIN Binaries AS B
//		ON A.BINARY_PK = B.ID
//	INNER JOIN Users AS C
//		ON A.USER_PK = C.ID
//	INNER JOIN Ratings AS D
//		ON A.RATING_PK = D.ID
//	END;`

// Return this package (ID)
func handle_packages_id(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user7:s$cret@tcp(127.0.0.1:3306)/testdb") // EMILE FIX PLS
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info(fmt.Sprintf("Received %s request", r.Method))
	headers := "Headers:\n"
	for key, value := range r.Header {
		headers += fmt.Sprintf("%s=%s\n", key, value)
	}
	logger.Info(headers)

	id := r.Header.Get("Id")
	if id == "" {
		logger.Info("\nNo Matching Value to key Id\n")
		return_error_packet(w, r)
		return
	}

	res, err := db.Query(`SELECT ID, NAME, VERSION FROM Registry;`)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	var meta PackageMetadata
	err = res.Scan(&meta.ID, &meta.Name, &meta.Version)
	res, err = db.Query(`BEGIN SELECT 
							B.BINARY_FILE, 
							A.URL 
							B.JS_PROGRAM
							FROM Registry AS A
							WHERE A.ID == id
							INNER JOIN Binaries AS B
								ON A.BINARY_PIK == B.ID
							END;`)
	defer res.Close()
	if err != nil {
		log.Fatal(err)
	}
	var packData PackageData
	// Need to append NULL for JSProgram
	err = res.Scan(&packData.Content, &packData.URL, &packData.JSProgram)
	totalPack := PackageModel{Metadata: &meta, Data: &packData}
	packJson, err := json.Marshal(totalPack)
	w.Write(packJson)
	w.WriteHeader(200)
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
