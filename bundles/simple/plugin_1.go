package main

import (
	"net/http"
)

// AddFooBarHeader adds custom "Foo: Bar" header to the request
func AddFooBarHeader(_ http.ResponseWriter, r *http.Request) {
	r.Header.Add("Foo", "Bar")
}

func main() {}
