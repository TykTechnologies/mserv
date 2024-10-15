package main

import (
	"net/http"
)

// AddFooBarHeader2 adds custom "Foo: Bar" header to the request
func AddFooBarHeader2(_ http.ResponseWriter, r *http.Request) {
	r.Header.Add("Foo", "Bar")
}
