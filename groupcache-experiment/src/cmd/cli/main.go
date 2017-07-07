package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"

	"api"
	"client"
)

var (
	listenAddr = flag.String("listen", "localhost:9001", "frontend listening address")
	backend    = flag.String("backend", "localhost:8080", "backend listening address")

	get = flag.Bool("get", false, "get a key-value pair")
	set = flag.Bool("set", false, "set a new key-value pair")

	key   = flag.String("k", "", "key to get/set")
	value = flag.String("v", "", "value to set")
)

func init() {
	flag.Parse()
}

func main() {
	switch true {
	case *get:
		getKV(*listenAddr, *key)
	case *set:
		setKV(*backend, *key, *value)
	default:
		flag.PrintDefaults()
	}
}

func getKV(addr, key string) {
	rc, err := rpc.DialHTTP("tcp", addr)
	checkerr(err)
	var reply api.ValueResult
	err = rc.Call("Frontend.Get", &api.Load{key}, &reply)
	checkerr(err)
	fmt.Println(string(reply.Value))
}

func setKV(addr, k, v string) {
	c := client.NewClient(addr)
	c.Set(k, v)
}

func checkerr(err interface{}) {
	if err != nil {
		log.Fatalln("error:", err)
	}
}
