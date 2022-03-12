package client

import (
	"encoding/xml"
)

type SOAPInterface interface {
	toXml() ([]byte, error)
}

type DynamicCategoryType uint8

const (
	DynamicCategoryType_ComputerModel DynamicCategoryType = 1
	DynamicCategoryType_Device        DynamicCategoryType = 2
	DynamicCategoryType_Application   DynamicCategoryType = 4
	DynamicCategoryType_Any           DynamicCategoryType = 255
)

type DynamicCategoryOriginType uint8

const (
	DynamicCategoryOriginType_Automatic DynamicCategoryOriginType = 1
	DynamicCategoryOriginType_Manual    DynamicCategoryOriginType = 2
	DynamicCategoryOriginType_Any       DynamicCategoryOriginType = 255
)

// DynamicCategory represents a dynamic category in WSUS.
//
// Only ID and Type are required.
//
// ID should be a valid UUID.
//
// Type should be one of the DynamicCategoryType constants.
//
// If supplied, Origin should be one of the DynamicCategoryOriginType constants.
//
// It is recommended to only supply the required fields.
type DynamicCategoryInterface struct {
	XMLName       xml.Name                  `xml:"apir:AddDynamicCategory"`
	ID            [16]byte                  `xml:"apir:id"`
	Name          string                    `xml:"apir:name"` // Optional
	Type          DynamicCategoryType       `xml:"apir:type"`
	Origin        DynamicCategoryOriginType `xml:"apir:origin omitempty"`
	IsSyncEnabled bool                      `xml:"apir:isSyncEnabled omitempty"`
	DiscoveryTime int64                     `xml:"apir:discoveryTime omitempty"`
	TargetId      int                       `xml:"apir:targetId omitempty"`
}

func (d DynamicCategoryInterface) toXml() ([]byte, error) {
	return xml.Marshal(d)
}

// DynamicCategories represents a list of dynamic categories.
//
// This is used to add multiple dynamic categories at once
type DynamicCategoriesInterface struct {
	XMLName    xml.Name                   `xml:"apir:AddDynamicCategories"`
	Categories []DynamicCategoryInterface `xml:"apir:categories"`
}

func (d DynamicCategoriesInterface) toXml() ([]byte, error) {
	temp, err := xml.Marshal(d)
	if err != nil {
		return nil, err
	}
	return []byte("<apir:AddDynamicCategories><apir:categories>" + string(temp) + "</apir:categories></apir:AddDynamicCategories>"), nil
}

func (d DynamicCategoriesInterface) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	rows, err := toGenericReadableRows(d.Categories)
	if err != nil {
		return err
	}
	return e.Encode(rows)
}

// NewLanguage represents a locale in WSUS.
type NewLanguageInterface struct {
	XMLName       xml.Name `xml:"apir:AddNewLanguage"`
	ID            int      `xml:"apir:newLanguage>apir:LanguageId"`
	ShortLanguage string   `xml:"apir:newLanguage>apir:ShortLanguage"` // Optional
	LongLanguage  string   `xml:"apir:newLanguage>apir:LongLanguage"`  // Optional
	Enabled       bool     `xml:"apir:newLanguage>apir:Enabled"`
}

func (n NewLanguageInterface) toXml() ([]byte, error) {
	return xml.Marshal(n)
}

// AutomaticUpdateApprovalRule is used to apply automatic update approval rules to WSUS.
type AutomaticUpdateApprovalRuleInterface struct {
	XMLName xml.Name `xml:"apir:ApplyAutomaticUpdateApprovalRule"`
	ID      int      `xml:"apir:ruleId"`
}

func (a AutomaticUpdateApprovalRuleInterface) toXml() ([]byte, error) {
	return xml.Marshal(a)
}

// CatalogSiteGetMetadataAndImport is used to get metadata and import update binaries.
type CatalogSiteGetMetadataAndImportInterface struct {
	XMLName             xml.Name `xml:"apir:CatalogSiteGetMetadataAndImport"`
	ID                  [16]byte `xml:"apir:updateId"`
	DownloadFileDigests []byte   `xml:"apir:downloadFileDigests>apir:base64Binary"` // Optional
}

