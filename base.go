package btstrpr

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
)

const rootTopTpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Starter Template for Bootstrap</title>
<link rel="stylesheet" href="{{ .BaseCSS }}" crossorigin="anonymous">
</head>
<body>`

const rootBottomTpl = `<script src="{{ .JQuery }}"></script>
<script src="{{ .BaseJS }}" crossorigin="anonymous"></script>
</body>
</html>`

type Bootstrap struct {
	baseCSS string
	jQuery  string
	baseJS  string
	body    Renderer
}

func (b *Bootstrap) Render(c context.Context) {
	t := template.Must(template.New("rootTop").Parse(rootTopTpl))
	var bufTop bytes.Buffer
	err := t.Execute(&bufTop, struct{ BaseCSS string }{b.baseCSS})
	if err != nil {
		panic(err)
	}
	fmt.Println(bufTop.String())

	if b.body != nil {
		b.body(c)
	}

	t = template.Must(template.New("rootBottom").Parse(rootBottomTpl))
	var bufBottom bytes.Buffer
	err = t.Execute(&bufBottom, struct{ JQuery, BaseJS string }{b.jQuery, b.baseJS})
	if err != nil {
		panic(err)
	}
	fmt.Println(bufBottom.String())

}

type Option func(*Bootstrap)

func BaseCSS(css string) Option {
	return func(b *Bootstrap) {
		b.baseCSS = css
	}
}
func JQuery(jQuery string) Option {
	return func(b *Bootstrap) {
		b.jQuery = jQuery
	}
}

func BaseJS(js string) Option {
	return func(b *Bootstrap) {
		b.baseJS = js
	}
}

func Body(body Renderer) Option {
	return func(b *Bootstrap) {
		b.body = body
	}
}

func New(opts ...Option) *Bootstrap {
	b := &Bootstrap{
		baseCSS: "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css",
		jQuery:  "https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js",
		baseJS:  "https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js",
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}
