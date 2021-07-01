package testing

import (
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/qaparser"

	"github.com/erda-project/erda-actions/actions/integration-test/1.0/internal/conf"
)

func Exec(cfg *conf.Conf) error {
	var (
		suites    []*apistructs.TestSuite
		suite     *apistructs.TestSuite
		err       error
		qaID      string
		utResults *apistructs.TestCallBackRequest
	)

	if suite, err = MavenTest(cfg); err != nil {
		return err
	}

	suites = append(suites, suite)

	if utResults, err = makeItResults(suites, cfg); err != nil {
		return err
	}

	if qaID, err = callback(utResults, cfg); err != nil {
		return err
	}

	return storeMetaFile(cfg, qaID)
}

// callback to qa
func callback(req *apistructs.TestCallBackRequest, cfg *conf.Conf) (string, error) {
	var result = struct {
		Success bool   `json:"success"`
		Data    string `json:"data"`
		Err     struct {
			Code    string                 `json:"code,omitempty"`
			Message string                 `json:"msg,omitempty"`
			Ctx     map[string]interface{} `json:"ctx,omitempty"`
		} `json:"err,omitempty"`
	}{}
	resp, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(cfg.DiceOpenapiPrefix).Path("/api/qa/actions/test-callback").
		Header("Authorization", cfg.DiceOpenapiToken).
		Header("Content-Type", "application/json").
		JSONBody(&req).Do().JSON(&result)

	if err != nil {
		return "", errors.Wrapf(err, "failed to report results to qa, req: %+v", req)
	}

	if !resp.IsOK() {
		return "", errors.Errorf("failed to report results to qa, code: %d, req: %+v, result: %+v",
			resp.StatusCode(), req, result)
	}

	// 一般不会发生
	if result.Err.Code != "" {
		return "", errors.Errorf("failed to report results to qa, (%+v)", result.Err)
	}

	logrus.Infof("successed to report results to qa, req: %+v, qaID: %s", req, result.Data)

	return result.Data, nil
}

func storeMetaFile(cfg *conf.Conf, qaID string) error {
	meta := apistructs.ActionCallback{
		Metadata: apistructs.Metadata{
			{
				Name:  "projectId",
				Value: strconv.FormatUint(cfg.ProjectID, 10),
			},
			{
				Name:  "AppId",
				Value: strconv.FormatUint(cfg.AppID, 10),
			},
			{
				Name:  "operatorId",
				Value: cfg.OperatorID,
			},
			{
				Name:  "commitId",
				Value: cfg.GittarCommit,
			},
			{
				Name:  apistructs.ActionCallbackQaID,
				Value: qaID,
				Type:  apistructs.ActionCallbackTypeLink,
			},
		},
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := filehelper.CreateFile(cfg.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}

func makeItResults(suites []*apistructs.TestSuite, cfg *conf.Conf) (*apistructs.TestCallBackRequest, error) {
	results := &apistructs.TestCallBackRequest{
		Results: &apistructs.TestResults{
			Extra: make(map[string]string),
		},
		Totals: &apistructs.TestTotals{
			Statuses: make(map[apistructs.TestStatus]int),
		},
	}

	results.Suites = suites
	calculateTotals(suites, results)

	var name string
	if cfg.Name != "" {
		name = cfg.Name
	}
	composeResults(results.Results, cfg, name)

	return results, nil
}

func calculateTotals(suites []*apistructs.TestSuite, totals *apistructs.TestCallBackRequest) {
	if totals.Totals == nil {
		totals.Totals = &apistructs.TestTotals{
			Statuses: make(map[apistructs.TestStatus]int),
		}
	}
	for _, s := range suites {
		to := &qaparser.Totals{totals.Totals}
		totals.Totals = to.Add(s.Totals).TestTotals
	}
}

func composeResults(results *apistructs.TestResults, cfg *conf.Conf, name string) error {
	if err := composeEnv(results, cfg); err != nil {
		return err
	}

	if name == "" {
		if len(results.CommitID) > 6 {
			results.Name = results.CommitID[:6]
		} else {
			results.Name = results.CommitID
		}
	} else {
		results.Name = name
	}

	results.Status = "FINISHED"
	results.Type = "IT"

	return nil
}

func composeEnv(results *apistructs.TestResults, cfg *conf.Conf) error {
	results.OperatorID = cfg.OperatorID
	results.OperatorName = cfg.OperatorName
	results.ApplicationID = int64(cfg.AppID)
	results.ProjectID = int64(cfg.ProjectID)
	results.ApplicationName = cfg.AppName
	results.BuildID = cfg.BuildID
	results.GitRepo = cfg.GittarRepo
	results.Branch = cfg.GittarBranch
	results.CommitID = cfg.GittarCommit
	results.Workspace = cfg.Workspace
	results.UUID = cfg.UUID

	return nil
}
