// Copyright 2017 Ulderico Cirello. All rights reserved.
// Copyright 2016 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !debug

package errors_test

import (
	"fmt"

	"cirello.io/errors"
)

func ExampleError() {
	err := errors.Str("network unreachable")

	// Single error.
	e1 := errors.E("Get", err)
	fmt.Println("\nSimple error:")
	fmt.Println(e1)

	// Nested error.
	fmt.Println("\nNested error:")
	e2 := errors.E("Read", e1)
	fmt.Println(e2)

	// Output:
	//
	// Simple error:
	// Get: network unreachable
	//
	// Nested error:
	// Read:
	//	Get: network unreachable
}

func ExampleMatch() {
	err := errors.Str("network unreachable")

	// Construct an error, one we pretend to have received from a test.
	got := errors.E("Get", err)

	// Now construct a reference error, which might not have all
	// the fields of the error from the test.
	expect := errors.E(err)

	fmt.Println("Match:", errors.Match(expect, got))

	// Output:
	//
	// Match: true
}
