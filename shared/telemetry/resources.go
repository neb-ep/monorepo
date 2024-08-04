package telemetry

import (
	"context"
	"sync"

	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var (
	resource     *sdkresource.Resource
	resourceOnce sync.Once
)

func initResource(config Service) *sdkresource.Resource {
	resourceOnce.Do(func() {
		extra, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithFromEnv(),
			sdkresource.WithTelemetrySDK(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
			sdkresource.WithAttributes(
				semconv.ServiceName(config.Name),
				semconv.ServiceVersion(config.Version),
				semconv.DeploymentEnvironment(config.Environment),
			),
		)

		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extra,
		)
	})

	return resource
}
