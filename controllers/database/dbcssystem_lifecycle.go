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
	"context"
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"

	databasev4 "github.com/oracle/oracle-database-operator/apis/database/v4"
	dbcsv4 "github.com/oracle/oracle-database-operator/commons/dbcssystem"
)

// maxLaunchAttempts bounds the terminate-and-relaunch loop a FAILED launch
// enters. Three covers the provider-side flake this exists for (a host-agent
// upgrade failing during first boot) without hiding a launch that fails for
// a reason of ours — after the third failure the CR stays FAILED with the
// last error in status.message.
const maxLaunchAttempts = 3

// dbSystemIDForDeletion resolves the OCID a deleting CR must terminate. The
// operator records the launched OCID in three places with different
// lifetimes: spec.id (in-memory during the launch; a declarative owner that
// re-applies the manifest strips it), the lastSuccessfulSpec annotation
// (written only once the launch reached AVAILABLE), and status.id (written by
// every status sync, including the FAILED one). Any of them is enough to
// terminate; none means nothing was ever launched.
func dbSystemIDForDeletion(dbcs *databasev4.DbcsSystem) (string, bool) {
	if dbcs == nil {
		return "", false
	}
	if dbcs.Spec.Id != nil && *dbcs.Spec.Id != "" {
		return *dbcs.Spec.Id, true
	}
	if last, err := dbcs.GetLastSuccessfulSpec(); err == nil && last != nil && last.Id != nil && *last.Id != "" {
		return *last.Id, true
	}
	if dbcs.Status.Id != nil && *dbcs.Status.Id != "" {
		return *dbcs.Status.Id, true
	}
	return "", false
}

// isGoneOrGoing reports lifecycle states a terminate request must not be sent
// for: the system is already being torn down, or is gone.
func isGoneOrGoing(state string) bool {
	switch strings.ToUpper(state) {
	case "TERMINATED", "TERMINATING":
		return true
	default:
		return false
	}
}

// launchAction is what the controller does after a failed launch attempt.
type launchAction int

const (
	// launchFail leaves the CR FAILED: the error is not a terminal OCI state
	// (spec validation, a rejected launch request, a wait timeout) or the
	// attempt budget is spent.
	launchFail launchAction = iota
	// launchRetry terminates the FAILED system and requeues so the next
	// reconcile launches again.
	launchRetry
)

// decideLaunchAction classifies a launch error against the attempt budget.
// attemptsSoFar counts failed launches recorded on the CR BEFORE this one.
// It returns the terminal-state error when the failure is one, so the caller
// has the OCID to terminate.
func decideLaunchAction(err error, attemptsSoFar int) (launchAction, *dbcsv4.LaunchTerminalStateError) {
	var terminal *dbcsv4.LaunchTerminalStateError
	if !errors.As(err, &terminal) {
		return launchFail, nil
	}
	if attemptsSoFar+1 >= maxLaunchAttempts {
		return launchFail, terminal
	}
	return launchRetry, terminal
}

// recordLaunchAttempt persists the retry bookkeeping on the status
// subresource: the attempt count, the phase the CR is in while the relaunch
// is pending, and a human-readable message. It patches against the latest
// copy so the write cannot race the status sync a parallel reconcile made.
func (r *DbcsSystemReconciler) recordLaunchAttempt(ctx context.Context, dbcs *databasev4.DbcsSystem, attempts int, state databasev4.LifecycleState, message string) error {
	latest := &databasev4.DbcsSystem{}
	if err := r.KubeClient.Get(ctx, client.ObjectKeyFromObject(dbcs), latest); err != nil {
		return fmt.Errorf("failed to fetch the latest DbcsSystem before recording the launch attempt: %w", err)
	}
	updated := latest.DeepCopy()
	updated.Status.LaunchAttempts = attempts
	updated.Status.State = state
	updated.Status.Message = message
	if err := r.KubeClient.Status().Patch(ctx, updated, client.MergeFrom(latest)); err != nil {
		return fmt.Errorf("failed to record the launch attempt: %w", err)
	}
	dbcs.Status.LaunchAttempts = attempts
	dbcs.Status.State = state
	dbcs.Status.Message = message
	return nil
}
