package client

import (
	"encoding/xml"
	"errors"
	"reflect"
	"strconv"
)

const soapHeader = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:apir="http://www.microsoft.com/SoftwareDistribution/Server/ApiRemotingWebService"><soapenv:Header/><soapenv:Body>`
const soapFooter = `</soapenv:Body></soapenv:Envelope>`

func wrapXML(data []byte) []byte {
	return []byte(soapHeader + string(data) + soapFooter)
}

type genericReadableRow struct {
	XMLName xml.Name      `xml:"apir:GenericReadableRow"`
	Values  []interface{} `xml:"apir:Values>apir:anyType"`
}

func toGenericReadableRows(rows interface{}) ([]genericReadableRow, error) {
	var interfaces []interface{}
	rv := reflect.ValueOf(rows)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			interfaces = append(interfaces, rv.Index(i).Interface())
		}
	} else {
		return nil, errors.New("not a slice")
	}
	var genericReadableRows []genericReadableRow
	for _, row := range interfaces {
		var values []interface{}
		fields := reflect.ValueOf(row).NumField()
		for i := 0; i < fields; i++ {
			value := reflect.ValueOf(row).Field(i)
			values = append(values, value.Interface())
		}
		genericReadableRows = append(genericReadableRows, genericReadableRow{Values: values})
	}
	return genericReadableRows, nil
}

func (row *readableRow) fromGenericReadableRow(target interface{}) error {
	rv := reflect.Indirect(reflect.ValueOf(target))
	if rv.Kind() != reflect.Struct {
		return errors.New("not a struct")
	}
	fields := rv.NumField()
	for i := 0; i < fields; i++ {
		if row.Values[i] == "" {
			continue
		}
		field := rv.Field(i)
		typeField := field.Type()
		if field.CanSet() {
			switch typeField.Kind() {
			case reflect.String:
				field.SetString(row.Values[i])
			case reflect.Int, reflect.Int64:
				if value, err := strconv.ParseInt(row.Values[i], 10, 64); err == nil {
					field.SetInt(value)
				} else {
					return err
				}
			case reflect.TypeOf(PatchingType(0)).Kind():
				if value, err := strconv.ParseInt(row.Values[i], 10, 64); err == nil {
					field.Set(reflect.ValueOf(PatchingType(value)))
				} else {
					return err
				}
			case reflect.TypeOf(UpdateDeploymentAction(0)).Kind():
				if value, err := strconv.ParseInt(row.Values[i], 10, 64); err == nil {
					field.Set(reflect.ValueOf(UpdateDeploymentAction(value)))
				} else {
					return err
				}
			case reflect.Uint:
				if value, err := strconv.ParseUint(row.Values[i], 10, 64); err == nil {
					field.SetUint(value)
				} else {
					return err
				}
			case reflect.Bool:
				if value, err := strconv.ParseBool(row.Values[i]); err == nil {
					field.SetBool(value)
				} else {
					return err
				}
			}
		} else {
			return errors.New("cannot set field")
		}
	}
	return nil
}

type readableRow struct {
	XMLName xml.Name `xml:"GenericReadableRow"`
	Values  []string `xml:"Values>anyType"`
}

type arrayOfGenericReadableRow struct {
	XMLName xml.Name      `xml:"ArrayOfGenericReadableRow"`
	Rows    []readableRow `xml:"GenericReadableRow"`
}

type soapBody struct {
	XMLName xml.Name
	Data    []byte `xml:",innerxml"`
}

type soapEnvelope struct {
	XMLName xml.Name
	Body    soapBody
}

type executeSPCountObsoleteUpdatesToCleanupResponse struct {
	XMLName xml.Name `xml:"ExecuteSPCountObsoleteUpdatesToCleanupResponse"`
	Count   int      `xml:"ExecuteSPCountObsoleteUpdatesToCleanupResult"`
}

// GetSPCountUpdatesToCleanupResponse returns the number of updates to cleanup or an error.
func GetSPCountObsoleteUpdatesToCleanupResponse(response []byte) (int, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return 0, err
	}
	var r executeSPCountObsoleteUpdatesToCleanupResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return 0, err
	}
	return r.Count, nil
}

type executeSPCountUpdatesToCompressResponse struct {
	XMLName xml.Name `xml:"ExecuteSPCountUpdatesToCompressResponse"`
	Count   int      `xml:"ExecuteSPCountUpdatesToCompressResult"`
}

// GetSPCountUpdatesToCompressResponse returns the number of updates to compress or an error.
func GetSPCountUpdatesToCompressResponse(response []byte) (int, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return 0, err
	}
	var r executeSPCountUpdatesToCompressResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return 0, err
	}
	return r.Count, nil
}

type executeSPGetAllComputersResponse struct {
	XMLName xml.Name                    `xml:"ExecuteSPGetAllComputersResponse"`
	Array   []arrayOfGenericReadableRow `xml:"ExecuteSPGetAllComputersResult>ArrayOfGenericReadableRow"`
}

// ITargetComputer represents a target computer in the WSUS database.
type ITargetComputer struct {
	ComputerID             string
	LastSyncTime           string
	LastReportedStatusTime string
	IPAddress              string
	FullDomainName         string
	OSMajorVersion         int
	OSMinorVersion         int
	OSBuildNumber          int
	OSServicePackMajor     int
	OSServicePackMinor     int
	OSLocale               string
	ComputerMake           string
	ComputerModel          string
	BiosVersion            string
	BiosName               string
	BiosReleaseDate        string
	ProcessorArchitecture  string
	RequestedTargetGroupID string
	LastInventoryTime      string
	AccountServerID        string
	LastSyncResult         int
	SuiteMask              int
	OldProductType         int
	NewProductType         int
	SystemMetrics          int
	ClientVersion          string
	OSFamily               string
	OSDescription          string
	OEM                    string
	DeviceType             string
	FirmwareVersion        string
	MobileOperator         string
}

// RequestedTargetGroups represents any requested groups for a target computer.
type RequestedTargetGroups struct {
	Computer string
	Group    string
}

// AssignedTargetGroups represents any assigned groups for a target computer.
type AssignedTargetGroups struct {
	Computer string
	Group    string
}

// GetSPGetAllComputersResponse returns a slice of ITargetComputer, RequestedTargetGroups, and AssignedTargetGroups - or an error.
func GetSPGetAllComputersResponse(response []byte) ([]ITargetComputer, []RequestedTargetGroups, []AssignedTargetGroups, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return nil, nil, nil, err
	}
	var r executeSPGetAllComputersResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, nil, nil, err
	}

	var computers []ITargetComputer
	for _, row := range r.Array[0].Rows {
		var computer ITargetComputer
		if err := row.fromGenericReadableRow(&computer); err != nil {
			return nil, nil, nil, err
		}
		computers = append(computers, computer)
	}

	var requestedTargetGroups []RequestedTargetGroups
	for _, row := range r.Array[1].Rows {
		var requestedTargetGroup RequestedTargetGroups
		if err := row.fromGenericReadableRow(&requestedTargetGroup); err != nil {
			return nil, nil, nil, err
		}
		requestedTargetGroups = append(requestedTargetGroups, requestedTargetGroup)
	}

	var assignedTargetGroups []AssignedTargetGroups
	for _, row := range r.Array[2].Rows {
		var assignedTargetGroup AssignedTargetGroups
		if err := row.fromGenericReadableRow(&assignedTargetGroup); err != nil {
			return nil, nil, nil, err
		}
		assignedTargetGroups = append(assignedTargetGroups, assignedTargetGroup)
	}

	return computers, requestedTargetGroups, assignedTargetGroups, nil
}

type executeSPGetAllDownstreamServersResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetAllDownstreamServersResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetAllDownstreamServersResult>ArrayOfGenericReadableRow"`
}

