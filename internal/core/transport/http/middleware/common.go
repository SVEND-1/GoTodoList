package middleware

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/response"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIdHeader = "X-Request-Id"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIdHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIdHeader, requestID)
			w.Header().Set(requestIdHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *logger.Log) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)
			l := log.With(
				zap.String("requestId", requestId),
				zap.String("url", r.URL.String()))

			ctx := logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHadler := response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					log.Error("panic recovered",
						zap.Any("panic", p),
						zap.String("stack", string(debug.Stack())),
					)
					responseHadler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			rw := response.NewResponseWriter(w)

			before := time.Now()
			log.Debug(">>> incoming HTTP request",
				zap.String("method", r.Method),
				zap.Time("time", before.UTC()))

			next.ServeHTTP(rw, r)

			log.Debug("<<< done HTTP request",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Since(before)),
			)

		})
	}
}
