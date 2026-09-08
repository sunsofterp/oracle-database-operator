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

package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	databasev4 "github.com/oracle/oracle-database-operator/apis/database/v4"
	dbcsv4 "github.com/oracle/oracle-database-operator/commons/dbcssystem"
)

func strp(s string) *string { return &s }

func withLastSuccessfulSpec(t *testing.T, dbcs *databasev4.DbcsSystem, id string) {
	t.Helper()
	spec := databasev4.DbcsSystemSpec{Id: strp(id)}
	raw, err := json.Marshal(spec)
	if err != nil {
		t.Fatal(err)
	}
	dbcs.SetAnnotations(map[string]string{"lastSuccessfulSpec": string(raw)})
}

func TestDbSystemIDForDeletionPrecedence(t *testing.T) {
	t.Parallel()

	t.Run("nothing recorded means nothing to terminate", func(t *testing.T) {
		dbcs := &databasev4.DbcsSystem{}
		if id, ok := dbSystemIDForDeletion(dbcs); ok || id != "" {
			t.Fatalf("got (%q, %v), want (\"\", false)", id, ok)
		}
		dbcs.Spec.Id = strp("")
		dbcs.Status.Id = strp("")
		if _, ok := dbSystemIDForDeletion(dbcs); ok {
			t.Fatal("empty strings must not count as an OCID")
		}
		if _, ok := dbSystemIDForDeletion(nil); ok {
			t.Fatal("nil CR must not resolve")
		}
	})

	t.Run("status.id alone is enough (FAILED launch, spec.id stripped)", func(t *testing.T) {
		dbcs := &databasev4.DbcsSystem{}
		dbcs.Status.Id = strp("ocid1.dbsystem.oc1.iad.status")
		id, ok := dbSystemIDForDeletion(dbcs)
		if !ok || id != "ocid1.dbsystem.oc1.iad.status" {
			t.Fatalf("got (%q, %v)", id, ok)
		}
	})

	t.Run("lastSuccessfulSpec beats status.id", func(t *testing.T) {
		dbcs := &databasev4.DbcsSystem{}
		dbcs.Status.Id = strp("ocid1.dbsystem.oc1.iad.status")
		withLastSuccessfulSpec(t, dbcs, "ocid1.dbsystem.oc1.iad.last")
		id, ok := dbSystemIDForDeletion(dbcs)
		if !ok || id != "ocid1.dbsystem.oc1.iad.last" {
			t.Fatalf("got (%q, %v)", id, ok)
		}
	})

	t.Run("spec.id beats everything", func(t *testing.T) {
		dbcs := &databasev4.DbcsSystem{}
		dbcs.Spec.Id = strp("ocid1.dbsystem.oc1.iad.spec")
		dbcs.Status.Id = strp("ocid1.dbsystem.oc1.iad.status")
		withLastSuccessfulSpec(t, dbcs, "ocid1.dbsystem.oc1.iad.last")
		id, ok := dbSystemIDForDeletion(dbcs)
		if !ok || id != "ocid1.dbsystem.oc1.iad.spec" {
			t.Fatalf("got (%q, %v)", id, ok)
		}
	})

	t.Run("a corrupt lastSuccessfulSpec falls through to status.id", func(t *testing.T) {
		dbcs := &databasev4.DbcsSystem{}
		dbcs.SetAnnotations(map[string]string{"lastSuccessfulSpec": "{not json"})
		dbcs.Status.Id = strp("ocid1.dbsystem.oc1.iad.status")
		id, ok := dbSystemIDForDeletion(dbcs)
		if !ok || id != "ocid1.dbsystem.oc1.iad.status" {
			t.Fatalf("got (%q, %v)", id, ok)
		}
	})
}

func TestIsGoneOrGoing(t *testing.T) {
	t.Parallel()
	for state, want := range map[string]bool{
		"TERMINATED": true, "terminating": true, "AVAILABLE": false, "FAILED": false, "": false,
	} {
		if got := isGoneOrGoing(state); got != want {
			t.Fatalf("%q: got %v want %v", state, got, want)
		}
	}
}

func TestDecideLaunchAction(t *testing.T) {
	t.Parallel()

	terminal := &dbcsv4.LaunchTerminalStateError{ID: "ocid1.dbsystem.oc1.iad.x", State: "FAILED", Err: errors.New("state FAILED")}

	t.Run("non-terminal errors are never retried", func(t *testing.T) {
		for _, err := range []error{errors.New("spec validation failed"), fmt.Errorf("wrapped: %w", errors.New("timeout"))} {
			if action, got := decideLaunchAction(err, 0); action != launchFail || got != nil {
				t.Fatalf("%v: got (%v, %v)", err, action, got)
			}
		}
	})

	t.Run("a terminal state is retried until the budget is spent", func(t *testing.T) {
		for attempts := 0; attempts+1 < maxLaunchAttempts; attempts++ {
			action, got := decideLaunchAction(fmt.Errorf("launch: %w", terminal), attempts)
			if action != launchRetry || got == nil || got.ID != terminal.ID {
				t.Fatalf("attempts=%d: got (%v, %v)", attempts, action, got)
			}
		}
		action, got := decideLaunchAction(terminal, maxLaunchAttempts-1)
		if action != launchFail || got == nil {
			t.Fatalf("last attempt: got (%v, %v), want (launchFail, the terminal error)", action, got)
		}
		if action, _ := decideLaunchAction(terminal, maxLaunchAttempts+5); action != launchFail {
			t.Fatal("over budget must not retry")
		}
	})
}
