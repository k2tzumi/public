package btstrpr

import (
	"bytes"
	"fmt"
	"html/template"
)

const rootTpl = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Starter Template for Bootstrap</title>
<link rel="stylesheet" href="{{ .BaseCSS }}" crossorigin="anonymous">
</head>
<body>{{ .Body }}
<script src="{{ .JQuery }}"></script>
<script src="{{ .BaseJS }}" crossorigin="anonymous"></script>
</body>
</html>`

type Bootstrap struct {
	baseCSS, jQuery, baseJS, body string
}

func (b *Bootstrap) Render() {
	t := template.Must(template.New("root").Parse(rootTpl))
	var buf bytes.Buffer
	err := t.Execute(&buf, struct {
		BaseCSS string
		Body    string
		JQuery  string
		BaseJS  string
	}{
		b.baseCSS,
		b.body,
		b.jQuery,
		b.baseJS,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(buf.String())
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
