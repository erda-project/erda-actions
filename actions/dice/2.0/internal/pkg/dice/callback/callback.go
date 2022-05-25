package callback

import (
	"time"
	"strconv"
	"fmt"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
	"github.com/erda-project/erda-actions/pkg/metawriter"
)

func BatchReportRuntimeInfo(conf *conf.Conf, rets map[string]*common.DeployResult) error {
	if err := metawriter.WriteKV(apistructs.ActionCallbackOperatorID, conf.OperatorID); err != nil {
		logrus.Errorf("failed to write operator id to meta: %v", err)
	}

	runtimeIds := make([]string, 0, len(rets))
	for appName, result := range rets {
		runtimeId := strconv.FormatUint(result.RuntimeId, 10)
		runtimeIds = append(runtimeIds, runtimeId)
		if result.RuntimeId == 0 {
			logrus.Warningf("application %s runtime id is emty, can not report to ci", appName)
			continue
		}

		// TODO: report all applicationName_runtimeId format, link to deployment order info
		switch utils.ConvertType(conf.ReleaseTye) {
		case common.TypeProjectRelease, common.TypeApplicationRelease:
			if err := metawriter.WriteKV(fmt.Sprintf("%s_%s", appName, apistructs.ActionCallbackRuntimeID), runtimeId); err != nil {
				logrus.Errorf("failed to write %s_%s to meta: %v", appName, apistructs.ActionCallbackRuntimeID, err)
			}
		default:
			if err := metawriter.WriteLink(apistructs.ActionCallbackRuntimeID, runtimeId); err != nil {
				logrus.Errorf("failed to write runtime %s link to meta: %v", runtimeId, err)
			}
			break
		}
	}

	return nil
}

// Callback runtime
func Callback(url string, runtimeId uint64, applicationId uint64, options interface{}, status string) {
	if len(url) == 0 {
		logrus.Info("no callback set, return directly")
		return
	}
	logrus.Infof("start to callback with url=%s, runtimeId=%d, applicationId=%d, status=%s, options=%v ...", url, runtimeId, applicationId, status, options)
	data := struct {
		ApplicationId uint64      `json:"applicationId"`
		RuntimeId     uint64      `json:"runtimeId"`
		Status        string      `json:"status"`
		Options       interface{} `json:"options,omitempty"`
	}{
		ApplicationId: applicationId,
		RuntimeId:     runtimeId,
		Status:        status,
		Options:       options,
	}
	dj, err := json.Marshal(&data)
	if err != nil {
		logrus.Error(err)
	}
	resp, body, errs := gorequest.New().Post(url).Send(string(dj)).Timeout(5 * time.Second).End()

	if errs != nil {
		logrus.Error(multierror.Append(errors.New("callback failed"), errs...))
	}
	if resp.StatusCode/100 != 2 {
		logrus.Error(errors.Errorf("callback failed, response code: %d, response body: %s",
			resp.StatusCode, body))
	}
	logrus.Infof("callback successfully! response code=%d", resp.StatusCode)
}
