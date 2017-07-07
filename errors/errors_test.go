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

// Upspin derived error.Kind
const (
	IO         Kind = "I/O error"
	Invalid    Kind = "Invalid"
	Permission Kind = "permission denied"
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
	e1 := E("Get", IO, err)

	// Nested error.
	e2 := E("Read", Other, e1)

	want := "Read: I/O error:: Get: network unreachable"
	if errorAsString(e2) != want {
		t.Errorf("expected %q; got %q", want, e2)
	}
}

func TestDoesNotChangePreviousError(t *testing.T) {
	err := E(Permission)
	err2 := E("I will NOT modify err", err)

	expected := "I will NOT modify err: permission denied"
	if errorAsString(err2) != expected {
		t.Fatalf("Expected %q, got %q", expected, err2)
	}
	kind := err.(*Error).Kind
	if kind != Permission {
		t.Fatalf("Expected kind %v, got %v", Permission, kind)
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
	{E("Op", Invalid, io.EOF), E("Op", Invalid, io.EOF), true},
	{E("Op", Invalid), E("Op", Invalid, io.EOF), true},
	{E("Op"), E("Op", Invalid, io.EOF), true},
	// Failure.
	{E(io.EOF), E(io.ErrClosedPipe), false},
	{E("Op1"), E("Op2"), false},
	{E(Invalid), E(Permission), false},
	{E("jane"), E("john"), false},
	{E("Op", Invalid, io.EOF, "jane"), E("Op", Invalid, io.EOF, "john"), false},
	{E("path1", Str("something")), E("path1"), false}, // Test nil error on rhs.
	// Nested *Errors.
	{E("Op1", "path1"), E("Op1", "john", E("Op2", "jane", "path1")), false},
	{E("Op1", E("path1")), E("Op1", "john", Str(E("Op2", "jane", "path1").Error())), false},
}

func TestMatch(t *testing.T) {
	for i, test := range matchTests {
		matched := Match(test.err1, test.err2)
		if matched != test.matched {
			t.Errorf("%d: Match(%q, %q)=%t; want %t", i, test.err1, test.err2, matched, test.matched)
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
