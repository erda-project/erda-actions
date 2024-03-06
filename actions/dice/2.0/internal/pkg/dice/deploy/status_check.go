package deploy

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda/pkg/http/httpclient"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/dice/callback"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
)

var (
	ErrCheckTimeout = errors.New("deployments status check timeout")
)

var (
	defaultTimeout       = (60 * 60) * 24
	defaultCheckInterval = 10 * time.Second
)

var (
	statusResult  = make(map[string]*common.DeploymentStatusRespData)
	lastStatusMap = sync.Map{}
)

func (d *deploy) StatusCheck(orderId string, result map[string]*common.DeployResult, timeout int) error {
	// Set default deployment timeout is 24h.
	minTimeoutSec := defaultTimeout
	if timeout < minTimeoutSec {
		timeout = minTimeoutSec
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	releaseType := utils.ConvertType(d.cfg.ReleaseTye)
	switch releaseType {
	case common.TypeProjectRelease:
		if err := d.checkBatchStatusLoop(ctx, orderId); err != nil {
			if err == ErrCheckTimeout {
				ErrTimeOutPrint()
				return nil
			}
			return err
		}
		return nil
	default:
		// loop check deployments status
		if err := d.statusCheckLoop(ctx, result); err != nil {
			if err == ErrCheckTimeout {
				ErrTimeOutPrint()
				return nil
			}
			return err
		}

		// batch execute callback
		d.exeCallBack(result)
		return d.store.BatchStoreMetaFile(statusResult)
	}
}

func (d *deploy) statusCheckLoop(ctx context.Context, drMap map[string]*common.DeployResult) error {
	ticker := time.NewTicker(defaultCheckInterval)
	for {
		select {
		case <-ctx.Done():
			return ErrCheckTimeout
		case <-ticker.C:
			isDeploying, err := d.check(drMap)
			if err != nil {
				logrus.Errorf("failed to check deploy %v", err)
				return err
			}
			// deploy done
			if !isDeploying {
				return nil
			}
			ticker.Reset(defaultCheckInterval)
		}
	}
}

func (d *deploy) checkBatchStatusLoop(ctx context.Context, orderId string) error {
	ticker := time.NewTicker(defaultCheckInterval)
	for {
		select {
		case <-ctx.Done():
			return ErrCheckTimeout
		case <-ticker.C:
			isDeploying, err := d.checkBatch(orderId)
			if err != nil {
				logrus.Errorf("failed to check deploy %v", err)
				return err
			}
			// deploy done
			if !isDeploying {
				return nil
			}
			ticker.Reset(defaultCheckInterval)
		}
	}
}

func (d *deploy) checkBatch(orderId string) (bool, error) {
	var err error
	resp, err := d.getDeploymentOrderStatus(orderId)
	if err != nil {
		logrus.Errorf("failed to get deployment order status, orderId: %s, err: %s", orderId, err)
		return false, err
	}

	resp.Print()
	if resp.Data.Status == "FAILED" {
		err = fmt.Errorf("get failed status, please check runtime log at deploy center")
	}
	return continueDeploying(resp.Data.Status), err
}

func (d *deploy) check(drMap map[string]*common.DeployResult) (bool, error) {
	// batch check deployment status
	var wg sync.WaitGroup
	wg.Add(len(drMap))

	for appName, dr := range drMap {
		// execute check deployment status
		go func(appName string, dr *common.DeployResult) {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("get application %s deployment status panic: %v", appName, err)
				}
			}()
			defer wg.Done()

			var dStatus *common.DeploymentStatusRespData
			if dr != nil {
				var err error
				// get deployment status by deployment id
				dStatus, err = d.getDeploymentStatus(dr.DeploymentId)
				if err != nil {
					logrus.Errorf("failed to get deployment status, deployment id: %d, err: %v",
						dr.DeploymentId, err)
				}
			}

			d.statusLock.Lock()
			defer d.statusLock.Unlock()
			// get deployment error will cause status == nil
			statusResult[appName] = dStatus
		}(appName, dr)
	}

	wg.Wait()

	var (
		needPrint, isErr, isDeploying bool
	)

	defer func() {
		// error message will store meta and report
		if isErr {
			err := d.store.BatchStoreMetaFile(statusResult)
			if err != nil {
				logrus.Errorf("failed to batch store meta file, err: %v", err)
			}
		}
	}()

	// analyse status check results
	for appName, data := range statusResult {
		if data == nil {
			logrus.Errorf("failed to check application %s status, deployment status is nil", appName)
			// reset status to deploying, recheck next loop
			needPrint, isDeploying = true, true
			continue
		}

		// update application last status
		status, _ := lastStatusMap.Load(appName)
		if data.Data.Status != status {
			needPrint = true
		}
		lastStatusMap.Store(appName, data.Data.Status)

		// check err message
		if len(data.Err.Message) != 0 {
			isErr = true
			needPrint = true
		}

		// parse deployment status from response data
		tmpIsDeploying, err := parseStatus(data)
		if err != nil {
			logrus.Debug(err)
			isErr = true
			return isDeploying, err
		}

		// if one of applications is deploying, group status is deploying
		if tmpIsDeploying {
			isDeploying = true
		}
	}

	// one of application deploy status change or process had error message
	if needPrint {
		utils.BatchPrintStatusCheckResult(statusResult)
	}

	return isDeploying, nil
}

