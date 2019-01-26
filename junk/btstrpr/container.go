package btstrpr

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
)

type htmlattributes [][2]string

// Container renders a div with container class. Panics if attributes are not
// pairs.
func Container(body Renderer, attributes ...interface{}) Renderer {
	if len(attributes)%2 != 0 {
		panic("attributes must always be pairs.")
	}
	attrs := htmlattributes{
		[2]string{"class", "container"},
	}
	for i := 0; i < len(attributes); i += 2 {
		attrs = append(attrs, [2]string{
			fmt.Sprint(attributes[i]),
			fmt.Sprint(attributes[i+1]),
		})
	}
	return func(c context.Context, out io.Writer) {
		fmt.Fprint(out, `<div `, renderAttrs(attrs), `>`)
		body(c, out)
		fmt.Fprintln(out, `</div>`)
	}
}

// FluidContainer renders a div with container-fluid class. Panics if attributes
// are not pairs.
func FluidContainer(body Renderer, attributes ...interface{}) Renderer {
	if len(attributes)%2 != 0 {
		panic("attributes must always be pairs.")
	}
	attrs := htmlattributes{
		[2]string{"class", "container-fluid"},
	}
	for i := 0; i < len(attributes); i += 2 {
		attrs = append(attrs, [2]string{
			fmt.Sprint(attributes[i]),
			fmt.Sprint(attributes[i+1]),
		})
	}
	return func(c context.Context, out io.Writer) {
		fmt.Fprint(out, `<div `, renderAttrs(attrs), `>`)
		body(c, out)
		fmt.Fprintln(out, `</div>`)
	}
}

func renderAttrs(attrs htmlattributes) string {
	var buf bytes.Buffer
	for _, attr := range attrs {
		buf.WriteString(attr[0])
		buf.WriteString("=")
		buf.WriteString(`"`)
		buf.WriteString(attr[1])
		buf.WriteString(`" `)
	}
	return strings.TrimSpace(buf.String())
}

// S stands for String - it is used to insert arbitrary text in the code. Does
// not do any sanitization.
func S(args ...interface{}) Renderer {
	return func(c context.Context, out io.Writer) {
		fmt.Fprint(out, args...)
	}
}

// Nil is a terminator for components that demand the existence of body
func Nil() Renderer {
	return func(context.Context, io.Writer) {}
}
