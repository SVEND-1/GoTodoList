package server

import (
	"TodoList/internal/core/logger"
	"TodoList/internal/core/transport/http/middleware"
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *logger.Log

	middleware []middleware.Middleware
}

func NewHTTPServer(config Config, log *logger.Log, middleware ...middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(), config: config,
		log: log, middleware: middleware,
	}
}

func (s *HTTPServer) RegisterAPIRouters(routers ...*APIVersionRouter) {
	for _, router := range routers {
		prefix := "/api/" + string(router.apiVersion)
		s.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (s *HTTPServer) Run(ctx context.Context) error {
	mux := middleware.ChainMiddleware(s.mux, s.middleware...)
	server := &http.Server{
		Addr:    s.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		s.log.Warn("Starting HTTP server", zap.String("addr", s.config.Addr))

		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("server error: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("Shutting down HTTP server", zap.Error(ctx.Err()))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Warn("HTTP server stopped")
	}
	return nil
}
