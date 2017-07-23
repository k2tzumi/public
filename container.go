package btstrpr

import (
	"context"
	"fmt"
)

// Container renders a div with container class
func Container(body Renderer) Renderer {
	return func(c context.Context) {
		fmt.Print(`<div class="container">`)
		body(c)
		fmt.Println(`</div>`)
	}
}

// FluidContainer renders a div with container-fluid class
func FluidContainer(body Renderer) Renderer {
	return func(c context.Context) {
		fmt.Print(`<div class="container-fluid">`)
		body(c)
		fmt.Println(`</div>`)
	}
}

// S stands for String - it is used to insert arbitrary text in the code. Does
// not do any sanitization.
func S(args ...interface{}) Renderer {
	return func(c context.Context) {
		fmt.Print(args...)
	}
}

// Nil is a terminator for components that demand the existence of body
func Nil() Renderer {
	return func(c context.Context) {}
}
