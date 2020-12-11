package main

import (
	"io/ioutil"
	"net/http"
)

func cacheRequestHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			cacheGET(w, r)
		case http.MethodPost:
			cachePOST(w, r)
		}
	})
}

func healthCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			healthGET(w, r)
		}
	})
}

// handles cache get request
func cacheGET(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cachePath):]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("get cache request without key"))
		logger.Info("get cache request without key.")
		return
	}
	value, exist := cache.Get(key)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("key not present in the cache"))
	}
	w.Write([]byte(value))
}

// handles cache post request
func cachePOST(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len(cachePath):]
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post cache request without key"))
		logger.Info("post cache request without key.")
		return
	}
	reqBody, _ := ioutil.ReadAll(r.Body)
	if len(reqBody) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("post cache request without value"))
		logger.Info("post cache request without value.")
		return
	}
	cache.Set(key, string(reqBody))
	w.WriteHeader(http.StatusCreated)
}

// handles healthcheck get request
func healthGET(w http.ResponseWriter, r *http.Request) {
	// simple check to see if endpoint is reachable
	w.Write([]byte("OK"))
}
