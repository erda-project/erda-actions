package appRun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/app-run/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/httpclient"
)

func handlerPipelineYmlName(ymlName string) string {

	if ymlName != "pipeline.yml" && !strings.HasPrefix(ymlName, ".dice/pipelines") {
		return ".dice/pipelines/" + ymlName
	}

	return ymlName
}

func handleAPIs() error {

	// get application name
	applications, err := getApplicationList()
	if err != nil {
		return err
	}
	var existApp *apistructs.ApplicationDTO
	for _, app := range applications {
		if strings.EqualFold(app.Name, conf.ActionApplicationName()) {
			existApp = &app
			break
		}
	}
	if existApp == nil {
		return fmt.Errorf("not find application name %s", conf.ActionApplicationName())
	}

	// start pipeline
	logrus.Infof("start run pipeline %s", conf.ActionPipelineYmlName())
	pipelineDTO, err := startPipeline(apistructs.PipelineCreateRequest{
		AppID:             existApp.ID,
		Branch:            conf.ActionBranch(),
		PipelineYmlSource: apistructs.PipelineYmlSourceGittar,
		PipelineYmlName:   handlerPipelineYmlName(conf.ActionPipelineYmlName()),
		Source:            apistructs.PipelineSourceDice,
		AutoRun:           true,
	})
	if err != nil {
		return err
	}
	logrus.Infof("end run application %s", conf.ActionPipelineYmlName())

	// watch pipeline done
	for {
		dto, err := pipelineSimpleDetail(PipelineDetailRequest{
			SimplePipelineBaseResult: true,
			PipelineID:               pipelineDTO.ID,
		})
		if err != nil {
			return err
		}

		if dto.Status.IsEndStatus() {
			return storeMetaFile(dto.ID, dto.Status.String())
		}

		time.Sleep(10 * time.Second)
	}
}

func getApplicationList() ([]apistructs.ApplicationDTO, error) {

	var resp apistructs.ApplicationListResponse
	var b bytes.Buffer

	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/applications")).
		Param("projectId", strconv.FormatUint(conf.ProjectId(), 10)).
		Param("pageSize", "9999").
		Param("pageNo", "1").
		Param("q", conf.ActionApplicationName()).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, status-code: %d, content-type: %s, raw bod: %s", response.StatusCode(), response.ResponseHeader("Content-Type"), b.String()))
	}

	if !resp.Success {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, error code: %s, error message: %s", resp.Error.Code, resp.Error.Msg))
	}

	if resp.Data.Total == 0 {
		return nil, nil
	}

	return resp.Data.List, nil
}

func startPipeline(req apistructs.PipelineCreateRequest) (*apistructs.PipelineDTO, error) {
	var resp apistructs.PipelineCreateResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path("/api/cicds").
		Header("Authorization", conf.DiceOpenapiToken()).
		JSONBody(&req).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("start pipeline error %s", err)
	}

	if !resp.Success {
		return nil, fmt.Errorf("start pipeline not success %s", resp.Error.Msg)
	}

	if !r.IsOK() {
		return nil, fmt.Errorf("start pipeline failed")
	}

	return resp.Data, nil
}

type PipelineDetailRequest struct {
	SimplePipelineBaseResult bool   `json:"simplePipelineBaseResult"`
	PipelineID               uint64 `json:"pipelineID"`
}

func pipelineSimpleDetail(req PipelineDetailRequest) (*apistructs.PipelineDetailDTO, error) {

	var resp apistructs.PipelineDetailResponse
	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiPublicUrl()).
		Path("/api/cicds/actions/pipeline-detail").
		Param("simplePipelineBaseResult", strconv.FormatBool(req.SimplePipelineBaseResult)).
		Param("pipelineId", strconv.FormatUint(req.PipelineID, 10)).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, status-code: %d, content-type: %s", response.StatusCode(), response.ResponseHeader("Content-Type")))
	}

	if !resp.Success {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, error code: %s, error message: %s", resp.Error.Code, resp.Error.Msg))
	}

	return resp.Data, nil
}

func storeMetaFile(pipelineID uint64, status string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "pipelineID",
				Value: strconv.FormatUint(pipelineID, 10),
			},
			{
				Name:  "status",
				Value: status,
			},
		},
	}

	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(conf.MetaFile(), string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}
