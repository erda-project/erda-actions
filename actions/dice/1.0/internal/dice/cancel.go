package dice

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
)

type CancelReq struct {
	RuntimeId int64  `json:"runtimeId"`
	Operator  string `json:"operator"`
}

func Cancel() error {
	var cfg conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Info(cfg)

	var (
		runtimeId int64
	)
	logrus.Info("now we are going to cancel the task...")
	deploymentIdStr, runtimeIdStr, err := getDiceInfo(cfg.WorkDir)
	if err != nil {
		return err
	}

	runtimeId, err = strconv.ParseInt(runtimeIdStr, 10, 64)
	if err != nil {
		return err
	}

	cReq := &CancelReq{
		RuntimeId: runtimeId,
		Operator:  fmt.Sprintf("%v", cfg.OperatorID),
	}

	err = cancelRequest(cReq, deploymentIdStr, &cfg)
	if err != nil {
		return errors.Wrapf(err, "failed to cancel deployment with req=%+v", cReq)
	}
	return nil
}

func getDiceInfo(wd string) (string, string, error) {
	fileValue, err := ioutil.ReadFile(filepath.Join(wd, "diceInfo"))
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to read file diceInfo")
	}

	if fileValue == nil {
		return "", "", errors.New("null diceInfo content")
	}

	diceList := strings.Split(string(fileValue), ",")
	if len(diceList) != 2 {
		return "", "", errors.New("failed to split diceInfo content")
	}

	deploymentIdInfo := strings.Split(diceList[0], "=")
	runtimeIdInfo := strings.Split(diceList[1], "=")
	if len(deploymentIdInfo) != 2 || len(runtimeIdInfo) != 2 {
		return "", "", errors.New("failed to get deploymentId and runtimeId")
	}

	return deploymentIdInfo[1], runtimeIdInfo[1], nil
}

func cancelRequest(cancelReq *CancelReq, deploymentId string, conf *conf) error {
	var diceResp DiceResponse
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).Path(fmt.Sprintf("/api/deployments/%s/actions/cancel", deploymentId)).
		Header(Authorization, conf.DiceOpenapiToken).JSONBody(&cancelReq).Do().JSON(&diceResp)
	if err != nil {
		return err
	}
	if !r.IsOK() {
		return errors.Errorf("failed to cancel dice deploy(id=%s). statusCode: %d, respCode=%s, message=%s, ctx=%v, DiceResponse:%+v",
			deploymentId, r.StatusCode(), diceResp.Err.Code, diceResp.Err.Message, diceResp.Err.Ctx, diceResp)
	}
	if !diceResp.Success {
		return errors.Errorf("failed to cancel dice deploy(id=%s). code=%s, message=%s, ctx=%v",
			deploymentId, diceResp.Err.Code, diceResp.Err.Message, diceResp.Err.Ctx)
	}

	logrus.Infof(">>> success to cancel dice deploy.")

	return nil
}
