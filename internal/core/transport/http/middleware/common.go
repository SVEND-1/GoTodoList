package middleware

import (
	"TodoList/internal/core/domain"
	core_errors "TodoList/internal/core/errors"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/cookies"
	"TodoList/internal/core/transport/http/response"
	"TodoList/internal/features/auth/authService/jwt"
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIdHeader = "X-Request-Id"

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowedOrigin := map[string]struct{}{
				"http://localhost:5050": {},
			}
			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigin[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PATCH, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

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

type contextKey struct{}

var (
	userIDKey = contextKey{}
	roleKey   = contextKey{}
)

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

func RoleFromContext(ctx context.Context) (domain.Role, bool) {
	role, ok := ctx.Value(roleKey).(domain.Role)
	return role, ok
}

type tokenParser interface {
	ParseAccessToken(tokenString string) (*jwt.Claims, error)
}

func Auth(tokenParser tokenParser) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := response.NewHTTPResponseHandler(log, w)

			token, err := cookies.GetTokenFromCookie(r, cookies.AccessTokenCookie)
			if err != nil {
				responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "missing access token")
				return
			}

			claims, err := tokenParser.ParseAccessToken(token)
			if err != nil {
				responseHandler.ErrorResponse(core_errors.ErrUnauthorized, "invalid or expired access token")
				return
			}

			ctx = context.WithValue(ctx, userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, roleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowedRoles ...domain.Role) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.FromContext(r.Context())
			responseHandler := response.NewHTTPResponseHandler(log, w)

			role, ok := RoleFromContext(r.Context())
			if !ok {
				responseHandler.ErrorResponse(core_errors.ErrForbidden, "role not found in context: Auth middleware must run before RequireRole")
				return
			}

			for _, allowed := range allowedRoles {
				if role == allowed {
					next.ServeHTTP(w, r)
					return
				}
			}

			responseHandler.ErrorResponse(core_errors.ErrForbidden, "insufficient permissions for this action")
		})
	}
}
