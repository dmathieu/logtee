package metrics

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegistryCount(t *testing.T) {
	registry := newRegistry()

	when := time.Now()
	registry.Counter("metric.test").Increment(when)

	counters := registry.counters
	assert.Equal(t, 1, len(counters))
	assert.Equal(t, "metric.test", counters["metric.test"].key)
	assert.Equal(t, 1, len(counters["metric.test"].values))
	assert.Equal(t, valueHolder{when: when, value: 1}, counters["metric.test"].values[0])

	registry.Counter("metric.test").Increment(when)
	registry.Counter("metric.foobar").Increment(when)

	assert.Equal(t, 2, len(counters))

	assert.Equal(t, 1, len(counters["metric.test"].values))
	assert.Equal(t, valueHolder{when: when, value: 2}, counters["metric.test"].values[0])

	assert.Equal(t, 1, len(counters["metric.foobar"].values))
	assert.Equal(t, valueHolder{when: when, value: 1}, counters["metric.foobar"].values[0])
}

func TestRegistryCountDifferentTimes(t *testing.T) {
	for _, test := range []struct {
		Name       string
		FirstTime  time.Time
		SecondTime time.Time
	}{
		{
			Name:       "when times have a large difference",
			FirstTime:  time.Now(),
			SecondTime: time.Now().Add(time.Hour),
		},
		{
			Name:       "when times have a small difference",
			FirstTime:  time.Now(),
			SecondTime: time.Now().Add(time.Second),
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			registry := newRegistry()

			registry.Counter("metric.test").Increment(test.FirstTime)
			registry.Counter("metric.test").Increment(test.SecondTime)
			registry.Counter("metric.test").Increment(test.SecondTime)

			counters := registry.counters
			assert.Equal(t, 1, len(counters))
			assert.Equal(t, "metric.test", counters["metric.test"].key)
			assert.Equal(t, 2, len(counters["metric.test"].values))
			assert.Equal(t, valueHolder{when: test.FirstTime, value: 1}, counters["metric.test"].values[0])
			assert.Equal(t, valueHolder{when: test.SecondTime, value: 2}, counters["metric.test"].values[1])
		})
	}
}

func TestRegistryCountSimilarTimes(t *testing.T) {
	registry := newRegistry()

	firstTime := time.Now()
	secondTime := time.Now().Add(time.Millisecond)

	registry.Counter("metric.test").Increment(firstTime)
	registry.Counter("metric.test").Increment(secondTime)
	registry.Counter("metric.test").Increment(secondTime)

	counters := registry.counters
	assert.Equal(t, 1, len(counters))
	assert.Equal(t, "metric.test", counters["metric.test"].key)
	assert.Equal(t, 1, len(counters["metric.test"].values))
	assert.Equal(t, valueHolder{when: firstTime, value: 3}, counters["metric.test"].values[0])
}

func TestRegistryCountCollectNoData(t *testing.T) {
	registry := newRegistry()
	counter := registry.Counter("metric.test")

	assert.Equal(t, int64(0), counter.Collect(time.Now()))
}

func TestRegistryCountCollect(t *testing.T) {
	registry := newRegistry()
	counter := registry.Counter("metric.test")

	when := time.Now()
	for i := 1; i <= 5; i++ {
		counter.Increment(when)
	}

	assert.Equal(t, int64(5), counter.Collect(time.Now().Add(time.Second*-2)))
}

func TestRegistryCountCollectOldData(t *testing.T) {
	registry := newRegistry()
	counter := registry.Counter("metric.test")

	now := time.Now()
	for i := 1; i <= 5; i++ {
		counter.Increment(now.Add(time.Second * -3))
	}

	for i := 1; i <= 5; i++ {
		counter.Increment(now)
	}

	assert.Equal(t, int64(5), counter.Collect(time.Now().Add(time.Second*-2)))
}
