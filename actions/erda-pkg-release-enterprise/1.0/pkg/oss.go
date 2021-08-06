package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (

	// const related to archive and release
	OssArchiveBucket         = "erda-release"
	OssArchivePath           = "archived-versions"
	OssPkgReleaseBucket      = "erda-release"
	OssPkgReleasePublicPath  = "erda"
	OssPkgReleasePrivatePath = "enterprise"

	// const related to oss acl
	//OssAclPublicReadWrite = "public-read-write"
	OssAclPublicRead = "public-read"
	//OssAclPrivate         = "private"
	//OssAclDefault         = "default"

	// const related to release type
	// just common release, has not test completely
	ReleaseCommon = "common"

	// released by sre, just released to test tools, avoid erda's templates reflecting
	ReleaseTools = "tools"

	// pkg released with ReleaseOffline used for erda cluster without internet
	ReleaseOffline = "offline"

	// go through completely testing, the quality of release pkg is guaranteed
	ReleaseCompletely = "completely"
)

// Oss object contains oss info used to auth with oss server
type Oss struct {
	OssEndPoint        string `json:"endpoint"`
	OssAccessKeyId     string `json:"accessKeyID"`
	OssAccessKeySecret string `json:"accessKeySecret"`
}

// OSS object release to erda-pkg-release action
type OSS struct {
	// oss Oss info, use to auth with oss
	oss *Oss

	// erdaVersion version of erda released package
	erdaVersion string

	// archiveBucket bucket in oss which erda archive to
	archiveBucket string

	// archiveBasePath base archive path in oss which erda archive to
	archiveBasePath string

	// releaseBucket bucket in oss which erda pkg release to
	releaseBucket string

	// releaseType type of erda release pkg, reference to ReleaseCommon | ReleaseTools | ReleaseCompletely | ReleaseOffline
	// reflecting path in oss of erda release package
	releaseType string

	// policy of release pkg, decide if osArch type as a dir
	osArch bool

	// actionReleaseBasePath related to erda-pkg-release-* action
	// judge whether the erda package release path is enterprise or erda
	actionReleaseBasePath string
}

// NewOSS get OSS object
func NewOSS(o *Oss, erdaVersion, releaseType, actionPath string, osArch bool) *OSS {
	return &OSS{
		oss:                   o,
		erdaVersion:           erdaVersion,
		archiveBucket:         OssArchiveBucket,
		archiveBasePath:       OssArchivePath,
		releaseBucket:         OssPkgReleaseBucket,
		releaseType:           releaseType,
		osArch:                osArch,
		actionReleaseBasePath: actionPath,
	}
}

// GetOss get oss info in object OSS
func (o *OSS) Oss() *Oss {
	return o.oss
}

// ErdaVersion get erdaVersion info in OSS
func (o *OSS) ErdaVersion() string {
	return o.erdaVersion
}

// ErdaVersion get releaseBucket info in OSS
func (o *OSS) ReleaseBucket() string {
	return o.releaseBucket
}

// ErdaVersion get releaseType info in OSS
func (o *OSS) ReleaseType() string {
	return o.releaseType
}

// ErdaVersion get osArch info in OSS
func (o *OSS) OsArch() bool {
	return o.osArch
}

// InitOssConfig init oss client's config in erda-pkg-release-* action
func (o *OSS) InitOssConfig() error {
	return o.oss.InitOssConfig()
}

func (o *OSS) GenArchivePath() string {
	return fmt.Sprintf("%s/%s", o.archiveBasePath, o.erdaVersion)
}

// GenReleasePath generate base release path
func (o *OSS) GenReleasePath(osArch, path string) string {
	// policy of release pkg, decide if osArch type as a dir
	if o.osArch {
		return fmt.Sprintf("%s/%s/%s/%s", o.releaseType, o.actionReleaseBasePath, osArch, path)
	}

	return fmt.Sprintf("%s/%s/%s", o.releaseType, o.actionReleaseBasePath, path)
}

// GenReleaseUrl generate erda release pkg url which can be used to get erda release package when access in browser
func (o *OSS) GenReleaseUrl(osArch, pkg string) string {
	return fmt.Sprintf("http://%s.%s/%s",
		o.releaseBucket, o.oss.OssEndPoint, o.GenReleasePath(osArch, pkg))
}

// ReleasePackage push erda release pkg to oss
func (o *OSS) ReleasePackage(releasePathInfo map[string]string) error {

	// upload release installing pkg of erda
	for osArch, pkgPath := range releasePathInfo {
		if !path.IsAbs(pkgPath) {
			return errors.Errorf("release pkg path is "+
				"not a absolute path: %s", pkgPath)
		}

		_, pkgName := path.Split(pkgPath)
		ossReleasePath := o.GenReleasePath(osArch, pkgName)

		if err := o.oss.UploadFile(pkgPath, o.releaseBucket, ossReleasePath, OssAclPublicRead); err != nil {
			return err
		}
	}

	return nil
}

