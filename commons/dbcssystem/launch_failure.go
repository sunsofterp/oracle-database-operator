/*
** Copyright (c) 2026 Oracle and/or its affiliates.
**
** The Universal Permissive License (UPL), Version 1.0
**
** Subject to the condition set forth below, permission is hereby granted to any
** person obtaining a copy of this software, associated documentation and/or data
** (collectively the "Software"), free of charge and under any and all copyright
** rights in the Software, and any and all patent rights owned or freely
** licensable by each licensor hereunder covering either (i) the unmodified
** Software as contributed to or provided by such licensor, or (ii) the Larger
** Works (as defined below), to deal in both
**
** (a) the Software, and
** (b) any piece of software and/or hardware listed in the lrgrwrks.txt file if
** one is included with the Software (each a "Larger Work" to which the Software
** is contributed by such licensors),
**
** without restriction, including without limitation the rights to copy, create
** derivative works of, display, perform, and distribute the Software and make,
** use, sell, offer for sale, import, export, have made, and have sold the
** Software and the Larger Work(s), and to sublicense the foregoing rights on
** either these or other terms.
**
** This license is subject to the following condition:
** The above copyright notice and either this complete permission notice or at
** a minimum a reference to the UPL must be included in all copies or
** substantial portions of the Software.
**
** THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
** IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
** FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
** AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
** LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
** OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
** SOFTWARE.
 */

package common

import (
	"fmt"
	"strings"
)

// LaunchTerminalStateError reports that a DB System this reconcile LAUNCHED
// reached a terminal lifecycle state (FAILED) instead of AVAILABLE. It carries
// the OCID so the caller can terminate the wreck and launch again: an OCI Base
// DB launch fails on the provider's side often enough (host-agent upgrade
// during first boot) that FAILED must be a retry, not the end of the CR.
type LaunchTerminalStateError struct {
	ID    string
	State string
	Err   error
}

func (e *LaunchTerminalStateError) Error() string {
	return fmt.Sprintf("DB System %s reached terminal state %s during launch: %v", e.ID, e.State, e.Err)
}

func (e *LaunchTerminalStateError) Unwrap() error { return e.Err }

// isTerminalLaunchState is the set of lifecycle states a launch cannot recover
// from on its own. TERMINATING/TERMINATED are included for completeness: a
// system torn down out-of-band while the launch wait was running is just as
// gone as one that FAILED.
func isTerminalLaunchState(state string) bool {
	switch strings.ToUpper(state) {
	case "FAILED", "TERMINATED", "TERMINATING":
		return true
	default:
		return false
	}
}

// classifyLaunchWaitError wraps the outcome of the post-launch state wait: a
// terminal state becomes a LaunchTerminalStateError (retryable by the
// controller); anything else (timeout, transport) passes through unchanged.
func classifyLaunchWaitError(id, state string, err error) error {
	if err == nil {
		return nil
	}
	if isTerminalLaunchState(state) {
		return &LaunchTerminalStateError{ID: id, State: strings.ToUpper(state), Err: err}
	}
	return err
}
