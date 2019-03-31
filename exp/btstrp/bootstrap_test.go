// Copyright 2019 github.com/ucirello and https://cirello.io. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to writing, software distributed
// under the License is distributed on a "AS IS" BASIS, WITHOUT WARRANTIES OR
// CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

package btstrp_test

import (
	"html/template"
	"log"
	"os"
	"testing"
	"text/template/parse"

	"cirello.io/exp/btstrp"
	"github.com/davecgh/go-spew/spew"
)

func TestParse(t *testing.T) {
	a, b := parse.Parse("hello", `
		{{ something }}
		{{ end }}
	`, "{{", "}}")
	spew.Dump(a, b)
}

func ExampleTree() {
	const pageTpl = `{{ template "Page" . }}`
	page := template.Must(btstrp.Template().Parse(pageTpl))
	err := page.Execute(os.Stdout, []interface{}{
		btstrp.Container{
			Variant: "danger",
			Body: btstrp.Button{
				Variant: "danger",
				Body:    btstrp.Inline{Body: "hellow"},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// <!doctype html>
	// <html lang="en">
	// <head> <meta charset="utf-8">  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"> <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
	// <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
	// <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
	// <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
	// </head>
	// <body><div class="container container-danger"><button type="button" class="btn btn-danger">hellow</button></div></body>
	// </html>
}

func Example() {
	const pageTpl = `<!doctype html>
<html lang="en">
	<head>{{template "Intro"}}</head>
	<body>
	{{template "ContainerStart" .}}
		{{ template "Button" .ButtonA }}
		{{ template "Button" .ButtonB }}
	{{template "ContainerEnd" .}}
	</body>
</html>`

	page := template.Must(btstrp.Template().Parse(pageTpl))

	err := page.Execute(os.Stdout, map[string]interface{}{
		"ButtonA": btstrp.Button{
			Type:    "reset",
			Body:    "reset me",
			Variant: "primary",
			Class:   "bold",
			Attributes: btstrp.Attributes{
				"onClick":     "javascript: alert('hello world')",
				"onMouseOver": "javascript: alert('hello world 2')",
			},
		},
		"ButtonB": btstrp.Button{
			Body:    "button B",
			Variant: "danger",
		},
		"ContainerB": btstrp.Container{
			Body: "containerB",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Output:
	// <!doctype html>
	// <html lang="en">
	// 	<head> <meta charset="utf-8">  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"> <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
	// <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
	// <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
	// <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
	// </head>
	// 	<body>
	// 	<div class="container">
	// 		<button type="reset" class="btn btn-primary bold" onClick="javascript: alert(&#39;hello world&#39;)" onMouseOver="javascript: alert(&#39;hello world 2&#39;)">reset me</button>
	// 		<button type="button" class="btn btn-danger">button B</button>
	// 	</div>
	// 	</body>
	// </html>
}
