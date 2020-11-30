package main

import (
	"log"
	"net/http"
	"strconv"
)

const (
	// base HTTP paths.
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion + "/"

	// path to cache.
	cachePath = apiBasePath + "cache/"
)

var (
	port  int
	cache map[string]string
)

func init() {
	cache = make(map[string]string)
	port = 8080
}

func main() {

	// api paths
	http.Handle(cachePath, cacheRequestHandler())

	log.Printf("starting server on :%d", port)

	strPort := ":" + strconv.Itoa(port)
	log.Fatal("ListenAndServe: ", http.ListenAndServe(strPort, nil))
}
