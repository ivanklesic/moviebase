package middleware

import (
	"io"
	stdlog "log"
	"net/http"
	"runtime/debug"
	"time"

	log "github.com/go-kit/log"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

type LoggerWrapper struct {
	Logger log.Logger
}

func (logWrap *LoggerWrapper) Log(keyvals ...interface{}) error {
	return logWrap.Logger.Log(keyvals...)
}

func LoggerInit(writeTo io.Writer) *LoggerWrapper {
	logger := log.NewLogfmtLogger(log.NewSyncWriter(writeTo))
	stdlog.SetOutput(log.NewStdlibAdapter(logger))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "loc", log.DefaultCaller)
	return &LoggerWrapper{Logger: logger}
}

func responseWriterWrapper(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

func LoggingMiddleware(logger *LoggerWrapper) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Log(
						"err", err,
						"trace", debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := responseWriterWrapper(w)
			next.ServeHTTP(wrapped, r)
			logger.Log(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}
		return http.HandlerFunc(fn)
	}
}
