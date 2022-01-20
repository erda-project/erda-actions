package utils

import (
	"encoding/json"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/sirupsen/logrus"
	"github.com/erda-project/erda-actions/pkg/log"
)

func BatchPrintStatusCheckResult(respMap map[string]*common.DeploymentStatusRespData) {
	log.AddNewLine(1)
	logrus.Infof(" check deploy status: ")
	for appName, r := range respMap {
		logrus.Infof(" %s", appName)
		if r == nil {
			logrus.Info("  result: failed to print deployment status, response data is nil")
			continue
		}
		logrus.Infof("  success: %v", r.Success)
		logrus.Infof("  deploymentID: %v", r.Data.DeploymentId)
		logrus.Infof("  status: %v", r.Data.Status)
		if r.Data.FailCause != "" {
			logrus.Infof("  failCause: %v", r.Data.FailCause)
		}
		if len(r.Data.ModuleErrMsg) > 0 {
			for k, v := range r.Data.ModuleErrMsg {
				if v != "" {
					logrus.Infof("  %s: %s", k, v)
				}
			}
		}
		if r.Data.Runtime != nil {
			b, err := json.MarshalIndent(r.Data.Runtime, "", " ")
			if err != nil {
				logrus.Errorf("fail to json marshal: err: %v", err)
			}
			logrus.Infof("  runtime: %s", string(b))
		}
		if r.Err.Code != "" {
			logrus.Infof("  err code: %s", r.Err.Code)
		}
		if r.Err.Message != "" {
			logrus.Infof("  err message: %s", r.Err.Message)
		}
		if r.Err.Ctx != nil {
			for k, v := range r.Err.Ctx {
				logrus.Infof("  err ctx %s: %v", k, v)
			}
		}
	}

	log.AddLineDelimiter(" ")
}
