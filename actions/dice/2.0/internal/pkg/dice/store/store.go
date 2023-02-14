package store

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/common"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/conf"
	"github.com/erda-project/erda-actions/actions/dice/2.0/internal/pkg/utils"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/metadata"
)

type Store interface {
	StoreDiceInfo(orderId string, dr map[string]*common.DeployResult) error
	BatchStoreMetaFile(statusResp map[string]*common.DeploymentStatusRespData) error
}

type Option func(s *store)

type store struct {
	cfg *conf.Conf
}

func New(opts ...Option) Store {
	s := store{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

func WithConf(c *conf.Conf) Option {
	return func(s *store) {
		s.cfg = c
	}
}

func (s *store) StoreDiceInfo(orderId string, dr map[string]*common.DeployResult) error {
	if len(dr) == 0 {
		return errors.New("no dice info need store")
	}

	content := []string{
		fmt.Sprintf("deployment_order_id=%s", orderId),
	}

	for k, v := range dr {
		if len(dr) == 1 {
			content = append(content, fmt.Sprintf("deploymentId=%d", v.DeploymentId))
			content = append(content, fmt.Sprintf("runtimeId=%d", v.RuntimeId))
			break
		}
		content = append(content, fmt.Sprintf("%s_deploymentId=%d", k, v.DeploymentId))
		content = append(content, fmt.Sprintf("%s_runtimeId=%d", k, v.RuntimeId))
	}

	err := utils.CreateFile(filepath.Join(s.cfg.WorkDir, "diceInfo"), strings.Join(content, ","), 0755)
	if err != nil {
		return errors.Wrap(err, "write file:diceInfo failed")
	}
	return nil
}

func (s *store) BatchStoreMetaFile(statusResp map[string]*common.DeploymentStatusRespData) error {
	if statusResp == nil {
		return s.storeMetaFile(metadata.Metadata{})
	}

	metaData := []metadata.MetadataField{
		{Name: "projectID", Value: strconv.FormatUint(s.cfg.ProjectID, 10)},
	}

	for _, resp := range statusResp {
		switch utils.ConvertType(s.cfg.ReleaseTye) {
		case common.TypeProjectRelease, common.TypeApplicationRelease:
			if len(resp.Data.ModuleErrMsg) == 0 {
				continue
			}
			for k, v := range resp.Data.ModuleErrMsg {
				if v == "" {
					continue
				}
				metaData = append(metaData, metadata.MetadataField{
					Name:  k,
					Value: v,
				})
			}
		default:
			metaData = append(metaData, metadata.Metadata{
				{Name: "appID", Value: strconv.FormatUint(s.cfg.AppID, 10)},
				{Name: "deploymentID", Value: strconv.Itoa(resp.Data.DeploymentId)},
			}...)
			break
		}
	}

	return s.storeMetaFile(metaData)
}

func (s *store) storeMetaFile(metaData metadata.Metadata) error {
	meta := apistructs.ActionCallback{
		Metadata: metaData,
	}
	b, err := json.Marshal(&meta)
	if err != nil {
		return err
	}
	if err := utils.CreateFile(s.cfg.MetaFile, string(b), 0644); err != nil {
		return errors.Wrap(err, "write file:metafile failed")
	}
	return nil
}
