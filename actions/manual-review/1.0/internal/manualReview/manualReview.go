package manualReview

import (
	"encoding/json"

	"fmt"
	"strconv"
	"time"

	"github.com/erda-project/erda-actions/actions/manual-review/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/httpclient"
	"github.com/erda-project/erda/pkg/httpserver/errorresp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func handleAPIs() error {
	//查询该任务之前是否建立过审核
	total, err := getTask(conf.TaskId())
	if err != nil {
		logrus.Warningf("not exist usecase test env info, total:%s, (%+v)", total, err)
	}
	if total != 0 {
		logrus.Errorf("manualReview failed")
	} else {
		//创建审核记录
		err := CreateReview()
		if err != nil {
			logrus.Warningf("create fail")
		}
		//创建用户审批权限
		for _, val := range conf.ProcessorId() {
			err = CreateReviewUser(val)
			if err != nil {
				logrus.Warningf("create fail")
			}
		}
		task_id := strconv.FormatUint(conf.TaskId(), 10)

		processor_id, _ := json.Marshal(conf.ProcessorId())

		ID, _ := getId(conf.TaskId())
		review_id := strconv.FormatUint(ID, 10)

		if err := storeMetaFile(task_id, string(processor_id), review_id); err != nil {
			return err
		}

		var approvalStatus string

		waitingTime := conf.WaitingTime()
		for true {
			time.Sleep(time.Duration(waitingTime) * time.Second)
			approvalStatus, _ = getReview(conf.TaskId())
			if approvalStatus == "Accept" {
				logrus.Errorf("manualReview succeed")
				break
			} else if approvalStatus == "Reject" {
				logrus.Warningf("create fail")
				err = fmt.Errorf("审批拒绝")
				return err
			}
		}
	}
	return nil
}

//通过taskId查询审核记录个数
func getTask(envID uint64) (int64, error) {
	// invoke
	var resp ManualReview
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/reviews/actions/%d", envID)).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)
	if err != nil {
		return 0, err
	}
	if !r.IsOK() {
		return 0, errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}

	return resp.Data.Total, nil
}

//通过taskId查询审核状态
func getReview(envID uint64) (string, error) {
	// invoke
	var resp ManualReview
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/reviews/actions/%d", envID)).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)
	if err != nil {
		return "", err
	}
	if !r.IsOK() {
		return "", errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}

	return resp.Data.Approvalstatus, nil
}

func getId(envID uint64) (uint64, error) {
	// invoke
	var resp ManualReview
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Get(conf.DiceOpenapiAddr()).
		Path(fmt.Sprintf("/api/reviews/actions/%d", envID)).
		Header("Authorization", conf.DiceOpenapiToken()).Do().JSON(&resp)
	if err != nil {
		return 0, err
	}
	if !r.IsOK() {
		return 0, errorresp.New(errorresp.WithCode(r.StatusCode(), resp.Error.Code), errorresp.WithMessage(resp.Error.Msg))
	}

	return resp.Data.Id, nil
}

//创建审核记录
func CreateReview() error {
	// invoke
	createReq := CreateReviewRequest{
		ProjectId:       conf.ProjectId(),
		BuildId:         conf.PipelineId(),
		ApplicationId:   conf.ApplicationId(),
		ApplicationName: conf.ApplicationName(),
		TaskId:          conf.TaskId(),
		OrgId:           conf.OrgId(),
		SponsorId:       conf.SponsorId(),
		CommitID:        conf.CommitId(),
		ProjectName:     conf.ProjectName(),
		BranchName:      conf.BranchName(),
		ApprovalStatus:  "WaitApprove",
	}
	var artifact CreateReviewResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path("/api/reviews/actions/review/approve").
		Header("Authorization", conf.DiceOpenapiToken()).
		JSONBody(&createReq).Do().JSON(&artifact)
	logrus.Error("pipipipipi:", artifact.Error)

	if !artifact.Success {
		return errorresp.New(errorresp.WithMessage(artifact.Error.Msg))
	}

	if err != nil {
		return err
	}
	if !r.IsOK() {
		return err
	}
	return nil
}

//创建审核
func CreateReviewUser(Operator string) error {
	createReq := CreateReviewUserRequest{
		TaskId:   conf.TaskId(),
		OrgId:    conf.OrgId(),
		Operator: Operator,
	}
	var artifact CreateReviewUserResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiAddr()).
		Path("/api/reviews/actions/user/create").
		Header("Authorization", conf.DiceOpenapiToken()).
		JSONBody(&createReq).Do().JSON(&artifact)
	logrus.Error("pipipipipi:", artifact.Error)

	if !artifact.Success {
		return errorresp.New(errorresp.WithMessage(artifact.Error.Msg))
	}

	if err != nil {
		return err
	}
	if !r.IsOK() {
		return err
	}
	return nil
}

func storeMetaFile(task_id string, processor_id string, review_id string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "task_id",
				Value: task_id,
			},
			{
				Name:  "processor_id",
				Value: processor_id,
			},
			{
				Name:  "review_id",
				Value: review_id,
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
