/*
Copyright 2018 github.com/ucirello

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pglock

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"testing"

	"github.com/lib/pq"
	"golang.org/x/xerrors"
)

func TestTypedError(t *testing.T) {
	errs := []struct {
		err          error
		expectedType error
	}{
		{fmt.Errorf("random error"), nil},
		{sql.ErrNoRows, ErrNotExist},
		{&net.OpError{}, ErrUnavailable},
		{&pq.Error{}, nil},
		{&pq.Error{Code: "40001"}, ErrFailedPrecondition},
	}
	for _, err := range errs {
		if err.err == nil {
			continue
		}
		typedErr := typedError(err.err, "")
		if err.expectedType == nil {
			if typedErr == nil {
				t.Errorf("unexpected nil error: %#v", typedErr)
			}
			continue
		} else if !xerrors.Is(typedErr, err.expectedType) {
			t.Errorf("untyped error found: %#v", typedErr)
		}
	}
}

func TestRetry(t *testing.T) {
	t.Run("type check", func(t *testing.T) {
		internal := xerrors.New("internal")
		c := &Client{
			log: &testLogger{t},
		}
		errs := []error{
			xerrors.Errorf("failed precondition: %w", ErrFailedPrecondition),
			xerrors.Errorf("other error: %w", internal),
		}
		err := c.retry(func() error {
			var err error
			err, errs = errs[0], errs[1:]
			return err
		})
		if !xerrors.Is(err, internal) {
			t.Fatal("unexpected error kind found")
		}
	})
	t.Run("max retries", func(t *testing.T) {
		c := &Client{
			log: log.New(ioutil.Discard, "", 0),
		}
		var retries int
		err := c.retry(func() error {
			retries++
			return xerrors.Errorf("failed precondition: %w", ErrFailedPrecondition)
		})
		if !xerrors.Is(err, ErrFailedPrecondition) {
			t.Fatal("unexpected error kind found")
		}
		if retries != maxRetries {
			t.Fatal("unexpected retries count found")
		}
		t.Log(retries, maxRetries)
	})
}

type testLogger struct {
	t *testing.T
}

func (t *testLogger) Println(v ...interface{}) {
	t.t.Log(v...)
}
