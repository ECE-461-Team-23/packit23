package main

/*
Provide GET functionality for the following endpoints:
/package/{id}
/package/{id}/rate
/package/byName/{name}

Provide POST functionality for the following endpoints:
/packages
/package/byRegEx
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Masterminds/semver"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

func print_req_info(r *http.Request) {
	fmt.Print(r.Method, r.URL, r.Header, r.URL.Query(), r.URL.RawPath, r.Body)
}

func verifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header["X-Authorization"] != nil {
			token, err := jwt.Parse(request.Header["X-Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				_, ok := token.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					fmt.Print("Error validating JWT")
					return_400_packet(writer, request)
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				fmt.Print("Error validating JWT")
				return_400_packet(writer, request)
				return
			}
			if token.Valid {
				endpointHandler(writer, request)
			} else {
				fmt.Print("Error validating JWT")
				return_400_packet(writer, request)
				return
			}
		}
	})
}

/*
// Get the packages from the registry
Missing:
- Pagination
*/
func handle_packages(w http.ResponseWriter, r *http.Request) {
	print_req_info(r)
	db, err := connect()
	//db, err := connect_test_db()
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer db.Close()
	// parse query for offset (pagination). if empty, return the first page of results
	query := r.URL.Query()
	offset := query.Get("offset")
	if offset == "" {
		offset = "1"
	}

	// parse body for versions to find
	var response_arr []PackageQuery
	var body PackagesBody
	var packages_metadata []PackageMetadata
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Print("\nError reading body of request\n")
		return_404_packet(w, r)
		return
	}
	for _, q := range body.Items {
		response_arr = append(response_arr, PackageQuery{q.Name, q.Version})
	}

	// make version field a range to look for
	for _, response := range response_arr {
		// check if response version field contains '-' character without surrounding whitespace, if it doesn't add it
		char_idx := strings.Index(response.Version, "-")
		if strings.Contains(response.Version, "-") && response.Version[char_idx-1] != ' ' && response.Version[char_idx+1] != ' ' {
			response.Version = strings.Replace(response.Version, "-", " - ", -1)
		}
		c, err := semver.NewConstraint(response.Version)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return
		}

		// query all versions of a package if found in db
		metadataList, err := getMetadataFromName(db, response.Name)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return
		}
		// check which version is in range
		for _, md := range metadataList {
			v, err := semver.NewVersion(md.Version)
			if err != nil {
				fmt.Print(err)
				return_500_packet(w, r)
				return
			}
			if c.Check(v) {
				packages_metadata = append(packages_metadata, md)
			}
		}
	}

	if packages_metadata == nil {
		fmt.Print("packages response is empty")
	}
	json.NewEncoder(w).Encode(packages_metadata)
	return_200_packet(w, r)
}

// Return this package (ID) from google cloud bucket
func handle_package_id(w http.ResponseWriter, r *http.Request) {
	print_req_info(r)
	db, err := connect()
	//db, err := connect_test_db()
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer db.Close()

	var meta PackageMetadata
	var bucket_object_name string
	var ret_package Package
	var data PackageData
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		fmt.Print("Empty {id} in path")
		return_404_packet(w, r)
		return
	}

	rows, err := db.Query("SELECT id, name, version FROM packages WHERE id = ?;", id)
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&meta.ID, &meta.Name, &meta.Version)
		if err != nil {
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
	}

	res, err := db.Query("SELECT rating_pk FROM packages WHERE id = ?;", id)
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer res.Close()

	// bucket object name is the same as the rating pk
	for res.Next() {
		err = res.Scan(&bucket_object_name)
		if err != nil {
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
	}

	b64contents, err := getBucketObject(os.Getenv("BUCKET_NAME"), bucket_object_name)
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}

	data.Content = string(b64contents)
	data.URL = ""
	data.JSProgram = ""
	ret_package.Metadata = &meta
	ret_package.Data = &data

	json.NewEncoder(w).Encode(ret_package)
	return_200_packet(w, r)
}

// Return the package rating with ID (id from path)
func handle_package_rate(w http.ResponseWriter, r *http.Request) {
	print_req_info(r)
	db, err := connect()
	//db, err := connect_test_db()
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}

	var ratings PackageRating
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		fmt.Print("id cannot be empty")
		return_404_packet(w, r)
		return
	}

	res, err := db.Query("SELECT A.busFactor, A.correctness, A.rampUp, A.responsiveMaintainer, A.licenseScore, A.goodPinningPractice, A.pullRequest, A.netScore FROM ratings AS A INNER JOIN packages AS B ON A.id = B.rating_pk WHERE B.id = ?;", id)
	if err != nil {
		fmt.Print(err)
		return_400_packet(w, r)
		return
	}

	for res.Next() {
		err = res.Scan(&ratings.BusFactor, &ratings.Correctness, &ratings.RampUp, &ratings.ResponsiveMaintainer, &ratings.LicenseScore, &ratings.GoodPinningPractice, &ratings.PullRequest, &ratings.NetScore)
		if err != nil {
			fmt.Print(err)
			return_400_packet(w, r)
			return
		}
	}

	json.NewEncoder(w).Encode(ratings)
	return_200_packet(w, r)
}

