// Copyright 2017 Ulderico Cirello. All rights reserved.
// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !debug

package errors

import (
	"io"
	"os"
	"os/exec"
	"testing"
)

func TestDebug(t *testing.T) {
	// Test with -tags debug to run the tests in debug_test.go
	cmd := exec.Command("go", "test", "-tags", "debug")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("external go test failed: %v", err)
	}
}

func TestSeparator(t *testing.T) {
	defer func(prev string) {
		Separator = prev
	}(Separator)
	Separator = ":: "

	// Same pattern as above.
	err := Str("network unreachable")

	// Single error. No user is set, so we will have a zero-length field inside.
	e1 := E("Get", err)

	// Nested error.
	e2 := E("Read", e1)

	want := "Read:: Get: network unreachable"
	if errorAsString(e2) != want {
		t.Errorf("expected %q; got %q", want, e2)
	}
}

func TestDoesNotChangePreviousError(t *testing.T) {
	err := E("permission denied")
	err2 := E("I will NOT modify err", err)

	expected := "I will NOT modify err:\n\tpermission denied"
	if errorAsString(err2) != expected {
		t.Fatalf("Expected %q, got %q", expected, err2)
	}
}

func TestNoArgs(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("E() did not panic")
		}
	}()
	_ = E()
}

type matchTest struct {
	err1, err2 error
	matched    bool
}

var matchTests = []matchTest{
	// Errors not of type *Error fail outright.
	{nil, nil, false},
	{io.EOF, io.EOF, false},
	{E(io.EOF), io.EOF, false},
	{io.EOF, E(io.EOF), false},
	// Success. We can drop fields from the first argument and still match.
	{E(io.EOF), E(io.EOF), true},
	{E("Op", io.EOF), E("Op", io.EOF), true},
	{E("Op"), E("Op", io.EOF), true},
	// Failure.
	{E(io.EOF), E(io.ErrClosedPipe), false},
	{E("Op1"), E("Op2"), false},
	{E("Op", io.EOF, "jane"), E("Op", io.EOF, "john"), false},
	{E("path1", Str("something")), E("path1"), false}, // Test nil error on rhs.
	// Nested *Errors.
	{E("Op1", E("path1")), E("Op1", E("Op2", "path1")), true},
	{E("Op1", "path1"), E("Op1", E("Op2", "path1")), false},
	{E("Op1", E("path1")), E("Op1", Str(E("Op2", "path1").Error())), false},
}

func TestMatch(t *testing.T) {
	for _, test := range matchTests {
		matched := Match(test.err1, test.err2)
		if matched != test.matched {
			t.Errorf("Match(%q, %q)=%t; want %t", test.err1, test.err2, matched, test.matched)
		}
	}
}

// errorAsString returns the string form of the provided error value.
// If the given string is an *Error, the stack information is removed
// before the value is stringified.
func errorAsString(err error) string {
	if e, ok := err.(*Error); ok {
		e2 := *e
		e2.stack = stack{}
		return e2.Error()
	}
	return err.Error()
}
