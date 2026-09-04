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

package v4

// revive:disable:exported,var-naming
// Legacy API field/type names are preserved for backward compatibility.

import (
	"encoding/json"
	"reflect"

	"github.com/oracle/oci-go-sdk/v65/database"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AutonomousDatabaseSpec defines the desired state of AutonomousDatabase
// Important: Run "make" to regenerate code after modifying this file
type AutonomousDatabaseSpec struct {
	// +kubebuilder:validation:Enum:="";Create;Sync;Update;Stop;Start;Terminate;Clone;Switchover;Failover
	Action    string                    `json:"action"`
	Details   AutonomousDatabaseDetails `json:"details,omitempty"`
	Clone     AutonomousDatabaseClone   `json:"clone,omitempty"`
	Wallet    WalletSpec                `json:"wallet,omitempty"`
	OciConfig OciConfigSpec             `json:"ociConfig,omitempty"`
	// +kubebuilder:default:=false
	HardLink *bool `json:"hardLink,omitempty"`
}

type AutonomousDatabaseDetails struct {
	AutonomousDatabaseBase `json:",inline"`
	Id                     *string `json:"id,omitempty"`
}

type AutonomousDatabaseClone struct {
	AutonomousDatabaseBase `json:",inline"`
	// +kubebuilder:validation:Enum:="FULL";"METADATA"
	CloneType database.CreateAutonomousDatabaseCloneDetailsCloneTypeEnum `json:"cloneType,omitempty"`
}

// AutonomousDatabaseBase defines the detail information of AutonomousDatabase, corresponding to oci-go-sdk/database/AutonomousDatabase
type AutonomousDatabaseBase struct {
	CompartmentId               *string `json:"compartmentId,omitempty"`
	AutonomousContainerDatabase AcdSpec `json:"autonomousContainerDatabase,omitempty"`
	DisplayName                 *string `json:"displayName,omitempty"`
	DbName                      *string `json:"dbName,omitempty"`
	// +kubebuilder:validation:Enum:="OLTP";"DW";"AJD";"APEX";"LH"
	DbWorkload database.AutonomousDatabaseDbWorkloadEnum `json:"dbWorkload,omitempty"`
	// +kubebuilder:validation:Enum:="LICENSE_INCLUDED";"BRING_YOUR_OWN_LICENSE"
	LicenseModel database.AutonomousDatabaseLicenseModelEnum `json:"licenseModel,omitempty"`
	// DatabaseEdition selects the Oracle Database edition for a BRING_YOUR_OWN_LICENSE
	// Autonomous Database Serverless (Standard Edition 2 or Enterprise Edition).
	// Ignored by OCI unless licenseModel is BRING_YOUR_OWN_LICENSE.
	// +kubebuilder:validation:Enum:="STANDARD_EDITION";"ENTERPRISE_EDITION"
	DatabaseEdition      database.AutonomousDatabaseDatabaseEditionEnum `json:"databaseEdition,omitempty"`
	DbVersion            *string                                        `json:"dbVersion,omitempty"`
	DataStorageSizeInTBs *int                                           `json:"dataStorageSizeInTBs,omitempty"`
	CpuCoreCount         *int                                           `json:"cpuCoreCount,omitempty"`
	// +kubebuilder:validation:Enum:="ECPU";"OCPU"
	ComputeModel         database.AutonomousDatabaseComputeModelEnum `json:"computeModel,omitempty"`
	ComputeCount         *float32                                    `json:"computeCount,omitempty"`
	OcpuCount            *float32                                    `json:"ocpuCount,omitempty"`
	AdminPassword        PasswordSpec                                `json:"adminPassword,omitempty"`
	IsAutoScalingEnabled *bool                                       `json:"isAutoScalingEnabled,omitempty"`
	IsDedicated          *bool                                       `json:"isDedicated,omitempty"`
	IsFreeTier           *bool                                       `json:"isFreeTier,omitempty"`

	// NetworkAccess
	IsAccessControlEnabled   *bool    `json:"isAccessControlEnabled,omitempty"`
	WhitelistedIps           []string `json:"whitelistedIps,omitempty"`
	SubnetId                 *string  `json:"subnetId,omitempty"`
	NsgIds                   []string `json:"nsgIds,omitempty"`
	PrivateEndpointLabel     *string  `json:"privateEndpointLabel,omitempty"`
	IsMtlsConnectionRequired *bool    `json:"isMtlsConnectionRequired,omitempty"`

	FreeformTags map[string]string `json:"freeformTags,omitempty"`

	// PeerDbId is the OCID of the Data Guard / disaster-recovery peer used by
	// the Switchover and Failover actions. A cross-region switchover must be
	// invoked on the STANDBY database (spec.details.id) with peerDbId set to
	// the PRIMARY's OCID; OCI rejects the request when it is sent to the primary.
	PeerDbId *string `json:"peerDbId,omitempty"`

	// DisasterRecovery, when set on a CR that has no spec.details.id, makes the
	// operator create this database as the cross-region disaster-recovery PEER of
	// SourceId (CreateCrossRegionDisasterRecoveryDetails) instead of a new
	// standalone database. The OCI client must target the PEER's (this CR's)
	// region — see spec.ociConfig / the region ConfigMap.
	DisasterRecovery *DisasterRecoverySpec `json:"disasterRecovery,omitempty"`
}

// DisasterRecoverySpec describes the cross-region disaster-recovery peer to create,
// corresponding to oci-go-sdk/database/CreateCrossRegionDisasterRecoveryDetails.
type DisasterRecoverySpec struct {
	// SourceId is the OCID of the primary (source) Autonomous Database.
	// +kubebuilder:validation:Required
	SourceId *string `json:"sourceId"`
	// Type selects the disaster-recovery technology: Autonomous Data Guard (ADG)
	// or backup-based (BACKUP_BASED). Required with SourceId: OCI rejects a
	// peer create without it, so the schema rejects it first.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum:="BACKUP_BASED";"ADG"
	Type database.DisasterRecoveryConfigurationDisasterRecoveryTypeEnum `json:"type"`
	// IsReplicateAutomaticBackups replicates the primary's automatic backups to the peer region.
	IsReplicateAutomaticBackups *bool `json:"isReplicateAutomaticBackups,omitempty"`
}

/************************
*	ACD specs
************************/
type K8sAcdSpec struct {
	Name *string `json:"name,omitempty"`
}

type OciAcdSpec struct {
	Id *string `json:"id,omitempty"`
}

// AcdSpec defines the spec of the target for backup/restore runs.
// The name could be the name of an AutonomousDatabase or an AutonomousDatabaseBackup
type AcdSpec struct {
	K8sAcd K8sAcdSpec `json:"k8sAcd,omitempty"`
	OciAcd OciAcdSpec `json:"ociAcd,omitempty"`
}

/************************
*	Secret specs
************************/
type K8sSecretSpec struct {
	Name *string `json:"name,omitempty"`
}

type OciSecretSpec struct {
	Id *string `json:"id,omitempty"`
}

type PasswordSpec struct {
	K8sSecret K8sSecretSpec `json:"k8sSecret,omitempty"`
	OciSecret OciSecretSpec `json:"ociSecret,omitempty"`
}

type WalletSpec struct {
	Name     *string      `json:"name,omitempty"`
	Password PasswordSpec `json:"password,omitempty"`
}

// AutonomousDatabaseStatus defines the observed state of AutonomousDatabase
type AutonomousDatabaseStatus struct {
	// Lifecycle State of the ADB
	LifecycleState database.AutonomousDatabaseLifecycleStateEnum `json:"lifecycleState,omitempty"`
	// Creation time of the ADB
	TimeCreated string `json:"timeCreated,omitempty"`
	// Expiring date of the instance wallet
	WalletExpiringDate string `json:"walletExpiringDate,omitempty"`
	// Connection Strings of the ADB
	AllConnectionStrings []ConnectionStringProfile `json:"allConnectionStrings,omitempty"`
	// Data Guard / disaster-recovery role of the ADB (PRIMARY, STANDBY, ...)
	Role string `json:"role,omitempty"`
	// OCIDs of the Data Guard / disaster-recovery peers of the ADB
	PeerDbIds []string `json:"peerDbIds,omitempty"`
	// Whether the ADB is the PRIMARY or the REMOTE side of a cross-region disaster-recovery pair
	DisasterRecoveryRegionType string `json:"disasterRecoveryRegionType,omitempty"`
	// Whether the ADB is in the PRIMARY_DG_REGION or the REMOTE_STANDBY_DG_REGION of a cross-region Data Guard pair
	DataguardRegionType string `json:"dataguardRegionType,omitempty"`
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

type TLSAuthenticationEnum string

const (
	tlsAuthenticationTLS  TLSAuthenticationEnum = "TLS"
	tlsAuthenticationMTLS TLSAuthenticationEnum = "Mutual TLS"
)

func GetTLSAuthenticationEnumFromString(val string) (TLSAuthenticationEnum, bool) {
	var mappingTLSAuthenticationEnum = map[string]TLSAuthenticationEnum{
		"TLS":        tlsAuthenticationTLS,
		"Mutual TLS": tlsAuthenticationMTLS,
	}

	enum, ok := mappingTLSAuthenticationEnum[val]
	return enum, ok
}

type ConnectionStringProfile struct {
	TLSAuthentication TLSAuthenticationEnum  `json:"tlsAuthentication,omitempty"`
	ConnectionStrings []ConnectionStringSpec `json:"connectionStrings"`
}

type ConnectionStringSpec struct {
	TNSName          string `json:"tnsName,omitempty"`
	ConnectionString string `json:"connectionString,omitempty"`
}

// AutonomousDatabase is the Schema for the autonomousdatabases API
// +kubebuilder:object:root=true
// +kubebuilder:resource:shortName="adb";"adbs"
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:JSONPath=".spec.details.displayName",name="Display Name",type=string
// +kubebuilder:printcolumn:JSONPath=".spec.details.dbName",name="Db Name",type=string
// +kubebuilder:printcolumn:JSONPath=".status.lifecycleState",name="State",type=string
// +kubebuilder:printcolumn:JSONPath=".spec.details.isDedicated",name="Dedicated",type=string
// +kubebuilder:printcolumn:JSONPath=".spec.details.computeCount",name="Compute Count",type=number
// +kubebuilder:printcolumn:JSONPath=".spec.details.dataStorageSizeInTBs",name="Storage (TB)",type=integer
// +kubebuilder:printcolumn:JSONPath=".spec.details.dbWorkload",name="Workload Type",type=string
// +kubebuilder:printcolumn:JSONPath=".status.timeCreated",name="Created",type=string
// +kubebuilder:storageversion
type AutonomousDatabase struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutonomousDatabaseSpec   `json:"spec,omitempty"`
	Status AutonomousDatabaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutonomousDatabaseList contains a list of AutonomousDatabase
type AutonomousDatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutonomousDatabase `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AutonomousDatabase{}, &AutonomousDatabaseList{})
}

// Implement conversion.Hub interface, which means any resource version can convert into v4
func (*AutonomousDatabase) Hub() {}

// UpdateStatusFromOCIADB updates the status subresource
func (adb *AutonomousDatabase) UpdateStatusFromOciAdb(ociObj database.AutonomousDatabase) {
	adb.Status.LifecycleState = ociObj.LifecycleState
	adb.Status.TimeCreated = FormatSDKTime(ociObj.TimeCreated)

	// Data Guard / cross-region disaster-recovery facts. These are observed
	// only; spec.details.disasterRecovery is a create-time intent and is never
	// written back from OCI.
	adb.Status.Role = string(ociObj.Role)
	if len(ociObj.PeerDbIds) != 0 {
		adb.Status.PeerDbIds = append([]string(nil), ociObj.PeerDbIds...)
	} else {
		adb.Status.PeerDbIds = nil
	}
	adb.Status.DisasterRecoveryRegionType = string(ociObj.DisasterRecoveryRegionType)
	adb.Status.DataguardRegionType = string(ociObj.DataguardRegionType)

	if *ociObj.IsDedicated {
		conns := make([]ConnectionStringSpec, len(ociObj.ConnectionStrings.AllConnectionStrings))
		for key, val := range ociObj.ConnectionStrings.AllConnectionStrings {
			conns = append(conns, ConnectionStringSpec{TNSName: key, ConnectionString: val})
		}

		adb.Status.AllConnectionStrings = []ConnectionStringProfile{
			{ConnectionStrings: conns},
		}
	} else {
		var mTLSConns []ConnectionStringSpec
		var tlsConns []ConnectionStringSpec

		var conns []ConnectionStringProfile

		for _, profile := range ociObj.ConnectionStrings.Profiles {
			if profile.TlsAuthentication == database.DatabaseConnectionStringProfileTlsAuthenticationMutual {
				mTLSConns = append(mTLSConns, ConnectionStringSpec{TNSName: *profile.DisplayName, ConnectionString: *profile.Value})
			} else {
				tlsConns = append(tlsConns, ConnectionStringSpec{TNSName: *profile.DisplayName, ConnectionString: *profile.Value})
			}
		}

		if len(mTLSConns) > 0 {
			conns = append(conns, ConnectionStringProfile{
				TLSAuthentication: tlsAuthenticationMTLS,
				ConnectionStrings: mTLSConns,
			})
		}

		if len(tlsConns) > 0 {
			conns = append(conns, ConnectionStringProfile{
				TLSAuthentication: tlsAuthenticationTLS,
				ConnectionStrings: tlsConns,
			})
		}

		adb.Status.AllConnectionStrings = conns
	}
}

// UpdateFromOciAdb updates the attributes using database.AutonomousDatabase object
func (adb *AutonomousDatabase) UpdateFromOciAdb(ociObj database.AutonomousDatabase, overwrite bool) (specChanged bool) {
	oldADB := adb.DeepCopy()

	/***********************************
	* update the spec
	***********************************/
	if overwrite || adb.Spec.Details.Id == nil {
		adb.Spec.Details.Id = ociObj.Id
	}
	if overwrite || adb.Spec.Details.CompartmentId == nil {
		adb.Spec.Details.CompartmentId = ociObj.CompartmentId
	}
	if overwrite || adb.Spec.Details.AutonomousContainerDatabase.OciAcd.Id == nil {
		adb.Spec.Details.AutonomousContainerDatabase.OciAcd.Id = ociObj.AutonomousContainerDatabaseId
	}
	if overwrite || adb.Spec.Details.DisplayName == nil {
		adb.Spec.Details.DisplayName = ociObj.DisplayName
	}
	if overwrite || adb.Spec.Details.DbName == nil {
		adb.Spec.Details.DbName = ociObj.DbName
	}
	if overwrite || adb.Spec.Details.DbWorkload == "" {
		adb.Spec.Details.DbWorkload = ociObj.DbWorkload
	}
	if overwrite || adb.Spec.Details.LicenseModel == "" {
		adb.Spec.Details.LicenseModel = ociObj.LicenseModel
	}
	if overwrite || adb.Spec.Details.DatabaseEdition == "" {
		adb.Spec.Details.DatabaseEdition = ociObj.DatabaseEdition
	}
	if overwrite || adb.Spec.Details.DbVersion == nil {
		adb.Spec.Details.DbVersion = ociObj.DbVersion
	}
	if overwrite || adb.Spec.Details.DataStorageSizeInTBs == nil {
		adb.Spec.Details.DataStorageSizeInTBs = ociObj.DataStorageSizeInTBs
	}
	if overwrite || adb.Spec.Details.CpuCoreCount == nil {
		adb.Spec.Details.CpuCoreCount = ociObj.CpuCoreCount
	}
	if overwrite || adb.Spec.Details.ComputeModel == "" {
		adb.Spec.Details.ComputeModel = ociObj.ComputeModel
	}
	if overwrite || adb.Spec.Details.OcpuCount == nil {
		adb.Spec.Details.OcpuCount = ociObj.OcpuCount
	}
	if overwrite || adb.Spec.Details.ComputeCount == nil {
		adb.Spec.Details.ComputeCount = ociObj.ComputeCount
	}
	if overwrite || adb.Spec.Details.IsAutoScalingEnabled == nil {
		adb.Spec.Details.IsAutoScalingEnabled = ociObj.IsAutoScalingEnabled
	}
	if overwrite || adb.Spec.Details.IsDedicated == nil {
		adb.Spec.Details.IsDedicated = ociObj.IsDedicated
	}
	if overwrite || adb.Spec.Details.IsFreeTier == nil {
		adb.Spec.Details.IsFreeTier = ociObj.IsFreeTier
	}
	if overwrite || adb.Spec.Details.FreeformTags == nil {
		// Special case: an emtpy map will be nil after unmarshalling while the OCI always returns an emty map.
		if len(ociObj.FreeformTags) != 0 {
			adb.Spec.Details.FreeformTags = ociObj.FreeformTags
		} else {
			adb.Spec.Details.FreeformTags = nil
		}
	}

	if overwrite || adb.Spec.Details.IsAccessControlEnabled == nil {
		adb.Spec.Details.IsAccessControlEnabled = ociObj.IsAccessControlEnabled
	}

	if overwrite || adb.Spec.Details.WhitelistedIps == nil {
		if len(ociObj.WhitelistedIps) != 0 {
			adb.Spec.Details.WhitelistedIps = ociObj.WhitelistedIps
		} else {
			adb.Spec.Details.WhitelistedIps = nil
		}
	}
	if overwrite || adb.Spec.Details.IsMtlsConnectionRequired == nil {
		adb.Spec.Details.IsMtlsConnectionRequired = ociObj.IsMtlsConnectionRequired
	}
	if overwrite || adb.Spec.Details.SubnetId == nil {
		adb.Spec.Details.SubnetId = ociObj.SubnetId
	}
	if overwrite || adb.Spec.Details.NsgIds == nil {
		if len(ociObj.NsgIds) != 0 {
			adb.Spec.Details.NsgIds = ociObj.NsgIds
		} else {
			adb.Spec.Details.NsgIds = nil
		}
	}
	if overwrite || adb.Spec.Details.PrivateEndpointLabel == nil {
		adb.Spec.Details.PrivateEndpointLabel = ociObj.PrivateEndpointLabel
	}

	/***********************************
	* update the status subresource
	***********************************/
	adb.UpdateStatusFromOciAdb(ociObj)

	return !reflect.DeepEqual(oldADB.Spec, adb.Spec)
}

// RemoveUnchangedDetails removes the unchanged fields in spec.details, and returns if the details has been changed.
func (adb *AutonomousDatabase) RemoveUnchangedDetails(prevSpec AutonomousDatabaseSpec) (bool, error) {

	changed, err := RemoveUnchangedFields(prevSpec.Details, &adb.Spec.Details)
	if err != nil {
		return changed, err
	}

	return changed, nil
}

// A helper function which is useful for debugging. The function prints out a structural JSON format.
func (adb *AutonomousDatabase) String() (string, error) {
	out, err := json.MarshalIndent(adb, "", "    ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
