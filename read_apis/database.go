package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/storage"
	"github.com/Masterminds/semver"
	sql_driver "github.com/go-sql-driver/mysql"
)

func connect_test_db() (*sql.DB, error) {
	db, err := sql.Open(
		"mysql",
		"db_user:oldpassword!!!@tcp(127.0.0.1:3306)/test_db",
	)
	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return db, nil
}

func connect() (*sql.DB, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Fatal Error in connect_connector.go: %s environment variable not set.", k)
		}
		return v
	}

	var (
		db_user                  = mustGetenv("DB_USER")
		db_pass                  = mustGetenv("DB_PASSWORD")
		db_name                  = mustGetenv("DB_NAME")
		instance_connection_name = mustGetenv("INSTANCE_CONNECTION_NAME")
		usePrivate               = os.Getenv("PRIVATE_IP")
	)

	d, err := cloudsqlconn.NewDialer(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %v", err)
	}
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}

	sql_driver.RegisterDialContext("cloudsqlconn",
		func(ctx context.Context, addr string) (net.Conn, error) {
			return d.Dial(ctx, instance_connection_name, opts...)
		})

	dbURI := fmt.Sprintf("%s:%s@cloudsqlconn(localhost:3306)/%s?parseTime=true",
		db_user, db_pass, db_name)

	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}
	return dbPool, nil
}

func getBucketObject(bucketName string, objectName string) ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	rc, err := client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", objectName, err)
	}
	defer rc.Close()
	bucketObject, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}
	return bucketObject, nil
}

// return all versions of package name in db
func getMetadataFromName(db *sql.DB, name string) ([]PackageMetadata, error) {
	var metadataList []PackageMetadata
	rows, err := db.Query("SELECT id, name, version FROM packages WHERE name = ?;", name)
	if err != nil {
		log.Print(err)
	}
	defer rows.Close()
	for rows.Next() {
		var md PackageMetadata
		if err := rows.Scan(&md.ID, &md.Name, &md.Version); err != nil {
			return nil, fmt.Errorf("version of package not found. rows.Scan: %v", err)
		}
		metadataList = append(metadataList, md)
	}
	return metadataList, nil
}

func get_paginated_packages_metadata(w http.ResponseWriter, r *http.Request, db *sql.DB, offset string, response_arr []PackageQuery) []PackageMetadata {
	var amount_per_page int = 10
	var max_returned_packages int = 50
	var packages_metadata []PackageMetadata
	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return nil
	}

	for i := (offset_int - 1) * amount_per_page; i < offset_int*amount_per_page; i++ {
		if i >= len(response_arr) {
			break
		}
		if i >= max_returned_packages {
			fmt.Print("Too many packages to query")
			return_413_packet(w, r)
			return nil
		}
		response := response_arr[i]
		// check if response version field contains '-' character without surrounding whitespace, if it doesn't add it
		char_idx := strings.Index(response.Version, "-")
		if strings.Contains(response.Version, "-") && response.Version[char_idx-1] != ' ' && response.Version[char_idx+1] != ' ' {
			response.Version = strings.Replace(response.Version, "-", " - ", -1)
		}
		c, err := semver.NewConstraint(response.Version)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return nil
		}

		// query all versions of a package if found in db
		metadataList, err := getMetadataFromName(db, response.Name)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return nil
		}
		// check which version is in range
		for _, md := range metadataList {
			v, err := semver.NewVersion(md.Version)
			if err != nil {
				fmt.Print(err)
				return_500_packet(w, r)
				return nil
			}
			if c.Check(v) {
				packages_metadata = append(packages_metadata, md)
			}
		}
	}
	return packages_metadata
}

func get_all_packages_metadata(w http.ResponseWriter, r *http.Request, db *sql.DB, response_arr []PackageQuery) []PackageMetadata {
	var max_returned_packages int = 50
	var packages_metadata []PackageMetadata
	for i, response := range response_arr {
		if i >= max_returned_packages {
			fmt.Print("Too many packages to query")
			return_413_packet(w, r)
			return nil
		}
		// check if response version field contains '-' character without surrounding whitespace, if it doesn't add it
		char_idx := strings.Index(response.Version, "-")
		if strings.Contains(response.Version, "-") && response.Version[char_idx-1] != ' ' && response.Version[char_idx+1] != ' ' {
			response.Version = strings.Replace(response.Version, "-", " - ", -1)
		}
		c, err := semver.NewConstraint(response.Version)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return nil
		}

		// query all versions of a package if found in db
		metadataList, err := getMetadataFromName(db, response.Name)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return nil
		}
		// check which version is in range
		for _, md := range metadataList {
			v, err := semver.NewVersion(md.Version)
			if err != nil {
				fmt.Print(err)
				return_500_packet(w, r)
				return nil
			}
			if c.Check(v) {
				packages_metadata = append(packages_metadata, md)
			}
		}
	}
	return packages_metadata
}
