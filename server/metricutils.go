package main

import (
	"expvar"
	"runtime"
	"strconv"
	"strings"
	"time"

	metric "github.com/zserge/metric"
)

func fmtReqCounterPath(url string, method string, code int) string {
	// formatting base path from path="/api/v1/cache/..."
	basePath := strings.Split(url, "/")[1:4]
	basePath = append(basePath, method)
	basePath = append(basePath, strconv.Itoa(code))
	return strings.Join(basePath, ":")
}

func fmtReqLatencyPath(url string, method string, code int) string {
	// formatting base path from path="/api/v1/cache/..."
	basePath := strings.Split(url, "/")[1:4]
	basePath = append(basePath, method)
	basePath = append(basePath, strconv.Itoa(code))
	basePath = append(basePath, "latency")
	return strings.Join(basePath, ":")
}

func requestCounterInc(url string, method string, code int) {
	expvar.Get(fmtReqCounterPath(url, method, code)).(metric.Metric).Add(1)
}

func requestLatencyUpdate(url string, method string, code int, latency time.Duration) {
	expvar.Get(fmtReqLatencyPath(url, method, code)).(metric.Metric).Add(latency.Seconds())
}

func startSystemHealthMetrics() {
	for range time.Tick(100 * time.Millisecond) {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)
		expvar.Get("go:numgoroutine").(metric.Metric).Add(float64(runtime.NumGoroutine()))
		expvar.Get("go:numcgocall").(metric.Metric).Add(float64(runtime.NumCgoCall()))
		expvar.Get("go:alloc").(metric.Metric).Add(float64(m.Alloc) / 1000000)
		expvar.Get("go:alloctotal").(metric.Metric).Add(float64(m.TotalAlloc) / 1000000)
	}
}
