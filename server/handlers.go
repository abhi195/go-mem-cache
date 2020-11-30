package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func cacheRequestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCacheHandler(w, r)
		case http.MethodPost:
			postCacheHandler(w, r)
		}
	})
}

// handles get requests.
func getCacheHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cachePath):]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("get cache request without key"))
		log.Print("get cache request without key.")
		return
	}
	value, exist := cache[key]
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not present in the cache"))
	}
	w.Write([]byte(value))
}

// handles post requests.
func postCacheHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cachePath):]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post cache request without key"))
		log.Print("post cache request without key.")
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	if len(reqBody) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post cache request without value"))
		log.Print("post cache request without value.")
		return
	}
	cache[key] = string(reqBody)
	w.WriteHeader(http.StatusCreated)
}
