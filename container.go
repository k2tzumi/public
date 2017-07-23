package btstrpr

import (
	"context"
	"fmt"
)

func Container(body Renderer) Renderer {
	return func(c context.Context) {
		fmt.Print(`<div class="container">`)
		body(c)
		fmt.Println(`</div>`)
	}
}

func FluidContainer(body Renderer) Renderer {
	return func(c context.Context) {
		fmt.Print(`<div class="container-fluid">`)
		body(c)
		fmt.Println(`</div>`)
	}
}

func S(args ...interface{}) Renderer {
	return func(c context.Context) {
		fmt.Print(args...)
	}
}
