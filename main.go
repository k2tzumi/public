package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"cirello.io/logmon/monitor"
	"cirello.io/logmon/tail"
)

func main() {
	threshold := flag.Uint64("threshold", 1000, "number of hits per second")
	duration := flag.Duration("duration", 2*time.Minute, "time period of analysis of threshold breaches before alert is triggered")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
	}()

	fn := flag.Arg(0)
	content, err := tail.Open(fn)
	exitonerr(err)

	m := monitor.Open(
		ctx, content,
		monitor.Threshold(*threshold),
		monitor.Timespan(*duration),
	)
	exitonerr(drawscreen(ctx, m))
	exitonerr(m.Err())
}

func exitonerr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