// PreparePatchRelease prepare erda release info to action
func (o *OSS) PreparePatchRelease() error {

	// download release from oss
	archivePath := o.GenArchivePath()
	if err := o.oss.DownloadDir("/tmp", o.archiveBucket, archivePath); err != nil {
		return errors.WithMessage(err, "cp release patch to /tmp/")
	}

	tars := []string{
		"erda-actions-enterprise.tar.gz",
		"erda-actions.tar.gz",
		"erda-addons-enterprise.tar.gz",
		"erda-addons.tar.gz",
	}

	// tar release
	for _, tar := range tars {
		if _, err := ExecCmd(os.Stdout, os.Stderr, fmt.Sprintf("/tmp/%s/extensions", o.erdaVersion),
			"tar", "-zxvf", tar); err != nil {
			return errors.WithMessage(err, fmt.Sprintf("decompress %s failed", tar))
		}
	}

	return nil
}

// OssRemotePath generate oss remote access path
func (o *Oss) OssRemotePath(bucket, path string) string {

	return fmt.Sprintf("oss://%s/%s", bucket, path)
}

// UploadFile upload file to oss with --force parameter
func (o *Oss) UploadFile(local, bucket, path, acl string) error {

	exists, err := FileExist(local)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("upload file %s to oss", local))
	}
	if !exists {
		return fmt.Errorf("the file %s waited to upload is not exists", local)
	}

	remote := o.OssRemotePath(bucket, path)

	if acl == "" {
		_, err = ExecCmd(os.Stdout, os.Stderr, "", "ossutil64", "cp", "-f", local, remote)
	} else {
		_, err = ExecCmd(os.Stdout, os.Stderr, "", "ossutil64",
			"cp", "-f", fmt.Sprintf("--acl=%s", acl), local, remote)
	}
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("upload file %s to %s failed", local, remote))
	}

	logrus.Infof("upload file %s to %s success", local, remote)
	return nil
}

// UploadDir upload dir to oss with --force parameter
func (o *Oss) UploadDir(dir, bucket, path, acl string) error {

	exists, err := IsDirExists(dir)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("upload dir %s to oss", dir))
	}
	if !exists {
		return fmt.Errorf("the dir %s waited to upload is not exists", dir)
	}

	remote := o.OssRemotePath(bucket, path)

	if acl == "" {
		_, err = ExecCmd(os.Stdout, os.Stderr, "", "ossutil64", "cp", "-rf", dir, path)
	} else {
		_, err = ExecCmd(os.Stdout, os.Stderr, "", "ossutil64", "cp", "-rf",
			fmt.Sprintf("--acl=%s", acl), dir, path)
	}
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("upload file %s to %s failed", dir, remote))
	}

	logrus.Infof("upload file %s to %s success", dir, remote)
	return nil
}

// DownloadFile download file to local with --force parameter
func (o *Oss) DownloadFile(local, bucket, path string) error {
	remote := o.OssRemotePath(bucket, path)

	_, err := ExecCmd(os.Stdout, os.Stderr, "", "ossutil64", "cp", "-f", remote, local)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("download file %s to %s failed", remote, local))
	}

	logrus.Infof("download file %s to %s success", remote, local)
	return nil
}

// DownloadDir download dir to local with --force parameter
func (o *Oss) DownloadDir(parent, bucket, path string) error {
	remote := o.OssRemotePath(bucket, path)

	_, err := ExecCmd(os.Stdout, os.Stderr, "", "ossutil64", "cp", "-rf", remote, parent)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("download dir %s from to %s failed", remote, parent))
	}

	logrus.Infof("download dir %s to %s success", remote, parent)
	return nil
}

// InitOssConfig init oss client's config to local
func (o *Oss) InitOssConfig() error {

	logrus.Info("start to init oss config...")
	// current user in action
	u, err := user.Current()
	if err != nil {
		return errors.WithMessage(err, "get current user when init oss config")
	}

	// oss config path
	home := u.HomeDir
	ossConfigPath := path.Join(home, ".ossutilconfig")

	// oss config
	ossConfig := fmt.Sprintf("[Credentials]\nlanguage=CH\nendpoint=%s\naccessKeyID="+
		"%s\naccessKeySecret=%s", o.OssEndPoint, o.OssAccessKeyId, o.OssAccessKeySecret)
	if err := ioutil.WriteFile(ossConfigPath, []byte(ossConfig), 0666); err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Info("init oss config success!!")

	return nil
}
