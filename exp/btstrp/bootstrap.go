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

// Package btstrp adds a set of html templates based on Bootstrap.
package btstrp // import "cirello.io/exp/btstrp"

import (
	"fmt"
	"html/template"
	"strings"
)

var base = template.Must(template.New("bootstraptpl").
	Option("missingkey=zero").
	Funcs(template.FuncMap{
		"isIter": func(pipeline interface{}) bool {
			switch pipeline.(type) {
			case []interface{}:
				return true
			}
			return false
		},
		"typeAssert": func(pipeline interface{}, t string) bool {
			return strings.HasSuffix(fmt.Sprintf("%T", pipeline), "."+t)
		},
		"isNotBootstrap": func(pipeline interface{}) bool {
			return !strings.HasPrefix(fmt.Sprintf("%T", pipeline), "btstrp.")
		},
	}).
	Parse(bootstraptplDefs))

// Template returns the preloaded bootstrap template.
func Template() *template.Template {
	return template.Must(base.Clone())
}

// Attributes extends the HTML tag of the bootstrap component.
type Attributes map[template.HTMLAttr]string

// Type defines the bootstrap elements that this template implementation knows
// how to convert in runtime.
type Type string

// All known element types.
const (
	ComponentInline    Type = "Inline"
	ComponentIntro     Type = "Intro"
	ComponentContainer Type = "Container"
	ComponentButton    Type = "Button"
)

var bootstraptplDefs = `{{define "CSS"}}<link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">{{end}}
{{define "JS"}}
<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
{{end}}
{{define "Meta"}} <meta charset="utf-8"> {{end}}
{{define "Responsive"}} <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"> {{end}}
{{define "Intro"}}{{template "Meta"}}{{template "Responsive"}}{{template "CSS"}}{{template "JS"}}{{end}}
{{define "ContainerStart"}}<div class="container{{with .Variant}} container-{{.}}{{end}}{{with .Class}} {{.}}{{end -}}"{{range $k, $v := .Attributes}} {{$k}}="{{- $v -}}"{{end}}>{{end}}
{{define "ContainerEnd"  }}</div>{{end}}
{{define "Container"     }}{{ template "ContainerStart" . }}{{template "Components" .Body }}{{ template "ContainerEnd" . }}{{end}}
{{define "RowStart"}}<div class="row{{with .Class}} {{.}}{{end -}}"{{range $k, $v := .Attributes}} {{$k}}="{{- $v -}}"{{end}}>{{end}}
{{define "RowEnd"  }}</div>{{end}}
{{define "Row"     }}{{ template "RowStart" . }}{{.Body}}{{ template "RowEnd" . }}{{end}}
{{define "ColStart"}}<div class="col{{with .Variant}} col-{{.}}{{end}}{{with .Class}} {{.}}{{end -}}"{{range $k, $v := .Attributes}} {{$k}}="{{- $v -}}"{{end}}>{{end}}
{{define "ColEnd"  }}</div>{{end}}
{{define "Col"     }}{{ template "ColStart" . }}{{.Body}}{{ template "ColEnd" . }}{{end}}
{{define "Button"}}<button type="{{- with .Type -}}{{.}}{{- else -}}button{{- end}}" class="btn{{with .Variant}} btn-{{.}}{{end}}{{with .Class}} {{.}}{{end -}}"{{range $k, $v := .Attributes}} {{$k}}="{{- $v -}}"{{end}}>{{template "Components" .Body}}</button>{{end}}
{{define "Inline"}}{{template "Components" .Body}}{{end}}
{{- define "Page" }}<!doctype html>
<html lang="en">
<head>{{block "Head" .}}{{template "Intro"}}{{end}}</head>
<body>{{block "Body" .}}{{template "Components" .}}{{end}}</body>
</html>
{{end -}}
{{- define "Component" }}
{{- if typeAssert       . "Inline"    }}{{template "Inline" . }}{{ end -}}
{{- if typeAssert       . "Intro"     }}{{template "Intro" . }}{{ end -}}
{{- if typeAssert       . "Container" }}{{template "Container" . }}{{ end -}}
{{- if typeAssert       . "Button"    }}{{template "Button" . }}{{ end -}}
{{- if isNotBootstrap   .             }}{{ . }}{{ end -}}
{{end -}}
{{- define "Components" }}
{{- if isIter . }}
{{- range . }}{{ template "Component" . }}{{end -}}
{{- else -}}{{ template "Component" . }}{{end -}}
{{end -}}
`

// Inline adds an inline component without any escaping.
type Inline struct {
	Body interface{}
}

// Intro adds the HTML preamble necessary for bootstrap.
type Intro struct {
}

// Container creates a container component.
type Container struct {
	Type       string
	Body       interface{}
	Variant    string
	Class      string
	Attributes map[template.HTMLAttr]string
}

// Button creates a button component.
type Button struct {
	Type       string
	Body       interface{}
	Variant    string
	Class      string
	Attributes map[template.HTMLAttr]string
}
