package client

import (
	"log"
	"net/rpc"

	"api"
)

type Client struct {
	rpc *rpc.Client
}

func NewClient(addr string) *Client {
	rpc, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Fatalln("error:", err)
	}
	return &Client{rpc}
}

func (c *Client) Get(key string) string {
	var reply api.ValueResult
	err := c.rpc.Call("Backend.Get", &api.Load{key}, &reply)
	if err != nil {
		log.Fatalln("error:", err)
	}
	return string(reply.Value)
}

func (c *Client) Set(key string, value string) {
	var reply int
	err := c.rpc.Call("Backend.Set", &api.Store{key, value}, &reply)
	if err != nil {
		log.Fatalln("error:", err)
	}
}
