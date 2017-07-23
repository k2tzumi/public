package btstrpr

import (
	"context"
	"fmt"
)

func Container(bodies ...Renderer) Renderer {
	return func(c context.Context) {
		fmt.Print(`<div class="container">`)
		for _, body := range bodies {
			body(c)
		}
		fmt.Println(`</div>`)
	}
}

func FluidContainer(bodies ...Renderer) Renderer {
	return func(c context.Context) {
		fmt.Print(`<div class="container-fluid">`)
		for _, body := range bodies {
			body(c)
		}
		fmt.Println(`</div>`)
	}
}

func S(args ...interface{}) Renderer {
	return func(c context.Context) {
		fmt.Print(args...)
	}
}
