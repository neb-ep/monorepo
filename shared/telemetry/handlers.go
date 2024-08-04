package telemetry

import (
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc/stats"
)

func NewServerHandler() stats.Handler {
	return otelgrpc.NewServerHandler()
}