func (c CatalogSiteGetMetadataAndImportInterface) toXml() ([]byte, error) {
	return xml.Marshal(c)
}

// InstallApprovalRule is used to create an install approval rule.
type InstallApprovalRuleInterface struct {
	XMLName xml.Name `xml:"apir:CreateInstallApprovalRule"`
	Name    string   `xml:"apir:name"`
}

func (c InstallApprovalRuleInterface) toXml() ([]byte, error) {
	return xml.Marshal(c)
}

// DeleteDynamicCategory is used to delete a dynamic category.
type DeleteDynamicCategoryInterface struct {
	XMLName xml.Name `xml:"apir:DeleteDynamicCategory"`
	ID      [16]byte `xml:"apir:id"`
}

func (d DeleteDynamicCategoryInterface) toXml() ([]byte, error) {
	return xml.Marshal(d)
}

// DeleteInstallApprovalRule is used to delete an install approval rule.
type DeleteInstallApprovalRuleInterface struct {
	XMLName xml.Name `xml:"apir:DeleteInstallApprovalRule"`
	ID      int      `xml:"apir:id"`
}

func (d DeleteInstallApprovalRuleInterface) toXml() ([]byte, error) {
	return xml.Marshal(d)
}

// ExecuteGetSigningCertificate is used to get the WSUS signing certificate.
type ExecuteGetSigningCertificateInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteGetSigningCertificate"`
}

func (e ExecuteGetSigningCertificateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteReplicaSPDeleteDeployment is used to delete a replica deployment.
type ExecuteReplicaSPDeleteDeploymentInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteReplicaSPDeleteDeployment"`
	ID      [16]byte `xml:"apir:id"`
}

func (e ExecuteReplicaSPDeleteDeploymentInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSetSelfSigningCertificate instructs the WSUS server to use a self-signed certificate.
type ExecuteSetSelfSigningCertificateInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSetSelfSigningCertificate"`
}

func (e ExecuteSetSelfSigningCertificateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSetSigningCertificate allows specification of a signing certificate.
// The certificate should base64 encoded PFX format.
type ExecuteSetSigningCertificateInterface struct {
	XMLName  xml.Name `xml:"apir:ExecuteSetSigningCertificate"`
	Cert     []byte   `xml:"apir:PFXFileContent"`
	Password []byte   `xml:"apir:passwordBytes"` // Optional, if private key is password protected
}

func (e ExecuteSetSigningCertificateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPAcceptEula is used to accept the EULA for an update.
//
// AdminName represents the name of the user who accepted the update.
// Should be in domain\username format if provided.
type ExecuteSPAcceptEulaInterface struct {
	XMLName   xml.Name `xml:"apir:ExecuteSPAcceptEula"`
	EulaID    [16]byte `xml:"apir:id"`
	AdminName string   `xml:"apir:adminName"`
	UpdateID  [16]byte `xml:"apir:updateId>apir:updateId"`
	Revision  int      `xml:"apir:updateId>apir:revision"`
}

func (e ExecuteSPAcceptEulaInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPAcceptEulaForReplicaDSS is used to accept the EULA for replica downstream servers.
type ExecuteSPAcceptEulaForReplicaDSSInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPAcceptEulaForReplicaDSS"`
	EulaID  [16]byte `xml:"apir:eulaId"`
}

func (e ExecuteSPAcceptEulaForReplicaDSSInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPAddComputerToTargetGroupAllowMultipleGroups is used to add a computer to a target group allowing multiple groups.
type ExecuteSPAddComputerToTargetGroupAllowMultipleGroupsInterface struct {
	XMLName       xml.Name `xml:"apir:ExecuteSPAddComputerToTargetGroupAllowMultipleGroups"`
	TargetGroupID [16]byte `xml:"apir:targetGroupId"`
	ComputerID    string   `xml:"apir:computerId"`
}

func (e ExecuteSPAddComputerToTargetGroupAllowMultipleGroupsInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCancelAllDownloads is used to cancel all current update downloads.
type ExecuteSPCancelAllDownloadsInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPCancelAllDownloads"`
}

func (e ExecuteSPCancelAllDownloadsInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCacnelDownload is used to cancel a current update download.
type ExecuteSPCancelDownloadInterface struct {
	XMLName  xml.Name `xml:"apir:ExecuteSPCancelDownload"`
	ID       [16]byte `xml:"apir:id>apir:UpdateId"`
	Revision int      `xml:"apir:id>apir:RevisionNumber"`
}

func (e ExecuteSPCancelDownloadInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCleanupObsoleteComputers deletes obsolete computers from WSUS.
type ExecuteSPCleanupObsoleteComputersInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPCleanupObsoleteComputers"`
}

func (e ExecuteSPCleanupObsoleteComputersInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCleanupUnneededContentFiles deletes unneeded content files from WSUS.
type ExecuteSPCleanupUnneededContentFilesInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPCleanupUnneededContentFiles"`
}

func (e ExecuteSPCleanupUnneededContentFilesInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCleanupUnneededContentFilesPreciseInterface deletes unneeded content files from WSUS, allowing for server targetting and local published files.
type ExecuteSPCleanupUnneededContentFilesPreciseInterface struct {
	XMLName                  xml.Name `xml:"apir:ExecuteSPCleanupUnneededContentFiles2"`
	ServerName               string   `xml:"apir:updateServerName"` // optional
	CleanLocalPublishedFiles bool     `xml:"apir:cleanLocalPublishedContentFiles"`
}

func (e ExecuteSPCleanupUnneededContentFilesPreciseInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCompressUpdate is used to enable compression for an update.
type ExecuteSPCompressUpdateInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPCompressUpdate"`
	ID      int      `xml:"apir:localUpdateId"`
}

func (e ExecuteSPCompressUpdateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCountObsoleteUpdatesToCleanup counts the number of obsolete updates to cleanup.
// Use GetSPCountObsoleteUpdatesToCleanupResponse on the response to get the result.
type ExecuteSPCountObsoleteUpdatesToCleanupInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPCountObsoleteUpdatesToCleanup"`
}

func (e ExecuteSPCountObsoleteUpdatesToCleanupInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCountUpdatesToCompress counts the number of updates to compress.
// Use GetSPCountUpdatesToCompressResponse on the response to get the result.
type ExecuteSPCountUpdatesToCompressInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPCountUpdatesToCompress"`
}

func (e ExecuteSPCountUpdatesToCompressInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCreateTargetGroup is used to create a target group.
type ExecuteSPCreateTargetGroupInterface struct {
	XMLName         xml.Name `xml:"apir:ExecuteSPCreateTargetGroup"`
	TargetGroupName string   `xml:"apir:name"`
	ParentGroupID   [16]byte `xml:"apir:parentGroupId"`
}

func (e ExecuteSPCreateTargetGroupInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPCreateTargetGroupPrecise is used to create a target group, allowing for specification of the ID of the group and whether or not the group
// can be added when the server is in replication mode.
type ExecuteSPCreateTargetGroupPreciseInterface struct {
	XMLName         xml.Name `xml:"apir:ExecuteSPCreateTargetGroup2"`
	TargetGroupName string   `xml:"apir:name"`
	ParentGroupID   [16]byte `xml:"apir:parentGroupId"`
	ID              [16]byte `xml:"apir:id"`
	FailIfReplica   bool     `xml:"apir:failIfReplica"`
}

func (e ExecuteSPCreateTargetGroupPreciseInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeclineExpiredUpdates is used to decline all expired updates.
type ExecuteSPDeclineExpiredUpdatesInterface struct {
	XMLName   xml.Name `xml:"apir:ExecuteSPDeclineExpiredUpdates"`
	AdminName string   `xml:"apir:adminName"`
}

func (e ExecuteSPDeclineExpiredUpdatesInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeclineSupercededUpdates is used to decline all superceded updates.
type ExecuteSPDeclineSupercededUpdatesInterface struct {
	XMLName   xml.Name `xml:"apir:ExecuteSPDeclineSupercededUpdates"`
	AdminName string   `xml:"apir:adminName"`
}

func (e ExecuteSPDeclineSupercededUpdatesInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeclineUpdate will decline the specified update.
type ExecuteSPDeclineUpdateInterface struct {
	XMLName       xml.Name `xml:"apir:ExecuteSPDeclineUpdate"`
	ID            [16]byte `xml:"apir:updateId"`
	AdminName     string   `xml:"apir:adminName"`
	FailIfReplica bool     `xml:"apir:failIfReplica"`
}

func (e ExecuteSPDeclineUpdateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeleteComputer deletes a computer from the WSUS inventory.
type ExecuteSPDeleteComputerInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPDeleteComputer"`
	ID      string   `xml:"apir:id"`
}

func (e ExecuteSPDeleteComputerInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeleteDeployment deletes a deployed WSUS server from the deployment group.
type ExecuteSPDeleteDeploymentInterface struct {
	XMLName   xml.Name `xml:"apir:ExecuteSPDeleteDeployment"`
	ID        [16]byte `xml:"apir:id"`
	AdminName string   `xml:"apir:adminName"`
}

func (e ExecuteSPDeleteDeploymentInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeleteDownstreamServer deletes a downstream server from the deployment group.
type ExecuteSPDeleteDownstreamServerInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPDeleteDownstreamServer"`
	ID      [16]byte `xml:"apir:id"`
}

func (e ExecuteSPDeleteDownstreamServerInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeleteTargetGroup deletes a target group.
type ExecuteSPDeleteTargetGroupInterface struct {
	XMLName       xml.Name `xml:"apir:ExecuteSPDeleteTargetGroup"`
	ID            [16]byte `xml:"apir:id"`
	AdminName     string   `xml:"apir:adminName"`
	FailIfReplica bool     `xml:"apir:failIfReplica"`
}

func (e ExecuteSPDeleteTargetGroupInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeleteUpdate removes a local update from the WSUS inventory.
//
// ID is the local update ID.
type ExecuteSPDeleteUpdateInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPDeleteUpdate"`
	ID      int      `xml:"apir:localUpdateID"`
}

func (e ExecuteSPDeleteUpdateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeleteUpdateByID removes an update from the WSUS inventory.
//
// ID is the update GUID.
type ExecuteSPDeleteUpdateByIDInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPDeleteUpdateByID"`
	ID      [16]byte `xml:"apir:updateID"`
}

func (e ExecuteSPDeleteUpdateByIDInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

type UpdateDeploymentAction int

const (
	UpdateDeploymentAction_Install UpdateDeploymentAction = iota
	UpdateDeploymentAction_Uninstall
	UpdateDeploymentAction_NotApproved
)

type UpdateDeploymentDownloadPriority int

const (
	UpdateDeploymentDownloadPriority_High UpdateDeploymentDownloadPriority = iota
	UpdateDeploymentDownloadPriority_Normal
	UpdateDeploymentDownloadPriority_Low
)

// ExecuteSPDeployUpdate allows for deployment of an update to a specified target group.
//
// Deadline is the deadline for the update to be deployed on a client (essentially a priority modifier).
// Deadline is a nanosecond timestamp.
//
// IsAssigned must be set to true.
type ExecuteSPDeployUpdateInterface struct {
	XMLName          xml.Name               `xml:"apir:ExecuteSPDeployUpdate1"`
	ID               [16]byte               `xml:"apir:updateId>UpdateId"`
	Revision         int                    `xml:"apir:updateId>Revision"`
	DeploymentAction UpdateDeploymentAction `xml:"apir:deploymentAction"`
	TargetGroup      [16]byte               `xml:"apir:targetGroupId"`
	Deadline         int64                  `xml:"apir:deadline"`
	AdminName        string                 `xml:"apir:adminName"`
	IsAssigned       bool                   `xml:"apir:isAssigned"`
}

func (e ExecuteSPDeployUpdateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPDeployUpdatePrecise allows for deployment of an update with extended options.
//
// GoLiveTime is the time at which the update will appear for clients.
//
// DownloadPriority is the priority of the update, the higher the number the higher the priority.
//
// DeploymentGUID denotes the GUID of the deployment.
// It is present when replicating deployments of a USS to DSS, as deployment GUIDs MUST be the same.
// If the deployment is not part of a USS/DSS replica synchronization, this field SHOULD be left default.
//
// TranslateSQLException is undocumented - leaving as false is recommended.
//
// FailIfReplica will prevent the action if the server is in replication.
//
// IsReplicaSync should always be set to false.
type ExecuteSPDeployUpdatePreciseInterface struct {
	XMLName               xml.Name                         `xml:"apir:ExecuteSPDeployUpdate2"`
	ID                    [16]byte                         `xml:"apir:updateId"`
	Revision              int                              `xml:"apir:revisionNumber"`
	DeploymentAction      UpdateDeploymentAction           `xml:"apir:deploymentAction"`
	TargetGroup           [16]byte                         `xml:"apir:targetGroupId"`
	AdminName             string                           `xml:"apir:adminName"`
	Deadline              int64                            `xml:"apir:deadline"`
	IsAssigned            bool                             `xml:"apir:isAssigned"`
	GoLiveTime            int64                            `xml:"apir:goLiveTime"`
	DownloadPriority      UpdateDeploymentDownloadPriority `xml:"apir:downloadPriority"`
	DeploymentGUID        [16]byte                         `xml:"apir:deploymentGUID"` // optional
	TranslateSQLException bool                             `xml:"apir:translateSqlException"`
	FailIfReplica         bool                             `xml:"apir:failIfReplica"`
	IsReplicaSync         bool                             `xml:"apir:isReplicaSync"`
}

func (e ExecuteSPDeployUpdatePreciseInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPGetAllComputers returns information for all computers in the WSUS database.
//
// Use GetSPGetAllComputersResponse to parse the response.
//
// Note: Group ID a0a08746-4dbe-4a37-9adf-9e7652c0b421 is the all computers group.
type ExecuteSPGetAllComputersInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPGetAllComputers"`
}

func (e ExecuteSPGetAllComputersInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPGetAllDownstreamServers returns information for all downstream servers in the WSUS database.
//
// Use GetSPGetAllDownstreamServersResponse to parse the response.
type ExecuteSPGetAllDownstreamServersInterface struct {
	XMLName               xml.Name `xml:"apir:ExecuteSPGetAllDownstreamServers"`
	ParentID              [16]byte `xml:"apir:parentServerId"` // optional
	IncludeNestedChildren bool     `xml:"apir:includeNestedChildren"`
}

func (e ExecuteSPGetAllDownstreamServersInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPGetAllLanguagesWithEnabledState returns information for all languages in the WSUS database.
//
// Use GetSPGetAllLanguagesWithEnabledStateResponse to parse the response.
type ExecuteSPGetAllLanguagesWithEnabledStateInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPGetAllLanguagesWithEnabledState"`
}

func (e ExecuteSPGetAllLanguagesWithEnabledStateInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}

// ExecuteSPGetAllTargetGroups returns information for all target groups in the WSUS database.
//
// Use GetSPGetAllTargetGroupsResponse to parse the response.
type ExecuteSPGetAllTargetGroupsInterface struct {
	XMLName xml.Name `xml:"apir:ExecuteSPGetAllTargetGroups"`
}

func (e ExecuteSPGetAllTargetGroupsInterface) toXml() ([]byte, error) {
	return xml.Marshal(e)
}
