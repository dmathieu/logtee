package main

import (
	"testing"
	"time"

	"github.com/dmathieu/logtee/alerts"
	"github.com/dmathieu/logtee/metrics"
	"github.com/stretchr/testify/assert"
)

func TestAlertRequestsTotal(t *testing.T) {
	alertTotalRequests = 100

	registry := alerts.NewRegistry()
	registerAlerts(registry)
	counter := metrics.RegistryProvider().Counter("requests.total")

	now := time.Now()
	since := now.Add(time.Minute * -2)
	al, err := registry.CheckAlerts(since)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(al))

	for i := int64(1); i <= alertTotalRequests; i++ {
		counter.Increment(now)
	}

	al, err = registry.CheckAlerts(since)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(al))

	future := now.Add(4 * time.Minute)
	al, err = registry.CheckAlerts(future)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(al))
}
