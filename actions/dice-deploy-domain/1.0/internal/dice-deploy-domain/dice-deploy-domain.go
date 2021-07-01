package dice_deploy_domains

import (
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/dice-deploy-service/1.0/dice-deploy-services"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/http/httpclient"
	"github.com/erda-project/erda/pkg/retry"
)

func Run() error {
	var cfg dice_deploy_services.Conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}
	logrus.Info(cfg)
	d := &dice_deploy_services.Dice{Conf: &cfg}
	result, err := Deploy(&cfg)
	if err != nil {
		return errors.Wrap(err, "deploy dice failed")
	}

	//Set default deployment timeout is 24h.
	timeout := cfg.TimeOut
	minTimeoutSec := (60 * 60) * 24
	if timeout < minTimeoutSec {
		timeout = minTimeoutSec
	}
	_, err = dice_deploy_services.CheckDeploymentLoop(d, result, "", time.Duration(timeout), Check)
	deployResult, Deployerr := dice_deploy_services.GetDeploymentStatus(result, &cfg)
	if Deployerr != nil {
		return Deployerr
	}
	if err != nil {
		dice_deploy_services.StoreMetaFileWithErr(&cfg, result.RuntimeId, result.DeploymentId, deployResult)
		return err
	}
	logrus.Infof("checkDeploymentLoop end storeMetaFile")
	return dice_deploy_services.StoreMetaFileWithErr(&cfg, result.RuntimeId, result.DeploymentId, deployResult)
}

func Deploy(conf *dice_deploy_services.Conf) (*dice_deploy_services.DeployResult, error) {
	var diceResp dice_deploy_services.DiceResponse
	err := retry.DoWithInterval(func() error {
		r, err := httpclient.New(httpclient.WithCompleteRedirect()).Post(conf.DiceOpenapiPrefix).Path(fmt.Sprintf("/api/deployments/%s/actions/deploy-domains", conf.DeploymentID)).
			Header(dice_deploy_services.Authorization, conf.DiceOpenapiToken).Do().JSON(&diceResp)
		if err != nil {
			return err
		}
		if !r.IsOK() {
			return errors.Errorf("create a dice deploy failed, statusCode: %d, diceResp:%+v",
				r.StatusCode(), diceResp)
		}

		if !diceResp.Success {
			return errors.Errorf("create dice deploy failed. code=%s, message=%s, ctx=%v",
				diceResp.Err.Code, diceResp.Err.Message, diceResp.Err.Ctx)
		}

		return nil
	}, 5, time.Second*3)
	if err != nil {
		logrus.Errorf("deploy to dice failed! response err:%v.", err)
		return nil, err
	}
	result := dice_deploy_services.DeployResult{}
	if err := mapstructure.Decode(diceResp.Data, &result); err != nil {
		return nil, errors.Wrapf(err, "mapstructure data=%+v", result)
	}

	return &result, nil
}

func Check(res *dice_deploy_services.DeployResult, conf *dice_deploy_services.Conf) (deploying bool, runtime interface{}, e error) {
	defer func() {
		if deploying && e == nil {
			Deploy(conf)
		}
	}()
	result, err := dice_deploy_services.GetDeploymentStatus(res, conf)
	if err != nil {
		deploying = false
		runtime = nil
		e = err
		return
	}
	if len(result.Data.MoudleErrMsg) > 0 {
		dice_deploy_services.StoreMetaFileWithErr(conf, res.RuntimeId, res.DeploymentId, result)
	}
	switch result.Data.Status {
	case "WAITING", "WAITAPPROVE", "INIT":
		deploying = true
		runtime = nil
		e = nil
		return
	case "DEPLOYING":
		switch result.Data.Phase {
		case "INIT", "ADDON_REQUESTING", "SCRIPT_APPLYING", "SERVICE_DEPLOYING", "DISCOVERY_REGISTER":
			logrus.Infof("continue deploying..., ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
				res.ApplicationId, res.RuntimeId, res.DeploymentId)
			deploying = true
			runtime = nil
			e = nil
			return
		case "COMPLETED":
			logrus.Info("deploy addons success")
			deploying = false
			runtime = result.Data.Runtime
			e = nil
			return
		}
		logrus.Infof("continue deploying..., ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d",
			res.ApplicationId, res.RuntimeId, res.DeploymentId)
		deploying = true
		runtime = nil
		e = nil
		return
	case "OK":
		logrus.Info("deploy success")
		deploying = false
		runtime = result.Data.Runtime
		e = nil
		return
	case "CANCELED":
		deploying = false
		runtime = nil
		e = &dice_deploy_services.DiceDeployError{"deployment canceled by dice"}
		return
	case "FAILED":
		deploying = false
		runtime = nil
		e = &dice_deploy_services.DiceDeployError{"deployment failed in dice, " + result.Data.FailCause}
		return
	}
	deploying = false
	runtime = nil
	e = errors.Errorf("deployment unknown %s in dice", result.Data.Status)
	return
}
