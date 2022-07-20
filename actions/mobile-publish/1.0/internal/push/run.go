package push

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-actions/actions/mobile-publish/1.0/internal/conf"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/envconf"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/erda-project/erda/pkg/metadata"
)

func Run() error {
	var cfg conf.Conf
	if err := envconf.Load(&cfg); err != nil {
		return err
	}

	release, err := GetRelease(cfg, cfg.ReleaseID)
	if err != nil {
		return err
	}
	relations, err := GetAppPublishItemRelations(cfg)
	if err != nil {
		return err
	}
	relation, ok := relations.Data[cfg.Workspace]
	if !ok {
		return fmt.Errorf("not found env config %s", cfg.Workspace)
	}
	readme := ""
	spec := ""
	if cfg.ReadmeFile != "" {
		readmeBytes, err := ioutil.ReadFile(cfg.ReadmeFile)
		if err != nil {
			return fmt.Errorf("read readme file err: %s", err)
		}
		readme = string(readmeBytes)
	}
	if cfg.SpecFile != "" {
		specBytes, err := ioutil.ReadFile(cfg.SpecFile)
		if err != nil {
			return fmt.Errorf("read spec file err: %s", err)
		}
		spec = string(specBytes)
	}
	var (
		version, buildID, packageName string
		resourceType                  apistructs.ResourceType
		h5VersionInfo                 apistructs.H5VersionInfo
	)
	// android和ios区分以后，release 时保证了 resource 只有一种，如果出现多个则是 release 时出现了bug
	for _, v := range release.Data.Resources {
		if isMobileResource(v.Type) {
			resourceType = v.Type
		}
		if _, ok := v.Meta["h5VersionInfo"]; ok {
			decodeBytes, err := base64.StdEncoding.DecodeString(v.Meta["h5VersionInfo"].(string))
			if err != nil {
				return err
			}
			if err = json.Unmarshal(decodeBytes, &h5VersionInfo); err != nil {
				return err
			}
		}
		if _, ok := v.Meta["buildID"]; ok {
			buildID = v.Meta["buildID"].(string)
		}
		if _, ok := v.Meta["version"]; ok {
			version = v.Meta["version"].(string)
		}
		if _, ok := v.Meta["packageName"]; ok {
			packageName = v.Meta["packageName"].(string)
		}
	}
	pushRequest := apistructs.CreatePublishItemVersionRequest{
		PackageName:   packageName,
		Version:       version,
		BuildID:       buildID,
		Public:        false,
		IsDefault:     false,
		Logo:          "",
		Desc:          "",
		Readme:        readme,
		Spec:          spec,
		ReleaseID:     release.Data.ReleaseID,
		MobileType:    resourceType,
		H5VersionInfo: h5VersionInfo,
		PublishItemID: relation.PublishItemID,
		Creator:       cfg.UserID,
	}
	logrus.Infof("publish release version:%s", pushRequest.Version)
	publishItemResponse, err := CreatePublishItemVersion(cfg, pushRequest)
	if err != nil {
		return err
	}

	// write metafile
	metaInfos := make([]metadata.MetadataField, 0, 1)
	metaInfos = append(metaInfos, metadata.MetadataField{
		Name:  apistructs.ActionCallbackPublisherID,
		Value: strconv.FormatUint(uint64(relation.PublisherID), 10),
	})
	metaInfos = append(metaInfos, metadata.MetadataField{
		Name:  apistructs.ActionCallbackPublishItemID,
		Value: strconv.FormatUint(uint64(relation.PublishItemID), 10),
	})
	metaInfos = append(metaInfos, metadata.MetadataField{
		Name:  apistructs.ActionCallbackPublishItemVersionID,
		Value: strconv.FormatUint(publishItemResponse.Data.ID, 10),
	})
	metaInfos = append(metaInfos, metadata.MetadataField{
		Name:  "type",
		Value: apistructs.PublishItemTypeMobile,
	})

	metaByte, _ := json.Marshal(apistructs.ActionCallback{Metadata: metaInfos})
	if err = filehelper.CreateFile(cfg.Metafile, string(metaByte), 0644); err != nil {
		logrus.Warnf("failed to write metafile, %v", err)
	}

	return nil
}

func isMobileResource(resourceType apistructs.ResourceType) bool {
	switch resourceType {
	case apistructs.ResourceTypeAndroid, apistructs.ResourceTypeIOS,
		apistructs.ResourceTypeH5, apistructs.ResourceTypeAndroidAppBundle:
		return true
	default:
		return false
	}
}
