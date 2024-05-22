// middleware/middleware.go
package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"golang.org/x/exp/slog"
)

// LoggingMiddleware logs the details of each request
func LoggingMiddleware(logger *slog.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			logger.Info("request completed",
				slog.String("method", r.Method),
				slog.String("uri", r.RequestURI),
				slog.Duration("duration", time.Since(start)))
		})
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// TracingMiddleware adds OpenTracing spans to each request
func TracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span, ctx := opentracing.StartSpanFromContext(r.Context(), "HTTP "+r.Method+" "+r.URL.Path)
		defer span.Finish()

		ext.HTTPMethod.Set(span, r.Method)
		ext.HTTPUrl.Set(span, r.URL.String())

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)

		statusCode, err := strconv.Atoi(w.Header().Get("Status-Code"))
		if err != nil {
			statusCode = 0
		}
		ext.HTTPStatusCode.Set(span, uint16(statusCode))
	})
}
