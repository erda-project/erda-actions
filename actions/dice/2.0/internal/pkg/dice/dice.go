package dice

import (
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/pkg/log"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/dice/deploy"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/dice/store"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/dice/callback"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/cancel"
)

func Run() error {
	// parse config
	cfg, err := conf.HandleConf()
	if err != nil {
		logrus.Errorf("failed to handle conf, err: %v", err)
		return err
	}

	// init
	s := store.New(store.WithConf(&cfg))
	d := deploy.New(
		deploy.WithConf(&cfg),
		deploy.WithStore(s),
	)

	// do deploy
	orderId, rets, err := d.Do()
	if err != nil {
		return errors.Wrap(err, "deploy failed")
	}

	// store dice info
	// default (drive by pipeline build release): deploymentId=uint64,runtimeId=uint64
	// application or project release (deploy release): applicationName_deploymentId=uint64,applicationName_runtimeId=uint64
	if err := s.StoreDiceInfo(orderId, rets); err != nil {
		logrus.Warning(err)
	}

	// report runtime info
	if err := callback.BatchReportRuntimeInfo(&cfg, rets); err != nil {
		logrus.Errorf("failed to batch report runtime info, err: %v", err)
	}

	// tips deploy info, link to current runtime last deployment log
	log.AddNewLine(1)
	// add msg because of frontend will match regularly, the regular is msg=\"(.+)\"
	for k, r := range rets {
		logrus.Infof("msg=\"application %s deploy starting... ##to_link:applicationId:%d,runtimeId:%d,deploymentId:%d\"",
			k, r.ApplicationId, r.RuntimeId, r.DeploymentId)
	}

	// check deployments status
	// TODO: check deployment order status instead
	if err := d.StatusCheck(rets, cfg.TimeOut); err != nil {
		logrus.Errorf("failed to check status, err: %v", err)
		cancel.ExeWithConfig(orderId, &cfg)
		return err
	}
	return nil
}
