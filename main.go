package main

import "net/http"

func main() {
	api.RegisterHandlers()

	http.ListenAndServe(":8080", nil)
}
