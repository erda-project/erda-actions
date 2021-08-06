package pkg

import (
	"encoding/json"
	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/pkg/filehelper"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// metafile keys
const (
	MetaErdaVersion = "erdaVersion"
	erdaPkgMapUrl   = "pkgMapUrl"

	// json string
	MetaReleaseInfoType = "ToolsPkgReleaseInfo"
)

// metafile object to operate metafile
type metafile struct {
	oss         *OSS
	metafile    string
	erdaVersion string
}

// NewMetafile get an new metafile object
func NewMetafile(oss *OSS, metafilePath string) *metafile {
	return &metafile{
		oss:         oss,
		metafile:    metafilePath,
		erdaVersion: oss.erdaVersion,
	}
}

// GetOss get oss info of metafile object
func (m *metafile) GetOss() *OSS {
	return m.oss
}

// GetMetafile get metafile path
func (m *metafile) GetMetafile() string {
	return m.metafile
}

// ErdaVersion get erda version
func (m *metafile) ErdaVersion() string {
	return m.erdaVersion
}

// WriteMetaFile write metafile
func (m *metafile) WriteMetaFile(releaseInfo map[string]string) error {

	logrus.Infof("start to write metafile")

	metaInfos := make([]apistructs.MetadataField, 0, 1)

	urlInfo := map[string]string{}

	// generate erda release package url info
	for osArch, pkg := range releaseInfo {
		urlInfo[osArch] = m.oss.GenReleaseUrl(osArch, pkg)
	}

	// serialize url info
	vJson, err := json.Marshal(urlInfo)
	if err != nil {
		return errors.WithMessage(err, "change release pkg url info to json string")
	}

	sJson := string(vJson)

	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  MetaErdaVersion,
		Value: m.erdaVersion,
	})
	metaInfos = append(metaInfos, apistructs.MetadataField{
		Name:  erdaPkgMapUrl,
		Type:  MetaReleaseInfoType,
		Value: sJson,
	})

	metaByte, _ := json.Marshal(apistructs.ActionCallback{Metadata: metaInfos})
	if err := filehelper.CreateFile(m.metafile, string(metaByte), 0644); err != nil {
		logrus.Warnf("failed to write metafile, %v", err)
	}

	logrus.Infof("write metafile success...")

	return nil
}
