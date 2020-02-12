package metrics

import (
	"sync"
	"time"
)

const comparisonFormat = "2006-01-02T15:04:05"

type valueHolder struct {
	when  time.Time
	value int64
}

// Counter allows storing metrics which will always be incremented by one
type Counter struct {
	key    string
	values []valueHolder
	mux    sync.Mutex
}

func newCounter(key string) *Counter {
	return &Counter{
		key: key,
	}
}

// Increment increments the counter by one
func (c *Counter) Increment(when time.Time) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if len(c.values) == 0 {
		c.values = append(c.values, valueHolder{when: when})
	}
	valueIndex := len(c.values) - 1

	if timesAreTooFarApart(c.values[valueIndex].when, when) {
		c.values = append(c.values, valueHolder{when: when})
		valueIndex = len(c.values) - 1
	}

	c.values[valueIndex].value = c.values[valueIndex].value + 1
}

// Collect takes a snapshot of the metric data and exports it to be analyzed
func (c *Counter) Collect(since time.Time) int64 {
	var value int64

	// Make a copy of the values so we can safely iterate through them while
	// acquiring a mutex for the least amount of time, and therefore avoid
	// blocking metrics aggregation for too long
	c.mux.Lock()
	values := make([]valueHolder, len(c.values))
	copy(values, c.values)
	c.mux.Unlock()

	for i := len(values) - 1; i >= 0; i-- {
		if values[i].when.Before(since) {
			break
		}

		value = value + values[i].value
	}

	return value
}

func timesAreTooFarApart(f, s time.Time) bool {
	return f.Format(comparisonFormat) != s.Format(comparisonFormat)
}
