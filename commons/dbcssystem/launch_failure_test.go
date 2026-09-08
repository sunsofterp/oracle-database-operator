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
	"errors"
	"testing"
)

func TestClassifyLaunchWaitErrorWrapsTerminalStates(t *testing.T) {
	t.Parallel()

	base := errors.New("DB System state FAILED is not matching expected state AVAILABLE")
	for _, state := range []string{"FAILED", "failed", "TERMINATED", "TERMINATING"} {
		err := classifyLaunchWaitError("ocid1.dbsystem.oc1.iad.x", state, base)
		var terminal *LaunchTerminalStateError
		if !errors.As(err, &terminal) {
			t.Fatalf("state %q: expected LaunchTerminalStateError, got %T (%v)", state, err, err)
		}
		if terminal.ID != "ocid1.dbsystem.oc1.iad.x" {
			t.Fatalf("state %q: ID = %q", state, terminal.ID)
		}
		if terminal.State != "FAILED" && terminal.State != "TERMINATED" && terminal.State != "TERMINATING" {
			t.Fatalf("state %q: normalized State = %q", state, terminal.State)
		}
		if !errors.Is(err, base) {
			t.Fatalf("state %q: wrapped error must unwrap to the wait error", state)
		}
	}
}

func TestClassifyLaunchWaitErrorPassesThroughNonTerminal(t *testing.T) {
	t.Parallel()

	base := errors.New("timed out after 4h0m0s waiting for DB System x to reach state AVAILABLE")
	for _, state := range []string{"PROVISIONING", "", "UPDATING"} {
		err := classifyLaunchWaitError("ocid1.dbsystem.oc1.iad.x", state, base)
		if err != base {
			t.Fatalf("state %q: expected the wait error unchanged, got %v", state, err)
		}
	}
	if err := classifyLaunchWaitError("ocid1.dbsystem.oc1.iad.x", "FAILED", nil); err != nil {
		t.Fatalf("nil error must stay nil, got %v", err)
	}
}
