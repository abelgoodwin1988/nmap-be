package main

import (
	"net/http"
)

func portsOpen(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("yhelothar"))
}

func handleRequests() {
	http.HandleFunc("/", portsOpen)

	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequests()
}
