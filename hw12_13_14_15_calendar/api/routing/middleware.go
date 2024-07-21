package routing

import (
	"net/http"
	"time"

	lg "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

// Перехватчик вызовов WriteHeader, для записи статуса ответа.
func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

// Перехватчик вызовов Write, для записи статуса ответа.
func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.status == 0 {
		lrw.status = http.StatusOK
	}
	return lrw.ResponseWriter.Write(b)
}

func LogRequest(log lg.Logger, handlerFunc http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w}

		handlerFunc.ServeHTTP(lrw, r)

		log.Info(
			"Processed request",
			"client_ip", r.RemoteAddr,
			"time", start.Format("02/Jan/2006:15:04:05 -0700"),
			"method", r.Method,
			"path", r.URL.Path,
			"version", r.Proto,
			"status", lrw.status,
			"latency", time.Since(start),
			"user_agent", r.UserAgent(),
		)
	})
}
