package api

import (
	"net/http"

	"github.com/abelgoodwin1988/nmap-be/internal/portscan"
)

// HandleRequests is the nexus for registration of route handling as well as creating the listener
func HandleRequests() {
	http.HandleFunc("/portscan", portscan.Get)

	http.ListenAndServe(":8080", nil)
}
