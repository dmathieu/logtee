package output

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/dmathieu/logtee/alerts"
	"github.com/dmathieu/logtee/metrics"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	refreshInterval = 2 * time.Second
	metricsWindow   = 10 * time.Second
	alertsWindow    = 2 * time.Minute
)

// Run refreshes the screen's content regularly with the metrics data
func Run(ctx context.Context, alertsRegistry *alerts.Registry) error {
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	err := ui.Init()
	if err != nil {
		return err
	}
	defer ui.Close()
	uiEvents := ui.PollEvents()

	err = updateContent(alertsRegistry)
	if err != nil {
		return err
	}

	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			}
		case <-ticker.C:
			err := updateContent(alertsRegistry)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func updateContent(alertsRegistry *alerts.Registry) error {
	since := time.Now().Add(-metricsWindow)
	snap, err := metrics.RegistryProvider().Collect(since)
	if err != nil {
		return err
	}

	tot, err := buildTotal(snap)
	if err != nil {
		return err
	}

	il, err := buildInvalidLogLines(snap)
	if err != nil {
		return err
	}

	sec, err := buildSections(snap)
	if err != nil {
		return err
	}

	stat, err := buildStatuses(snap)
	if err != nil {
		return err
	}

	al, err := buildAlerts(alertsRegistry)
	if err != nil {
		return err
	}

	ui.Render(tot, il, sec, stat, al)
	return nil
}

func buildTotal(snap metrics.Snapshot) (ui.Drawable, error) {
	p := widgets.NewParagraph()
	p.Title = "Total Requests"
	p.Text = fmt.Sprintf("%d", snap.FetchCounter("requests.total"))
	p.TextStyle.Fg = ui.ColorRed
	p.TextStyle.Modifier = ui.ModifierBold
	p.SetRect(0, 0, 25, 4)
	return p, nil
}

func buildInvalidLogLines(snap metrics.Snapshot) (ui.Drawable, error) {
	p := widgets.NewParagraph()
	p.Title = "Invalid Log Lines"
	p.Text = fmt.Sprintf("%d", snap.FetchCounter("logtee.invalid_line"))
	p.TextStyle.Fg = ui.ColorYellow
	p.SetRect(0, 4, 25, 8)
	return p, nil
}

func buildSections(snap metrics.Snapshot) (ui.Drawable, error) {
	l := widgets.NewList()
	l.Title = "Sections"
	l.Rows = []string{}
	l.TextStyle.Fg = ui.ColorYellow
	l.SetRect(25, 0, 50, 8)

	sections := snap.FetchAllCounters("requests.section.")
	for _, v := range sections {
		section := strings.Split(v.Key, ".")[2]
		l.Rows = append(l.Rows, fmt.Sprintf("%s: %d", section, v.Value))
	}

	return l, nil
}

func buildStatuses(snap metrics.Snapshot) (ui.Drawable, error) {
	l := widgets.NewList()
	l.Title = "Statuses"
	l.Rows = []string{}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(50, 0, 75, 8)

	statuses := snap.FetchAllCounters("requests.status.")
	for _, v := range statuses {
		status := strings.Split(v.Key, ".")[2]
		l.Rows = append(l.Rows, fmt.Sprintf("%s: %d", status, v.Value))
	}

	return l, nil
}

func buildAlerts(registry *alerts.Registry) (ui.Drawable, error) {
	since := time.Now().Add(-alertsWindow)
	alerts, err := registry.CheckAlerts(since)
	if err != nil {
		return nil, err
	}

	l := widgets.NewList()
	l.Title = fmt.Sprintf("%d triggered alerts", len(alerts))
	l.Rows = []string{}
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetRect(0, 8, 150, 16)

	for _, v := range alerts {
		l.Rows = append(l.Rows, v.Message)
	}

	return l, nil
}