// DownstreamServer represents the information stored in WSUS about DSS.
type DownstreamServer struct {
	DomainName     string
	ID             string
	LastSyncTime   string
	ParentID       string
	LastRollupTime string
	Version        string
	IsReplica      bool
}

// GetSPAllDownstreamServersResponse returns a slice of DownstreamServer - or an error.
func GetSPGetAllDownstreamServersResponse(response []byte) ([]DownstreamServer, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return nil, err
	}
	var r executeSPGetAllDownstreamServersResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var downstreamServers []DownstreamServer
	for _, row := range r.Rows {
		var downstreamServer DownstreamServer
		if err := row.fromGenericReadableRow(&downstreamServer); err != nil {
			return nil, err
		}
		downstreamServers = append(downstreamServers, downstreamServer)
	}

	return downstreamServers, nil
}

type executeSPGetAllLanguagesWithEnabledStateResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetAllLanguagesWithEnabledStateResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetAllLanguagesWithEnabledStateResult>ArrayOfGenericReadableRow"`
}

// LanguageState represents the state of a language in the WSUS database.
type LanguageState struct {
	ShortLanguageName string
	Enabled           bool
	UssEnabled        bool
}

// GetSPGetAllLanguagesWithEnabledStateResponse returns a slice of LanguageState - or an error.
func GetSPGetAllLanguagesWithEnabledStateResponse(response []byte) ([]LanguageState, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return nil, err
	}
	var r executeSPGetAllLanguagesWithEnabledStateResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var languageStates []LanguageState
	for _, row := range r.Rows {
		var languageState LanguageState
		if err := row.fromGenericReadableRow(&languageState); err != nil {
			return nil, err
		}
		languageStates = append(languageStates, languageState)
	}

	return languageStates, nil
}

type getSPGetAllTargetGroupsResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetAllTargetGroupsResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetAllTargetGroupsResult>ArrayOfGenericReadableRow"`
}

// TargetGroup represents the information stored in WSUS about a target group.
type TargetGroup struct {
	Type          string
	Name          string
	Guid          string
	OrderValue    int
	GroupPriority int
	ParentGroupID string
}

// GetSPGetAllTargetGroupsResponse returns a slice of TargetGroup - or an error.
func GetSPGetAllTargetGroupsResponse(response []byte) ([]TargetGroup, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return nil, err
	}
	var r getSPGetAllTargetGroupsResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var targetGroups []TargetGroup
	for _, row := range r.Rows {
		var targetGroup TargetGroup
		if err := row.fromGenericReadableRow(&targetGroup); err != nil {
			return nil, err
		}
		targetGroups = append(targetGroups, targetGroup)
	}

	return targetGroups, nil
}

type getSPGetApprovedUpdatesMetaDataResponse struct {
	XMLName xml.Name                    `xml:"ExecuteSPGetApprovedUpdatesMetaDataResponse"`
	Array   []arrayOfGenericReadableRow `xml:"ExecuteSPGetApprovedUpdatesMetaDataResult>ArrayOfGenericReadableRow"`
}

// UpdateMetaData represents the information stored in WSUS about an update.
type UpdateMetaData struct {
	UpdateID       string
	RevisionNumber int
	RevisionID     int
	XML            string

	// If xml is not set, this will be set.
	// Base64 encoded data.
	// Compressed using in-memory CAB files.
	// https://docs.microsoft.com/en-us/windows/win32/msi/cabinet-files
	XMLCompressed string

	LocalUpdateID int
}

type PatchingType int

const (
	PatchingTypeNone      PatchingType = iota
	PatchingTypeContained PatchingType = iota
	PatchingTypeExpress   PatchingType = iota
	PatchingTypeDelta     PatchingType = iota
)

// UpdateFile represents the information stored in WSUS about an update file.
type UpdateFile struct {
	RevisionID   int
	FileName     string
	Modified     int64
	HostedOnMU   bool
	Size         string
	FileSize     int
	PatchingType PatchingType
	IsEula       bool
}

// UpdateApproval represents the information stored in WSUS about an update approval.
type UpdateApproval struct {
	DeploymentTime string // utc
	Deployed       bool
	ActionID       UpdateDeploymentAction
	GoLiveTime     string // utc
	Deadline       string
	AdminName      string
	DeploymentGUID string
	IsAssigned     bool
	UpdateID       string
	RevisionNumber int
	TargetGroupID  string
}

// GetSPGetApprovedUpdatesMetaDataResponse returns a slice of UpdateMetaData, UpdateFiles, and UpdateApprovals - or an error.
func GetSPGetApprovedUpdatesMetaDataResponse(response []byte) ([]UpdateMetaData, []UpdateFile, []UpdateApproval, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(response, &e); err != nil {
		return nil, nil, nil, err
	}
	var r getSPGetApprovedUpdatesMetaDataResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, nil, nil, err
	}
	var updates []UpdateMetaData
	for _, row := range r.Array[0].Rows {
		var update UpdateMetaData
		if err := row.fromGenericReadableRow(&update); err != nil {
			return nil, nil, nil, err
		}
		updates = append(updates, update)
	}

	var updateFiles []UpdateFile
	for _, row := range r.Array[1].Rows {
		var updateFile UpdateFile
		if err := row.fromGenericReadableRow(&updateFile); err != nil {
			return nil, nil, nil, err
		}
		updateFiles = append(updateFiles, updateFile)
	}

	var updateApprovals []UpdateApproval
	for _, row := range r.Array[2].Rows {
		var updateApproval UpdateApproval
		if err := row.fromGenericReadableRow(&updateApproval); err != nil {
			return nil, nil, nil, err
		}
		updateApprovals = append(updateApprovals, updateApproval)
	}

	return updates, updateFiles, updateApprovals, nil
}

type getSPGetCategoriesResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetCategoriesResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetCategoriesResult>ArrayOfGenericReadableRow"`
}

// Category represents an update category
type Category struct {
	LocalUpdateID          int
	UpdateID               [16]byte
	CategoryType           string
	ProhibitsSubcategories bool
	ProhibitsUpdates       bool
	CategoryIndex          int
	DisplayOrder           int
	Title                  string
	Description            string
	ReleaseNotes           string
	Received               int64
	UpdateSource           int
}

// GetSPGetCategoriesResponse returns a slice of Category - or an error.
func GetSPGetCategoriesResponse(data []byte) ([]Category, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	var r getSPGetCategoriesResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var categories []Category
	for _, row := range r.Rows {
		var category Category
		if err := row.fromGenericReadableRow(&category); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

type getSPGetCategoryByIDResponse struct {
	XMLName xml.Name    `xml:"ExecuteSPGetCategoryByIDResponse"`
	Row     readableRow `xml:"ExecuteSPGetCategoryByIDResult>GenericReadableRow"`
}

// GetSPGetCategoryByIDResponse returns a Category - or an error.
func GetSPGetCategoryByIDResponse(data []byte) (*Category, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	var r getSPGetCategoryByIDResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var category Category
	if err := r.Row.fromGenericReadableRow(&category); err != nil {
		return nil, err
	}

	return &category, nil
}

type getSPGetChildTargetGroupsResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetChildTargetGroupsResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetChildTargetGroupsResult>ArrayOfGenericReadableRow"`
}

// GetSPGetChildTargetGroupsResponse returns a slice of TargetGroup - or an error.
func GetSPGetChildTargetGroupsResponse(data []byte) ([]TargetGroup, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	var r getSPGetChildTargetGroupsResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var targetGroups []TargetGroup
	for _, row := range r.Rows {
		var targetGroup TargetGroup
		if err := row.fromGenericReadableRow(&targetGroup); err != nil {
			return nil, err
		}
		targetGroups = append(targetGroups, targetGroup)
	}

	return targetGroups, nil
}

type getSPGetClientsWithRecentNameChangeResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetClientsWithRecentNameChangeResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetClientsWithRecentNameChangeResult>ArrayOfGenericReadableRow"`
}

// GetSPGetClientsWithRecentNameChangeResponse returns a slice of FQDNs - or an error.
func GetSPGetClientsWithRecentNameChangeResponse(data []byte) ([]string, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	var r getSPGetClientsWithRecentNameChangeResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var clients []string
	for _, row := range r.Rows {
		var client string
		if err := row.fromGenericReadableRow(&client); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

type getSPGetComponentWithErrorsResponse struct {
	XMLName xml.Name      `xml:"ExecuteSPGetComponentWithErrorsResponse"`
	Rows    []readableRow `xml:"ExecuteSPGetComponentWithErrorsResult>ArrayOfGenericReadableRow"`
}

// GetSPGetComponentWithErrorsResponse returns a slice of component names - or an error.
func GetSPGetComponentWithErrorsResponse(data []byte) ([]string, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	var r getSPGetComponentWithErrorsResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var clients []string
	for _, row := range r.Rows {
		var client string
		if err := row.fromGenericReadableRow(&client); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

type getSPGetComputerByIDResponse struct {
	XMLName xml.Name    `xml:"ExecuteSPGetComputerByIDResponse"`
	Row     readableRow `xml:"ExecuteSPGetComputerByIDResult>GenericReadableRow"`
}

// GetSPGetComputerByIDResponse returns a ITargetComputer - or an error.
func GetSPGetComputerByIDResponse(data []byte) (*ITargetComputer, error) {
	var e soapEnvelope
	if err := xml.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	var r getSPGetComputerByIDResponse
	if err := xml.Unmarshal(e.Body.Data, &r); err != nil {
		return nil, err
	}
	var computer ITargetComputer
	if err := r.Row.fromGenericReadableRow(&computer); err != nil {
		return nil, err
	}

	return &computer, nil
}
