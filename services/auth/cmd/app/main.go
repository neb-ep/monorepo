package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/neb-ep/monorepo/shared/logger"
	"github.com/neb-ep/monorepo/shared/telemetry"
	authv1pb "github.com/neb-ep/shared/contracts/protos/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// -------------------- CONFIGS (BEGIN) --------------------
type Config struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

// --------------------  CONFIGS (END)  --------------------

// --------------------   API (BEGIN)   --------------------
type api struct {
	authv1pb.UnimplementedAuthServiceServer
}

// SignIn implements auth.AuthServiceServer.
func (a *api) SignIn(context.Context, *authv1pb.SignInRequest) (*authv1pb.SignInResponse, error) {
	panic("unimplemented")
}

// SignUp implements auth.AuthServiceServer.
func (a *api) SignUp(context.Context, *authv1pb.SignUpRequest) (*authv1pb.SignUpResponse, error) {
	panic("unimplemented")
}

// --------------------    API (END)   --------------------

type AppConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

func main() {
	// Common
	ctx := context.Background()

	// Config
	conf := AppConfig{
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
	}

	srv := grpc.NewServer(
		grpc.StatsHandler(telemetry.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				slog.Error("recovered from panic", "panic", p, "stack", debug.Stack())
				return status.Errorf(codes.Internal, "%s", p)
			})),
		),
	)

	authv1pb.RegisterAuthServiceServer(srv, &api{})

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
