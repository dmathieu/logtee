package alerts

import (
	"testing"
	"time"

	"github.com/dmathieu/logtee/metrics"
	"github.com/stretchr/testify/assert"
)

func TestRegisterAlert(t *testing.T) {
	p := NewRegistry()
	assert.Equal(t, 0, len(p.alerts))

	p.Register(func(metrics.Snapshot, time.Time) (string, bool) {
		return "A test alert", true
	})

	assert.Equal(t, 1, len(p.alerts))
}

func TestCheckAlerts(t *testing.T) {
	p := NewRegistry()
	p.Register(func(metrics.Snapshot, time.Time) (string, bool) {
		return "A matching alert", true
	})
	p.Register(func(metrics.Snapshot, time.Time) (string, bool) {
		return "A non-matching alert", false
	})

	alerts, err := p.CheckAlerts(time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(alerts))
	assert.Equal(t, "A matching alert", alerts[0].Message)
	triggeredAt := alerts[0].TriggeredAt

	alerts, err = p.CheckAlerts(time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1, len(alerts))
	assert.Equal(t, triggeredAt, alerts[0].TriggeredAt)
}

func TestCheckResolvingAlert(t *testing.T) {
	p := NewRegistry()
	p.Register(func(snap metrics.Snapshot, ta time.Time) (string, bool) {
		totalRequests := snap.FetchCounter("requests.total")

		if totalRequests == 1 {
			return "Alert triggered!", true
		}

		return "", false
	})
	counter := metrics.RegistryProvider().Counter("requests.total")

	alerts, err := p.CheckAlerts(time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(alerts))

	counter.Increment(time.Now())
	alerts, err = p.CheckAlerts(time.Now().Add(time.Minute * -1))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(alerts))
	assert.Equal(t, "Alert triggered!", alerts[0].Message)

	counter.Increment(time.Now())
	alerts, err = p.CheckAlerts(time.Now().Add(time.Minute * -1))
	assert.NoError(t, err)
	assert.Equal(t, 0, len(alerts))
}
