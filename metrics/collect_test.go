package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollect(t *testing.T) {
	registry := newRegistry()
	counter := registry.Counter("metric.test")
	secondCounter := registry.Counter("metric.foobar")

	when := time.Now()
	for i := 1; i <= 5; i++ {
		counter.Increment(when)
	}
	for i := 1; i <= 10; i++ {
		secondCounter.Increment(when)
	}

	coll, err := registry.Collect(time.Now().Add(time.Second * -2))
	assert.NoError(t, err)
	assert.Equal(t, Snapshot{
		Counters: []SnapshotMetric{
			SnapshotMetric{Key: "metric.foobar", Value: 10},
			SnapshotMetric{Key: "metric.test", Value: 5},
		},
	}, coll)
}