func parseStatus(resp *common.DeploymentStatusRespData) (bool, error) {
	var err error
	switch resp.Data.Status {
	case "CANCELED":
		err = &common.DeployErrResponse{
			Msg: "deployment canceled by dice",
		}
	case "FAILED":
		err = &common.DeployErrResponse{
			Msg: "deployment failed in dice, " + resp.Data.FailCause,
		}
	}
	return continueDeploying(resp.Data.Status), err
}

func continueDeploying(status string) bool {
	switch status {
	case "WAITING", "WAITAPPROVE", "INIT":
		return true
	case "DEPLOYING":
		return true
	case "OK":
		return false
	case "CANCELED":
		return false
	case "FAILED":
		return false
	default:
		return false
	}
}

func (d *deploy) getDeploymentStatus(deploymentId uint64) (*common.DeploymentStatusRespData, error) {
	var result common.DeploymentStatusRespData
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(d.cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/deployments/%d/status", deploymentId)).
		Header("Authorization", d.cfg.DiceOpenapiToken).Do().JSON(&result)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() {
		return nil, errors.Errorf("deploy to dice failed, statusCode: %d", r.StatusCode())
	}
	if !result.Success {
		return nil, errors.Errorf("failed to request dice deploy status. code=%s, message=%s, ctx=%v",
			result.Err.Code, result.Err.Message, result.Err.Ctx)
	}
	return &result, nil
}

func (d *deploy) getDeploymentOrderStatus(orderId string) (*common.DeploymentOrderStatusRespData, error) {
	var result common.DeploymentOrderStatusRespData
	r, err := httpclient.New(httpclient.WithCompleteRedirect()).Get(d.cfg.DiceOpenapiPrefix).
		Path(fmt.Sprintf("/api/deployment-orders/%s", orderId)).
		Header("Authorization", d.cfg.DiceOpenapiToken).Do().JSON(&result)
	if err != nil {
		return nil, err
	}
	if !r.IsOK() {
		return nil, errors.Errorf("deploy to dice failed, statusCode: %d", r.StatusCode())
	}
	if !result.Success {
		return nil, errors.Errorf("failed to request deploy order status. code=%s, message=%s, ctx=%v",
			result.Err.Code, result.Err.Message, result.Err.Ctx)
	}
	return &result, nil
}

func (d *deploy) exeCallBack(result map[string]*common.DeployResult) {
	cbTarget := d.cfg.Callback

	if cbTarget == "" {
		logrus.Info("no callback set, return directly")
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(result))

	for appName, dr := range result {
		go func(appName string, dr *common.DeployResult) {
			defer func() {
				if err := recover(); err != nil {
					logrus.Errorf("execute application %s callback panic, err: %v", appName, err)
				}
			}()
			defer wg.Done()

			var (
				runtime interface{}
				status  = common.CallbackStatusSuccess
			)

			sr, ok := statusResult[appName]
			if ok {
				if _, err := parseStatus(sr); err != nil {
					status = common.CallbackStatusFailed
				} else {
					runtime = sr.Data.Runtime
				}
			} else {
				status = common.CallbackStatusFailed
			}

			callback.Callback(cbTarget, dr.RuntimeId, dr.ApplicationId, runtime, status)
		}(appName, dr)
	}

	wg.Wait()
}

func ErrTimeOutPrint() {
	logrus.Infof("Deploying timeout, you can: ")
	logrus.Infof("   1. increase timeout in pipeline.yml")
	logrus.Infof("   2. try again ")
}
