package testplan

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/loop"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/testplan/1.0/internal/conf"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type UsecaseIDs struct {
	usecaseIDs []uint64
}

func Run() error {
	if err := conf.Load(); err != nil {
		return err
	}

	return execTestPlan()
}

func execTestPlan() error {
	usecaseIDs := &UsecaseIDs{}
	// invoke
	params := make(url.Values)
	params.Add("projectId", strconv.FormatUint(conf.ProjectID(), 10))
	params.Add("projectTestEnvID", strconv.FormatUint(conf.ProjectTestEnvID(), 10))
	var body bytes.Buffer
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/pmp/api/testing/testplan/%d/run-usecase-apis", conf.TestPlanID())).
		Header("Authorization", conf.DiceOpenapiToken()).
		Params(params).
		JSONBody(usecaseIDs).Do().Body(&body)
	if err != nil {
		return errors.Errorf("failed to run test plan, err:%v, url:%s, body: %s",
			err, conf.DiceOpenapiAddr(), body.String())
	}
	if !r.IsOK() {
		return errors.Errorf("fialed to tun test plan, statusCode: %d, body: %s.", r.StatusCode(), body.String())
	}

	logrus.Info("正在执行测试计划，请等待...")

	// 获取最新的 pipeline
	pipelineReq := apistructs.PipelinePageListRequest{
		PageNum:         1,
		PageSize:        1,
		Sources:         []apistructs.PipelineSource{apistructs.PipelineSourceAPITest},
		MustMatchLabels: map[string][]string{"testPlanID": {"48"}},
		YmlNames:        []string{fmt.Sprintf("%s-%d.yml", apistructs.PipelineSourceAPITest, conf.ProjectID())},
	}
	var pageResp apistructs.PipelinePageListResponse
	httpResp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).Path("/api/pipelines").
		Header("Authorization", conf.DiceOpenapiToken()).
		Params(pipelineReq.UrlQueryString()).
		Do().JSON(&pageResp)
	if err != nil {
		return errors.Errorf("failed to get latest pipeline, err:%v, url:%s, body: %+v",
			err, conf.DiceOpenapiAddr(), pageResp)
	}
	if !httpResp.IsOK() || !pageResp.Success {
		return errors.Errorf("failed to get latest pipeline, statusCode: %d, body: %s.",
			httpResp.StatusCode(), pageResp.Error)
	}

	pipelineInfo := pageResp.Data.Pipelines
	if len(pipelineInfo) == 0 {
		return errors.Errorf("nil pipeline")
	}

	// 每隔1秒，轮询获取测试计划执行结果
	var pipelineResp apistructs.PipelineDetailResponse

	l := loop.New(loop.WithInterval(time.Second))
	err = l.Do(func() (bool, error) {
		httpResp, err = httpclient.New(httpclient.WithCompleteRedirect()).
			Get(conf.DiceOpenapiAddr()).Path(fmt.Sprintf("/api/apitests/pipeline/%d", pipelineInfo[0].ID)).
			Header("Authorization", conf.DiceOpenapiToken()).
			Do().JSON(&pipelineResp)
		if err != nil {
			return true, errors.Errorf("failed to get pipeline info, err:%v, url:%s, body: %+v, pipeline: %d",
				err, conf.DiceOpenapiAddr(), pipelineResp, pipelineInfo[0].ID)
		}
		if !httpResp.IsOK() || !pipelineResp.Success {
			return true, errors.Errorf("failed to get pipeline info, statusCode: %d, body: %+v",
				httpResp.StatusCode(), pipelineResp)
		}

		logrus.Infof("Check pipeline result, status: %s", pipelineResp.Data.Status)

		if pipelineResp.Data.Status.IsEndStatus() {
			if pipelineResp.Data.Status.IsFailedStatus() {
				return true, errors.Errorf("pipeline status error, status: %s", pipelineResp.Data.Status)
			}
			logrus.Infof("Finish to run pipeline, status: %s", pipelineResp.Data.Status)
			return true, nil
		}

		return false, nil
	})

	return err
}
