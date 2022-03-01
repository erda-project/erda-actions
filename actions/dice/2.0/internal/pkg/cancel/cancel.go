package cancel

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
)

func ExeWithConfig(orderId string, cfg *conf.Conf) {
	cReq := &common.CancelRequest{
		DeploymentOrderId: orderId,
		Operator:          fmt.Sprintf("%v", cfg.OperatorID),
		Force:             true,
	}

	if err := cancelRequest(cReq, cfg); err != nil {
		logrus.Errorf("failed to cancel deployment with req=%+v, err: %v", cReq, err)
	}
}

func Cancel() error {
	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}

	logrus.Info("now we are going to cancel the task...")

	orderId, err := getDiceInfo(cfg.WorkDir)
	if err != nil {
		logrus.Errorf("failed to get dice info, err: %v", err)
		return err
	}

	cReq := &common.CancelRequest{
		DeploymentOrderId: orderId,
		Operator:          fmt.Sprintf("%v", cfg.OperatorID),
		Force:             true,
	}

	if err := cancelRequest(cReq, &cfg); err != nil {
		return errors.Wrapf(err, "failed to cancel deployment with req=%+v", cReq)
	}
	return nil
}

func getDiceInfo(workdir string) (string, error) {
	fileValue, err := ioutil.ReadFile(filepath.Join(workdir, "diceInfo"))
	if err != nil {
		return "", errors.Wrapf(err, "failed to read file diceInfo")
	}

	if fileValue == nil {
		return "", errors.New("null diceInfo content")
	}

	diceList := strings.Split(string(fileValue), ",")
	if len(diceList) < 1 {
		return "", errors.New("failed to split diceInfo content")
	}

	orderId := strings.Split(diceList[0], "=")
	if len(orderId) != 2 {
		return "", errors.New("failed to get deploymentOrderId")
	}

	return orderId[1], nil
}

func cancelRequest(cancelReq *common.CancelRequest, conf *conf.Conf) error {
	var diceResp common.Response
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).
		Post(conf.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/deployment-orders/%s/actions/cancel", cancelReq.DeploymentOrderId)).
		Header(common.Authorization, conf.DiceOpenapiToken).
		JSONBody(&cancelReq).Do().JSON(&diceResp)
	if err != nil {
		return err
	}
	if !r.IsOK() {
		return errors.Errorf("failed to cancel dice deploy(id=%s). statusCode: %d, respCode=%s, message=%s, ctx=%v, "+
			"DiceResponse:%+v", cancelReq.DeploymentOrderId, r.StatusCode(), diceResp.Err.Code, diceResp.Err.Message,
			diceResp.Err.Ctx, diceResp)
	}
	if !diceResp.Success {
		return errors.Errorf("failed to cancel dice deploy(id=%s). code=%s, message=%s, ctx=%v",
			cancelReq.DeploymentOrderId, diceResp.Err.Code, diceResp.Err.Message, diceResp.Err.Ctx)
	}

	logrus.Infof(">>> success to cancel dice deploy.")

	return nil
}
