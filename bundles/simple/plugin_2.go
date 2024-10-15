package main

import (
	"net/http"
)

// AddHelloWorldHeader adds custom "Hello: World" header to the request
func AddHelloWorldHeader(_ http.ResponseWriter, r *http.Request) {
	r.Header.Add("Hello", "World")
}
