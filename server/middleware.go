package main

import (
	"net/http"

	httpsnoop "github.com/felixge/httpsnoop"
)

func loggingMiddlewareHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(h, w, r)
		requestCounterInc(r.URL.Path, r.Method, m.Code)
		requestLatencyUpdate(r.URL.Path, r.Method, m.Code, m.Duration)
		logger.Infof(
			"%s %s (code=%d dt=%s written=%d)",
			r.Method,
			r.URL,
			m.Code,
			m.Duration,
			m.Written,
		)
	})
}
