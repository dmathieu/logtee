package alerts

import (
	"sync"
	"time"

	"github.com/dmathieu/logtee/metrics"
)

// Alert defines the condition for each alert and an explanation for it
type Alert func(metrics.Snapshot, time.Time) (string, bool)

// TriggeredAlert is an alert which has been checked positively
type TriggeredAlert struct {
	Message     string
	TriggeredAt time.Time
	alertID     int
}

// Registry allows storing and retrieving all alerts
type Registry struct {
	alerts          []Alert
	triggeredAlerts []TriggeredAlert
	mux             sync.Mutex
}

// NewRegistry creates a new instances of a registry
func NewRegistry() *Registry {
	return &Registry{}
}

// Register registers an alert to be checked later on
func (r *Registry) Register(a Alert) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.alerts = append(r.alerts, a)
}

// CheckAlerts checks the condition of each alert with the current snapshot and
// returns all the ones that match
func (r *Registry) CheckAlerts(since time.Time) ([]TriggeredAlert, error) {
	alerts := []TriggeredAlert{}
	snap, err := metrics.RegistryProvider().Collect(since)
	if err != nil {
		return nil, err
	}

	for i, a := range r.alerts {
		ta := TriggeredAlert{TriggeredAt: time.Now(), alertID: i}
		for _, t := range r.triggeredAlerts {
			if t.alertID == i {
				ta = t
				break
			}
		}

		m, ok := a(snap, ta.TriggeredAt)
		if ok {
			ta.Message = m
			alerts = append(alerts, ta)
		}
	}

	r.triggeredAlerts = alerts
	return alerts, nil
}
