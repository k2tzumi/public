package btstrpr

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var layoutTests = []struct {
	b *Bootstrap
	c context.Context
}{
	{New(
		BaseCSS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css"),
		JQuery("https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"),
		BaseJS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"),
	), context.Background()},
	{New(
		BaseCSS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css"),
		JQuery("https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"),
		BaseJS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"),
		Body(
			Container(S("hello world")),
		),
	), context.Background()},
	{New(
		BaseCSS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css"),
		JQuery("https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"),
		BaseJS("https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js"),
		Body(
			Container(
				S("hello world"),
				"style", "margin: 0",
			),
		),
	), context.Background()},
}

func TestRender(t *testing.T) {
	for i, layout := range layoutTests {
		var got bytes.Buffer
		layout.b.Render(layout.c, &got)

		fn := filepath.Join("golden", fmt.Sprint(i, ".html"))
		expected, err := ioutil.ReadFile(fn)
		if err != nil {
			t.Fatal(err)
		}

		if result := bytes.Compare(got.Bytes(), expected); result != 0 {
			t.Error(fn, "error")
			t.Log("got:", got.String())
			t.Log("len:", len(got.String()))
			t.Log("expected:", string(expected))
			t.Log("len:", len(string(expected)))
		}
	}
}
