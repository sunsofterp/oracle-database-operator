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
	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	databasev4 "github.com/oracle/oracle-database-operator/apis/database/v4"
)

// applyLaunchNetwork copies the network-placement fields of the spec onto a
// fresh-create LaunchDbSystemDetails. These fields have always existed on
// DbSystemDetails but were only honored on the clone paths (which read them
// back from an existing DB system); a CR-driven create silently dropped them,
// so a DB system could not be pinned to a private IP, a time zone, fault
// domains, or — crucially for a database that must accept connections only
// from a known set of proxies — network security groups.
//
// Empty values are left unset so OCI applies its defaults, exactly as before.
func applyLaunchNetwork(details *database.LaunchDbSystemDetails, spec *databasev4.DbSystemDetails) {
	if details == nil || spec == nil {
		return
	}
	if len(spec.NsgIds) > 0 {
		details.NsgIds = append([]string(nil), spec.NsgIds...)
	}
	if len(spec.BackupNetworkNsgIds) > 0 {
		details.BackupNetworkNsgIds = append([]string(nil), spec.BackupNetworkNsgIds...)
	}
	if spec.TimeZone != "" {
		details.TimeZone = common.String(spec.TimeZone)
	}
	if spec.PrivateIp != "" {
		details.PrivateIp = common.String(spec.PrivateIp)
	}
	if len(spec.FaultDomains) > 0 {
		details.FaultDomains = append([]string(nil), spec.FaultDomains...)
	}
	if spec.BackupSubnetId != "" {
		details.BackupSubnetId = common.String(spec.BackupSubnetId)
	}
}
