package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dmathieu/logtee/logtee"
	"github.com/dmathieu/logtee/metrics"
)

var (
	requestCountMetric = metrics.RegistryProvider().Counter("requests.total")
	sectionMetrics     = map[string]*metrics.Counter{}
	statusMetrics      = map[string]*metrics.Counter{}
)

// Handler allows receiving a `logtee.LogLine` and emits the appropriate metrics for it
type Handler struct {
}

// Handle is received whenever a log line is received, and emits metrics
func (h Handler) Handle(l logtee.LogLine) error {
	requestCountMetric.Increment(l.Time)
	logSection(l)

	strStatus := strconv.Itoa(l.Status)
	if _, ok := statusMetrics[strStatus]; !ok {
		statusMetrics[strStatus] = metrics.RegistryProvider().Counter(fmt.Sprintf("requests.status.%s", strStatus))
	}
	statusMetrics[strStatus].Increment(l.Time)

	return nil
}

func logSection(l logtee.LogLine) {
	section := strings.Split(l.Path, "/")[1]
	if section == "" {
		section = "home"
	}

	if _, ok := sectionMetrics[section]; !ok {
		sectionMetrics[section] = metrics.RegistryProvider().Counter(fmt.Sprintf("requests.section.%s", section))
	}

	sectionMetrics[section].Increment(l.Time)
}
