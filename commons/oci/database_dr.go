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
	"errors"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/database"

	dbv4 "github.com/oracle/oracle-database-operator/apis/database/v4"
)

// isCrossRegionDRPeer reports whether the CR asks for this database to be
// created as the cross-region disaster-recovery peer of another database:
// the CR must not already name a database (spec.details.id — once the peer
// exists the operator writes its OCID back and disasterRecovery stays as a
// record of how it was created), and both sourceId and type must be set
// (validateDisasterRecovery rejects the half-specified case before this is
// consulted).
func isCrossRegionDRPeer(adb *dbv4.AutonomousDatabase) bool {
	if id := adb.Spec.Details.Id; id != nil && *id != "" {
		return false
	}
	dr := adb.Spec.Details.DisasterRecovery
	return dr != nil && dr.SourceId != nil && *dr.SourceId != "" && dr.Type != ""
}

// validateDisasterRecovery rejects a spec.details.disasterRecovery that names
// a sourceId without a type (or a type without a sourceId): OCI would answer
// 400 to the former and the latter would silently create a STANDALONE
// database, so neither may reach the create path.
func validateDisasterRecovery(adb *dbv4.AutonomousDatabase) error {
	dr := adb.Spec.Details.DisasterRecovery
	if dr == nil {
		return nil
	}
	hasSource := dr.SourceId != nil && *dr.SourceId != ""
	switch {
	case hasSource && dr.Type == "":
		return errors.New("spec.details.disasterRecovery.type is required when sourceId is set (BACKUP_BASED or ADG)")
	case !hasSource && dr.Type != "":
		return errors.New("spec.details.disasterRecovery.sourceId is required when type is set")
	}
	return nil
}

// crossRegionDRDetails maps spec.details of a CR onto the OCI request body that
// creates a cross-region disaster-recovery peer of spec.details.disasterRecovery.sourceId.
//
// Only the fields the OCI API accepts on a peer are mapped. The peer inherits
// the primary's admin password, workload, version and edition-defining
// properties from the source, so no admin password is sent.
func crossRegionDRDetails(adb *dbv4.AutonomousDatabase) database.CreateCrossRegionDisasterRecoveryDetails {
	details := adb.Spec.Details
	dr := details.DisasterRecovery

	out := database.CreateCrossRegionDisasterRecoveryDetails{
		CompartmentId: details.CompartmentId,
		DisplayName:   details.DisplayName,
		DbName:        details.DbName,

		CpuCoreCount:         details.CpuCoreCount,
		ComputeModel:         database.CreateAutonomousDatabaseBaseComputeModelEnum(details.ComputeModel),
		ComputeCount:         details.ComputeCount,
		OcpuCount:            details.OcpuCount,
		DataStorageSizeInTBs: details.DataStorageSizeInTBs,
		LicenseModel:         database.CreateAutonomousDatabaseBaseLicenseModelEnum(details.LicenseModel),
		DatabaseEdition:      database.AutonomousDatabaseSummaryDatabaseEditionEnum(details.DatabaseEdition),

		IsAccessControlEnabled:   details.IsAccessControlEnabled,
		WhitelistedIps:           details.WhitelistedIps,
		SubnetId:                 details.SubnetId,
		NsgIds:                   details.NsgIds,
		PrivateEndpointLabel:     details.PrivateEndpointLabel,
		IsMtlsConnectionRequired: details.IsMtlsConnectionRequired,

		FreeformTags: details.FreeformTags,
	}

	if dr != nil {
		out.SourceId = dr.SourceId
		out.RemoteDisasterRecoveryType = dr.Type
		out.IsReplicateAutomaticBackups = dr.IsReplicateAutomaticBackups
	}

	return out
}

// switchoverRequest builds the Switchover request for adbOCID. peerDbId is
// sent only when non-empty. For a cross-region pair the request must be made
// on the STANDBY (adbOCID) with peerDbId set to the PRIMARY's OCID.
func switchoverRequest(adbOCID string, peerDbId string) database.SwitchoverAutonomousDatabaseRequest {
	retryPolicy := common.DefaultRetryPolicy()

	request := database.SwitchoverAutonomousDatabaseRequest{
		AutonomousDatabaseId: common.String(adbOCID),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: &retryPolicy,
		},
	}
	if peerDbId != "" {
		request.PeerDbId = common.String(peerDbId)
	}
	return request
}

// failoverRequest builds the Failover request for adbOCID. peerDbId is sent
// only when non-empty; see switchoverRequest for the cross-region rule.
func failoverRequest(adbOCID string, peerDbId string) database.FailOverAutonomousDatabaseRequest {
	retryPolicy := common.DefaultRetryPolicy()

	request := database.FailOverAutonomousDatabaseRequest{
		AutonomousDatabaseId: common.String(adbOCID),
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: &retryPolicy,
		},
	}
	if peerDbId != "" {
		request.PeerDbId = common.String(peerDbId)
	}
	return request
}
