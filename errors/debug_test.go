// Copyright 2017 Ulderico Cirello. All rights reserved.
// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build debug

package errors_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"cirello.io/errors"
)

var errorLines = strings.Split(strings.TrimSpace(`
	.*/cirello.io/errors/debug_test.go:\d+: cirello.io/errors_test.func1:
	.*/cirello.io/errors/debug_test.go:\d+: ...T.func2:
	.*/cirello.io/errors/debug_test.go:\d+: ...func3:
	.*/cirello.io/errors/debug_test.go:\d+: ...func4: op:
	valid.UserName: bad-username
`), "\n")

var errorLineREs = make([]*regexp.Regexp, len(errorLines))

func init() {
	for i, s := range errorLines {
		errorLineREs[i] = regexp.MustCompile(fmt.Sprintf("^%s$", s))
	}
}

// Test that the error stack includes all the function calls between where it
// was generated and where it was printed. It should not include the name
// of the function in which the Error method is called. It should coalesce
// the call stacks of nested errors into one single stack, and present that
// stack before the other error values.
func TestDebug(t *testing.T) {
	got := printErr(t, func1())
	lines := strings.Split(got, "\n")
	for i, re := range errorLineREs {
		if i >= len(lines) {
			// Handled by line number check.
			break
		}
		if !re.MatchString(lines[i]) {
			t.Errorf("error does not match at line %v, got:\n\t%q\nwant:\n\t%q", i, lines[i], re)
		}
	}
	// Check number of lines after checking the lines themselves,
	// as the content check will likely be more illuminating.
	if got, want := len(lines), len(errorLines); got != want {
		t.Errorf("got %v lines of errors, want %v", got, want)
	}
}

func printErr(t *testing.T, err error) string {
	return err.Error()
}

func func1() error {
	var t T
	return t.func2()
}

type T struct{}

func (T) func2() error {
	return errors.E("op", func3())
}

func func3() error {
	return func4()
}

func func4() error {
	return errors.E("valid.UserName", errors.Str("bad-username"))
}
