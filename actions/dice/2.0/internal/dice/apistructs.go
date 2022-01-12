package dice

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
)

const (
	Authorization              = "Authorization"
	DeploymentOrderRequestPath = "/api/deployment-orders"
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

type CreateDeploymentOrderRequest struct {
	Type      string `json:"type,omitempty"`
	ReleaseId string `json:"releaseId"`
	Workspace string `json:"workspace,omitempty"`
	AutoRun   bool   `json:"autoRun"`
}

func (d *CreateDeploymentOrderRequest) print() {
	log.AddNewLine(1)
	logrus.Infof("request deploy body: ")
	logrus.Infof(" releaseId: %s", d.ReleaseId)
	logrus.Infof(" type: %s", d.Type)
	logrus.Infof(" worspace: %s", d.Workspace)
	logrus.Infof(" autoRun: %v", d.AutoRun)
	log.AddLineDelimiter(" ")
}

type CreateDeploymentOrderResponse struct {
	Response
	Data struct {
		DeploymentOrderId string                    `json:"id"`
		Deployments       map[uint64]DeploymentInfo `json:"deployments"`
	} `json:"data"`
}

type DeploymentInfo struct {
	DeploymentID  int64 `json:"deploymentId"`
	ApplicationID int64 `json:"applicationId"`
	RuntimeID     int64 `json:"runtimeId"`
}

type DeploymentStatusRespData struct {
	Response
	Data struct {
		DeploymentId int               `json:"deploymentId"`
		Status       string            `json:"status"`
		FailCause    string            `json:"failCause"`
		ModuleErrMsg map[string]string `json:"lastMessage"`
		Runtime      interface{}       `json:"runtime"`
	} `json:"data"`
}

func (r *DeploymentStatusRespData) Print() {
	log.AddNewLine(1)
	logrus.Infof(" check deploy status: ")
	logrus.Infof(" success: %v", r.Success)
	logrus.Infof(" deploymentID: %v", r.Data.DeploymentId)
	logrus.Infof(" status: %v", r.Data.Status)
	if r.Data.FailCause != "" {
		logrus.Infof(" failCause: %v", r.Data.FailCause)
	}
	if len(r.Data.ModuleErrMsg) > 0 {
		for k, v := range r.Data.ModuleErrMsg {
			if v != "" {
				logrus.Infof(" %s: %s", k, v)
			}
		}
	}
	if r.Data.Runtime != nil {
		b, err := json.MarshalIndent(r.Data.Runtime, "", " ")
		if err != nil {
			logrus.Errorf("fail to json marshal: err: %v", err)
		}
		logrus.Infof(" runtime: %s", string(b))
	}
	if r.Err.Code != "" {
		logrus.Infof(" err code: %s", r.Err.Code)
	}
	if r.Err.Message != "" {
		logrus.Infof(" err message: %s", r.Err.Message)
	}
	if r.Err.Ctx != nil {
		for k, v := range r.Err.Ctx {
			logrus.Infof(" err ctx %s: %v", k, v)
		}
	}

	log.AddLineDelimiter(" ")
}

type DeployResult struct {
	DeploymentId  int64
	ApplicationId int64
	RuntimeId     int64
}

type DeployErrResponse struct {
	s string
}

func (d *DeployErrResponse) Error() string {
	return d.s
}
