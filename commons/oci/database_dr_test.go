/*
** Copyright (c) 2022 Oracle and/or its affiliates.
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

package oci

import (
	"reflect"
	"testing"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	dbv4 "github.com/oracle/oracle-database-operator/apis/database/v4"
)

const (
	testStandbyOCID = "ocid1.autonomousdatabase.oc1.phx.standby"
	testPrimaryOCID = "ocid1.autonomousdatabase.oc1.iad.primary"
)

func TestIsCrossRegionDRPeer(t *testing.T) {
	fullDR := &dbv4.DisasterRecoverySpec{SourceId: common.String(testPrimaryOCID), Type: database.DisasterRecoveryConfigurationDisasterRecoveryTypeBackupBased}
	cases := []struct {
		name string
		id   *string
		dr   *dbv4.DisasterRecoverySpec
		want bool
	}{
		{name: "nil disasterRecovery", dr: nil, want: false},
		{name: "disasterRecovery without sourceId", dr: &dbv4.DisasterRecoverySpec{Type: database.DisasterRecoveryConfigurationDisasterRecoveryTypeAdg}, want: false},
		{name: "empty sourceId", dr: &dbv4.DisasterRecoverySpec{SourceId: common.String("")}, want: false},
		{name: "sourceId set without type", dr: &dbv4.DisasterRecoverySpec{SourceId: common.String(testPrimaryOCID)}, want: false},
		{name: "sourceId and type set", dr: fullDR, want: true},
		{name: "sourceId and type set but id already populated", id: common.String(testStandbyOCID), dr: fullDR, want: false},
		{name: "sourceId and type set with empty id", id: common.String(""), dr: fullDR, want: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			adb := &dbv4.AutonomousDatabase{}
			adb.Spec.Details.Id = tc.id
			adb.Spec.Details.DisasterRecovery = tc.dr
			if got := isCrossRegionDRPeer(adb); got != tc.want {
				t.Fatalf("isCrossRegionDRPeer() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestValidateDisasterRecovery(t *testing.T) {
	cases := []struct {
		name    string
		dr      *dbv4.DisasterRecoverySpec
		wantErr bool
	}{
		{name: "nil disasterRecovery", dr: nil, wantErr: false},
		{name: "empty struct", dr: &dbv4.DisasterRecoverySpec{}, wantErr: false},
		{name: "sourceId without type", dr: &dbv4.DisasterRecoverySpec{SourceId: common.String(testPrimaryOCID)}, wantErr: true},
		{name: "type without sourceId", dr: &dbv4.DisasterRecoverySpec{Type: database.DisasterRecoveryConfigurationDisasterRecoveryTypeAdg}, wantErr: true},
		{name: "empty sourceId with type", dr: &dbv4.DisasterRecoverySpec{SourceId: common.String(""), Type: database.DisasterRecoveryConfigurationDisasterRecoveryTypeAdg}, wantErr: true},
		{name: "both set", dr: &dbv4.DisasterRecoverySpec{SourceId: common.String(testPrimaryOCID), Type: database.DisasterRecoveryConfigurationDisasterRecoveryTypeAdg}, wantErr: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			adb := &dbv4.AutonomousDatabase{}
			adb.Spec.Details.DisasterRecovery = tc.dr
			err := validateDisasterRecovery(adb)
			if (err != nil) != tc.wantErr {
				t.Fatalf("validateDisasterRecovery() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestCrossRegionDRDetails(t *testing.T) {
	adb := &dbv4.AutonomousDatabase{}
	adb.Spec.Details.CompartmentId = common.String("ocid1.compartment.oc1..peer")
	adb.Spec.Details.DisplayName = common.String("erp-standby")
	adb.Spec.Details.DbName = common.String("erpstby")
	adb.Spec.Details.LicenseModel = database.AutonomousDatabaseLicenseModelBringYourOwnLicense
	adb.Spec.Details.DatabaseEdition = database.AutonomousDatabaseDatabaseEditionStandardEdition
	adb.Spec.Details.ComputeModel = database.AutonomousDatabaseComputeModelEcpu
	adb.Spec.Details.ComputeCount = common.Float32(4)
	adb.Spec.Details.DataStorageSizeInTBs = common.Int(1)
	adb.Spec.Details.IsAccessControlEnabled = common.Bool(true)
	adb.Spec.Details.WhitelistedIps = []string{"10.0.0.0/8"}
	adb.Spec.Details.SubnetId = common.String("ocid1.subnet.oc1.phx.peer")
	adb.Spec.Details.NsgIds = []string{"ocid1.networksecuritygroup.oc1.phx.peer"}
	adb.Spec.Details.PrivateEndpointLabel = common.String("erpstby")
	adb.Spec.Details.IsMtlsConnectionRequired = common.Bool(false)
	adb.Spec.Details.FreeformTags = map[string]string{"tenant": "acme"}
	// Create-time only fields that must NOT leak into the peer request.
	adb.Spec.Details.AdminPassword.K8sSecret.Name = common.String("admin-password")
	adb.Spec.Details.DisasterRecovery = &dbv4.DisasterRecoverySpec{
		SourceId:                    common.String(testPrimaryOCID),
		Type:                        database.DisasterRecoveryConfigurationDisasterRecoveryTypeAdg,
		IsReplicateAutomaticBackups: common.Bool(true),
	}

	got := crossRegionDRDetails(adb)

	want := database.CreateCrossRegionDisasterRecoveryDetails{
		CompartmentId:               common.String("ocid1.compartment.oc1..peer"),
		SourceId:                    common.String(testPrimaryOCID),
		RemoteDisasterRecoveryType:  database.DisasterRecoveryConfigurationDisasterRecoveryTypeAdg,
		IsReplicateAutomaticBackups: common.Bool(true),
		DisplayName:                 common.String("erp-standby"),
		DbName:                      common.String("erpstby"),
		LicenseModel:                database.CreateAutonomousDatabaseBaseLicenseModelBringYourOwnLicense,
		DatabaseEdition:             database.AutonomousDatabaseSummaryDatabaseEditionStandardEdition,
		ComputeModel:                database.CreateAutonomousDatabaseBaseComputeModelEcpu,
		ComputeCount:                common.Float32(4),
		DataStorageSizeInTBs:        common.Int(1),
		IsAccessControlEnabled:      common.Bool(true),
		WhitelistedIps:              []string{"10.0.0.0/8"},
		SubnetId:                    common.String("ocid1.subnet.oc1.phx.peer"),
		NsgIds:                      []string{"ocid1.networksecuritygroup.oc1.phx.peer"},
		PrivateEndpointLabel:        common.String("erpstby"),
		IsMtlsConnectionRequired:    common.Bool(false),
		FreeformTags:                map[string]string{"tenant": "acme"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("crossRegionDRDetails() mismatch\n got: %#v\nwant: %#v", got, want)
	}
	if got.AdminPassword != nil {
		t.Fatalf("crossRegionDRDetails() must not send an admin password for a DR peer, got %q", *got.AdminPassword)
	}

	// The DR body must be usable where the operator sends create requests.
	var _ database.CreateAutonomousDatabaseBase = got
}

func TestCrossRegionDRDetailsBackupBasedMinimal(t *testing.T) {
	adb := &dbv4.AutonomousDatabase{}
	adb.Spec.Details.CompartmentId = common.String("ocid1.compartment.oc1..peer")
	adb.Spec.Details.DisasterRecovery = &dbv4.DisasterRecoverySpec{
		SourceId: common.String(testPrimaryOCID),
		Type:     database.DisasterRecoveryConfigurationDisasterRecoveryTypeBackupBased,
	}

	got := crossRegionDRDetails(adb)

	if got.SourceId == nil || *got.SourceId != testPrimaryOCID {
		t.Fatalf("SourceId = %v, want %q", got.SourceId, testPrimaryOCID)
	}
	if got.RemoteDisasterRecoveryType != database.DisasterRecoveryConfigurationDisasterRecoveryTypeBackupBased {
		t.Fatalf("RemoteDisasterRecoveryType = %q, want BACKUP_BASED", got.RemoteDisasterRecoveryType)
	}
	if got.IsReplicateAutomaticBackups != nil {
		t.Fatalf("IsReplicateAutomaticBackups should stay nil when unset, got %v", *got.IsReplicateAutomaticBackups)
	}
	// Unset enum-typed spec fields must map to the SDK zero value so OCI applies its defaults.
	if got.LicenseModel != "" || got.DatabaseEdition != "" || got.ComputeModel != "" {
		t.Fatalf("unset enums leaked into the request: licenseModel=%q databaseEdition=%q computeModel=%q",
			got.LicenseModel, got.DatabaseEdition, got.ComputeModel)
	}
}

func TestSwitchoverRequestPeerDbId(t *testing.T) {
	cases := []struct {
		name     string
		peerDbId string
		want     *string
	}{
		{name: "local Data Guard: no peerDbId", peerDbId: "", want: nil},
		{name: "cross-region: peerDbId is the primary", peerDbId: testPrimaryOCID, want: common.String(testPrimaryOCID)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := switchoverRequest(testStandbyOCID, tc.peerDbId)

			if got.AutonomousDatabaseId == nil || *got.AutonomousDatabaseId != testStandbyOCID {
				t.Fatalf("AutonomousDatabaseId = %v, want %q", got.AutonomousDatabaseId, testStandbyOCID)
			}
			if !reflect.DeepEqual(got.PeerDbId, tc.want) {
				t.Fatalf("PeerDbId = %v, want %v", got.PeerDbId, tc.want)
			}
			if got.RequestMetadata.RetryPolicy == nil {
				t.Fatal("RetryPolicy must be set, matching the other ADB requests")
			}
		})
	}
}

func TestFailoverRequestPeerDbId(t *testing.T) {
	cases := []struct {
		name     string
		peerDbId string
		want     *string
	}{
		{name: "local Data Guard: no peerDbId", peerDbId: "", want: nil},
		{name: "cross-region: peerDbId is the primary", peerDbId: testPrimaryOCID, want: common.String(testPrimaryOCID)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := failoverRequest(testStandbyOCID, tc.peerDbId)

			if got.AutonomousDatabaseId == nil || *got.AutonomousDatabaseId != testStandbyOCID {
				t.Fatalf("AutonomousDatabaseId = %v, want %q", got.AutonomousDatabaseId, testStandbyOCID)
			}
			if !reflect.DeepEqual(got.PeerDbId, tc.want) {
				t.Fatalf("PeerDbId = %v, want %v", got.PeerDbId, tc.want)
			}
			if got.RequestMetadata.RetryPolicy == nil {
				t.Fatal("RetryPolicy must be set, matching the other ADB requests")
			}
		})
	}
}
