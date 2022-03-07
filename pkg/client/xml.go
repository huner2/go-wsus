package client

import (
	"encoding/xml"
)

const soapHeader = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:apir="http://www.microsoft.com/SoftwareDistribution/Server/ApiRemotingWebService"><soapenv:Header/><soapenv:Body>`
const soapFooter = `</soapenv:Body></soapenv:Envelope>`

func wrapXML(data []byte) []byte {
	return []byte(soapHeader + string(data) + soapFooter)
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
