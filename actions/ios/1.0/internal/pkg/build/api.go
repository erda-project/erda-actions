package build

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/http/httpclient"

	"github.com/erda-project/erda-actions/actions/ios/1.0/internal/conf"
)

func QueryTask(cfg conf.Conf, taskID string) (*apistructs.RunnerTask, error) {
	var resp apistructs.QueryRunnerTaskResponse
	request := httpclient.New(httpclient.WithCompleteRedirect()).Get(cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/runner/tasks/%s", taskID)).
		Header("Authorization", cfg.CiOpenapiToken)
	httpResp, err := request.Do().JSON(&resp)
	if err != nil {
		return nil, err
	}

	if !httpResp.IsOK() {
		return nil, errors.Errorf("failed to query runner task, status code: %d body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return nil, errors.Errorf(resp.Error.Msg)
	}
	return &resp.Data, nil
}

func UpdateTaskStatus(cfg conf.Conf, taskID int64, status string) error {
	request := httpclient.New(httpclient.WithCompleteRedirect()).Put(cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/runner/tasks/%d", taskID)).
		Header("Authorization", cfg.CiOpenapiToken)
	httpResp, err := request.JSONBody(map[string]string{
		"status": status}).Do().DiscardBody()
	if err != nil {
		return err
	}

	if !httpResp.IsOK() {
		return errors.Errorf("failed to query runner task, status code: %d body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}

	return nil
}

func CreateTask(cfg conf.Conf, req *apistructs.CreateRunnerTaskRequest) (int64, error) {
	var resp apistructs.CreateRunnerTaskResponse
	request := httpclient.New(httpclient.WithCompleteRedirect()).Post(cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/runner/tasks")).
		Header("Authorization", cfg.CiOpenapiToken)
	httpResp, err := request.JSONBody(req).Do().JSON(&resp)
	if err != nil {
		return 0, err
	}

	if !httpResp.IsOK() {
		return 0, errors.Errorf("failed to create runner task, status code: %d body:%s", httpResp.StatusCode(), string(httpResp.Body()))
	}
	if !resp.Success {
		return 0, errors.Errorf(resp.Error.Msg)
	}
	return resp.Data, nil
}
