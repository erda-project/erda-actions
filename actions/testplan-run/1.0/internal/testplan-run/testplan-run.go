package testplan_run

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/testplan-run/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/encoding/jsonparse"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/strutil"
)

func handleAPIs() error {

	logrus.Info("start execution")
	pipelineDTO, err := ExecuteDiceAutotestTestPlans(conf.TestPlan(), conf.Cms())
	if err != nil {
		logrus.Info("Failed execution of the plan")
		return err
	}
	logrus.Info("test plan is being executed")
	logrus.Info("pipeline status ", pipelineDTO.Status)

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
		logrus.Info("pipeline status ", pipelineDTO.Status)

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
			err = storeMetaFile(dto.ID, dto.Status.String(), runtimeIDs)
			if err != nil {
				err = fmt.Errorf("上报执行信息失败")
				if conf.IsContinueExecution() {
					return nil
				}
				return err
			}

			if dto.Status.IsFailedStatus() {
				err = fmt.Errorf("执行失败")
				if conf.IsContinueExecution() {
					return nil
				}
				return err
			}

			return nil
		}

		time.Sleep(10 * time.Second)
	}
}

// 执行计划
func ExecuteDiceAutotestTestPlans(TestPlanID uint64, ns string) (*apistructs.PipelineDTO, error) {
	// invoke
	var req apistructs.AutotestExecuteTestPlansRequest
	req.TestPlan.ID = TestPlanID
	req.ClusterName = conf.DiceClusterName()
	req.ConfigManageNamespaces = ns
	var resp apistructs.AutotestExecuteTestPlansResponse
	response, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/autotests/testplans/%v/actions/execute", req.TestPlan.ID)).
		Header("Authorization", conf.DiceOpenapiToken()).
		JSONBody(&req).Do().JSON(&resp)

	if err != nil {
		return nil, fmt.Errorf("failed to request (%s)", err.Error())
	}

	if !response.IsOK() {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, header: %s, msg: %s", response.Headers(), response.Body()))
	}

	if !resp.Success {
		return nil, fmt.Errorf(fmt.Sprintf("failed to request, error code: %s, error message: %s", resp.Error.Code, resp.Error.Msg))
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