// return the package history with package name from path (all versions)
// const mysqlFormat = "2006-01-02 15:04:05"
// const timeFormat = "2006-01-02T15:04:05Z"
func handle_package_byname(w http.ResponseWriter, r *http.Request) {
	print_req_info(r)
	db, err := connect()
	//db, err := connect_test_db()
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer db.Close()

	var ret []PackageHistoryEntry
	var metadataList []PackageMetadata
	var times []string
	vars := mux.Vars(r)
	name := vars["name"]
	if name == "" {
		fmt.Print("name cannot be empty")
		return_404_packet(w, r)
		return
	}

	// get registry entry from name
	rows, err := db.Query("SELECT id, name, version, upload_time FROM packages WHERE name = ?;", name)
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}

	defer rows.Close()
	// get all versions of named package
	for rows.Next() {
		var timevar string
		var md PackageMetadata
		if err := rows.Scan(&md.ID, &md.Name, &md.Version, &timevar); err != nil {
			// package with name not found
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
		if err != nil {
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
		t, err := time.Parse("2006-01-02T15:04:05Z", timevar)
		if err != nil {
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
		timevar = t.Format(time.RFC3339)
		metadataList = append(metadataList, md)
		times = append(times, timevar)
	}

	// iterate through versions of package and get rest of history
	for i, md := range metadataList {
		var history PackageHistoryEntry
		history.User = &User{Name: "test", IsAdmin: false}
		history.Date, err = time.Parse("2006-01-02T15:04:05Z", times[i])
		if err != nil {
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
		history.PackageMetadata = &md
		history.Action = "Uploaded"
		ret = append(ret, history)
	}

	json.NewEncoder(w).Encode(ret)
	return_200_packet(w, r)
}

// return a list of package metadata from package names that match the regex
func handle_package_byregex(w http.ResponseWriter, r *http.Request) {
	print_req_info(r)
	db, err := connect()
	//db, err := connect_test_db()
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer db.Close()

	//grab RegEx from body
	var body PackageRegExBody
	var retList []RegExReturn
	var listoflists [][]PackageMetadata
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		fmt.Print("\nError reading body of request\n")
		return_404_packet(w, r)
		return
	}

	// for list of names that match regex, get metadata and append to list of metadata
	rows, err := db.Query("SELECT id, name, version FROM packages WHERE name REGEXP ?;", body.RegEx)
	if err != nil {
		fmt.Print(err)
		return_500_packet(w, r)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var md PackageMetadata
		var cont bool = false
		if err := rows.Scan(&md.ID, &md.Name, &md.Version); err != nil {
			fmt.Print(err)
			return_500_packet(w, r)
			return
		}
		for _, md_list := range listoflists {
			// might need to check first iteration
			// all temps will have same name so only check one
			if md_list[0].Name == md.Name {
				cont = true
				break
			}
		}

		if !cont {
			mdl, err := getMetadataFromName(db, md.Name)
			if err != nil {
				fmt.Print("Error getting metadata from name")
				return_500_packet(w, r)
				return
			}
			listoflists = append(listoflists, mdl)
		}
	}

	for i, md_list := range listoflists {
		var sortedVersions []*semver.Version
		for _, md := range md_list {
			// create a semantic version for each version
			sv, err := semver.NewVersion(md.Version)
			if err != nil {
				fmt.Print("Error creating semantic version")
				return_413_packet(w, r)
				return
			}
			sortedVersions = append(sortedVersions, sv)
		}

		// check if versions are sorted
		sort.Sort(semver.Collection(sortedVersions))

		// exact version found
		if len(sortedVersions) == 1 {
			var ret RegExReturn
			ret.Version = md_list[0].Version
			ret.Name = md_list[i].Name
			retList = append(retList, ret)
		} else {
			// get begin and end of list
			begin := sortedVersions[0]
			end := sortedVersions[len(sortedVersions)-1]
			begin_split := strings.Split(begin.String(), ".")
			end_split := strings.Split(end.String(), ".")
			var ret RegExReturn

			// for bounded range (up to major version can change):
			// if first group IS different
			if begin_split[0] != end_split[0] {
				ret.Version = begin.String() + "-" + end.String()
				ret.Name = md_list[i].Name
			} else if begin_split[1] != end_split[1] {
				// for caret range (only major must match, up to minor version can change):
				ret.Version = "^" + begin.String()
				ret.Name = md_list[i].Name
			} else if begin_split[2] != end_split[2] {
				// for tilde range (major and minor must match, patch version can change):
				ret.Version = "~" + begin.String()
				ret.Name = md_list[i].Name
			}
			retList = append(retList, ret)
		}
	}
	json.NewEncoder(w).Encode(retList)
	return_200_packet(w, r)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/packages", verifyJWT(handle_packages))
	router.HandleFunc("/package/{id}", verifyJWT(handle_package_id))
	router.HandleFunc("/package/{id}/rate", verifyJWT(handle_package_rate))
	router.HandleFunc("/package/byName/{name}", verifyJWT(handle_package_byname))
	router.HandleFunc("/package/byRegEx", verifyJWT(handle_package_byregex))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	handleRequests()
}
