package output

import (
	"context"
	"testing"
	"time"

	"github.com/dmathieu/logtee/alerts"
	"github.com/dmathieu/logtee/metrics"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	metrics.RegistryProvider().Counter("requests.total").Increment(time.Now())
	refreshInterval = time.Millisecond
	go func() {
		err := Run(ctx, alerts.NewRegistry())
		assert.NoError(t, err)
	}()

	time.Sleep(2 * time.Millisecond)
	cancel()
}
