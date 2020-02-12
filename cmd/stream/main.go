package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/dmathieu/logtee/alerts"
	"github.com/dmathieu/logtee/handler"
	"github.com/dmathieu/logtee/logtee"
	"github.com/dmathieu/logtee/metrics"
	"github.com/dmathieu/logtee/output"
	"github.com/oklog/run"
)

var (
	logFile            string
	alertTotalRequests int64
)

func registerAlerts(registry *alerts.Registry) {
	registry.Register(func(snap metrics.Snapshot, triggeredAt time.Time) (string, bool) {
		totalRequests := snap.FetchCounter("requests.total")
		if totalRequests < alertTotalRequests {
			return "", false
		}
		ta := triggeredAt.Format(time.RFC3339)

		return fmt.Sprintf("High traffic generated an alert - hits = %d, triggered at %s", totalRequests, ta), true
	})
}

func main() {
	flag.StringVar(&logFile, "file", "/tmp/access.log", "Path to the log file")
	flag.Int64Var(&alertTotalRequests, "alert-total-requests", 10, "Total requests threshold after which to send alert")
	flag.Parse()

	var g run.Group
	alertsRegistry := alerts.NewRegistry()
	registerAlerts(alertsRegistry)

	s, err := logtee.NewStreamer(logFile, logtee.NewParser(handler.Handler{}))
	if err != nil {
		log.Fatal(err)
	}

	g.Add(func() error {
		return s.Run()
	}, func(error) {
		s.Close()
	})

	ctx, cancel := context.WithCancel(context.Background())
	g.Add(func() error {
		return output.Run(ctx, alertsRegistry)
	}, func(error) {
		cancel()
	})

	err = g.Run()
	if err != nil {
		log.Fatal(err)
	}
}
