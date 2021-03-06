package api

import (
	"net/http"
)

// HandleRequests is the nexus for registration of route handling as well as creating the listener
func HandleRequests() {
	http.HandleFunc("/portscan", scanHandler)

	http.ListenAndServe(":8080", nil)
}
