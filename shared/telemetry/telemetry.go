package telemetry

import (
	"fmt"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
)

var errRuntimeStart = "failed to start opentelemetry runtime with reduced min read mem stats interval: %w"

func StartWithMinimumReadMemStatsInterval(d time.Duration) error {
	if err := runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(time.Second),
	); err != nil {
		return fmt.Errorf(errRuntimeStart, err)
	}

	return nil
}
