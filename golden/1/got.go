package main

import (
	"context"

	"cirello.io/btstrpr"
)

func main() {

	b := btstrpr.New(
		btstrpr.BaseCSS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css"),
		btstrpr.JQuery("https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"),
		btstrpr.BaseJS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"),
		btstrpr.Body(
			btstrpr.Container(btstrpr.S("hello world")),
		),
	)

	b.Render(context.Background())
}
