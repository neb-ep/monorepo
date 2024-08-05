package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

const (
	attrServiceName        = "service_name"
	attrServiceVersion     = "service_version"
	attrServiceEnvironment = "service_environment"
)

type Config struct {
	Name        string     `yaml:"name"`
	Version     string     `yaml:"version"`
	Environment string     `yaml:"environment"`
	Level       slog.Level `yaml:"level"`
}

func InitLogger(config Config) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     config.Level,
		AddSource: true,
	})
	slog.SetDefault(slog.New(handler).With(
		attrServiceName, config.Name,
		attrServiceVersion, config.Version,
		attrServiceEnvironment, config.Environment),
	)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return logging.UnaryServerInterceptor(
		logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
			switch lvl {
			case logging.LevelDebug:
				slog.Debug(msg, fields...)
			case logging.LevelInfo:
				slog.Info(msg, fields...)
			case logging.LevelWarn:
				slog.Warn(msg, fields...)
			case logging.LevelError:
				slog.Error(msg, fields...)
			default:
				panic(fmt.Sprintf("unknown level %v", lvl))
			}
		}),
		logging.WithFieldsFromContext(func(ctx context.Context) logging.Fields {
			if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
				return logging.Fields{"traceId", span.TraceID().String()}
			}
			return nil
		}),
	)
}
