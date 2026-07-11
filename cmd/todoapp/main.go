package main

import (
	"TodoList/internal/core/config"
	"TodoList/internal/core/logger"
	"TodoList/internal/core/repository/pool/postgres/core_pgx"
	"TodoList/internal/core/transport/http/middleware"
	"TodoList/internal/core/transport/http/server"
	"TodoList/internal/features/tasks/taskRepository"
	"TodoList/internal/features/tasks/taskService"
	"TodoList/internal/features/tasks/taskTransport/taskApi"
	"TodoList/internal/features/users/userRepository"
	"TodoList/internal/features/users/userService"
	"TodoList/internal/features/users/userTransport/userApi"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	config := config.NewConfigMust()
	time.Local = config.TimeZone

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log, err := logger.NewLog(logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to initialize logger,ex=", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Debug("Initialing postgres connection pool")
	pool, err := core_pgx.NewPool(ctx, core_pgx.NewConfigMust())
	if err != nil {
		log.Fatal("failed to init postgres connection poll", zap.Error(err))
	}
	defer pool.Close()

	log.Debug("Initialing feature", zap.String("feature", "user"))
	userPostgresRepository := userRepository.NewUserRepository(pool)
	userService := userService.NewUserService(userPostgresRepository)
	userTransportHttp := userApi.NewUserController(userService)

	log.Debug("Initializing feature", zap.String("feature", "task"))
	taskPostgresRepository := taskRepository.NewTaskRepository(pool)
	taskService := taskService.NewTaskService(taskPostgresRepository)
	taskTransportHttp := taskApi.NewTaskController(taskService)

	log.Debug("Initializing http server")
	httpServer := server.NewHTTPServer(
		server.NewConfigMust(), log,
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.Trace(),
		middleware.Panic(),
	)

	apiVersionRouter := server.NewAPIVersionRouter(server.ApiVersion1)

	apiVersionRouter.RegisterRouters(userTransportHttp.Routers()...)
	apiVersionRouter.RegisterRouters(taskTransportHttp.Routers()...)

	httpServer.RegisterAPIRouters(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		log.Error("HTTP server run, err=", zap.Error(err))
	}
}
