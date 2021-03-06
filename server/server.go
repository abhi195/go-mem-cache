package main

import (
	"net/http"
	"strconv"

	memcache "github.com/abhi195/go-mem-cache"
	log "github.com/sirupsen/logrus"
)

const (
	// base HTTP paths.
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion + "/"

	// health check path
	healthCheckPath = "/health"

	// path to cache.
	cachePath = apiBasePath + "cache/"
)

var (
	port   int
	cache  *memcache.MemCache
	logger *log.Logger
)

func init() {
	cache = memcache.New()
	port = 8080
	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{})
}

func main() {

	// handlers
	ch := cacheRequestHandler()
	hh := healthCheckHandler()

	// httpsnooping wrapper handlers
	cwh := loggingMiddlewareHandler(ch)
	hwh := loggingMiddlewareHandler(hh)

	// handling api paths
	http.Handle(cachePath, cwh)
	http.Handle(healthCheckPath, hwh)

	logger.Infof("Starting server on :%d", port)

	strPort := ":" + strconv.Itoa(port)
	logger.Fatal("ListenAndServe: ", http.ListenAndServe(strPort, nil))
}
