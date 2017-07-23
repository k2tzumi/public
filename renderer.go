package btstrpr

import (
	"context"
	"io"
)

type Renderer func(io.Writer, context.Context)
