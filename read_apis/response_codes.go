package main

import "net/http"

func return_200_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func return_400_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("400 - There is missing field(s) in the PackageQuery/AuthenticationToken or it is formed improperly, or the AuthenticationToken is invalid."))
}

func return_404_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - Not found"))
}

func return_413_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	w.Write([]byte("413 - Too many packages returned"))
}

func return_500_packet(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal error"))
}
