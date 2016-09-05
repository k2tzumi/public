package main // import "cirello.io/bloomfilterd/cmd/bloomfilterd"

import (
	"fmt"
	"os"
	"os/signal"

	"cirello.io/bloomfilterd/internal/http"
	"cirello.io/suture"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	var supervisor suture.Supervisor
	http := http.New()
	supervisor.Add(http)
	supervisor.ServeBackground()

	<-c
	fmt.Println("stopping...")
	supervisor.Stop()
}
