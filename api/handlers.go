package api

import "net/http"

// RegsiterHandlers registers handlers with the http package
func RegisterHandlers() {
	http.HandleFunc("/", DiscoverPortsWithNMap)
}

func ()