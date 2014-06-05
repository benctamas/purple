package purple

import (
	"encoding/xml"
)

type serverBase struct {
	Type                     string `xml:"type,attr"`
	Hostname                 string `xml:"hostname"`
	Port                     int    `xml:"port"`
	SocketType               string `xml:"socketType"`
	Username                 string `xml:"username"`
	Authentication           string `xml:"authentication"`
	AddThisServer            bool   `xml:"addThisServer,omitempty"`
	UseGlobalPrefferedServer bool   `xml:"useGlobalPrefferedServer,omitempty"`
	Restriction              string `xml:"restriction,omitempty"`
}

type Pop3Config struct {
	XMLName               xml.Name `xml:"pop3"`
	LeaveMessagesOnServer bool     `xml:"leaveMessagesOnServer"`
}

type IncomingServer struct {
	XMLName xml.Name `xml:"incomingServer"`
	serverBase
	Pop3 Pop3Config `xml:"pop3,omitempty"`
}

type OutgoingServer struct {
	XMLName xml.Name `xml:"outgoingServer"`
	serverBase
}

type Description struct {
	XMLName xml.Name `xml:"descr"`
	Lang    string   `xml:"lang,attr"`
	Text    string   `xml:",chardata"`
}

type Documentation struct {
	XMLName      xml.Name      `xml:"documentation"`
	Url          string        `xml:"url,attr"`
	Descriptions []Description `xml:"descr"`
}

type EnableDoc struct {
	XMLName      xml.Name `xml:"enable"`
	VisitUrl     string   `xml:"visiturl,attr"`
	Instructions []string `xml:"instruction"`
}

type EmailProvider struct {
	XMLName          xml.Name         `xml:"emailProvider"`
	Id               string           `xml:"id,attr"`
	Domains          []string         `xml:"domain"`
	DisplayName      string           `xml:"displayName"`
	DisplayShortName string           `xml:"displayShortName"`
	IncomingServers  []IncomingServer `xml:"incomingServer"`
	OutgoingServers  []OutgoingServer `xml:"outgoingServer"`
	Enable           EnableDoc        `xml:"enable,omitempty"`
	Documentations   []Documentation  `xml:"documentation,omitempty"`
}

type ClientConfig struct {
	XMLName        xml.Name        `xml:"clientConfig"`
	Version        string          `xml:"version,attr"`
	EmailProviders []EmailProvider `xml:"emailProvider"`
}

type ConfigMap map[string]*EmailProvider
