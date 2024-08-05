package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/neb-ep/monorepo/services/auth/internal/app"
	"github.com/neb-ep/monorepo/services/auth/internal/service"
	"github.com/neb-ep/monorepo/services/auth/internal/storage"
	"github.com/neb-ep/monorepo/shared/database"
	"github.com/neb-ep/monorepo/shared/logger"
	"github.com/neb-ep/monorepo/shared/telemetry"
	authv1 "github.com/neb-ep/shared/contracts/protos/auth/v1"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Config struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

func main() {
	// Common
	ctx := context.Background()

	// Config
	conf := Config{
		Name:        "auth-service",
		Version:     "v0.0.1",
		Environment: "devel",
	}

	// Logger
	logger.InitLogger(logger.Config{
		Name:        conf.Name,
		Version:     conf.Version,
		Environment: conf.Environment,
		Level:       slog.LevelDebug,
	})

	// Telemetry
	tp, err := telemetry.NewTraceProvider(ctx, telemetry.Service{
		Name:        conf.Name,
		Version:     conf.Version,
		Environment: conf.Environment,
	})
	if err != nil {
		slog.Error("failed telemetry.NewTraceProvider", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("failed trace provider shutdown", "error", err)
		}
	}()

	mp, err := telemetry.NewMeterProvider(ctx, telemetry.Service{
		Name:        conf.Name,
		Version:     conf.Version,
		Environment: conf.Environment,
	})
	if err != nil {
		slog.Error("failed telemetry.NewMeterProvider", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := mp.Shutdown(ctx); err != nil {
			slog.Error("failed meter provider shutdown", "error", err)
		}
	}()

	// gRPC service
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		slog.Error("failed to start listen", "error", err)
		os.Exit(1)
	}

	srv := grpc.NewServer(
		grpc.StatsHandler(telemetry.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				slog.Error("recovered from panic", "panic", p, "stack", debug.Stack())
				return status.Errorf(codes.Internal, "%s", p)
			})),
			logger.UnaryServerInterceptor(),
		),
	)

	slog.Debug("create new postgres connection pool")
	pool, err := database.NewPostgres(ctx, &database.Config{
		ConnectionURL: "postgres://devel:devel@localhost:5432/devel?sslmode=disable",
	}, slog.Default())
	if err != nil {
		slog.Error("failed to create new postgres connection pool", "error", err)
		os.Exit(1)
	}
	defer func() {
		slog.Info("closing connection pool")
		pool.Close()
	}()

	storage := storage.NewStorage(pool)
	service := service.NewService(storage)
	authv1.RegisterAuthServiceServer(srv, app.NewApi(service))
	reflection.Register(srv)

	ctx, cancelFunc := signal.NotifyContext(ctx,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGPIPE,
	)
	defer cancelFunc()

	slog.Info("service staring to listen at tcp", "addr", lis.Addr().String())

	go func() {
		if err := srv.Serve(lis); err != nil {
			slog.Error("failed to starting serve", "error", err)
		}
	}()

	<-ctx.Done()

	srv.GracefulStop()
	slog.Info("service stopped")
}
