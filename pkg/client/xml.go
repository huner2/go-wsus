package client

const soapHeader = `<soapenv:Envelope xlmns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:apir="http://www.microsoft.com/SoftwareDistribution/Server/ApiRemotingWebService"><soapenv:Header/><soapenv:Body>`
const soapFooter = `</soapenv:Body></soapenv:Envelope>`

func wrapXML(data []byte) []byte {
	return []byte(soapHeader + string(data) + soapFooter)
}
