package common

import (
	"github.com/mark3labs/mcp-go/mcp"
)

const (
	Authorization      = "Authorization"
	ReleaseRequestPath = "/api/releases"
)

const (
	AnnotationMcpDescription = "mcp.erda.cloud/description"
	AnnotationMcpConnectURI  = "mcp.erda.cloud/connect-uri"

	LabelMcpName          = "mcp.erda.cloud/name"
	LabelMcpVersion       = "mcp.erda.cloud/version"
	LabelMcpIsPublished   = "mcp.erda.cloud/is-published"
	LabelMcpIsDefault     = "mcp.erda.cloud/is-default"
	LabelMcpTransportType = "mcp.erda.cloud/transport-type"
	LabelMcpServicePort   = "mcp.erda.cloud/service-port"
)

// erda standard response struct

type Response struct {
	Success bool `json:"success"`
	Err     Err  `json:"err,omitempty"`
}

type Err struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"msg,omitempty"`
	Ctx     map[string]interface{} `json:"ctx,omitempty"`
}

type GetReleaseResponse struct {
	Response
	Data struct {
		DiceYaml string `json:"diceyml"`
	} `json:"data"`
}

type MCPServerRegisterRequest struct {
	Name             string     `json:"name"`
	Description      string     `json:"description"`
	Version          string     `json:"version"`
	Endpoint         string     `json:"endpoint"`
	TransportType    string     `json:"transportType"`
	ServerConfig     string     `json:"serverConfig"`
	Tools            []mcp.Tool `json:"tools"`
	IsPublished      bool       `json:"isPublished,omitempty"`
	IsDefaultVersion bool       `json:"isDefault,omitempty"`
}

type Service struct {
	Host  string `json:"host"`
	Ports []int  `json:"ports"`
}

type ServiceInfo struct {
	Services map[string]Service
}
