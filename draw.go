package main

import (
	"bytes"
	"context"
	"fmt"
	"sort"
	"text/tabwriter"

	"cirello.io/logmon/monitor"
	ui "gopkg.in/gizak/termui.v2"
)

type statsLoader interface {
	Stats() map[string]monitor.Stat
}

type alertsWatcher interface {
	Alerts() []monitor.Alert
}

type logWatcher interface {
	statsLoader
	alertsWatcher
}

func drawscreen(ctx context.Context, lw logWatcher) error {
	err := ui.Init()
	if err != nil {
		return err
	}
	defer ui.Close()

	go func() {
		<-ctx.Done()
		ui.StopLoop()
	}()

	sections := ui.NewPar("Hot Sections\nLoading...\n")
	sections.Height = ui.TermHeight()
	sections.Y = 1
	sections.Text = renderedStats(lw)

	alerts := ui.NewPar("Alerts\n")
	alerts.Height = ui.TermHeight()
	alerts.Y = 1

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(8, 0, sections),
			ui.NewCol(4, 0, alerts),
		),
	)
	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/sys/kbd/q", func(ui.Event) {
		ui.StopLoop()
	})

	// ui.DefaultEvtStream.Merge("timer", ui.NewTimerCh(10*time.Second))
	// ui.Handle("/timer/10s", func(e ui.Event) {
	// 	sections.Text = renderedStats(lw)
	// 	sections.Height = ui.TermHeight()
	// 	alerts.Height = ui.TermHeight()
	// 	ui.Body.Align()
	// 	ui.Render(ui.Body)
	// })
	ui.Handle("/timer/1s", func(e ui.Event) {
		sections.Text = renderedStats(lw)
		sections.Height = ui.TermHeight()
		alerts.Height = ui.TermHeight()
		ui.Body.Align()
		ui.Render(ui.Body)

		alerts.Text = renderedAlerts(lw)
		sections.Height = ui.TermHeight()
		alerts.Height = ui.TermHeight()
		ui.Body.Align()
		ui.Render(ui.Body)
	})
	ui.Loop()

	return nil
}

func renderedStats(l statsLoader) string {
	stats := l.Stats()
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "Sections")

	w := tabwriter.NewWriter(&buf, 10, 8, 1, ' ', 0)
	fmt.Fprint(w, "\t", "hits", "\t", "2xx", "\t", "3xx", "\t", "4xx", "\t", "5xx", "\n")

	var ranking []rank
	for k, v := range stats {
		ranking = append(ranking, rank{k, v.Hits})
	}
	sort.Sort(sort.Reverse(byHits(ranking)))

	for _, i := range ranking {
		section, stat := i.section, stats[i.section]
		fmt.Fprint(w, section, "\t", stat.Hits, "\t", stat.Status2xx, "\t", stat.Status3xx, "\t", stat.Status4xx, "\t", stat.Status5xx, "\n")
	}
	w.Flush()

	return buf.String()
}

func renderedAlerts(a alertsWatcher) string {
	alerts := a.Alerts()

	var msg string
	for _, alert := range alerts {
		msg = fmt.Sprintln("High traffic generated an alert - hits =", alert.Hits, "triggered at", alert.Start) + msg
		if !alert.End.IsZero() {
			msg = fmt.Sprintln("High traffic generated an alert - hits =", alert.Hits, "recovered at", alert.End) + msg
		}
	}

	return "Alerts\n" + msg
}

type rank struct {
	section string
	hits    uint64
}

type byHits []rank

func (a byHits) Len() int           { return len(a) }
func (a byHits) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byHits) Less(i, j int) bool { return a[i].hits < a[j].hits }
