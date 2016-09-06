package main // import "cirello.io/bloomfilterd/cmd/bloomfilterd"

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"cirello.io/bloomfilterd/internal/http"
	"cirello.io/bloomfilterd/internal/storage"
	"cirello.io/suture"
	"github.com/coreos/etcd/raft/raftpb"
)

func main() {
	cluster := flag.String("cluster", "http://127.0.0.1:9021", "comma separated cluster peers")
	id := flag.Int("id", 1, "node ID")
	listen := flag.String("listen", ":9121", "listen")
	join := flag.Bool("join", false, "join an existing cluster")
	flag.Parse()

	fmt.Println(*cluster, *id, *listen, *join)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	propose := make(chan string)
	defer close(propose)
	confChange := make(chan raftpb.ConfChange)
	defer close(confChange)

	var supervisor suture.Supervisor
	http := http.New(
		propose,
		confChange,
		http.Storage(storage.Memory),
		http.Listen(*listen),
	)

	supervisor.Add(http)
	supervisor.ServeBackground()

	<-c
	fmt.Println("stopping...")
	supervisor.Stop()
}
