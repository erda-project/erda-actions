package callback

import (
	"time"
	"strconv"
	"fmt"
	"strings"
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/retry"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/http/httputil"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
)

func BatchReportRuntimeInfo(conf *conf.Conf, rets map[string]*common.DeployResult) error {
	runtimeIds := make([]string, 0)

	cb := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{Name: apistructs.ActionCallbackOperatorID, Value: conf.OperatorID},
		},
		PipelineID:     conf.PipelineBuildID,
		PipelineTaskID: conf.PipelineTaskID,
	}

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
			cb.Metadata = append(cb.Metadata, apistructs.MetadataField{
				Name:  fmt.Sprintf("%s_%s", appName, apistructs.ActionCallbackRuntimeID),
				Value: runtimeId,
			})
		default:
			cb.Metadata = append(cb.Metadata, apistructs.MetadataField{
				Name:  apistructs.ActionCallbackRuntimeID,
				Value: runtimeId,
				Type:  apistructs.ActionCallbackTypeLink,
			})
			break
		}
	}

	b, err := json.Marshal(&cb)
	if err != nil {
		return err
	}

	var cbReq apistructs.PipelineCallbackRequest
	cbReq.Type = string(apistructs.PipelineCallbackTypeOfAction)
	cbReq.Data = b

	if err = retry.DoWithInterval(func() error {
		var resp apistructs.PipelineCallbackResponse
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).
			Post(conf.DiceOpenapiPrefix).
			Path("/api/pipelines/actions/callback").
			Header("Authorization", conf.DiceOpenapiToken).
			Header(httputil.UserHeader, conf.UserID).
			Header(httputil.InternalHeader, conf.InternalClient).
			JSONBody(&cbReq).
			Do().
			JSON(&resp)
		if err != nil {
			return err
		}
		if !r.IsOK() || !resp.Success {
			return errors.Errorf("status-code %d, resp %#v", r.StatusCode(), resp)
		}
		logrus.Infof("report runtimeID %s to pipeline platform successfully!", strings.Join(runtimeIds, ","))
		return nil
	}, 3, time.Millisecond*500); err != nil {
		logrus.Infof("report runtimeID to pipeline platform failed! err: %v", err)
		return err
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
