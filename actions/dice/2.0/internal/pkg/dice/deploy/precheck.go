package deploy

import (
	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
)

func paramsPreCheck(c *conf.Conf) error {
	if c.AssignedWorkspace != "" && !utils.VerifyWorkspace(c.AssignedWorkspace) {
		return errors.New("invalid assigned workspace")
	}

	if c.ReleaseTye != "" && c.AssignedWorkspace == "" {
		return errors.New("workspace is required when release type is specified.")
	}

	if c.ReleaseID != "" && c.ReleaseName != "" {
		return errors.New("only one of release id or release name can be specified.")
	}

	if c.ReleaseName != "" && c.ReleaseTye == "" {
		return errors.New("release type is required when release name is specified.")
	}

	if utils.ConvertType(c.ReleaseTye) == common.TypeApplicationRelease && c.ApplicationName == "" {
		return errors.New("application name is required when release type is application.")
	}

	return nil
}
