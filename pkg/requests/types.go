package requests

import (
	"encoding/xml"

	"github.com/google/uuid"
)

type SOAPRequest []byte

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
// Only ID, Name, and Type are required.
// ID should be a valid UUID.
// Type should be one of the DynamicCategoryType constants.
// If supplied, Origin should be one of the DynamicCategoryOriginType constants.
// It is recommended to only supply the required fields.
type DynamicCategory struct {
	XMLName       xml.Name                  `xml:"apir:AddDynamicCategory"`
	ID            uuid.UUID                 `xml:"apir:id"`
	Name          string                    `xml:"apir:name"`
	Type          DynamicCategoryType       `xml:"apir:type"`
	Origin        DynamicCategoryOriginType `xml:"apir:origin omitempty"`
	IsSyncEnabled bool                      `xml:"apir:isSyncEnabled omitempty"`
	DiscoveryTime int64                     `xml:"apir:discoveryTime omitempty"`
	TargetId      int                       `xml:"apir:targetId omitempty"`
}
