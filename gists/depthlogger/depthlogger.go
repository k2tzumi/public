package main // by Brad Fitzpatrick https://play.golang.org/p/PM0Fx-o6Drz

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func init() {
	log.SetOutput(depthLogger{})
}

func main() {
	log.Printf("in main")
	defer log.Printf("out of main")
	foo()
}

func foo() {
	log.Printf("in foo")
	defer log.Printf("out of foo")
	bar()
}

func bar() {
	log.Printf("in foo")
	defer log.Printf("out of foo")
}

type depthLogger struct{}

func (depthLogger) Write(p []byte) (int, error) {
	fmt.Fprintf(os.Stderr, "%s%s%s", p[:19], strings.Repeat(" ", depth()*2), p[19:])
	return len(p), nil
}

func depth() int {
	var p [50]uintptr
	return runtime.Callers(9, p[:])
}
