package build

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-actions/actions/lib-publish/1.0/internal/pkg/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
)

const (
	specYmlFile string = "spec.yml"
	readmeFile  string = "README.md"
)

func Run() error {
	// 1. 根据 appID 获取 publisherItemID
	// 2. 创建 publisherItemVersion
	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}

	// 切换工作目录
	if cfg.Context != "" {
		fmt.Fprintf(os.Stdout, "change workding directory to: %s\n", cfg.Context)
		if err := os.Chdir(cfg.Context); err != nil {
			return err
		}
	}

	// 检查 spec.yml & README.md 文件是否存在
	readmeContent, err := ioutil.ReadFile(readmeFile)
	if err != nil {
		readmeContent, err = ioutil.ReadFile(strings.ToLower(readmeFile))
		if err != nil {
			return err
		}
	}
	specContent, err := ioutil.ReadFile(specYmlFile)
	if err != nil {
		return err
	}

	// 检查 spec.yml 文件结构
	var spec apistructs.Spec
	if err := yaml.Unmarshal(specContent, &spec); err != nil {
		return err
	}
	if !checkSpecValid(spec) {
		return errors.Errorf("invalid format %v", specYmlFile)
	}

	relations, err := GetAppPublishItemRelations(cfg)
	if err != nil {
		return err
	}
	relation, ok := relations.Data[cfg.Workspace]
	if !ok {
		return fmt.Errorf("not found env config %s", cfg.Workspace)
	}

	pushRequest := apistructs.CreatePublishItemVersionRequest{
		Version:       spec.Version,
		PublishItemID: relation.PublishItemID,
		AppID:         cfg.AppID,
		Creator:       cfg.UserID,
		Readme:        string(readmeContent),
		Spec:          string(specContent),
		Public:        spec.Public,
		IsDefault:     spec.IsDefault,
	}
	publishItemResponse, err := CreatePublishItemVersion(cfg, pushRequest)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "library %s publish succ, version: %s\n", spec.Name, spec.Version)

	metaInfos := make([]apistructs.MetadataField, 0, 1)
	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  apistructs.ActionCallbackPublisherID,
		Value: strconv.FormatUint(uint64(relation.PublisherID), 10),
	})
	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  apistructs.ActionCallbackPublishItemID,
		Value: strconv.FormatUint(uint64(relation.PublishItemID), 10),
	})
	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  apistructs.ActionCallbackPublishItemVersionID,
		Value: strconv.FormatUint(publishItemResponse.Data.ID, 10),
	})
	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  "type",
		Value: apistructs.PublishItemTypeLIBRARY,
	})

	metaByte, _ := json.Marshal(apistructs.ActionCallback{Metadata: metaInfos})
	if err = filehelper.CreateFile(cfg.Metafile, string(metaByte), 0644); err != nil {
		logrus.Warnf("failed to write metafile, %v", err)
	}
	return nil
}

func checkSpecValid(spec apistructs.Spec) bool {
	if spec.Name == "" || spec.Type == "" || spec.Version == "" {
		return false
	}
	return true
}
