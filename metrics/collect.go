package metrics

import (
	"sort"
	"strings"
	"time"
)

// SnapshotMetric is a snapshot of a single metric
type SnapshotMetric struct {
	Key   string
	Value int64
}

// Snapshot is a snapshot of the metrics data when it was collected
type Snapshot struct {
	Counters []SnapshotMetric
}

// FetchCounter finds the value for the specified counter
func (s Snapshot) FetchCounter(key string) int64 {
	for _, v := range s.Counters {
		if v.Key == key {
			return v.Value
		}
	}

	return 0
}

// FetchAllCounters finds all the counters that start with the specified value
func (s Snapshot) FetchAllCounters(match string) []SnapshotMetric {
	metrics := []SnapshotMetric{}

	for _, v := range s.Counters {
		if strings.HasPrefix(v.Key, match) {
			metrics = append(metrics, v)
		}
	}

	return metrics
}

// Collect takes a snapshot of the current metrics data and exports them so
// they can be analyzed
func (r *Registry) Collect(since time.Time) (Snapshot, error) {
	snap := Snapshot{
		Counters: []SnapshotMetric{},
	}

	for _, c := range r.counters {
		snap.Counters = append(
			snap.Counters,
			SnapshotMetric{
				Key:   c.key,
				Value: c.Collect(since),
			})
	}
	sort.Slice(snap.Counters, func(i, j int) bool {
		return snap.Counters[i].Value > snap.Counters[j].Value
	})

	return snap, nil
}
