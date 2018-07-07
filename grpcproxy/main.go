package main

func main() {
	hubAddr, hub := startHub()
	go startGateway(hub)
	proxy(hubAddr)
}
