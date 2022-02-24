package common

import (
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/pkg/log"
)

const (
	Authorization              = "Authorization"
	DeploymentOrderRequestPath = "/api/deployment-orders"
	TypeApplicationRelease     = "APPLICATION_RELEASE"
	TypeProjectRelease         = "PROJECT_RELEASE"
	SourcePipeline             = "PIPELINE"
)

const (
	CallbackStatusFailed  = "Failed"
	CallbackStatusSuccess = "Success"
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
	Type                string `json:"type"`
	ReleaseId           string `json:"releaseId"`
	ReleaseName         string `json:"releaseName"`
	ProjectId           uint64 `json:"projectId"`
	ApplicationName     string `json:"applicationName"`
	Workspace           string `json:"workspace"`
	AutoRun             bool   `json:"autoRun"`
	DeployWithoutBranch bool   `json:"deployWithoutBranch"`
	Source              string `json:"source"`
}

func (d *CreateDeploymentOrderRequest) Print() {
	log.AddNewLine(1)
	logrus.Infof("request deploy body: ")
	logrus.Infof(" worspace: %s", d.Workspace)
	logrus.Infof(" autoRun: %v", d.AutoRun)
	logrus.Infof(" source: %s", SourcePipeline)

	switch d.Type {
	case TypeApplicationRelease, TypeProjectRelease:
		logrus.Infof(" type: %s", d.Type)
		logrus.Infof(" releaseName: %s", d.ReleaseName)
		if d.Type == TypeApplicationRelease {
			logrus.Infof(" application_name: %s", d.ApplicationName)
		} else {
			logrus.Infof(" projectId: %d", d.ProjectId)
		}
	case "":
		logrus.Infof(" releaseId: %s", d.ReleaseId)
		logrus.Infof(" deployWithoutBranch: %v", d.DeployWithoutBranch)
	}

	log.AddLineDelimiter(" ")
}

type CreateDeploymentOrderResponse struct {
	Response
	Data struct {
		DeploymentOrderId string                    `json:"id"`
		Deployments       map[string]DeploymentInfo `json:"deployments"`
	} `json:"data"`
}

type DeploymentInfo struct {
	DeploymentID  uint64 `json:"deploymentId"`
	ApplicationID uint64 `json:"applicationId"`
	RuntimeID     uint64 `json:"runtimeId"`
}

type DeploymentOrderStatusRespData struct {
	Response
	Data struct {
		BatchSize    uint64 `json:"batchSize"`
		CurrentBatch uint64 `json:"currentBatch"`
		Status       string `json:"status"`
	} `json:"data"`
}

func (d *DeploymentOrderStatusRespData) Print() {
	log.AddNewLine(1)
	logrus.Infof("response deploy status body: ")
	logrus.Infof(" batchSize: %d", d.Data.BatchSize)
	logrus.Infof(" currentBatch: %d", d.Data.CurrentBatch)
	logrus.Infof(" status: %s", d.Data.Status)
	log.AddLineDelimiter(" ")
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
	DeploymentId  uint64
	ApplicationId uint64
	RuntimeId     uint64
}

type DeployErrResponse struct {
	Msg string
}

func (d *DeployErrResponse) Error() string {
	return d.Msg
}

type CancelRequest struct {
	DeploymentOrderId string
	Operator          string
	Force             bool `json:"force"`
}
