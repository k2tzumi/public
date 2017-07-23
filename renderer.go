package btstrpr

import (
	"context"
	"io"
)

// Renderer represents components that yield any output.
type Renderer func(context.Context, io.Writer)
