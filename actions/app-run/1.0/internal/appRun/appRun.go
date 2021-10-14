package appRun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/app-run/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
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

	logrus.Infof("start run pipeline %s", conf.ActionPipelineYmlName())
	// start pipeline
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
	logrus.Infof("end run pipeline %s", conf.ActionPipelineYmlName())

	logrus.Infof("wait pipeline done %s", pipelineDTO.ID)

	err = storePipelineInfo(pipelineDTO.ID)
	if err != nil {
		return err
	}

	// watch pipeline done
	for {
		dto, err := pipelineSimpleDetail(PipelineDetailRequest{
			SimplePipelineBaseResult: true,
			PipelineID:               pipelineDTO.ID,
		})
		if err != nil {
			fmt.Printf(" get pipelineSimpleDetail error %v \n", err)
			time.Sleep(10*time.Second)
			continue
		}

		if dto.Status.IsEndStatus() {
			// get detail info
			dto, err := pipelineSimpleDetail(PipelineDetailRequest{
				PipelineID: pipelineDTO.ID,
			})
			if err != nil {
				fmt.Printf(" get pipelineDetail error %v \n", err)
				time.Sleep(10*time.Second)
				continue
			}

			logrus.Infof("pipeline %s was done status %v", pipelineDTO.ID, dto.Status.String())

			runtimeIDs := getDiceTaskRuntimeIDs(dto)

			return storeMetaFile(dto.ID, dto.Status.String(), runtimeIDs)
		}

		time.Sleep(10 * time.Second)
	}
}

func getDiceTaskRuntimeIDs(dto *apistructs.PipelineDetailDTO) []string {

	var runtimeIDs []string
	if dto == nil || dto.PipelineStages == nil {
		return runtimeIDs
	}
	for _, stage := range dto.PipelineStages {
		for _, task := range stage.PipelineTasks {
			if task.Type != "dice" || task.Result.Metadata == nil {
				continue
			}
			for _, meta := range task.Result.Metadata {
				if meta.Name == "runtimeID" {
					runtimeIDs = append(runtimeIDs, meta.Value)
				}
			}
		}
	}

	return runtimeIDs
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

func storePipelineInfo(pipelineID uint64) error {
	content := fmt.Sprint("pipelineID=", pipelineID)
	err := filehelper.CreateFile(filepath.Join(conf.WorkDir(), "pipelineInfo"), content, 0755)
	if err != nil {
		return errors.Wrap(err, "write file:pipelineInfo failed")
	}
	return nil
}

func storeMetaFile(pipelineID uint64, status string, runtimeID []string) error {
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

	if len(runtimeID) > 0 {
		meta.Metadata = append(meta.Metadata, apistructs.MetadataField{
			Name:  "runtimeIDs",
			Value: jsonparse.JsonOneLine(strutil.DedupSlice(runtimeID)),
		})

		meta.Metadata = append(meta.Metadata, apistructs.MetadataField{
			Name:  "runtimeID",
			Value: runtimeID[0],
		})
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
