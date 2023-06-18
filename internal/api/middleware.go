package api

import (
	"net/http"
	"time"

	"github.com/v1tbrah/promcli"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (a *API) registerMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeStart := time.Now()

		writer := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(writer, r)

		a.promCli.ObserveRequestDurationSeconds(r.URL.Path, writer.statusCode, time.Since(timeStart))

		a.promCli.IncRequestResultCount(r.URL.Path, promcli.LabelTotal)
		if writer.statusCode == http.StatusInternalServerError {
			a.promCli.IncRequestResultCount(r.URL.Path, promcli.LabelError)
		}
	})
}
