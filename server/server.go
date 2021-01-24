package main

import (
	"expvar"
	"net/http"
	"strconv"

	memcache "github.com/abhi195/go-mem-cache"
	log "github.com/sirupsen/logrus"
	metric "github.com/zserge/metric"
)

const (
	// base HTTP paths.
	apiVersion  = "v1"
	apiBasePath = "/api/" + apiVersion + "/"

	// path to cache.
	cachePath = apiBasePath + "cache/"

	// metrics path
	metricPath = "/metrics"
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

	// publish metrics
	initMetrics()

	// cache request handler
	h := cacheRequestHandler()
	// httpsnooping wrapper handler
	wh := loggingMiddlewareHandler(h)
	// handling api paths
	http.Handle(cachePath, wh)
	http.Handle(metricPath, metric.Handler(metric.Exposed))

	logger.Infof("Starting server on :%d", port)

	strPort := ":" + strconv.Itoa(port)
	logger.Fatal("ListenAndServe: ", http.ListenAndServe(strPort, nil))
}

func initMetrics() {

	// httpserver api counters
	expvar.Publish(fmtReqCounterPath(cachePath, http.MethodGet, http.StatusOK), metric.NewCounter("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqCounterPath(cachePath, http.MethodGet, http.StatusBadRequest), metric.NewCounter("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqCounterPath(cachePath, http.MethodGet, http.StatusNotFound), metric.NewCounter("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqCounterPath(cachePath, http.MethodPost, http.StatusCreated), metric.NewCounter("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqCounterPath(cachePath, http.MethodPost, http.StatusBadRequest), metric.NewCounter("120s1s", "15m10s", "1h1m"))

	// httpserver latency histograms
	expvar.Publish(fmtReqLatencyPath(cachePath, http.MethodGet, http.StatusOK), metric.NewHistogram("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqLatencyPath(cachePath, http.MethodGet, http.StatusBadRequest), metric.NewHistogram("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqLatencyPath(cachePath, http.MethodGet, http.StatusNotFound), metric.NewHistogram("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqLatencyPath(cachePath, http.MethodPost, http.StatusCreated), metric.NewHistogram("120s1s", "15m10s", "1h1m"))
	expvar.Publish(fmtReqLatencyPath(cachePath, http.MethodPost, http.StatusBadRequest), metric.NewHistogram("120s1s", "15m10s", "1h1m"))

	// some Go internal metrics
	expvar.Publish("go:numgoroutine", metric.NewGauge("2m1s", "15m30s", "1h1m"))
	expvar.Publish("go:numcgocall", metric.NewGauge("2m1s", "15m30s", "1h1m"))
	expvar.Publish("go:alloc", metric.NewGauge("2m1s", "15m30s", "1h1m"))
	expvar.Publish("go:alloctotal", metric.NewGauge("2m1s", "15m30s", "1h1m"))

	// start Go internal metrics reporting
	go startSystemHealthMetrics()
}
