package main

import (
	"net/http"
)

func main() {
	api := &api{
		store: NewStore(),
	}

	err := http.ListenAndServe(":8080", http.HandlerFunc(api.Serve))
	if err != nil {
		panic(err)
	}
}
