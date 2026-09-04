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
	"reflect"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/database"

	databasev4 "github.com/oracle/oracle-database-operator/apis/database/v4"
)

func TestApplyLaunchNetworkCopiesEverySetField(t *testing.T) {
	spec := &databasev4.DbSystemDetails{
		NsgIds:              []string{"ocid1.networksecuritygroup.oc1..a"},
		BackupNetworkNsgIds: []string{"ocid1.networksecuritygroup.oc1..b"},
		TimeZone:            "UTC",
		PrivateIp:           "10.21.0.5",
		FaultDomains:        []string{"FAULT-DOMAIN-1"},
		BackupSubnetId:      "ocid1.subnet.oc1..backup",
	}
	var d database.LaunchDbSystemDetails
	applyLaunchNetwork(&d, spec)
	if !reflect.DeepEqual(d.NsgIds, spec.NsgIds) || !reflect.DeepEqual(d.BackupNetworkNsgIds, spec.BackupNetworkNsgIds) || !reflect.DeepEqual(d.FaultDomains, spec.FaultDomains) {
		t.Errorf("slices not copied: %+v", d)
	}
	if d.TimeZone == nil || *d.TimeZone != "UTC" || d.PrivateIp == nil || *d.PrivateIp != "10.21.0.5" || d.BackupSubnetId == nil || *d.BackupSubnetId != spec.BackupSubnetId {
		t.Errorf("scalars not copied: %+v", d)
	}
	// Copies, not aliases: mutating the spec must not change the request.
	spec.NsgIds[0] = "changed"
	if d.NsgIds[0] == "changed" {
		t.Error("NsgIds aliased the spec slice")
	}
}

func TestApplyLaunchNetworkLeavesUnsetFieldsNil(t *testing.T) {
	var d database.LaunchDbSystemDetails
	applyLaunchNetwork(&d, &databasev4.DbSystemDetails{})
	if d.NsgIds != nil || d.BackupNetworkNsgIds != nil || d.TimeZone != nil || d.PrivateIp != nil || d.FaultDomains != nil || d.BackupSubnetId != nil {
		t.Errorf("unset spec fields must stay nil so OCI defaults apply: %+v", d)
	}
	applyLaunchNetwork(nil, nil) // must not panic
}
