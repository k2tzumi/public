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
	"net"

	"github.com/lib/pq"
	"golang.org/x/xerrors"
)

// ErrNotExist qualifies the wrapping error with the kind NotExist.
var ErrNotExist = xerrors.New("not exist")

// ErrUnavailable qualifies the wrapping error with the kind Unavailable.
var ErrUnavailable = xerrors.New("not exist")

// ErrFailedPrecondition qualifies the wrapping error with the kind
// FailedPrecondition.
var ErrFailedPrecondition = xerrors.New("failed precondition")

// ErrNotPostgreSQLDriver is returned when an invalid database connection is
// passed to this locker client.
var ErrNotPostgreSQLDriver = xerrors.New("this is not a PostgreSQL connection")

// ErrNotAcquired indicates the given lock is already enforce to some other
// client.
var ErrNotAcquired = xerrors.New("cannot acquire lock")

// ErrLockAlreadyReleased indicates that a release call cannot be fulfilled
// because the client does not hold the lock
var ErrLockAlreadyReleased = xerrors.New("lock is already released")

// ErrLockNotFound is returned for get calls on missing lock entries.
var ErrLockNotFound = xerrors.Errorf("lock not found: %w", ErrNotExist)

// Validation errors
var (
	ErrDurationTooSmall = xerrors.New("Heartbeat period must be no more than half the length of the Lease Duration, " +
		"or locks might expire due to the heartbeat thread taking too long to update them (recommendation is to make it much greater, for example " +
		"4+ times greater)")
)

func typedError(err error, msg string) error {
	const serializationErrorCode = "40001"
	if err == nil {
		return nil
	} else if err == sql.ErrNoRows {
		return xerrors.Errorf("%w: %w", xerrors.Errorf(msg+": %w", err), ErrNotExist)
	} else if _, ok := err.(*net.OpError); ok {
		return xerrors.Errorf("%w: %w", xerrors.Errorf(msg+": %w", err), ErrUnavailable)
	} else if e, ok := err.(*pq.Error); ok && e.Code == serializationErrorCode {
		return xerrors.Errorf("%w: %w", xerrors.Errorf(msg+": %w", err), ErrFailedPrecondition)
	}
	return err
}
